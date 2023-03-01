#!/usr/bin/env python

import os

from collections import namedtuple

from . import helper

OPERATIONS = [
    "add", "sub", "mul", "inv", "exp"
]
# Define a named tuple for the payload
Payload = namedtuple('Payload', ['curves', 'fields', 'operations'])



def build_command_gnark(payload, count):
    """
    Build the command to invoke the gnark ZKP-framework given the payload
    """
    if payload.fields is not None and payload.curves is not None:
        commands = [f"./gnark arithmetic --input={inp} --field={field} --operation={op} --curve={curve} --count={count}\n"
                    for curve in payload.curves
                    for field in payload.fields
                    for op, input_path in payload.operations.items()
                    for inp in helper.get_all_input_files(input_path)
                    ]

        # Join the commands into a single string
        command = "".join(commands)
        # Prepend the command to change the working directory to the gnark directory
        command = f"cd {helper.GNARK_DIR}; {command}"
    else:
        raise ValueError("Missing payload fields for arithmetic mode")
    return command

def build_command_circom(payload, count):
    """
    Build the command to invoke the circom ZKP-framework given the payload
    """
    for c in payload.curves: 
        if c not in helper.CIRCOM_CURVES:
            raise ValueError(f"Curve {c} not in {helper.CIRCOM_CURVES}")
    for f in payload.fields:
        if f not in helper.ARITHMETIC_FIELDS:
            raise ValueError(f"Field {f} not in {helper.ARITHMETIC_FIELDS}")
    commands = [
        "{script} {curve} {field} {operation} {count} {input_path} {benchmark}\n".format(
            script=helper.CIRCOM_ARITHMETICS_SCRIPT,
            curve=curve,
            field=field,
            operation=operation,
            count=count,
            input_path=inp,
            benchmark=os.path.join(helper.CIRCOM_BENCHMAKR_DIR, "circom_arithmetics.csv")
        )
        for operation, input_path in payload.operations.items()
        for inp in helper.get_all_input_files(input_path)
        for curve in payload.curves
        for field in payload.fields
    ]
    command = "".join(commands)
    return command


def default_case():
    raise ValueError("Framework not integrated into the benchmarking framework!")


# List ZKP-frameworks in the zk-Harness
projects = {
    "gnark":    build_command_gnark,
    "circom":   build_command_circom
}


def build_command(project, payload, count):
    """
    Build the command to execute the given project with the given payload.
    Input: project (e.g. gnark) + payload (config.json)
    """
    commands = projects.get(project, default_case)(payload, count)
    return commands


def get_arithmetic_payload(config):
    """
    Extract the payload for category "circuit" given a config.json
    """
    # Extract the relevant fields from the configuration data
    payload = config.get('payload')
    if payload is None:
        raise KeyError("Payload does not exist in arithmetic config")

    curves = payload.get('curves')
    if curves is None:
        raise KeyError("curves field does not exist in arithmetic payload")
    if len(curves) == 0:
        raise ValueError("curves field is empty")

    fields = payload.get('fields')
    if fields is None:
        raise KeyError("fields field does not exist in arithmetic payload")
    if len(fields) == 0:
        raise ValueError("fields field is empty")

    operations = payload.get('operations')
    if operations is None:
        raise KeyError("operations field does not exist in arithmetic payload")
    if len(operations.keys()) == 0:
        raise ValueError("operations field is empty")
    for op in operations.keys():
        if op not in OPERATIONS:
            raise ValueError(f"operation '{op}' not in {OPERATIONS}")

    input_path = []
    for op in payload['operations'].values():
        inp = op.get("input_path")
        if inp is None:
            raise KeyError(f"input_path does not exist to '{op}' operation")
        input_path.append(inp)

    # Map operation names onto input paths
    ops = dict(zip(operations.keys(), input_path))
    
    # Return a new instance of the named tuple with the extracted values
    return Payload(curves, fields, ops)
