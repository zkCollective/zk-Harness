import argparse
import json
import subprocess

from . import process_circuit
from . import process_arithmetic

def circuit_processing(project, config, count):
    # Extract relevant fields from config, build & execute command
    payload = process_circuit.get_circuit_payload(config)
    commands = process_circuit.build_command(project, payload, count)
    subprocess.run(commands, shell=True, check=True)

def arithmetic_processing(project, config, count):
    # Extract relevant fields from config, build & execute command
    payload = process_arithmetic.get_arithmetic_payload(config)
    commands = process_arithmetic.build_command(project, payload, count)
    subprocess.run(commands, shell=True, check=True)

def default_case():
    raise ValueError("Benchmark category not integrated into the benchmarking framework!")

# TODO - Add other modes (arithmetic & curves)
categories = {
    "circuit": circuit_processing,
    "arithmetic": arithmetic_processing
}

def parse_config(config_path):

    with open(config_path, 'r') as f:
        config = json.load(f)

    project = config['project']
    category = config['category']
    count = config['count']
    categories.get(category, default_case)(project, config, count)

if __name__ == '__main__':
    parser = argparse.ArgumentParser()
    parser.add_argument('--config', help='Path to configuration file')
    args = parser.parse_args()
    parse_config(args.config)
