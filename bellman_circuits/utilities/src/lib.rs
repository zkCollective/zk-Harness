use serde::Serialize;
use std::{env, fs::File};
use std::io::Read;
use std::process;
use psutil;
use std::fs;
use bellman::groth16::Proof;
use bls12_381::Bls12;

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

pub fn measure_size_in_bytes(proof: &Proof<Bls12>) -> usize {
    // TODO: Should we serialize the proof in another format?
    // Serialize data and save it to a temporary file
    let temp_file_path = "temp_file.bin";

    let mut file = File::create(temp_file_path).expect("Could not create temp file");
    proof.write(&mut file).expect("Could not write proof to file");

    // Measure the size of the file
    let file_size: usize = fs::metadata(&temp_file_path).expect("Cannot read the size of temp file").len() as usize;

    // Convert file size to MB
    let size_in_mb = file_size;

    // Remove the temporary file
    fs::remove_file(&temp_file_path).expect("Cannot remove temp file");

    return size_in_mb;
}