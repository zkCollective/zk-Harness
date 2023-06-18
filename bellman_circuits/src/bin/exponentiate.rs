use bellman_circuits::circuits::exponentiate;
use clap::{Parser};
use rust_utils::{
    read_file_contents,
};
use bellman_utils::{BinaryArgs, f_setup, f_verify, f_prove};
use bellman::gadgets::multipack;
use bls12_381::Scalar;
use ff::PrimeField;

fn main() {
    // Parse command line arguments
    let args = BinaryArgs::parse();

    // Read and parse input from the specified JSON file
    let input_str = read_file_contents(args.input);

    // Get data from config
    let (x_64, e, y_64) = exponentiate::get_exponentiate_data(input_str);

    // Create Scalar from some values
    let x = Scalar::from(x_64);
    let y = Scalar::from(y_64);

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
        let x_bits = multipack::bytes_to_bits_le(&x.to_repr().as_ref());
        let y_bits = multipack::bytes_to_bits_le(&y.to_repr().as_ref());
        let inputs: Vec<Scalar> = [multipack::compute_multipacking(&x_bits), multipack::compute_multipacking(&y_bits)].concat();
        let params_file = args.params.expect("Missing params argument");
        let proof_file = args.proof.expect("Missing proof argument");
        f_verify(params_file, proof_file, inputs)
    } else {
        panic!("Invalid phase (should be setup, prove, or verify)");
    }


}
