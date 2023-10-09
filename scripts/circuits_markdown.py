import argparse
import json

def generate_markdown_table(json_data):
    # Extract circuits and frameworks from the JSON data
    circuits = json_data["circuits"]
    frameworks = json_data["frameworks"]

    # Create the header row of the table
    header_row = "|          |"
    for circuit in circuits:
        header_row += f" {circuit} {' ' * (18 - len(circuit))} |"
    header_row += "\n"

    # Create the separator row of the table
    separator_row = "| -------- |"
    for _ in circuits:
        separator_row += " ------------------- |"
    separator_row += "\n"

    # Create the rows for each framework
    framework_rows = ""
    for framework in frameworks:
        framework_name = framework["name"]
        framework_url = framework["url"]

        framework_row = f"| [{framework_name}]({framework_url}) |"

        for circuit in circuits:
            circuit_entry = next(
                (
                    c
                    for c in framework["circuits"]
                    if c["name"] == circuit
                ),
                None,
            )

            if circuit_entry:
                custom = circuit_entry["implementation"]["custom"]
                url = circuit_entry["implementation"]["url"]
                if custom:
                    framework_row += " :heavy_check_mark: (custom) |"
                else:
                    framework_row += f" :heavy_check_mark: ([implementation]({url})) |"
            else:
                framework_row += " :x: |"

        framework_row += "\n"
        framework_rows += framework_row

    # Generate the markdown table
    markdown_table = header_row + separator_row + framework_rows

    return markdown_table

def main():
    # Parse the command-line arguments
    parser = argparse.ArgumentParser(description="Generate a markdown table from JSON data.")
    parser.add_argument("input_file", help="Path to the input JSON file")
    args = parser.parse_args()

    # Read the JSON data from the input file
    with open(args.input_file) as file:
        json_data = json.load(file)

    # Generate the markdown table
    markdown_table = generate_markdown_table(json_data)

    # Print the markdown table
    print(markdown_table)

if __name__ == "__main__":
    main()
