#!/usr/bin/env python

import os

from collections import namedtuple

from . import helper

OPERATIONS = [
    "scalar-multiplication", "multi-scalar-multiplication", "pairing"
]
# Define a named tuple for the payload
Payload = namedtuple('Payload', ['curves', 'groups', 'operations'])



def build_command_gnark(payload, count):
    """
    Build the command to invoke the gnark ZKP-framework given the payload
    """
    if payload.groups is not None and payload.curves is not None:
        commands = [f"./gnark ec --input={inp} --group={group} --operation={op} --curve={curve} --count={count}\n"
                    for curve in payload.curves
                    for group in payload.groups
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

def build_command_snarkjs(payload, count):
    """
    Build the command to invoke the snarkjs ZKP-framework given the payload
    """
    for c in payload.curves: 
        if c not in helper.CIRCOM_CURVES:
            raise ValueError(f"Curve {c} not in {helper.CIRCOM_CURVES}")
    for f in payload.groups:
        if f not in helper.GROUPS:
            raise ValueError(f"Field {f} not in {helper.GROUPS}")
    # It will add two commands for pairing if both g1 and g2 are provided
    commands = [
        "{script} {curve} {group} {operation} {count} {input_path} {benchmark}\n".format(
            script=helper.SNARKJS_EC_SCRIPT,
            curve=curve,
            group=group,
            operation=operation,
            count=count,
            input_path=inp,
            benchmark=os.path.join(helper.SNARKJS_BENCHMAKR_DIR, "snarkjs_ec.csv")
        )
        for operation, input_path in payload.operations.items()
        for inp in helper.get_all_input_files(input_path)
        for curve in payload.curves
        for group in payload.groups
    ]
    command = "".join(commands)
    return command


def default_case():
    raise ValueError("Framework not integrated into the benchmarking framework!")


# List ZKP-frameworks in the zk-Harness
projects = {
    "gnark":    build_command_gnark,
    "snarkjs":   build_command_snarkjs
}


def build_command(project, payload, count):
    """
    Build the command to execute the given project with the given payload.
    Input: project (e.g. gnark) + payload (config.json)
    """
    commands = projects.get(project, default_case)(payload, count)
    return commands


def get_ec_payload(config):
    """
    Extract the payload for category "ec" given a config.json
    """
    # Extract the relevant fields from the configuration data
    payload = config.get('payload')
    if payload is None:
        raise KeyError("Payload does not exist in ec config")

    curves = payload.get('curves')
    if curves is None:
        raise KeyError("curves field does not exist in ec payload")
    if len(curves) == 0:
        raise ValueError("curves field is empty")

    groups = payload.get('groups')
    if groups is None:
        raise KeyError("groups field does not exist in ec payload")

    operations = payload.get('operations')
    if operations is None:
        raise KeyError("operations field does not exist in ec payload")
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
    return Payload(curves, groups, ops)
