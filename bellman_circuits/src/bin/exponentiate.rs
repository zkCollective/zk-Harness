// use bellman_circuits::benches::benchmark_circuit; // Assuming this is the path to the bench_proof function
use bellman_circuits::circuits::exponentiate;
use clap::{Parser};
use rust_utils::{
    get_memory,
    read_file_contents,
    save_results,
};
use bellman_utils::measure_size_in_bytes;
use bellman::groth16;
use bellman::gadgets::multipack;
use bls12_381::{Bls12, Scalar};
use rand::rngs::OsRng;
use ff::PrimeField;

#[derive(Parser, Debug)]
#[clap(
    name = "MemoryBenchExponentiate",
    about = "MemoryBenchExponentiate CLI is a CLI Application to Benchmark memory consumption of Exponentiate",
    version = "0.0.1"
)]

struct Args {
    #[arg(short, long)]
    input: String,

    #[arg(short, long)]
    output: String,
}

fn main() {
    // Parse command line arguments
    let args = Args::parse();

    // Read and parse input from the specified JSON file
    let input_str = read_file_contents(args.input);

     // Get data from config
     let (x_64, e, y_64) = exponentiate::get_exponentiate_data(input_str);

     // Create Scalar from some values
     let x = Scalar::from(x_64);
     let y = Scalar::from(y_64);

      // Public inputs are x and y
    let x_bits = multipack::bytes_to_bits_le(&x.to_repr().as_ref());
    let y_bits = multipack::bytes_to_bits_le(&y.to_repr().as_ref());
    let inputs = [multipack::compute_multipacking(&x_bits), multipack::compute_multipacking(&y_bits)].concat();


    // Define the circuit
    let circuit = exponentiate::ExponentiationCircuit {
        x: Some(x),
        e: e,
        y: Some(y),
    };

    // Get the initial memory usage
    let initial_rss = get_memory();

    // Generate Parameters
    let params = groth16::generate_random_parameters::<Bls12, _, _>(circuit.clone(), &mut OsRng).unwrap();

    // Prepare the verification key
    let pvk = groth16::prepare_verifying_key(&params.vk);

    // Get the memory usage after setup
    let setup_rss = get_memory();

    // Create a Groth16 proof with our parameters
    let proof = groth16::create_random_proof(circuit, &params, &mut OsRng).unwrap();

    // Get the memory usage after proof generation
    let proof_rss = get_memory();

    // Verify the proof
    let _ = groth16::verify_proof(&pvk, &proof, &inputs);

    // Get the memory usage after proof verification
    let verify_rss = get_memory();

    // Measure the proof size
    let proof_size = measure_size_in_bytes(&proof);

    // Save the results
    save_results(initial_rss, setup_rss, proof_rss, verify_rss, proof_size, args.output);
}
