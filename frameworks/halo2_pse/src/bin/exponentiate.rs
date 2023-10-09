use halo2_proofs::halo2curves::bn256::Fr;
use halo2_pse_circuits::circuits::exponentiate::ExponentiationCircuit;
use halo2_pse_circuits::circuits::exponentiate::get_exponentiation_data;
use utilities::f_prove;
use utilities::f_setup;
use utilities::f_verify;
use utilities::read_file_contents;
use utilities::BinaryArgs;
use clap::Parser;

fn main () {
    let args = BinaryArgs::parse();
    let input_str = read_file_contents(args.input);
    let (k, e_value, x, e, y) = get_exponentiation_data(input_str);

    if args.phase == "setup" {
        let circuit = ExponentiationCircuit {
            row: e_value,
        };
        let params_file = args.params.expect("Missing params argument");
        let vk_file = args.vk.expect("Missing vk argument");
        let pk_file = args.pk.expect("Missing pk argument");
        f_setup(k, circuit, params_file, vk_file, pk_file);
    } else if args.phase == "prove" {
        let circuit = ExponentiationCircuit {
            row: e_value,
        };
        let public_input: &[&[Fr]] = &[&[x, e, y]];
        let params_file = args.params.expect("Missing params argument");
        let pk_file = args.pk.expect("Missing pk argument");
        let proof_file = args.proof.expect("Missing proof argument");
        f_prove(circuit, params_file, pk_file, public_input, proof_file);
    } else if args.phase == "verify" {
        let public_input: &[&[Fr]] = &[&[x, e, y]];
        let params_file = args.params.expect("Missing params argument");
        let vk_file = args.vk.expect("Missing vk argument");
        let proof_file = args.proof.expect("Missing proof argument");
        f_verify::<ExponentiationCircuit>(params_file, vk_file, proof_file, public_input)
    } else {
        panic!("Invalid phase (should be setup, prove, or verify)");
    }
}