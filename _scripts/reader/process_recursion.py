#!/usr/bin/env python

import os

from collections import namedtuple

from . import helper

OPERATIONS = [
    "compile", "setup", "witness", "prove", "verify"
]

def build_command_gnark_recursion(payload, count):
    """
    Build the command to invoke the gnark ZKP-framework given the payload
    """
    
    if payload.outerBackend is not None and payload.curve is not None:
        commands = [f"./gnark recursion --circuit={circ} --algo={op} --curve=bls12_377 --input={inp} --count={count}\n"
                    for circ, input_path in payload.circuit.items()
                    for inp in helper.get_all_input_files(input_path)
                    for op in payload.operation]

        # Join the commands into a single string
        command = "".join(commands)
        print(command)
        # Prepend the command to change the working directory to the gnark directory
        command = f"cd {helper.GNARK_DIR}; {command}"
    else:
        raise ValueError("Missing payload fields for circuit mode")
    return command

def default_case():
    raise ValueError("Framework not integrated into the benchmarking framework!")

# List ZKP-frameworks in the zk-Harness
projects = {
    "gnark":    build_command_gnark_recursion
}

def build_command(project, payload, count):
    """
    Build the command to execute the given project with the given payload.
    Input: project (e.g. gnark) + payload (config.json)
    """
    commands = projects.get(project, default_case)(payload, count)
    return commands


def get_recursion_payload(config):
    """
    Extract the payload for category "circuit" given a config.json
    """
    # Extract the relevant fields from the configuration data
    payload = config.get('payload')
    if payload is None:
        raise KeyError("Payload does not exist in circuit config")

    innerBackend = payload.get('innerBackend')
    if innerBackend is None:
        raise KeyError("backend field does not exist in circuit payload")
    if len(innerBackend) == 0:
        raise ValueError("backend field is empty")

    curve = payload.get('innerCurve')
    if curve is None:
        raise KeyError("curves field does not exist in circuit payload")
    if len(curve) == 0:
        raise ValueError("curves field is empty")
    
    outerBackend = payload.get('outerBackend')
    if outerBackend is None:
        raise KeyError("curves field does not exist in circuit payload")
    if len(outerBackend) == 0:
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
    Payload = namedtuple('Payload', ['innerBackend', 'outerBackend', 'curve', 'circuit', 'operation', 'input_path'])

    # Return a new instance of the named tuple with the extracted values
    return Payload(innerBackend, outerBackend, curve, circuit, operation, input_path)