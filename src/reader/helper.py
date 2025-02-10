import os
import platform


FRAMEWORKS_DIR = "frameworks"

class Paths():
    _instance = None
    MAIN_DIR = None
    BENCHMARKS_DIR = None
    CIRCUITS_BENCH = None
    MATH_BENCH = None

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

    def __new__(cls, machine=None, *args, **kwargs):
        if not cls._instance:
            assert machine is not None
            cls._instance = super().__new__(cls, *args, **kwargs)
            # GENERAL PATHS
            cls._instance.MAIN_DIR = os.path.abspath(os.path.join(os.path.dirname(__file__), '..', '..'))
            cls._instance.BENCHMARKS_DIR = os.path.join(cls._instance.MAIN_DIR, "benchmarks", machine)
            cls._instance.CIRCUITS_BENCH = os.path.join(cls._instance.MAIN_DIR, "benchmarks", "circuit", machine)
            cls._instance.MATH_BENCH = os.path.join(cls._instance.MAIN_DIR, "benchmarks", "math", machine)
            # BELLMAN PATHS
            cls._instance.BELLMAN = os.path.join(cls._instance.MAIN_DIR, FRAMEWORKS_DIR, "bellman_circuits")
            cls._instance.BELLMAN_BENCH = os.path.join(cls._instance.CIRCUITS_BENCH, "bellman")
            cls._instance.BELLMAN_BENCH_MEMORY = os.path.join(cls._instance.BELLMAN_BENCH, "memory")
            cls._instance.BELLMAN_BENCH_JSON = os.path.join(cls._instance.BELLMAN_BENCH, "jsons")
            # STARKY PATHS
            cls._instance.STARKY = os.path.join(cls._instance.MAIN_DIR, FRAMEWORKS_DIR, "starky_circuits")
            cls._instance.STARKY_BENCH = os.path.join(cls._instance.CIRCUITS_BENCH, "starky")
            cls._instance.STARKY_BENCH_MEMORY = os.path.join(cls._instance.STARKY_BENCH, "memory")
            cls._instance.STARKY_BENCH_JSON = os.path.join(cls._instance.STARKY_BENCH, "jsons")
            # GNARK PATHS
            cls._instance.GNARK_DIR = os.path.join(cls._instance.MAIN_DIR, FRAMEWORKS_DIR, "gnark")
            cls._instance.GNARK_BENCH = os.path.join(cls._instance.CIRCUITS_BENCH, "gnark")
            cls._instance.GNARK_BENCH_MEMORY = os.path.join(cls._instance.GNARK_BENCH, "memory")
            # HALO2_PSE PATHS
            cls._instance.HALO2_PSE = os.path.join(cls._instance.MAIN_DIR, FRAMEWORKS_DIR, "halo2_pse")
            cls._instance.HALO2_PSE_BENCH = os.path.join(cls._instance.CIRCUITS_BENCH, "halo2_pse")
            cls._instance.HALO2_PSE_BENCH_MEMORY = os.path.join(cls._instance.HALO2_PSE_BENCH, "memory")
            cls._instance.HALO2_PSE_BENCH_JSON = os.path.join(cls._instance.HALO2_PSE_BENCH, "jsons")
            # CIRCOM PATHS
            cls._instance.CIRCOM_DIR = os.path.join(cls._instance.MAIN_DIR, FRAMEWORKS_DIR, "circom")
            cls._instance.CIRCOM_BENCHMAKR_DIR = os.path.join(cls._instance.CIRCUITS_BENCH, "circom")
            cls._instance.CIRCOM_SCRIPT = os.path.join(cls._instance.CIRCOM_DIR, "scripts", "run_circuit.sh")
            cls._instance.CIRCOM_CIRCUITS_DIR = os.path.join(cls._instance.CIRCOM_DIR, "circuits", "benchmarks")
            cls._instance.CIRCOM_PTAU = os.path.join(cls._instance.CIRCOM_DIR, "phase1", "powersOfTau28_final.ptau")
        return cls._instance


### GENERAL ###
MEMORY_CMD = "/usr/bin/time"


def get_all_input_files(input_path, abspath=False):
    """
    Given a input_path return the full path of the file or if it is a directory
    return the full paths of all JSON files in this directory
    """
    if isinstance(input_path, list):
        files = []
        for p in input_path:
            files.extend(get_all_input_files(p, abspath))
        return files
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
