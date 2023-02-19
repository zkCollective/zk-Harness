#!/usr/bin/env python

import os

from collections import namedtuple

OPERATIONS = [
    "compile", "setup", "witness", "prove", "verify"
]
# GENERAL PATHS
MAIN_DIR = os.path.abspath(os.path.join(os.path.dirname(__file__), '..', '..'))
BENCHMARKS_DIR = os.path.join(MAIN_DIR, "benchmarks")
# GNARK PATHS
GNARK_DIR = os.path.join(MAIN_DIR, "gnark")
# CIRCOM PATHS
CIRCOM_DIR = os.path.join(MAIN_DIR, "circom")
CIRCOM_BENCHMAKR_DIR = os.path.join(BENCHMARKS_DIR, "circom")
CIRCOM_SCRIPT = os.path.join(CIRCOM_DIR, "scripts", "run_circuit.sh")
CIRCOM_CIRCUITS_DIR = os.path.join(CIRCOM_DIR, "circuits", "benchmarks")
CIRCOM_PTAU = os.path.join(CIRCOM_DIR, "phase1", "powersOfTau28_hez_final_16.ptau")


def get_all_input_files(input_path):
    """
    Given a input_path return the full path of the file or if it is a directory
    return the full paths of all JSON files in this directory
    """
    if not os.path.exists(input_path):
        raise ValueError(f"Input: {input_path} does not exist")
    if os.path.isfile(input_path):
        if not input_path.endswith(".json"):
            raise ValueError(f"Input: {input_path} is not a JSON file")
        return [os.path.abspath(input_path)]
    # input_path is a directory
    files = []
    # NOTE this operation is not recursive 
    for f in os.listdir(input_path):
        file = os.path.join(input_path, f)
        if os.path.isfile(file) and file.endswith(".json"):
            files.append(os.path.abspath(file))
    if len(files) == 0:
        raise ValueError(f"Input: no input file detected in {input_path}")
    return files


def build_command_gnark(payload):
    """
    Build the command to invoke the gnark ZKP-framework given the payload
    """

    if payload.backend is not None and payload.curves is not None:
        commands = [f"./gnark {backend} --circuit={circ} --algo={op} --curve={curve} --input={input_path}\n"
                    for backend in payload.backend
                    for curve in payload.curves
                    for circ, input_path in payload.circuit.items()
                    for op in payload.operation]

        # Join the commands into a single string
        command = "".join(commands)
        # Prepend the command to change the working directory to the gnark directory
        command = f"cd {GNARK_DIR}; {command}"
    else:
        raise ValueError("Missing payload fields for circuit mode")
    return command


def build_command_circom(payload):
    """
    Build the command to invoke the circom ZKP-framework given the payload
    """
    
    if len(payload.backend) != 1 or payload.backend[0] != "groth16":
        raise ValueError("Circom benchmark only supports groth16 backend")
    if len(payload.curves) != 1 or payload.curves[0] != "bn128":
        raise ValueError("Circom benchmark only suppports bn128 curve")
    # TODO handle diffent operations
    commands = []
    for circuit, input_path in payload.circuit.items():
        # TODO check if circuit exists
        for inp in get_all_input_files(input_path):
            command = "{script} {circuit_file} {input_path} {ptau} {benchmark}\n".format(
                script=CIRCOM_SCRIPT,
                circuit_file=os.path.join(CIRCOM_CIRCUITS_DIR, circuit, "circuit.circom"),
                input_path=inp,
                ptau=CIRCOM_PTAU,
                benchmark=os.path.join(CIRCOM_BENCHMAKR_DIR, "circom_" + circuit + ".csv")
            )
            commands.append(command)
    command = "".join(commands)
    return command


def default_case():
    raise ValueError("Framework not integrated into the benchmarking framework!")


# List ZKP-frameworks in the zk-Harness
projects = {
    "gnark":    build_command_gnark,
    "circom":   build_command_circom
}


def build_command(project, payload):
    """
    Build the command to execute the given project with the given payload.
    Input: project (e.g. gnark) + payload (config.json)
    """
    commands = projects.get(project, default_case)(payload)
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
    for c in payload['circuits'].values():
        inp = c.get("input_path")
        if inp is None:
            raise KeyError(f"input_path does not exist to '{c}' circuit")
        input_path.append(inp)

    # Map circuit names onto input paths
    circuit = dict(zip(circuits, input_path))
    
    # Define a named tuple for the payload
    Payload = namedtuple('Payload', ['backend', 'curves', 'circuit', 'operation', 'input_path'])

    # Return a new instance of the named tuple with the extracted values
    return Payload(backend, curves, circuit, operation, input_path)
