#!/usr/bin/env python

import os
import sys
import subprocess

from collections import namedtuple

from . import helper

OPERATIONS = [
    "compile", "setup", "witness", "prove", "verify"
]


def build_command_gnark(payload, count):
    """
    Build the command to invoke the gnark ZKP-framework given the payload
    """
    initial_cmd = f"cd {helper.Paths().GNARK_DIR} && "
    
    os.makedirs(helper.Paths().GNARK_BENCH, exist_ok=True)    
    os.makedirs(helper.Paths().GNARK_BENCH_MEMORY, exist_ok=True)    
    if payload.backend is not None and payload.curves is not None:

        # FIXME: Update that code when implementing more fine-grained benchmarking
        # in gnark
        commands = [f"{initial_cmd} ./gnark {backend} --circuit={circ} --algo={op} --curve={curve} --input={inp} --count=1 --outputPath={helper.Paths().GNARK_BENCH}/{backend}_{circ}.csv; \n"
                    for backend in payload.backend
                    for curve in payload.curves
                    for circ, input_path in payload.circuit.items()
                    for inp in helper.get_all_input_files(input_path)
                    for op in payload.operation
                    for _ in range(0,count)]

        # Builder command memory
        command_binary = f"{initial_cmd} ./build_memory.sh;"
        # Create /tmp folder if non-existent
        command_check_tmp = f"{initial_cmd} mkdir -p ./tmp;"

        # Memory commands
        commands_memory = [
            (
                os.makedirs(f"{helper.Paths().GNARK_BENCH_MEMORY}/{modified_inp}", exist_ok=True),
                f"{initial_cmd} {helper.get_memory_command()} ./{backend}_memory_{op} \
                    --circuit={circ} \
                    --curve={curve} \
                    --input={inp} \
                    --count=1 \
                    2> {helper.Paths().GNARK_BENCH_MEMORY}/{modified_inp}/gnark_{backend}_{circ}_memory_{op}.txt \
                    > /dev/null || true; \n"
            )[1]
            for backend in payload.backend
            for curve in payload.curves
            for circ, input_path in payload.circuit.items()
            for inp in helper.get_all_input_files(input_path)
            for modified_inp in [inp.replace('input/circuit/', '').replace('.json', '')]
            for op in payload.operation
        ]

        commands_memory.append("cd ../../;")

        commands_merge = [
            "python3 src/parsers/csv_parser.py --memory_folder {memory_folder}/{input_name} --time_filename {gnark_bench_folder}/{backend}_{circuit}.csv --circuit {circuit}; \n".format(
                memory_folder=helper.Paths().GNARK_BENCH_MEMORY,
                input_name=inp.replace('input/circuit/', '').replace('.json', ''),
                gnark_bench_folder=helper.Paths().GNARK_BENCH,
                backend=backend,
                circuit=circ
            )
            for backend in payload.backend
            for circ, input_path in payload.circuit.items()
            for op in payload.operation
            for inp in helper.get_all_input_files(input_path)
        ]

        # Join the commands into a single string
        pre_command = "".join(commands + commands_memory + commands_merge)
        
        command = f"cd {helper.Paths().GNARK_DIR}; \
                    {command_binary} \
                    {command_check_tmp} \
                    {pre_command}\n"
    else:
        raise ValueError("Missing payload fields for circuit mode")
    return command


def build_command_circom_snarkjs(payload, count):
    return build_command_circom(payload, count)


def build_command_circom_rapidsnark(payload, count):
    return build_command_circom(payload, count, rapidsnark=True)


def build_command_circom(payload, count, rapidsnark=False):
    """
    Build the command to invoke the circom ZKP-framework given the payload

    If rapidsnark is true it will benchmark both provers
    NOTE: rapidsnark only works in intel x64
    """
    os.makedirs(helper.Paths().CIRCOM_BENCHMAKR_DIR, exist_ok=True)    
    # TODO - Add count to command creation
    if len(payload.backend) != 1 or payload.backend[0] != "groth16":
        raise ValueError("Circom benchmark only supports groth16 backend")
    if len(payload.curves) != 1 or payload.curves[0] != "bn128":
        raise ValueError("Circom benchmark only suppports bn128 curve")
    # TODO handle diffent operations
    commands = []
    for circuit, input_path in payload.circuit.items():
        for _ in range(0, count):
            # TODO check if circuit exists
            for inp in helper.get_all_input_files(input_path):
                command = "{script} {circuit_file} {circuit_name} {input_path} {ptau} {benchmark} tmp {template_vars} {rapidsnark}\n".format(
                    script=helper.Paths().CIRCOM_SCRIPT,
                    circuit_file=os.path.join(helper.Paths().CIRCOM_CIRCUITS_DIR, circuit, "circuit.circom"),
                    circuit_name=circuit,
                    input_path=inp,
                    ptau=helper.Paths().CIRCOM_PTAU,
                    benchmark=os.path.join(helper.Paths().CIRCOM_BENCHMAKR_DIR, "circom_" + circuit + ".csv"),
                    template_vars="".join(payload.template_vars.get(circuit, [""])),
                    rapidsnark="1" if rapidsnark else ""
                )
                commands.append(command)
    command = "".join(commands)
    return command


# TODO - This currently uses the halo2 criterion rust parser
def build_command_bellman(payload, count):
    """
    Build the command to invoke the bellman ZKP-library given the payload
    """
    initial_cmd = f"cd {helper.Paths().BELLMAN} && "

    os.makedirs(helper.Paths().BELLMAN_BENCH, exist_ok=True)    
    os.makedirs(helper.Paths().BELLMAN_BENCH_JSON, exist_ok=True)    
    os.makedirs(helper.Paths().BELLMAN_BENCH_MEMORY, exist_ok=True)    
    # TODO - Add count to command creation
    if len(payload.backend) != 1 or payload.backend[0] != "bellman":
        raise ValueError("Bellman benchmark only supports groth16 backend")
    if len(payload.curves) != 1 or payload.curves[0] != "bls12_381":
        raise ValueError("Bellman benchmark only supports bls12_381 curve")
    # TODO handle diffent operations (i.e., algorithms)
    commands = []
    for circuit, input_path in payload.circuit.items():
        for inp in helper.get_all_input_files(input_path):
            commands.append(f"cd {helper.Paths().BELLMAN}; ")
            output_bench = os.path.join(
                helper.Paths().BELLMAN_BENCH_JSON,
                circuit + "_bench_" + os.path.basename(inp)
            )
            input_file = os.path.join("..", "..", inp)
            command_bench: str = "{initial_cmd} RUSTFLAGS=-Awarnings INPUT_FILE={input_file} CIRCUIT={circuit} cargo criterion --message-format=json --bench {bench} 1> {output}; ".format(
                initial_cmd = initial_cmd,
                circuit=circuit,
                input_file=input_file,
                bench="benchmark_circuit",
                output=output_bench
            )
            commands.append(command_bench)
            # Memory commands
            os.makedirs(f"{helper.Paths().BELLMAN_BENCH_MEMORY}/{inp}", exist_ok=True)
            # Altough each operation need only a subset of the arguments we pass
            # all of them for simplicity
            os.makedirs(os.path.join(helper.Paths().BELLMAN, "tmp"), exist_ok=True)
            for op in payload.operation:
                cargo_cmd = "cargo run --bin {circuit} --release -- --input {inp} --phase {phase} --params {params} --proof {proof}".format(
                    circuit=circuit,
                    inp=input_file,
                    phase=op,
                    params=os.path.join("tmp", "params"),
                    proof=os.path.join("tmp", "proof"),
                )
                commands.append(
                    "{initial_cmd} RUSTFLAGS=-Awarnings {memory_cmd} {cargo} 2> {time_file} > /dev/null; ".format(
                        initial_cmd=initial_cmd,
                        memory_cmd=helper.get_memory_command(),
                        cargo=cargo_cmd,
                        time_file=f"{helper.Paths().BELLMAN_BENCH_MEMORY}/{inp}/bellman_{circuit}_memory_{op}.txt"
                    )
                )
            commands.append("cd ../../; ")
            out = os.path.join(
                helper.Paths().BELLMAN_BENCH,
                "bellman_bls12_381_" + circuit + ".csv"
            )
            python_command = "python3"
            try:
                subprocess.run([python_command, "--version"], stdout=subprocess.DEVNULL, stderr=subprocess.DEVNULL)
            except (subprocess.CalledProcessError, FileNotFoundError):
                try:
                    python_command = "python3"
                    subprocess.run([python_command, "--version"], stdout=subprocess.DEVNULL, stderr=subprocess.DEVNULL)
                except (subprocess.CalledProcessError, FileNotFoundError):
                    print("Neither Python nor Python3 are installed or accessible. Please install or check your path settings.")
                    sys.exit(1)

            transform_command = "{python} src/parsers/criterion_rust_parser.py --framework bellman --category circuit --backend bellman --curve bls12_381 --input {inp} --criterion_json {bench} --proof {proof} --output_csv {out}; ".format(
                python=python_command,
                inp=inp,
                bench=output_bench,
                proof=os.path.join(helper.Paths().BELLMAN, "tmp", "proof"),
                out=out
            )
            commands.append(transform_command)
            time_merge = "python3 src/parsers/csv_parser_rust.py --memory_folder {memory_folder} --time_filename {time_filename} --circuit {circuit}; ".format(
                memory_folder=os.path.join(helper.Paths().BELLMAN_BENCH_MEMORY, inp),
                time_filename=out,
                circuit=circuit
            )
            commands.append(time_merge)

    # Join the commands into a single string
    command = "".join(commands)
    return command

def build_command_starky(payload, count):
    """
    Build the command to invoke the starky ZKP-library given the payload
    """
    
    initial_cmd = f"cd {helper.Paths().STARKY} && "

    os.makedirs(helper.Paths().BELLMAN_BENCH, exist_ok=True)    
    os.makedirs(helper.Paths().BELLMAN_BENCH_MEMORY, exist_ok=True)  

    if not os.path.exists(helper.Paths().STARKY_BENCH_JSON):
        try:
            os.makedirs(helper.Paths().STARKY_BENCH_JSON)
        except OSError as e:
            # This can happen when the process doesn't have write permissions or another error
            print("Error: Creating directory. " +  folder_path)
            raise
    # TODO - Add count to command creation
    if len(payload.backend) != 1 or payload.backend[0] != "starky":
        raise ValueError("Starky benchmark only supports starky backend")
    # TODO - Solution for Starks - don't use curve, rename parameter / other option?
    if len(payload.curves) != 1 or payload.curves[0] != "goldilocks":
        raise ValueError("Starky benchmark only supports goldilocks field")
          
    # TODO handle diffent operations (i.e., algorithms)
    commands = []
    for circuit, input_path in payload.circuit.items():
        for inp in helper.get_all_input_files(input_path):
            commands.append(f"cd {helper.Paths().STARKY}; ")
            output_bench = os.path.join(
                helper.Paths().STARKY_BENCH_JSON,
                circuit + "_bench_" + os.path.basename(inp)
            )
            input_file = os.path.join("..", "..", inp)
            command_bench: str = "{initial_cmd} RUSTFLAGS=-Awarnings INPUT_FILE={input_file} CIRCUIT={circuit} cargo criterion --message-format=json --bench {bench} 1> {output}; ".format(
                initial_cmd=initial_cmd,
                circuit=circuit,
                input_file=input_file,
                bench="benchmark_circuit",
                output=output_bench
            )
            commands.append(command_bench)
            # Memory commands
            os.makedirs(f"{helper.Paths().STARKY_BENCH_MEMORY}/{inp}", exist_ok=True)
            # Altough each operation need only a subset of the arguments we pass
            # all of them for simplicity
            os.makedirs(os.path.join(helper.Paths().STARKY, "tmp"), exist_ok=True)
            for op in payload.operation:
                cargo_cmd = "cargo run --bin {circuit}_{path} --release -- --input {inp} --proof {proof}".format(
                    initial_cmd=initial_cmd,
                    circuit=circuit,
                    inp=input_file,
                    path=op,
                    proof=os.path.join("tmp", "proof"),
                )
                commands.append(
                    "{initial_cmd} RUSTFLAGS=-Awarnings {memory_cmd} {cargo} 2> {time_file} > /dev/null; ".format(
                        initial_cmd=initial_cmd,
                        memory_cmd=helper.get_memory_command(),
                        cargo=cargo_cmd,
                        time_file=f"{helper.Paths().STARKY_BENCH_MEMORY}/{inp}/starky_{circuit}_memory_{op}.txt"
                    )
                )
            commands.append("cd ../../; ")
            out = os.path.join(
                helper.Paths().STARKY_BENCH,
                "starky_goldilocks_" + circuit + ".csv"
            )
            python_command = "python3"
            try:
                subprocess.run([python_command, "--version"], stdout=subprocess.DEVNULL, stderr=subprocess.DEVNULL)
            except (subprocess.CalledProcessError, FileNotFoundError):
                try:
                    python_command = "python3"
                    subprocess.run([python_command, "--version"], stdout=subprocess.DEVNULL, stderr=subprocess.DEVNULL)
                except (subprocess.CalledProcessError, FileNotFoundError):
                    print("Neither Python nor Python3 are installed or accessible. Please install or check your path settings.")
                    sys.exit(1)

            transform_command = "{python} src/parsers/criterion_rust_parser.py --framework starky --category circuit --backend starky --curve goldilocks --input {inp} --criterion_json {bench} --proof {proof} --output_csv {out}; ".format(
                python=python_command,
                inp=inp,
                bench=output_bench,
                proof=os.path.join(helper.Paths().STARKY, "tmp", "proof"),
                out=out
            )
            commands.append(transform_command)
            time_merge = "python3 src/parsers/csv_parser_rust.py --memory_folder {memory_folder} --time_filename {time_filename} --circuit {circuit}; ".format(
                memory_folder=os.path.join(helper.Paths().STARKY_BENCH_MEMORY, inp),
                time_filename=out,
                circuit=circuit
            )
            commands.append(time_merge)

    # Join the commands into a single string
    command = "".join(commands)
    return command



def build_command_halo2_pse(payload, count):
    """
    Build the command to invoke the halo2 PSE ZKP-library given the payload
    """
    initial_cmd = f"cd {helper.Paths().HALO2_PSE} && "

    os.makedirs(helper.Paths().HALO2_PSE_BENCH, exist_ok=True)    
    os.makedirs(helper.Paths().HALO2_PSE_BENCH_JSON, exist_ok=True)    
    os.makedirs(helper.Paths().HALO2_PSE_BENCH_MEMORY, exist_ok=True)    

    # TODO - Add count to command creation

    if len(payload.backend) != 1 or payload.backend[0] != "halo2":
        raise ValueError("PSE Halo2 benchmark only supports halo2 backend")
    if len(payload.curves) != 1 or payload.curves[0] != "bn256":
        raise ValueError("PSE Halo2 benchmark only suppports bn256 curve")
    # TODO handle diffent operations (i.e., algorithms)
    commands = []
    for circuit, input_path in payload.circuit.items():
        for inp in helper.get_all_input_files(input_path):
            commands.append(f"cd {helper.Paths().HALO2_PSE}; ")
            output_bench = os.path.join(
                helper.Paths().HALO2_PSE_BENCH_JSON,
                circuit + "_bench_" + os.path.basename(inp)
            )
            input_file = os.path.join("..", "..", inp)
            command_bench: str = "{initial_cmd} RUSTFLAGS=-Awarnings INPUT_FILE={input_file} cargo criterion --message-format=json --bench {bench} 1> {output}; ".format(
                initial_cmd=initial_cmd,
                input_file=input_file,
                bench=circuit + "_bench",
                output=output_bench
            )
            commands.append(command_bench)
            # Memory commands
            os.makedirs(f"{helper.Paths().HALO2_PSE_BENCH_MEMORY}/{inp}", exist_ok=True)
            # Altough each operation need only a subset of the arguments we pass
            # all of them for simplicity
            os.makedirs(os.path.join(helper.Paths().HALO2_PSE, "tmp"), exist_ok=True)
            for op in payload.operation:
                cargo_cmd = "cargo run --bin {circuit} --release -- --input {inp} --phase {phase} --params {params} --vk {vk} --pk {pk} --proof {proof}".format(
                    circuit=circuit,
                    inp=input_file,
                    phase=op,
                    params=os.path.join("tmp", "params"),
                    vk=os.path.join("tmp", "vk"),
                    pk=os.path.join("tmp", "pk"),
                    proof=os.path.join("tmp", "proof"),
                )
                commands.append(
                    "{initial_cmd} RUSTFLAGS=-Awarnings {memory_cmd} {cargo} 2> {time_file} > /dev/null; ".format(
                        initial_cmd=initial_cmd,
                        memory_cmd=helper.get_memory_command(),
                        cargo=cargo_cmd,
                        time_file=f"{helper.Paths().HALO2_PSE_BENCH_MEMORY}/{inp}/halo2_{circuit}_memory_{op}.txt"
                    )
                )
            commands.append("cd ../../; ")
            out = os.path.join(
                helper.Paths().HALO2_PSE_BENCH,
                "halo2_pse_bn256_" + circuit + ".csv"
            )
            transform_command: str = "python3 src/parsers/criterion_rust_parser.py --framework halo2_pse --category circuit --backend halo2 --curve bn256 --input {inp} --criterion_json {bench} --proof {proof} --output_csv {out}; ".format(
                inp=inp,
                bench=output_bench,
                out=out,
                proof=os.path.join(helper.Paths().HALO2_PSE, "tmp", "proof")
            )
            commands.append(transform_command)
            time_merge = "python3 src/parsers/csv_parser_rust.py --memory_folder {memory_folder} --time_filename {time_filename} --circuit {circuit}; ".format(
                memory_folder=os.path.join(helper.Paths().HALO2_PSE_BENCH_MEMORY, inp),
                time_filename=os.path.join(helper.Paths().HALO2_PSE_BENCH, f"halo2_pse_bn256_{circuit}.csv"),
                circuit=circuit
            )
            commands.append(time_merge)
    # Join the commands into a single string
    command = "".join(commands)

    return command

def default_case(_payload, _count):
    raise ValueError("Framework not integrated into the benchmarking framework!")


# List ZKP-frameworks in the zk-Harness
projects = {
    "gnark":                build_command_gnark,
    "circom/snarkjs":       build_command_circom_snarkjs,
    "circom/rapidsnark":    build_command_circom_rapidsnark,
    "bellman":              build_command_bellman,
    "starky":               build_command_starky,
    "halo2_pse":            build_command_halo2_pse
}


def build_command(project, payload, count):
    """
    Build the command to execute the given project with the given payload.
    Input: project (e.g. gnark) + payload (config.json)
    """
    commands = projects.get(project, default_case)(payload, count)
    return commands


def get_circuit_payload(config):
    """
    Extract the payload for category "circuit" given a config.json
    """
    # Extract the relevant fields from the configuration data
    payload = config.get('payload')
    if payload is None:
        raise KeyError("Payload does not exist in circuit config")

    backend = payload.get('backend')
    if backend is None:
        raise KeyError("backend field does not exist in circuit payload")
    if len(backend) == 0:
        raise ValueError("backend field is empty")

    curves = payload.get('curves')
    if curves is None:
        raise KeyError("curves field does not exist in circuit payload")
    if len(curves) == 0:
        raise ValueError("curves field is empty")

    circuits = payload.get('circuits')
    if circuits is None:
        raise KeyError("circuits field does not exist in circuit payload")
    circuits = config['payload']['circuits'].keys()
    if len(circuits) == 0:
        raise ValueError("circuits field is empty")

    # FIXME use operation instead of algorithm
    operation = payload.get('algorithm')
    if operation is None:
        raise KeyError("operation field does not exist in circuit payload")
    if len(operation) == 0:
        raise ValueError("operation field is empty")
    for op in operation:
        if op not in OPERATIONS:
            raise ValueError(f"operation '{op}' not in {OPERATIONS}")

    input_path = []
    template_vars = {}
    for c_name, c in payload['circuits'].items():
        inp = c.get("input_path")
        if inp is None:
            raise KeyError(f"input_path does not exist to '{c}' circuit")
        if "template_vars" in c:
            template_vars[c_name] = c["template_vars"]
        input_path.append(inp)

    # Map circuit names onto input paths
    circuit = dict(zip(circuits, input_path))
    
    # Define a named tuple for the payload
    Payload = namedtuple('Payload', ['backend', 'curves', 'circuit', 'operation', 'input_path', 'template_vars'])

    # Return a new instance of the named tuple with the extracted values
    return Payload(backend, curves, circuit, operation, input_path, template_vars)
