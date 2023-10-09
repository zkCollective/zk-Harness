// Extern crate declarations
extern crate rand;
extern crate criterion;

use starky_utils;
use criterion::{Criterion};
use starky::sha256::{Sha2CompressionStark, Sha2StarkCompressor};
use starky::{
    prover::prove,
    util::to_u32_array_be,
    verifier::verify_stark_proof,
};
use plonky2::plonk::config::{GenericConfig, PoseidonGoldilocksConfig};
use plonky2::util::timing::TimingTree;
use plonky2::field::types::Field;

// TODO - Generic bench_stark function.
// Currently not possible as the Stark Trait does not have a generic S::new() function

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
            let _ = prove::<F, C, S, D>(stark, &config, trace.clone(), [], &mut TimingTree::default()).unwrap();
        })
    });

    // Generate Proof for verification
    let proof = prove::<F, C, S, D>(stark, &config, trace, [], &mut TimingTree::default()).unwrap();

    // 3. Verify algorith, uses config & proof
    group.bench_function("verify", |b| {
        b.iter(|| { 
            verify_stark_proof(stark, proof.clone(), &config).unwrap();
        })
    });
}

fn fibonacci<F: Field>(n: usize, x0: F, x1: F) -> F {
    (0..n).fold((x0, x1), |x, _| (x.1, x.0 + x.1)).1
}

pub fn bench_fibonacci(c: &mut Criterion, num_rows_input: i32) -> Result<(), anyhow::Error> {
    const D: usize = 2;
    type C = PoseidonGoldilocksConfig;
    type F = <C as GenericConfig<D>>::F;
    type S = starky_circuits::circuits::fibonacci::FibonacciStark<F, D>;
    let config = starky_utils::secure_config();

    let mut group = c.benchmark_group("fibonacci");

    let num_rows = 1 << num_rows_input;
    let public_inputs = [F::ZERO, F::ONE,fibonacci(num_rows - 1, F::ZERO, F::ONE)];
    let stark = S::new(num_rows);
    
    // 1. Witness Generation
    group.bench_function("setup", |b| {
        b.iter(|| { 
            let _  = stark.generate_trace(public_inputs[0], public_inputs[1]);
        })
    });

    let trace = stark.generate_trace(public_inputs[0], public_inputs[1]);

    // 2. Compute the proof
    group.bench_function("proof", |b| {
        b.iter(|| { 
            let _  = prove::<F, C, S, D>(
                stark,
                &config,
                trace.clone(),
                public_inputs,
                &mut TimingTree::default(),
            );
        })
    });

    let proof = prove::<F, C, S, D>(
        stark,
        &config,
        trace.clone(),
        public_inputs,
        &mut TimingTree::default(),
    )?;

    group.bench_function("verify", |b| {
        b.iter(|| { 
            let _ = verify_stark_proof(stark, proof.clone(), &config);
        })
    });

    verify_stark_proof(stark, proof, &config)?;

    Ok(())
}

fn exponentiate<F: Field>(n: usize, x: F) -> F {
    (0..n).fold((F::ONE, x), |acc, _| (acc.1, acc.1 * x)).1
}

pub fn bench_exponentiate(c: &mut Criterion, input_str: String) -> Result<(), anyhow::Error> {
    const D: usize = 2;
    type C = PoseidonGoldilocksConfig;
    type F = <C as GenericConfig<D>>::F;
    type S = starky_circuits::circuits::exponentiate::ExponentiateStark<F, D>;
    let config = starky_utils::secure_config();

    // Get data from config
    let (x, e, y) = starky_circuits::circuits::exponentiate::get_exponentiate_data::<PoseidonGoldilocksConfig, 2>(input_str);

    let mut group = c.benchmark_group("exponentiate");

    let num_rows = e as usize;
    let public_inputs = [x, F::from_canonical_usize(num_rows), exponentiate(num_rows - 1, x)];
    let stark = S::new(num_rows);
    
    // 1. Witness Generation
    group.bench_function("witness", |b| {
        b.iter(|| { 
            let _  = stark.generate_trace(public_inputs[0], public_inputs[1], public_inputs[2]);
        
        })
    });

    let trace = stark.generate_trace(public_inputs[0], public_inputs[1], public_inputs[2]);

    // 2. Compute the proof
    group.bench_function("prove", |b| {
        b.iter(|| { 
            let _  = prove::<F, C, S, D>(
                stark,
                &config,
                trace.clone(),
                public_inputs,
                &mut TimingTree::default(),
            );
        })
    });

    let proof = prove::<F, C, S, D>(
        stark,
        &config,
        trace.clone(),
        public_inputs,
        &mut TimingTree::default(),
    )?;

    group.bench_function("verify", |b| {
        b.iter(|| { 
            let _ = verify_stark_proof(stark, proof.clone(), &config);
        })
    });

    verify_stark_proof(stark, proof, &config)?;

    Ok(())
}

fn main() {
    let mut criterion = Criterion::default()
        .configure_from_args()
        .sample_size(10);

    // FIXME - Parese input from file
    let input_file_str = starky_utils::read_file_from_env_var("INPUT_FILE".to_string());

    let circuit_str = starky_utils::read_env_variable("CIRCUIT".to_string());

    // TODO - SHA256 & Fibonacci file parsing
    match circuit_str.as_str() {
        // "sha256" => bench_sha256(&mut criterion, num_hashes),
        // "fibonacci" => bench_fibonacci(&mut criterion, num_hashes).unwrap(),
        "exponentiate" => bench_exponentiate(&mut criterion, input_file_str).unwrap(),
        _ => println!("Unsupported circuit"),
    }

    criterion.final_summary();
}