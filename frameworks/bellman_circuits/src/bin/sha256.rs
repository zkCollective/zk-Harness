// use bellman_circuits::benches::benchmark_circuit; // Assuming this is the path to the bench_proof function
use bellman_circuits::circuits::sha256;
use clap::{Parser};
use bellman::gadgets::multipack;
use sha2::{Digest, Sha256};
use bellman_utils::{BinaryArgs, f_setup, f_verify, f_prove, read_file_contents};

fn main() {
    // Parse command line arguments
    let args = BinaryArgs::parse();

    // Read and parse input from the specified JSON file
    let input_str = read_file_contents(args.input);

    let (preimage_length, preimage) = sha256::get_sha256_data(input_str);

    if args.phase == "setup" {
        let circuit = sha256::Sha256Circuit { 
            preimage: Some(vec![0; preimage.len()]), 
            preimage_length: preimage.len() 
        };
        let params_file = args.params.expect("Missing params argument");
        f_setup(circuit, params_file);
    } else if args.phase == "prove" {
        let circuit = sha256::Sha256Circuit {
            preimage: Some(preimage),
            preimage_length: preimage_length
        };
        let params_file = args.params.expect("Missing params argument");
        let proof_file = args.proof.expect("Missing proof argument");
        f_prove(circuit, params_file, proof_file);
    } else if args.phase == "verify" {
        // Compute the SHA256 hash of the preimage
        let hash = &Sha256::digest(&preimage);
        // Pack the hash as inputs for proof verification
        let hash_bits = multipack::bytes_to_bits_le(&hash);
        let inputs = multipack::compute_multipacking(&hash_bits);
        let params_file = args.params.expect("Missing params argument");
        let proof_file = args.proof.expect("Missing proof argument");
        f_verify(params_file, proof_file, inputs)
    } else {
        panic!("Invalid phase (should be setup, prove, or verify)");
    }
}
