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

                assert stage in ["setup", "prove", "verify"], "Invalid stage"

                mean = data.get("mean", {})
                assert mean.get("unit") == "ns", "Invalid unit"

                mean_time = mean.get("estimate", 0)
                count = sum(data.get("iteration_count", []))

                stages[stage] = {"mean": mean_time, "count": count}

    return circuit_name, stages


def compute_memory_usage(mem_proof_json, stages):
    with open(mem_proof_json, "r") as f:
        data = json.load(f)

    initial_rss = data.get("initial_rss")
    setup_rss = data.get("setup_rss")
    proof_rss = data.get("proof_rss")
    verify_rss = data.get("verify_rss")
    proof_size = data.get("proof_size")

    assert initial_rss is not None, "Missing 'initial_rss'"
    assert setup_rss is not None, "Missing 'setup_rss'"
    assert proof_rss is not None, "Missing 'proof_rss'"
    assert verify_rss is not None, "Missing 'verify_rss'"
    assert proof_size is not None, "Missing 'proof_size'"

    # FIXME That is definitely wrong
    setup_mem = setup_rss 
    proof_mem = proof_rss + initial_rss - setup_mem
    verify_mem = verify_rss + initial_rss - proof_mem

    stages["setup"]["ram"] = setup_mem
    stages["prove"]["ram"] = proof_mem
    stages["prove"]["proofSize"] = proof_size
    stages["verify"]["ram"] = verify_mem

    return stages


def save_csv(framework, category, backend, curve, circuit_name, input_path, stages, output_csv):
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
        "ram(mb)",
        "time(ms)",
        "proofSize",
        "nbPhysicalCores",
        "nbLogicalCores",
        "count",
        "cpu"
    ]

    write_header = not os.path.exists(output_csv)

    with open(output_csv, "a") as f:
        
        writer = csv.DictWriter(f, fieldnames=header)

        if write_header:
            writer.writeheader()

        for stage, data in stages.items():
            ram = data.get("ram", "")
            ram = int(ram / (1024 * 1024)) if isinstance(ram, int) else "" # convert bytes to mb
            mean = int(data.get("mean", "") / 1_000_000) # convert ns to ms
            mean = 1 if mean == 0 else mean

            row = {
                "framework": framework,
                "category": category,
                "backend": backend,
                "curve": curve,
                "circuit": circuit_name,
                "input": input_path,
                "operation": stage,
                # FIXME
                "nbConstraints": 1,
                # FIXME
                "nbSecret": 1,
                # FIXME
                "nbPublic": 1,
                "ram(mb)": ram,
                "time(ms)": mean, 
                "proofSize": int(data.get("proofSize", "")) if "proofSize" in data else "",
                # FIXME
                "nbPhysicalCores": 1,
                # FIXME
                "nbLogicalCores": 1,
                "count": data.get("count", ""),
                # FIXME
                "cpu": "MacOSX"
            }

            writer.writerow(row)

def combine_jsons_to_csv(framework, category, backend, curve, input_path, criterion_json, mem_proof_json, output_csv):
    circuit_name, stages = parse_criterion_json(criterion_json)
    stages = compute_memory_usage(mem_proof_json, stages)
    save_csv(framework, category, backend, curve, circuit_name, input_path, stages, output_csv)


if __name__ == "__main__":
    parser = argparse.ArgumentParser()
    parser.add_argument("--framework", required=True, help="Framework name")
    parser.add_argument("--category", required=True, help="Category")
    parser.add_argument("--backend", required=True, help="Backend")
    parser.add_argument("--curve", required=True, help="Curve")
    parser.add_argument("--input", required=True, help="Input")
    parser.add_argument("--criterion_json", required=True, help="Path to the criterion JSON file")
    parser.add_argument("--mem_proof_json", required=True, help="Path to the memory proof JSON file")
    parser.add_argument("--output_csv", required=True, help="Path to the output CSV file")
    args = parser.parse_args()

    combine_jsons_to_csv(
        args.framework, args.category, args.backend, args.curve,
        args.input, args.criterion_json, args.mem_proof_json, args.output_csv)
