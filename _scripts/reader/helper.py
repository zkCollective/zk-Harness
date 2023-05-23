import os


### PATHS ###
# GENERAL PATHS
MAIN_DIR = os.path.abspath(os.path.join(os.path.dirname(__file__), '..', '..'))
BENCHMARKS_DIR = os.path.join(MAIN_DIR, "benchmarks")
# GNARK PATHS
GNARK_DIR = os.path.join(MAIN_DIR, "gnark")
# CIRCOM PATHS
CIRCOM_DIR = os.path.join(MAIN_DIR, "circom")
CIRCOM_BENCHMAKR_DIR = os.path.join(BENCHMARKS_DIR, "circom")
SNARKJS_BENCHMAKR_DIR = os.path.join(BENCHMARKS_DIR, "snarkjs")
RAPIDSNARK_BENCHMAKR_DIR = os.path.join(BENCHMARKS_DIR, "rapidsnark")
CIRCOM_SCRIPT = os.path.join(CIRCOM_DIR, "scripts", "run_circuit.sh")
CIRCOM_CIRCUITS_DIR = os.path.join(CIRCOM_DIR, "circuits", "benchmarks")
CIRCOM_PTAU = os.path.join(CIRCOM_DIR, "phase1", "powersOfTau28_hez_final_16.ptau")
# SNARKJS PATHS
SNARKJS_DIR = os.path.join(MAIN_DIR, "snarkjs")
SNARKJS_ARITHMETICS_SCRIPT = os.path.join(SNARKJS_DIR, "scripts", "arithmetics.js")
SNARKJS_EC_SCRIPT = os.path.join(SNARKJS_DIR, "scripts", "curves.js")
# SNARKJS PATHS
RAPIDSNARK_DIR = os.path.join(MAIN_DIR, "rapidsnark")
RAPIDSNARK_ARITHMETICS_SCRIPT = os.path.join(RAPIDSNARK_DIR, "scripts", "arithmetics.js")
RAPIDSNARK_EC_SCRIPT = os.path.join(RAPIDSNARK_DIR, "scripts", "curves.js")
### GENERAL ###
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
