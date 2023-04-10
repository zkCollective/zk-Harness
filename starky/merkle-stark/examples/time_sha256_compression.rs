use std::env;

use log::{debug, Level, LevelFilter};
// use merkle_stark::serialization::Buffer;
use merkle_stark::{
    config::StarkConfig,
    prover::prove,
    serialization::Buffer,
    sha256_stark::{Sha2CompressionStark, Sha2StarkCompressor},
    stark::Stark,
    util::to_u32_array_be,
    verifier::verify_stark_proof,
};
use plonky2::plonk::config::{GenericConfig, PoseidonGoldilocksConfig};
use plonky2::util::timing::TimingTree;

const D: usize = 2;
type C = PoseidonGoldilocksConfig;
type F = <C as GenericConfig<D>>::F;
type S = Sha2CompressionStark<F, D>;

fn main() {
    let args: Vec<String> = env::args().collect();
    let num_hashes = args[1].parse::<i32>().unwrap();
    println!(
        "\n============== num hashes {} =======================================",
        num_hashes
    );

    let mut builder = env_logger::Builder::from_default_env();
    builder.format_timestamp(None);
    builder.filter_level(LevelFilter::Debug);
    builder.try_init().unwrap();

    let mut compressor = Sha2StarkCompressor::new();
    let zero_bytes = [0; 32];
    for _ in 0..num_hashes {
        let left = to_u32_array_be::<8>(zero_bytes);
        let right = to_u32_array_be::<8>(zero_bytes);

        compressor.add_instance(left, right);
    }

    let mut timing = TimingTree::new("stark", Level::Debug);
    timing.push("prove", Level::Debug);
    timing.push("gen trace", Level::Debug);
    let trace = compressor.generate();
    timing.pop();
    println!("trace len {} width {}", trace[0].len(), trace.len());

    let config = StarkConfig::standard_fast_config();

    debug!("Num Columns: {}", S::COLUMNS);
    let stark = S::new();
    let proof = prove::<F, C, S, D>(stark, &config, trace, [], &mut timing).unwrap();
    let mut buffer = Buffer::new(Vec::new());
    let _ = buffer.write_stark_proof_with_public_inputs(&proof);
    println!("proof size {}\n", buffer.bytes().len());
    timing.pop();

    timing.push("verify", Level::Debug);
    verify_stark_proof(stark, proof, &config).unwrap();
    timing.pop();
    timing.print();
}
