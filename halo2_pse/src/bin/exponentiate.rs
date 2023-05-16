use halo2_proofs::halo2curves::bn256::Fr;
use halo2_pse_circuits::circuits::exponentiate::ExponentiationCircuit;
use halo2_pse_circuits::circuits::exponentiate::get_exponentiation_data;
use utilities::get_memory;
use utilities::measure_size_in_bytes;
use utilities::prove_circuit;
use utilities::read_file_contents;
use utilities::save_results;
use utilities::setup_circuit;
use utilities::verify_circuit;
use clap::Parser;

#[derive(Parser, Debug)]
#[command(author, version, about, long_about = None)]
struct Args {
    #[arg(short, long)]
    input: String,

    #[arg(short, long)]
    output: String,
}

fn main () {
    let args = Args::parse();
    let input_str = read_file_contents(args.input);
    let (k, e_value, x, e, y) = get_exponentiation_data(input_str);

    let initial_rss = get_memory();

    let circuit = ExponentiationCircuit {
        row: e_value,
    };
    let (params, vk, pk) = setup_circuit(k, circuit.clone());

    let setup_rss = get_memory();

    let public_input: &[&[Fr]] = &[&[x, e, y]];
    let proof = prove_circuit(&params, &pk, circuit, public_input);

    let proof_rss = get_memory();

    verify_circuit::<ExponentiationCircuit>(&params, &vk, &proof, public_input);

    let verify_rss = get_memory();

    let proof_size = measure_size_in_bytes(&proof);

    save_results(initial_rss, setup_rss, proof_rss, verify_rss, proof_size, args.output);
}
