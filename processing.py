#!/usr/bin/env python3
"""
This script recursively reads the logs from a specific path and produces 
pandas dataframes for circuits, arithmetics, and elliptic curves logs.

The script works as follows:
    1. Detects CSV files
    2. Parse logs from CSV files into LogRow objects
    3. Save the results to a Result object
    4. Retrieve pandas dataframes for each benchmark category

Deps:
    pip install pandas 
"""
import argparse
import logging
import os
import csv
import inspect

from dataclasses import dataclass
from collections import defaultdict

import pandas as pd


logging.basicConfig(
    format='%(asctime)s - %(levelname)s - %(message)s', 
    datefmt='%H:%M:%S', 
    level=logging.INFO
)


@dataclass
class CircuitsResults:
    nr_zkp: int
    nr_backends: int
    nr_circuits: int
    multiple_curves_barcharts: list


class Result:
    def __init__(self, rows):
        self.rows = rows

    def get_circuit_rows_as_df(self):
        rows = []
        headers = CircuitLogRow.get_headers()
        headers.remove("category")
        # TODO Merge rows that measure the exact same thing in a single row
        # by taking the mean to all related values
        # TODO should we keep a counter of the number of benchs?
        same_rows = defaultdict(list)
        for row in self.rows:
            if isinstance(row, CircuitLogRow):
                # If two rows have the same values for the following columns
                # then we should merge them.
                sig = row.get_static_rows()
                same_rows[sig].append(row)
        for v in same_rows.values():
            row = v[0] if len(v) == 1 else CircuitLogRow.merge_rows(v)
            rows.append(row.get_row())
        return pd.DataFrame(rows, columns=headers)


class LogRow:
    def __init__(self, framework):
        self.framework = framework

    @staticmethod
    def new_rows(reader, headers):
        """Read all lines and return a list with parsed rows.
        """
        rows = []
        cls = None
        for row in reader:
            if cls is None:
                if row[1] == "arithmetic":
                    raise NotImplementedError("Arithmetic not implemented")
                elif row[1] == "ec":
                    raise NotImplementedError("EC not implemented")
                elif row[1] == "circuit":
                    cls = CircuitLogRow
                else:
                    raise Exception("category (column 2) should be arithmetic, ec, or circuit")
                # Headers mapping
                def mapping(s):
                    mappings = [
                        ("input", "input_path"),
                        ("nbConstraints", "nb_constraints"),
                        ("nbSecret", "nb_secret"),
                        ("nbPublic", "nb_public"),
                        ("nbPhysicalCores", "nb_physical_cores"),
                        ("nbLogicalCores", "nb_logical_cores"),
                        ("ram(mb)", "ram"),
                        ("time(ms)", "time")
                    ]
                    for i, t in mappings:
                        s = s.replace(i, t)
                    return s
                # Check headers
                headers = list(map(mapping, headers))
                if headers != cls.get_headers():
                    raise Exception("Wrong headers:\nExpected: {}\nFound: {}".format(
                        cls.get_headers(), headers
                    ))
            rows.append(cls(*row))
        return rows

    @classmethod
    def get_headers(cls):
        header = inspect.getfullargspec(cls.__init__).args
        return list(filter(lambda x: x != 'self', header))

    def get_row(self):
        return [getattr(self, h) for h in self.get_headers() if h != "category"]

    def get_static_rows(self):
        """Return the values that should remain same across all executions
        in the same machine as string
        """
        raise NotImplementedError

    @classmethod
    def merge_rows(cls):
        """Merge multiple rows by returning the mean value for each column
        that can be variable
        """


class CircuitLogRow(LogRow):
    # We need catrgory to easily verify that we pass the correct number of args.
    def __init__(
        self, framework, category, backend, curve, circuit, input_path, operation,
        nb_constraints, nb_secret, nb_public, ram, time, nb_physical_cores,
        nb_logical_cores, cpu
    ):
        super().__init__(framework)
        # TODO sanity checks
        # Check for caps
        self.backend =  backend
        self.curve = curve
        self.circuit = circuit
        self.input_path = input_path
        self.operation = operation
        self.nb_constraints = int(nb_constraints)
        self.nb_secret = int(nb_secret)
        self.nb_public = int(nb_public)
        self.ram = int(ram)
        self.time = int(time)
        self.nb_physical_cores = int(nb_physical_cores)
        self.nb_logical_cores = int(nb_logical_cores)
        self.cpu = cpu

    def get_static_rows(self):
        return (f"{self.framework},{self.backend},{self.curve},{self.circuit},"
                f"{self.input_path},{self.operation},{self.nb_constraints},"
                f"{self.nb_secret},{self.nb_public},{self.nb_physical_cores},"
                f"{self.nb_logical_cores}")

    @classmethod
    def merge_rows(cls, rows):
        # TODO we can perform a sanity check to check if the non-variable values
        # are the same across all rows.
        time_values = []
        ram_values = []
        for row in rows:
            time_values.append(row.time)
            ram_values.append(row.ram)
        time = int(sum(time_values)/len(time_values))
        ram = int(sum(ram_values)/len(ram_values))
        row = rows[0]
        row.time = time
        row.ram = ram
        return row


def parse_logs(log_files):
    res = []
    for lf in log_files:
        logging.info(f"Parse {lf}")
        try:
            with open(lf, 'r') as fp:
                reader = csv.reader(fp)
                header = next(reader)
                if header[0] != "framework":
                    raise Exception(
                        "First row should contain the header and first column should be framework"
                    )
                # Check if the given headers match the expected headers
                res.extend(LogRow.new_rows(reader, header))
        except NotImplementedError as e:
            logging.error(f"Cannot parse file: {lf}")
            print("Exception:")
            print(e)
            #sys.exit(-1)
    return Result(res)


def detect_files_to_process(path):
    """Look for csv files in directory recursively
    """
    res = []
    for root,d_names,f_names in os.walk(path):
        for d in d_names:
            logging.info(f"Process directory {d}")
        res.extend([os.path.join(root, f) for f in f_names if f.endswith('.csv')])
    return res


def parse_arguments():
    """Parse the command-line arguements.
    """
    parser = argparse.ArgumentParser(
        description='Parse logs and generate report'
    )
    parser.add_argument(
        "logs", help="Path that contains the logs (it will search recursively)"
    )
    return parser.parse_args()


def analyse_logs(logs, level=logging.INFO):
    logging.getLogger().setLevel(level)
    # Detect files
    logging.info(f"Process {logs}")
    log_files = detect_files_to_process(logs)
    logging.info(f"Files to process: {len(log_files)}")

    # Parse logs
    logging.info(f"Parse log files")
    logs = parse_logs(log_files)

    # Circuits analyses
    circuits_df = logs.get_circuit_rows_as_df()
    return circuits_df


def main():
    args = parse_arguments()
    # TODO print some stats
    analyse_logs(args.logs)


if __name__ == "__main__":
    main()
