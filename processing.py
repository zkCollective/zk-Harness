#!/usr/bin/env python3
"""
This script recursively reads the logs from a specific path and generates 
a report with the results.

The script works as follows:
    1. Detects CSV files
    2. Parse logs from CSV files into LogRow objects
    3. Save the results to a Result object
    4. Retrieve pandas dataframes for each benchmark category
    5. Perfom the following analyses per category:
        A. Arithmetic
        B. EC
        C. Circuit 
            i. 

Deps:
    pip install altair-viewer altair pandas matplotlib
"""
import argparse
from copy import deepcopy
import logging
import os
import csv
import inspect
import sys

from dataclasses import dataclass

import pandas as pd
import matplotlib.pyplot as plt
import altair as alt
import datapane as dp
import plotly.graph_objects as go
import plotly.express as px
import ipywidgets as widgets

GITHUB_REPO = "https://github.com/XXX/YYY"
HTML_HEADER = """
<h1 style="text-align:center;font-family:Georgia;font-variant:small-caps;font-size: 70px;color:#0F1419;">
ZKP Libraries Benchmarking
</h1>
<div><div style ="float:right;font-family: Courier New, monospace;">
View Source on <a href="{}">Github</a>
</div></div>
""".format(GITHUB_REPO)

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


def produce_html(output, stats_group, circuits_results):
    app = dp.Report(
            HTML_HEADER,
            stats_group,
            # FIXME
            dp.Plot(circuits_results.multiple_curves_barcharts[0], responsive=True)
    )
    app.save(path = output,
        formatting=dp.ReportFormatting(
        bg_color="#EEE"
    ))


def analyze_circuits(df):
    # General Statistics
    unique_frameworks = df['framework'].unique()
    unique_backends = df['backend'].unique()
    unique_circuits = df['circuit'].unique()

    # TODO FIXME
    # Filter the dataframe to only include rows where framework = 'gnark', circuit = 'cubic', and backend = 'groth16'
    gnark_cubic_groth16_df = df[(df['framework'] == 'gnark') & (df['circuit'] == 'cubic') & (df['backend'] == 'groth16')]

    # Create a bar chart using Plotly
    fig = go.Figure()
    colors = ['blue', 'orange', 'green']
    for i, op in enumerate(gnark_cubic_groth16_df['operation'].unique()):
        curve_values = gnark_cubic_groth16_df[gnark_cubic_groth16_df['operation'] == op]['curve'].tolist()
        time_values = gnark_cubic_groth16_df[gnark_cubic_groth16_df['operation'] == op]['time'].tolist()
        fig.add_trace(go.Bar(x=curve_values,
                             y=time_values,
                             name=op,
                             marker_color=colors[i]))
    fig.update_layout(title='Time vs Curve for Gnark operations on Cubic circuit with Groth16 backend',
                      xaxis_title='Curve',
                      yaxis_title='Time',
                      barmode='relative',
                      bargap=0.1)
    # For every framework supporting multiple curves, get a barchart to compare
    # them
    return CircuitsResults(
        len(unique_frameworks), len(unique_backends), len(unique_circuits), 
        [fig]
    )


class Result:
    def __init__(self, rows):
        self.rows = rows

    def get_circuit_rows_as_df(self):
        headers = CircuitLogRow.get_headers()
        rows = [
            r.row for r in self.rows if isinstance(r, CircuitLogRow)
        ]
        return pd.DataFrame(rows, columns=headers)


class LogRow:
    def __init__(self, framework):
        self.framework = framework

    @staticmethod
    def new(args):
        if args[1] == "arithmetic":
            raise NotImplementedError
        if args[1] == "ec":
            raise NotImplementedError
        if args[1] == "circuit":
            return CircuitLogRow(*args)
        raise Exception("category (column 2) should be arithmetic, ec, or circuit")

    @classmethod
    def get_headers(cls):
        header = inspect.getfullargspec(cls.__init__).args
        return list(filter(lambda x: x not in ('self', '_category'), header))


class CircuitLogRow(LogRow):
    # We need catrgory to easily verify that we pass the correct number of args.
    def __init__(
        self, framework, _category, backend, curve, circuit, input_path, operation,
        nb_constraints, nb_secret, nb_public, ram, time, nb_physical_cores,
        nb_logical_cores, cpu
    ):
        super().__init__(framework)
        # TODO sanity checks
        # Check for caps
        # The order should be the exact same order we use in parameters
        self.row = [
            framework, backend, curve, circuit, input_path, operation,
            int(nb_constraints), int(nb_secret), int(nb_public), int(ram), 
            int(time), int(nb_physical_cores), int(nb_logical_cores), cpu
        ]


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
                for row in reader:
                    res.append(LogRow.new(row))
        except Exception as e:
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
    parser.add_argument(
        "-o",
        "--output",
        default="index.html",
        help="Html file to save the report (default: index.html)"
    )
    return parser.parse_args()


def main():
    args = parse_arguments()

    # Detect files
    logging.info(f"Process {args.logs}")
    log_files = detect_files_to_process(args.logs)
    logging.info(f"Files to process: {len(log_files)}")

    # Parse logs
    logging.info(f"Parse log files")
    logs = parse_logs(log_files)

    # Circuits analyses
    circuits_df = logs.get_circuit_rows_as_df()
    circuits_results = analyze_circuits(circuits_df)

    stats_group = dp.Group(
        dp.BigNumber(heading="ZKP Languages / Libraries", value=circuits_results.nr_zkp),
        dp.BigNumber(heading="Nr. of Backends", value=circuits_results.nr_backends),
        dp.BigNumber(heading="Nr. of Test Circuits", value=circuits_results.nr_circuits),
        columns=3,
    )
    
    produce_html(args.output, stats_group, circuits_results)

if __name__ == "__main__":
    main()
