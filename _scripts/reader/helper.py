import os
import platform


class Paths():
    _instance = None
    MAIN_DIR = None
    BENCHMARKS_DIR = None

    BELLMAN = None
    BELLMAN_BENCH = None
    BELLMAN_BENCH_JSON = None
    BELLMAN_BENCH_MEMORY = None

    GNARK_DIR = None
    GNARK_BENCH = None
    GNARK_BENCH_MEMORY = None

    HALO2_PSE = None
    HALO2_PSE_BENCH = None
    HALO2_PSE_BENCH_JSON = None
    HALO2_PSE_BENCH_MEMORY = None

    CIRCOM_DIR = None
    CIRCOM_BENCHMAKR_DIR = None
    SNARKJS_BENCHMAKR_DIR = None
    RAPIDSNARK_BENCHMAKR_DIR = None
    CIRCOM_SCRIPT = None
    CIRCOM_CIRCUITS_DIR = None
    CIRCOM_PTAU = None
    # SNARKJS PATHS
    SNARKJS_DIR = None
    SNARKJS_ARITHMETICS_SCRIPT = None
    SNARKJS_EC_SCRIPT = None
    # SNARKJS PATHS
    RAPIDSNARK_DIR = None
    RAPIDSNARK_ARITHMETICS_SCRIPT = None
    RAPIDSNARK_EC_SCRIPT = None

    def __new__(cls, machine=None, *args, **kwargs):
        if not cls._instance:
            assert machine is not None
            cls._instance = super().__new__(cls, *args, **kwargs)
            # GENERAL PATHS
            cls._instance.MAIN_DIR = os.path.abspath(os.path.join(os.path.dirname(__file__), '..', '..'))
            cls._instance.BENCHMARKS_DIR = os.path.join(cls._instance.MAIN_DIR, "benchmarks", machine)
            # BELLMAN PATHS
            cls._instance.BELLMAN = os.path.join(cls._instance.MAIN_DIR, "bellman_circuits")
            cls._instance.BELLMAN_BENCH = os.path.join(cls._instance.BENCHMARKS_DIR, "bellman")
            cls._instance.BELLMAN_BENCH_MEMORY = os.path.join(cls._instance.BELLMAN_BENCH, "memory")
            cls._instance.BELLMAN_BENCH_JSON = os.path.join(cls._instance.BELLMAN_BENCH, "jsons")
            # GNARK PATHS
            cls._instance.GNARK_DIR = os.path.join(cls._instance.MAIN_DIR, "gnark")
            cls._instance.GNARK_BENCH = os.path.join(cls._instance.BENCHMARKS_DIR, "gnark")
            cls._instance.GNARK_BENCH_MEMORY = os.path.join(cls._instance.GNARK_BENCH, "memory")
            # HALO2_PSE PATHS
            cls._instance.HALO2_PSE = os.path.join(cls._instance.MAIN_DIR, "halo2_pse")
            cls._instance.HALO2_PSE_BENCH = os.path.join(cls._instance.BENCHMARKS_DIR, "halo2_pse")
            cls._instance.HALO2_PSE_BENCH_MEMORY = os.path.join(cls._instance.HALO2_PSE_BENCH, "memory")
            cls._instance.HALO2_PSE_BENCH_JSON = os.path.join(cls._instance.HALO2_PSE_BENCH, "jsons")
            # CIRCOM PATHS
            cls._instance.CIRCOM_DIR = os.path.join(cls._instance.MAIN_DIR, "circom")
            cls._instance.CIRCOM_BENCHMAKR_DIR = os.path.join(cls._instance.BENCHMARKS_DIR, "circom")
            cls._instance.SNARKJS_BENCHMAKR_DIR = os.path.join(cls._instance.BENCHMARKS_DIR, "snarkjs")
            cls._instance.RAPIDSNARK_BENCHMAKR_DIR = os.path.join(cls._instance.BENCHMARKS_DIR, "rapidsnark")
            cls._instance.CIRCOM_SCRIPT = os.path.join(cls._instance.CIRCOM_DIR, "scripts", "run_circuit.sh")
            cls._instance.CIRCOM_CIRCUITS_DIR = os.path.join(cls._instance.CIRCOM_DIR, "circuits", "benchmarks")
            cls._instance.CIRCOM_PTAU = os.path.join(cls._instance.CIRCOM_DIR, "phase1", "powersOfTau28_hez_final_16.ptau")
            # SNARKJS PATHS
            cls._instance.SNARKJS_DIR = os.path.join(cls._instance.MAIN_DIR, "snarkjs")
            cls._instance.SNARKJS_ARITHMETICS_SCRIPT = os.path.join(cls._instance.SNARKJS_DIR, "scripts", "arithmetics.js")
            cls._instance.SNARKJS_EC_SCRIPT = os.path.join(cls._instance.SNARKJS_DIR, "scripts", "curves.js")
            # SNARKJS PATHS
            cls._instance.RAPIDSNARK_DIR = os.path.join(cls._instance.MAIN_DIR, "rapidsnark")
            cls._instance.RAPIDSNARK_ARITHMETICS_SCRIPT = os.path.join(cls._instance.RAPIDSNARK_DIR, "scripts", "arithmetics.js")
            cls._instance.RAPIDSNARK_EC_SCRIPT = os.path.join(cls._instance.RAPIDSNARK_DIR, "scripts", "curves.js")
        return cls._instance


### GENERAL ###
MEMORY_CMD = "/usr/bin/time"
ARITHMETIC_FIELDS = ["base", "scalar"]
GROUPS = ["g1", "g2"]
# CIRCOM CURVES
CIRCOM_CURVES = ["bn128", "bls12_381"]


def get_all_input_files(input_path, abspath=False):
    """
    Given a input_path return the full path of the file or if it is a directory
    return the full paths of all JSON files in this directory
    """
    if not os.path.exists(input_path):
        raise ValueError(f"Input: {input_path} does not exist")
    if os.path.isfile(input_path):
        if not input_path.endswith(".json"):
            raise ValueError(f"Input: {input_path} is not a JSON file")
        return [os.path.abspath(input_path)] if abspath else [input_path]
    # input_path is a directory
    files = []
    # NOTE this operation is not recursive 
    for f in os.listdir(input_path):
        file = os.path.join(input_path, f)
        if os.path.isfile(file) and file.endswith(".json"):
            files.append(os.path.abspath(file) if abspath else file)
    if len(files) == 0:
        raise ValueError(f"Input: no input file detected in {input_path}")
    return files


def get_memory_command():
    system = platform.system()
    options = ""
    if system == 'Darwin':
        options = "-h -l"
    elif system == 'Linux':
        options = "-v"
    else:
        raise Exception("Unsupported operating system")
    return MEMORY_CMD + " " + options
