import argparse
import json

def generate_markdown_table(json_data):
    # Create the header row of the table
    header_row = "|           | Language | Curves/Fields | Frameworks |\n"
    header_row += "| --------- | -------- | ------------- | ---------- |\n"

    # Create the rows for each entry
    entry_rows = ""
    for entry in json_data:
        library_name = entry["name"]
        library_url = entry["url"]
        language = entry["language"]
        curves_fields = ", ".join(entry["curves/fields"])
        frameworks = []
        if "framework" in entry:
            for framework in entry["framework"]:
                framework_name = framework["name"]
                framework_url = framework["url"]
                frameworks.append(f"[{framework_name}]({framework_url})")
        frameworks_str = ", ".join(frameworks)

        entry_row = f"| [{library_name}]({library_url}) | {language} | {curves_fields} | {frameworks_str} |\n"
        entry_rows += entry_row

    # Generate the markdown table
    markdown_table = header_row + entry_rows

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
