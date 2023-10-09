"""
CSV Parser Script
----------------

This script is used to parse a CSV file containing benchmark data and update it with RAM values obtained from memory files.

Usage:
------
python3 _scripts/parsers/csv_parser.py --memory_folder <memory_folder_path> --time_filename <csv_file_path> --circuit <circuit_name>

Arguments:
----------
--memory_folder: The folder path containing the memory files.
--time_filename: The path to the CSV file to be parsed and updated.
--circuit: The name of the circuit to filter rows in the CSV file.

The script reads the CSV file, searches for rows matching the specified circuit name, extracts the RAM values from the corresponding memory files,
and adds the RAM values to the "ram" column in the CSV file. The updated CSV file will be saved with the changes.

Example Usage:
--------------
python3 _scripts/parsers/csv_parser.py --memory_folder benchmarks/halo2_pse/memory/_input/circuit/exponentiate/input_100.json --time_filename benchmarks/halo2_pse/halo2_pse_bn256_exponentiate.csv --circuit exponentiate

"""
import argparse
import csv
import os
import re

def extract_ram_from_file(filename):
    with open(filename, 'r') as file:
        content = file.read()
        match = re.search(r'(\d+)\s+maximum resident set size', content)
        if match:
            return match.group(1)
        else:
            match = re.search('Maximum resident set size \(kbytes\):\s+(\d+)', content)
            if match: 
                kb = match.group(1)
                return int(kb) * 1024  # convert kb to bytes
    return ''

def parse_csv(csv_filename, memory_folder, circuit):
    csv_rows = []
    with open(csv_filename, 'r') as file:
        csv_reader = csv.DictReader(file)
        for row in csv_reader:
            input_name_from_memory = os.path.basename(memory_folder).split(".")[0]
            input_name_from_row = os.path.basename(row['input']).split('.')[0]
            if row['circuit'] == circuit and input_name_from_memory == input_name_from_row:
                files =  os.listdir(memory_folder)
                memory_filename = next((f for f in files if row["operation"] in f), None)
                if memory_filename is None:
                    print(f"No file found for operation: {row['operation']}")
                    ram = ''
                else:
                    memory_file = os.path.join(memory_folder, memory_filename)
                    ram = extract_ram_from_file(memory_file)
                row['ram'] = ram
            csv_rows.append(row)

    # Write updated rows back to the CSV file
    with open(csv_filename, 'w', newline='') as file:
        fieldnames = csv_rows[0].keys()
        csv_writer = csv.DictWriter(file, fieldnames)
        csv_writer.writeheader()
        csv_writer.writerows(csv_rows)

def main():
    parser = argparse.ArgumentParser(description='CSV Parser')
    parser.add_argument('--memory_folder', type=str, help='Folder containing memory files')
    parser.add_argument('--time_filename', type=str, help='CSV file path')
    parser.add_argument('--circuit', type=str, help='Circuit name')
    args = parser.parse_args()

    parse_csv(args.time_filename, args.memory_folder, args.circuit)

if __name__ == '__main__':
    main()
