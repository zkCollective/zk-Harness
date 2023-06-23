use starky::serialization::Buffer;
use clap::Parser;

use starky_utils;
use starky::{
    prover::prove
};
use plonky2::plonk::config::{GenericConfig, PoseidonGoldilocksConfig};
use plonky2::util::timing::TimingTree;
use plonky2::field::types::Field;

use starky_circuits::circuits::exponentiate::{
    exponentiate,
    get_exponentiate_data
};


#[derive(Parser, Debug)]
#[clap(
    name = "StarkyMemoryBenchExponentiateWitness",
    about = "StarkyMemoryBenchExponentiateWitness CLI is a CLI Application to Benchmark memory consumption of Exponentiate in starky",
    version = "0.0.1"
)]

struct Args {
    #[arg(short, long)]
    input: String,
}

fn main() {
    // Parse command line arguments
    let args = Args::parse();

    const D: usize = 2;
    type C = PoseidonGoldilocksConfig;
    type F = <C as GenericConfig<D>>::F;
    type S = starky_circuits::circuits::exponentiate::ExponentiateStark<F, D>;
    let config = starky_utils::standard_fast_config();

    // Read and parse input from the specified JSON file
    let input_str = starky_utils::read_file_contents(args.input);
    
    // Get data from config
    let (x, e, _y) = get_exponentiate_data::<PoseidonGoldilocksConfig, 2>(input_str);
    
    // Compute Trace
    let num_rows = 1 << e;
    let public_inputs = [x, F::from_canonical_usize(num_rows), exponentiate(num_rows - 1, x)];
    let stark = S::new(num_rows);

    let trace = stark.generate_trace(public_inputs[0], public_inputs[1], public_inputs[2]);

    // Compute Proof
    let proof = prove::<F, C, S, D>(
        stark,
        &config,
        trace.clone(),
        public_inputs,
        &mut TimingTree::default(),
    ).unwrap();

    // Serialization of Proof
    let mut buffer = Buffer::new(Vec::new());
    let _ = buffer.write_stark_proof_with_public_inputs(&proof);
    // println!("proof size {}\n", buffer.bytes().len());

}
