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

    def get_rows_as_df(self, cls):
        rows = []
        headers = cls.get_headers()
        headers.remove("category")
        # TODO Merge rows that measure the exact same thing in a single row
        # by taking the mean to all related values
        # TODO should we keep a counter of the number of benchs?
        same_rows = defaultdict(list)
        for row in self.rows:
            # Get only rows of cls
            if isinstance(row, cls):
                # If two rows have the same values for the static columns
                # then we should merge them.
                sig = row.get_static_rows()
                same_rows[sig].append(row)
        for v in same_rows.values():
            row = v[0] if len(v) == 1 else cls.merge_rows(v)
            rows.append(row.get_row())
        df = pd.DataFrame(rows, columns=headers)
        cls.check_count(df)
        return df


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
                    cls = ArithmeticLogRow
                elif row[1] == "ec":
                    cls = ECLogRow
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
                        ("proofSize", "proof"),
                        ("ram(mb)", "ram"),
                        ("time(ms)", "time"),
                        ("p(bitlength)", "p"),
                        ("time(ns)", "time")
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
    def merge_rows(cls, rows):
        """Merge multiple rows by returning the mean value for each column
        that can be variable
        """
        # TODO we can perform a sanity check to check if the non-variable values
        # are the same across all rows.
        time_values = []
        ram_values = []
        count_values = []
        for row in rows:
            time_values.append(row.time)
            ram_values.append(row.ram)
            count_values.append(row.count)
        time = int(sum(time_values)/len(time_values))
        ram = int(sum(ram_values)/len(ram_values))
        count = int(sum(count_values))
        row = rows[0]
        row.time = time
        row.ram = ram
        row.count = count
        return row

    @classmethod
    def check_count(cls, df):
        raise NotImplementedError()


class CircuitLogRow(LogRow):
    # We need catrgory to easily verify that we pass the correct number of args.
    def __init__(
        self, framework, category, backend, curve, circuit, input_path, operation,
        nb_constraints, nb_secret, nb_public, ram, time, proof, nb_physical_cores,
        nb_logical_cores, count, cpu
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
        self.proof = int(proof) if proof != '' else 0
        self.nb_physical_cores = int(nb_physical_cores)
        self.nb_logical_cores = int(nb_logical_cores)
        self.count = int(count)
        self.cpu = cpu

    def get_static_rows(self):
        return (f"{self.framework},{self.backend},{self.curve},{self.circuit},"
                f"{self.input_path},{self.operation},{self.nb_constraints},"
                f"{self.nb_secret},{self.nb_public},"
                f"{self.nb_physical_cores},{self.nb_logical_cores},{self.count},"
                f"{self.cpu}")

    @classmethod
    def check_count(cls, df):
        dff = df[['circuit', 'input_path', 'operation', 'count']]
        dff_grouped = dff.groupby(['circuit', 'input_path', 'operation'])
        # This check is not required if we use criterion or other such framework
        #for name, group in dff_grouped:
        #    if len(group['count'].unique()) != 1:
        #        raise AssertionError(f"Each experiment in group `{name}` should have the same count")


class ArithmeticLogRow(LogRow):
    # We need category to easily verify that we pass the correct number of args.
    def __init__(
        self, framework, category, curve, field, operation, input_path, ram, time,
        nb_physical_cores, nb_logical_cores, count, cpu
    ):
        super().__init__(framework)
        # TODO sanity checks
        # Check for caps
        self.curve =  curve
        self.field =  field
        self.operation = operation
        self.input_path = input_path
        self.ram = int(ram)
        self.time = int(time)
        self.nb_physical_cores = int(nb_physical_cores)
        self.nb_logical_cores = int(nb_logical_cores)
        self.count = int(count)
        self.cpu = cpu

    def get_static_rows(self):
        return (f"{self.framework},{self.curve},{self.field},{self.operation},"
                f"{self.input_path},{self.nb_physical_cores},"
                f"{self.nb_logical_cores},{self.count},{self.cpu}")


    @classmethod
    def check_count(cls, df):
        dff = df[['field', 'input_path', 'operation', 'count']]
        dff_grouped = dff.groupby(['field', 'input_path', 'operation'])
        for name, group in dff_grouped:
            if len(group['count'].unique()) != 1:
                raise AssertionError(f"Each experiment in group `{name}` should have the same count")


class ECLogRow(LogRow):
    # We need category to easily verify that we pass the correct number of args.
    def __init__(
        self, framework, category, curve, operation, input_path, ram, time,
        nb_physical_cores, nb_logical_cores, count, cpu
    ):
        super().__init__(framework)
        # TODO sanity checks
        # Check for caps
        self.curve =  curve
        self.operation = operation
        self.input_path = input_path
        self.ram = int(ram)
        self.time = int(time)
        self.nb_physical_cores = int(nb_physical_cores)
        self.nb_logical_cores = int(nb_logical_cores)
        self.count = int(count)
        self.cpu = cpu

    def get_static_rows(self):
        return (f"{self.framework},{self.curve},{self.operation},"
                f"{self.input_path},{self.nb_physical_cores},"
                f"{self.nb_logical_cores},{self.count},{self.cpu}")

    @classmethod
    def check_count(cls, df):
        dff = df[['input_path', 'operation', 'count']]
        dff_grouped = dff.groupby(['input_path', 'operation'])
        for name, group in dff_grouped:
            if len(group['count'].unique()) != 1:
                raise AssertionError(f"Each experiment in group `{name}` should have the same count")


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
    circuits_df = logs.get_rows_as_df(CircuitLogRow)
    arithmetics_df = logs.get_rows_as_df(ArithmeticLogRow)
    ec_df = logs.get_rows_as_df(ECLogRow)
    return circuits_df, arithmetics_df, ec_df


def main():
    args = parse_arguments()
    # TODO print some stats
    analyse_logs(args.logs)


if __name__ == "__main__":
    main()
else:
    current_dir = os.path.abspath(os.path.join(os.path.dirname(__file__), '..'))
    benchmarks_dir = current_dir + "/benchmarks"
    circuits_df, arithmetics_df, ec_df = analyse_logs(benchmarks_dir, logging.ERROR)
