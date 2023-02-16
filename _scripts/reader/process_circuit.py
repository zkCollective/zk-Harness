#!/usr/bin/env python

import argparse
import json
import subprocess
import os

from collections import namedtuple

def build_command_gnark(payload):
    """
    Build the command to invoke the gnark ZKP-framework given the payload
    """

    # Get the absolute path to the gnark directory
    gnark_dir = os.path.abspath(os.path.join(os.path.dirname(__file__), '..', '..', 'gnark'))

    if payload.backend is not None and payload.curves is not None:
        commands = [f"./gnark {backend} --circuit={circ} --algo={algo} --curve={curve} --input={input_path}\n"
                    for backend in payload.backend
                    for curve in payload.curves
                    for circ, input_path in payload.circuit.items()
                    for algo in payload.algo]

        # Join the commands into a single string
        command = "".join(commands)

        # Prepend the command to change the working directory to the gnark directory
        command = f"cd {gnark_dir}; {command}"
    else:
        raise ValueError("Missing payload fields for circuit mode")
    return command


def default_case():
    raise ValueError("Framework not integrated into the benchmarking framework!")


# List ZKP-frameworks in the zk-Harness
projects = {
    "gnark":    build_command_gnark
    # "circom":   circom_processing
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
    backend = config['payload']['backend']
    curves = config['payload']['curves']
    circuits = config['payload']['circuits'].keys()
    algo = config['payload']['algorithm']
    input_path = [list(config['payload']['circuits'][c].values())[0] for c in circuits]

    # Map circuit names onto input paths
    circuit = dict(zip(circuits, input_path))
    
    # Define a named tuple for the payload
    Payload = namedtuple('Payload', ['backend', 'curves', 'circuit', 'algo', 'input_path'])

    # Return a new instance of the named tuple with the extracted values
    return Payload(backend, curves, circuit, algo, input_path)

        
