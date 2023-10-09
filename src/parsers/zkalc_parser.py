import os
import argparse
import json
import csv


CRITERION_BENCHES = [
    #"ark-curves.json", #TODO
    "arkworks.json",
    "blstrs.json",
    #"curve25519-dalek.json",
    "pasta_curves.json",     
    "zkcrypto.json"
]
CSV_ARITHMETICS_HEADER = [
    "framework",
    "category",
    "curve",
    "field",
    "operation",
    "time",
    "count"
]
CSV_CURVES_HEADER = [
    "framework",
    "category",
    "curve",
    "operation",
    "input",
    "time",
    "count"
]

ARITHMETIC_OPERATIONS = ["add", "mul", "ivn"]
EC_OPERATIONS = ["msm", "pairing"]


def parse_criterion(file_path, framework):
    arithmetics_results = [CSV_ARITHMETICS_HEADER]
    ec_results = [CSV_CURVES_HEADER]
    with open(file_path, "r") as file:
        for line in file:
            data = json.loads(line)
            if data["reason"] != "benchmark-complete":
                continue

            # Parse ids
            id_parts = data["id"].split(" ")[0].split("/")
            curve = id_parts[0]
            operation = id_parts[1]
            field = id_parts[2] if len(id_parts) > 2 else ""
            if field == "" and "_" in operation:
                field = operation.split("_")[1]
                operation = operation.split("_")[0]
            op_input = id_parts[3] if len(id_parts) > 3 else ""

            # Check if unit in ns
            assert data["mean"]["unit"] == "ns"
            mean = int(data["mean"]["estimate"])
            count = len(data["iteration_count"])

            if operation in ARITHMETIC_OPERATIONS and field == "ff":
                arithmetics_results.append(
                    [framework,"arithmetics",curve,field,operation,mean,count]
                )

            if operation in EC_OPERATIONS:
                if field != "":
                    operation = field + "-" + operation
                operation = operation.lower()
                ec_results.append(
                    [framework,"ec",curve,operation,op_input,mean,count]
                )
    return arithmetics_results, ec_results


def process_json_files(input_dir, output_dir):
    # Check if input directory exists
    if not os.path.exists(input_dir):
        print(f"Input directory '{input_dir}' does not exist.")
        return

    # Safely create output directory if it doesn't exist
    os.makedirs(output_dir, exist_ok=True)
    arithmetics_dir = os.path.join(output_dir, "arithmetics")
    ec_dir = os.path.join(output_dir, "ec")

    # Iterate over JSON files in the input directory
    for filename in os.listdir(input_dir):
        if filename.endswith(".json"):
            file_path = os.path.join(input_dir, filename)
            print(f"Processing {file_path}")

            framework = filename.replace(".json", "")

            # Processing logic for each JSON file
            if filename in CRITERION_BENCHES:
                arithmetics, ec = parse_criterion(file_path, framework)
            else:
                print(f"Skip: no parser available for {filename}")
                continue

            if len(arithmetics) > 1:
                os.makedirs(arithmetics_dir, exist_ok=True)
                target = os.path.join(arithmetics_dir, framework + ".csv")
                with open(target, 'w') as f:
                    writer = csv.writer(f)
                    writer.writerows(arithmetics)

            if len(ec) > 1:
                os.makedirs(ec_dir, exist_ok=True)
                target = os.path.join(ec_dir, framework + ".csv")
                with open(target, 'w') as f:
                    writer = csv.writer(f)
                    writer.writerows(ec)
                

    print("Processing complete.")


def main():
    # Create the argument parser
    parser = argparse.ArgumentParser(description="Process JSON files.")

    # Add the required arguments
    parser.add_argument("--input", required=True, help="Input directory path.")
    parser.add_argument("--output", required=True, help="Output directory path.")

    # Parse the command-line arguments
    args = parser.parse_args()

    # Get the input and output directory paths from the command-line arguments
    input_dir = args.input
    output_dir = args.output

    # Call the function to process JSON files
    process_json_files(input_dir, output_dir)


if __name__ == "__main__":
    main()
