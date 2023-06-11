#![feature(generic_const_exprs)]

// Extern crate declarations
extern crate rand;
extern crate criterion;

use starky_utils;
use rust_utils::{read_env_variable};
use criterion::{Criterion};
use starky::sha256::{Sha2CompressionStark, Sha2StarkCompressor};
use starky::{
    prover::prove,
    util::to_u32_array_be,
    verifier::verify_stark_proof,
};
use plonky2::plonk::config::{GenericConfig, PoseidonGoldilocksConfig};
use plonky2::util::timing::TimingTree;

pub fn bench_sha256(c: &mut Criterion, num_hashes: i32){
    const D: usize = 2;
    type C = PoseidonGoldilocksConfig;
    type F = <C as GenericConfig<D>>::F;
    type S = Sha2CompressionStark<F, D>;
    let config = starky_utils::secure_config();

    let mut group = c.benchmark_group("sha256");
    
    // 1. Witness generation - generation of the trace
    group.bench_function("setup", |b| {
        b.iter(|| { 
            let mut compressor = Sha2StarkCompressor::new();
            let zero_bytes = [0; 32];
            for _ in 0..num_hashes {
                let left = to_u32_array_be::<8>(zero_bytes);
                let right = to_u32_array_be::<8>(zero_bytes);

                compressor.add_instance(left, right);
            }
            let _ = compressor.generate::<F>();
        })
    });

    // Generate Trace for proving
    // What does this generate for e.g. num_hashes = 512?
    // If you run this loop 512 times, it would produce a total of 512 * 64 = 32kB of pre-Image.
    let mut compressor = Sha2StarkCompressor::new();
    let zero_bytes = [0; 32];
    for _ in 0..num_hashes {
        let left = to_u32_array_be::<8>(zero_bytes);
        let right = to_u32_array_be::<8>(zero_bytes);

        compressor.add_instance(left, right);
    }
    let trace = compressor.generate::<F>();

    // 2. Prove algorithm, uses config
    let stark = S::new();
    group.bench_function("prove", |b| {
        b.iter(|| { 
            let _ = prove::<F, C, S, D>(stark, &config, trace.clone(), [], &mut &mut TimingTree::default()).unwrap();
        })
    });

    // Generate Proof for verification
    let proof = prove::<F, C, S, D>(stark, &config, trace, [], &mut &mut TimingTree::default()).unwrap();

    // 3. Verify algorith, uses config & proof
    group.bench_function("verify", |b| {
        b.iter(|| { 
            verify_stark_proof(stark, proof.clone(), &config).unwrap();
        })
    });
}

fn main() {
    let mut criterion = Criterion::default()
        .configure_from_args()
        .sample_size(10);

    // FIXME - Parese input from file
    // let input_file_str = read_file_from_env_var("INPUT_FILE".to_string());
    // let input_file_str = "Empty".to_string();

    let circuit_str = read_env_variable("CIRCUIT".to_string());

    let num_hashes_string = read_env_variable("NUM_HASHES".to_string());
    let num_hashes: i32 = num_hashes_string.parse().unwrap();

    match circuit_str.as_str() {
        "sha256" => bench_sha256(&mut criterion, num_hashes),
        _ => println!("Unsupported circuit"),
    }

    criterion.final_summary();
}