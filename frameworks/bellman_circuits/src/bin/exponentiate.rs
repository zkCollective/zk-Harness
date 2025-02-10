use bellman_circuits::circuits::exponentiate;
use clap::{Parser};
use bellman_utils::{BinaryArgs, f_setup, f_verify, f_prove, read_file_contents};
use bellman::gadgets::multipack;
use bls12_381::Scalar;
use ff::PrimeField;

fn main() {
    // Parse command line arguments
    let args = BinaryArgs::parse();

    // Read and parse input from the specified JSON file
    let input_str = read_file_contents(args.input);

    // Get data from config
    let (x, e, y) = exponentiate::get_exponentiate_data(input_str);

    // Public input is empty
    let inputs = Vec::new();

    if args.phase == "setup" {
        let circuit = exponentiate::ExponentiationCircuit {
            x: Some(x),
            e: e,
            y: Some(y),
        };
        let params_file = args.params.expect("Missing params argument");
        f_setup(circuit, params_file);
    } else if args.phase == "prove" {
        let circuit = exponentiate::ExponentiationCircuit {
            x: Some(x),
            e: e,
            y: Some(y),
        };
        let params_file = args.params.expect("Missing params argument");
        let proof_file = args.proof.expect("Missing proof argument");
        f_prove(circuit, params_file, proof_file);
    } else if args.phase == "verify" {
        // Public inputs are x and y
        let params_file = args.params.expect("Missing params argument");
        let proof_file = args.proof.expect("Missing proof argument");
        f_verify(params_file, proof_file, inputs)
    } else {
        panic!("Invalid phase (should be setup, prove, or verify)");
    }


}
