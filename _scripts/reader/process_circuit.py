#!/usr/bin/env python

import argparse
import json
import subprocess
import os

from collections import namedtuple

def build_command_gnark(payload):
    """
    Build the command to execute gnark for the given payload of type "circuit".
    
    Args:
    - payload: A namedtuple containing the following fields:
        - backend (str or list): The name of the backend to use or a list of names of backends to use.
        - curves (str or list): The name of the curve to use or a list of names of curves to use.
        - circuit (str): The name of the circuit to use.
        - algo (str or list): The name of the algorithm to use or a list of names of algorithms to use.
    
    Returns:
    - command (str): A string containing the command to execute gnark with the given payload.
    """

    # Get the absolute path to the gnark directory
    gnark_dir = os.path.abspath(os.path.join(os.path.dirname(__file__), '..', '..', 'gnark'))

    if payload.backend is not None and payload.curves is not None:
        commands = [f"./gnark {backend} --circuit={payload.circuit} --algo={algo} --curve={curve}\n"
                    for backend in payload.backend
                    for curve in payload.curves
                    for algo in payload.algo]

        # Join the commands into a single string
        command = "".join(commands)

        # Prepend the command to change the working directory to the gnark directory
        command = f"cd {gnark_dir}; {command}"
    else:
        raise ValueError("Missing payload fields for circuit mode")
    return command


def default_case():
    """
    Raise a ValueError with the message "Framework not integrated into the benchmarking framework!".
    """
    raise ValueError("Framework not integrated into the benchmarking framework!")


# List of implemented projects in the zk-Harness
projects = {
    "gnark":    build_command_gnark
    # "circom":   circom_processing
}


def build_command(project, payload):
    """
    Build the command to execute the given project with the given payload.
    
    Args:
    - project (str): The name of the project to execute.
    - payload: A namedtuple containing the following fields:
        - backend (str or list): The name of the backend to use or a list of names of backends to use.
        - curves (str or list): The name of the curve to use or a list of names of curves to use.
        - circuit (str): The name of the circuit to use.
        - algo (str or list): The name of the algorithm to use or a list of names of algorithms to use.
    
    Returns:
    - command (str): A string containing the command to execute the given project with the given payload.
    """
    commands = projects.get(project, default_case)(payload)
    return commands


def get_payload(config):
    """
    Extract the relevant fields from the given configuration data and return them as a named tuple.
    
    Args:
    - config: A dictionary containing the configuration data.
    
    Returns:
    - payload (namedtuple): A namedtuple containing the following fields:
        - backend (str or list): The name of the backend to use or a list of names of backends to use.
        - curves (str or list): The name of the curve to use or a list of names of curves to use.
        - circuit (str): The name of the circuit to use.
        - algo (str or list): The name of the
        algorithm to use or a list of names of algorithms to use.
    """
    # Extract the relevant fields from the configuration data
    backend = config['payload']['backend']
    curves = config['payload']['curves']
    circuit = list(config['payload']['circuits'].keys())[0]
    algo = config['payload']['algorithm']

    # Define a named tuple for the payload
    Payload = namedtuple('Payload', ['backend', 'curves', 'circuit', 'algo'])

    # Return a new instance of the named tuple with the extracted values
    return Payload(backend, curves, circuit, algo)

        
