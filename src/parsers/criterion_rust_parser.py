"""
NOTE: Currently this is specific to halo2_pse and only for circuits.
"""
import argparse
import json
import csv
import os


def parse_criterion_json(json_file):
    circuit_name = ""
    stages = {}

    with open(json_file, "r") as f:
        for line in f:
            data = json.loads(line)

            reason = data.get("reason")
            if reason == "benchmark-complete":
                id_parts = data.get("id").split("/")
                assert len(id_parts) == 2, "Invalid ID format"
                circuit_name, stage = id_parts

                assert stage in ["setup", "witness", "prove", "verify"], "Invalid stage"

                mean = data.get("mean", {})
                assert mean.get("unit") == "ns", "Invalid unit"

                mean_time = mean.get("estimate", 0)
                count = sum(data.get("iteration_count", []))

                stages[stage] = {"mean": mean_time, "count": count}

    return circuit_name, stages


def get_file_size(file_path):
    try:
        size = os.path.getsize(file_path)
        return size
    except OSError as e:
        print(f"Error: {e}")
        return None


def save_csv(framework, category, backend, curve, circuit_name, input_path, stages, output_csv, proof_size):
    header = [
        "framework",
        "category",
        "backend",
        "curve",
        "circuit",
        "input",
        "operation",
        "nbConstraints",
        "nbSecret",
        "nbPublic",
        "ram",
        "time",
        "proofSize",
        "count"
    ]

    write_header = not os.path.exists(output_csv)

    with open(output_csv, "a") as f:
        
        writer = csv.DictWriter(f, fieldnames=header)

        if write_header:
            writer.writeheader()

        for stage, data in stages.items():
            ram = data.get("ram", "")
            ram = int(ram) if ram is not None and isinstance(ram, int) else ""
            mean = data.get("mean", "")
            mean = int(mean / 1_000_000) if mean is not None else ""  # convert ns to ms
            mean = 1 if mean == 0 else mean

            proofSize = int(proof_size) if proof_size is not None else ""

            row = {
                "framework": framework,
                "category": category,
                "backend": backend,
                "curve": curve,
                "circuit": circuit_name,
                "input": input_path,
                "operation": stage,
                # FIXME
                "nbConstraints": "",
                # FIXME
                "nbSecret": "",
                # FIXME
                "nbPublic": "",
                "ram": ram,
                "time": mean, 
                "proofSize": proofSize if stage == "prove" else "",
                "count": data.get("count", "")
            }

            writer.writerow(row)

def combine_jsons_to_csv(framework, category, backend, curve, input_path, criterion_json, output_csv, proof_file):
    circuit_name, stages = parse_criterion_json(criterion_json)
    if proof_file:
        proof_size = get_file_size(proof_file)
    else:
        proof_size = None
    save_csv(framework, category, backend, curve, circuit_name, input_path, stages, output_csv, proof_size)


if __name__ == "__main__":
    parser = argparse.ArgumentParser()
    parser.add_argument("--framework", required=True, help="Framework name")
    parser.add_argument("--category", required=True, help="Category")
    parser.add_argument("--backend", required=True, help="Backend")
    parser.add_argument("--curve", required=True, help="Curve")
    parser.add_argument("--input", required=True, help="Input")
    parser.add_argument("--criterion_json", required=True, help="Path to the criterion JSON file")
    parser.add_argument("--proof", required=False, help="Path to the proof file")
    parser.add_argument("--output_csv", required=True, help="Path to the output CSV file")
    args = parser.parse_args()

    combine_jsons_to_csv(
        args.framework, args.category, args.backend, args.curve,
        args.input, args.criterion_json, args.output_csv, args.proof)
