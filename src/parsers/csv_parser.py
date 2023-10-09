import argparse
import pandas as pd
import os

def parse_memory_txt(circuit_folder, circuit):
    memory_data = []
    for filename in os.listdir(circuit_folder):
        if filename.endswith(".txt"):
            with open(os.path.join(circuit_folder, filename), 'r') as file:
                lines = file.readlines()
                for line in lines:
                    if "maximum resident set size" in line:
                        ram = int(line.split()[0])
                        operation = filename.split("memory_")[-1].split("_")[0].replace(".txt", "")
                        input_file = "input/circuit/" + circuit + "/input_" + circuit_folder.split("_")[-1] + ".json" 
                        memory_data.append([input_file, operation, ram])
    df_memory = pd.DataFrame(memory_data, columns=["input", "operation", "ramReal"])
    return df_memory

def read_time_csv(filename):
    df_time = pd.read_csv(filename)
    return df_time

def merge_memory_time(df_memory, df_time):
    memory_dict = df_memory.set_index(['input', 'operation'])['ramReal'].to_dict()
    df_time['ramReal'] = df_time.set_index(['input', 'operation']).index.map(memory_dict.get)
    return df_time

def save_df(dataframe, filename):
    dataframe.to_csv(filename, index=False)

def combine_memory_time(circuit_folder, filename, circuit):
    df_memory = parse_memory_txt(circuit_folder, circuit)
    df_time = read_time_csv(filename)
    df_merged = merge_memory_time(df_memory, df_time)
    save_df(df_merged, filename)

if __name__ == "__main__":
    parser = argparse.ArgumentParser()
    parser.add_argument("--circuit", required=True, help="Circuit")
    parser.add_argument("--memory_folder", required=True, help="Memory Folder")
    parser.add_argument("--time_filename", required=True, help="Time Filename")
    args = parser.parse_args()
    combine_memory_time(args.memory_folder, args.time_filename, args.circuit)