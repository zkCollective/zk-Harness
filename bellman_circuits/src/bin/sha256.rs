// use bellman_circuits::benches::benchmark_circuit; // Assuming this is the path to the bench_proof function
use bellman_circuits::circuits::sha256;
use clap::{Parser};
use rust_utils::{
    get_memory,
    read_file_contents,
    save_results,
};
use bellman_utils::measure_size_in_bytes;
use bellman::groth16;
use bellman::gadgets::multipack;
use bls12_381::Bls12;
use rand::rngs::OsRng;
use sha2::{Digest, Sha256};

#[derive(Parser, Debug)]
#[clap(
    name = "MemoryBenchSha256",
    about = "MemoryBenchSha256 CLI is a CLI Application to Benchmark memory consumption of SHA-256",
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

    let (preimage_length, preimage) = sha256::get_sha256_data(input_str);

    // Get the initial memory usage
    let initial_rss = get_memory();

    // Generate circuit parameters
    let params = {
        let c = sha256::Sha256Circuit { 
            preimage: Some(vec![0; preimage.len()]), 
            preimage_length: preimage.len() 
        };
        groth16::generate_random_parameters::<Bls12, _, _>(c, &mut OsRng).unwrap()
    };

    // Prepare the verification key
    let pvk = groth16::prepare_verifying_key(&params.vk);

    // Get the memory usage after setup
    let setup_rss = get_memory();

    // Compute the SHA256 hash of the preimage
    let hash = &Sha256::digest(&preimage);

    // Create a circuit instance
    let circuit = sha256::Sha256Circuit {
        preimage: Some(preimage),
        preimage_length: preimage_length
    };

    // Create a Groth16 proof with our parameters
    let proof = groth16::create_random_proof(circuit, &params, &mut OsRng).unwrap();

    // Get the memory usage after proof generation
    let proof_rss = get_memory();

    // Pack the hash as inputs for proof verification
    let hash_bits = multipack::bytes_to_bits_le(&hash);
    let inputs = multipack::compute_multipacking(&hash_bits);

    // Verify the proof
    assert!(groth16::verify_proof(&pvk, &proof, &inputs).is_ok());

    // Get the memory usage after proof verification
    let verify_rss = get_memory();

    // Measure the proof size
    let proof_size = measure_size_in_bytes(&proof);

    // Save the results
    save_results(initial_rss, setup_rss, proof_rss, verify_rss, proof_size, args.output);
}
