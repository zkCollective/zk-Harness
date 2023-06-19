use halo2_proofs::halo2curves::bn256::Fr;
use halo2_pse_circuits::circuits::sha256::Sha256Circuit;
use halo2_pse_circuits::circuits::sha256::get_sha256_data;  
use utilities::f_prove;
use utilities::f_setup;
use utilities::f_verify;
use utilities::read_file_contents;
use utilities::BinaryArgs;
use clap::Parser;

fn main () {
    let args = BinaryArgs::parse();
    let input_str = read_file_contents(args.input);
    let (k, sha_data) = get_sha256_data(input_str);

    if args.phase == "setup" {
        let circuit = Sha256Circuit {
            sha_data: sha_data,
        };
        let params_file = args.params.expect("Missing params argument");
        let vk_file = args.vk.expect("Missing vk argument");
        let pk_file = args.pk.expect("Missing pk argument");
        f_setup(k, circuit, params_file, vk_file, pk_file);
    } else if args.phase == "prove" {
        let circuit = Sha256Circuit {
            sha_data: sha_data,
        };
        let public_input: &[&[Fr]] = &[];
        let params_file = args.params.expect("Missing params argument");
        let pk_file = args.pk.expect("Missing pk argument");
        let proof_file = args.proof.expect("Missing proof argument");
        f_prove(circuit, params_file, pk_file, public_input, proof_file);
    } else if args.phase == "verify" {
        let public_input: &[&[Fr]] = &[];
        let params_file = args.params.expect("Missing params argument");
        let vk_file = args.vk.expect("Missing vk argument");
        let proof_file = args.proof.expect("Missing proof argument");
        f_verify::<Sha256Circuit>(params_file, vk_file, proof_file, public_input)
    } else {
        panic!("Invalid phase (should be setup, prove, or verify)");
    }
}