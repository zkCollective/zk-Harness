use std::{env, fs::File};
use std::process;
use std::io::Read;
use std::fs;
use bellman::{Circuit, groth16};
use bellman::groth16::{Proof, Parameters};
use bls12_381::{Bls12, Scalar};
use clap::Parser;
use rand::rngs::OsRng;


#[derive(Parser, Debug)]
#[command(author, version, about, long_about = None)]
pub struct BinaryArgs {
    #[arg(short, long)]
    pub input: String,

    #[arg(short, long)]
    pub phase: String,

    #[arg(short, long)]
    pub params: Option<String>,

    #[arg(short, long)]
    pub proof: Option<String>,
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

pub fn read_file_contents(file_name: String) -> String {
    let mut file = File::open(file_name).expect("Cannot load file");
    let mut file_str = String::new();
    file.read_to_string(&mut file_str).expect("Cannot read file");
    return file_str;
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

pub fn save_params(params_file: String, params: Parameters<Bls12>) {
    let mut file = File::create(&params_file).expect("Failed to create file");
    // Write the init_params to the file
    params.write(&mut file).expect("Failed to write params to file");
}

pub fn load_params(params_file: String) -> Parameters<Bls12> {
    let mut file = File::open(&params_file).expect("Failed to open file");
    Parameters::read(&mut file, true).expect("Failed to read params from file")
}

pub fn save_proof(proof_file: String, proof: Proof<Bls12>) {
    let mut file = File::create(&proof_file).expect("Failed to create file");
    // Write the proof to the file
    proof.write(&mut file).expect("Failed to write proof to file");
}

pub fn load_proof(proof_file: String) -> Proof<Bls12> {
    let mut file = File::open(&proof_file).expect("Failed to open file");
    Proof::read(&mut file).expect("Failed to read proof from file")
}

pub fn f_setup<C: Circuit<Scalar> + Clone>(circuit: C, params_file: String) {
    let params = groth16::generate_random_parameters::<Bls12, _, _>(circuit.clone(), &mut OsRng).unwrap();
    save_params(params_file, params);
}

pub fn f_prove<C: Circuit<Scalar> + Clone>(circuit: C, params_file: String, proof_file: String) {
    let params = load_params(params_file);
    let proof = groth16::create_random_proof(circuit, &params, &mut OsRng).unwrap();
    save_proof(proof_file, proof);
}

pub fn f_verify(params_file: String, proof_file: String, public_input: Vec<Scalar>) {
    let params = load_params(params_file);
    let pvk = groth16::prepare_verifying_key(&params.vk);
    let proof = load_proof(proof_file);
    assert!(groth16::verify_proof(&pvk, &proof, &public_input).is_ok());
}
