use serde::Serialize;
use std::{env, fs::File};
use std::io::Read;
use std::process;
use psutil;

#[derive(Serialize)]
struct Results {
    initial_rss: u64,
    setup_rss: u64,
    proof_rss: u64,
    verify_rss: u64,
    proof_size: usize,
}

pub fn read_file_contents(file_name: String) -> String {
    let mut file = File::open(file_name).expect("Cannot load file");
    let mut file_str = String::new();
    file.read_to_string(&mut file_str).expect("Cannot read file");
    return file_str;
}

pub fn read_file_from_env_var(env_var_name: String) -> String {
    let input_file = env::var(env_var_name.clone()).unwrap_or_else(|_| {
        println!("Please set the {} environment variable to point to the input file", env_var_name);
        process::exit(1);
    });
    return read_file_contents(input_file);
}

pub fn read_env_variable(env_var_name: String) -> String {
    let variable_str = env::var(env_var_name.clone()).unwrap_or_else(|_| {
        println!("Please set the {} environment variable", env_var_name);
        process::exit(1);
    });
    return variable_str;
}


/// Return current RSS memory in bytes
pub fn get_memory() -> u64 {
    let current_process = psutil::process::Process::current().expect("Cannot get current process");
    let mem = current_process.memory_info().expect("Cannot get memory info");
    return mem.rss();
}

pub fn save_results(
    initial_rss: u64,
    setup_rss: u64,
    proof_rss: u64,
    verify_rss: u64,
    proof_size: usize,
    output_file_str: String,
) {
    // Create a Results struct with the provided data
    let results = Results {
        initial_rss,
        setup_rss,
        proof_rss,
        verify_rss,
        proof_size,
    };

    // Serialize the Results struct to JSON
    let serialized = serde_json::to_string(&results).expect("Could not serialize results");

    // Create a file and write the JSON data
    std::fs::write(
        output_file_str,
        serialized
    ).expect("Could not write to file");
}