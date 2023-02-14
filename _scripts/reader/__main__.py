import argparse
import json
import subprocess

from . import process_circuit

def circuit_processing(project, config):
    """
    Run circuit processing for the specified project and configuration data.
    
    Extracts the relevant fields from the configuration data, builds a command
    for the specified project and fields, and executes the command using
    subprocess.
    """
    # Extract the relevant fields from the configuration data
    payload = process_circuit.get_payload(config)

    # Build the command for the specified project and fields
    commands = process_circuit.build_command(project, payload)
    
    # Execute the command using subprocess
    subprocess.run(commands, shell=True, check=True)

def default_case():
    """
    Raise a ValueError if the specified mode is not integrated into the benchmarking framework.
    """
    raise ValueError("Mode not integrated into the benchmarking framework!")

# List of processing functions for each mode
# TODO - Add other modes (arithmetic & curves)
modes = {
    "circuit": circuit_processing
}

def parse_config(config_path):
    """
    Parse the configuration data from the specified path and execute the appropriate processing function.
    
    Loads the JSON data from the configuration file, selects the appropriate processing function
    based on the mode field in the configuration data, and executes the function with the specified
    project and configuration data.
    """
    # Load the JSON data from the configuration file
    with open(config_path, 'r') as f:
        config = json.load(f)

    # Select the appropriate processing function based on the mode field in the configuration data
    # TODO - Parse the input path!
    project = config['project']
    mode = config['mode']
    modes.get(mode, default_case)(project, config)

if __name__ == '__main__':
    # Define the command line arguments
    parser = argparse.ArgumentParser()
    parser.add_argument('--config', help='Path to configuration file')

    # Parse the command line arguments
    args = parser.parse_args()

    # Parse the configuration data and execute the appropriate processing function
    parse_config(args.config)
