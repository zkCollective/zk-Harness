use clap::Parser;
use starky_utils;
use plonky2::plonk::config::{GenericConfig, PoseidonGoldilocksConfig};
use plonky2::field::types::Field;
use starky_circuits::circuits::exponentiate::{
    exponentiate,
    get_exponentiate_data
};
use rust_utils::{
    read_file_contents,
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

    // Read and parse input from the specified JSON file
    let input_str = read_file_contents(args.input);
    
    // Get data from config
    let (x, e, _y) = get_exponentiate_data::<PoseidonGoldilocksConfig, 2>(input_str);
    
    let num_rows = 1 << e;
    let public_inputs = [x, F::from_canonical_usize(num_rows), exponentiate(num_rows - 1, x)];
    let stark = S::new(num_rows);

    let _ = stark.generate_trace(public_inputs[0], public_inputs[1], public_inputs[2]);

    // TODO - Serialization of Trace for memory consumption
}
