// Extern crate declarations
extern crate rand;
extern crate criterion;

// Use statements
use bellman_circuits::circuits::{sha256, exponentiate};
use rand::thread_rng;
use bellman::{Circuit, groth16};
use criterion::{Criterion, BenchmarkGroup};
use criterion::measurement::Measurement;
use sha2::{Digest, Sha256};
use bellman::gadgets::multipack;
use bls12_381::{Bls12, Scalar};
use rust_utils::{read_file_from_env_var, read_env_variable};
use ff::PrimeField;
use bellman::gadgets::test::TestConstraintSystem;

// Benchmark for a given circuit
pub fn bench_circuit<M: Measurement, C: Circuit<Scalar> + Clone + 'static>(
    c: &mut BenchmarkGroup<'_, M>,
    circuit: C,
    public_inputs: Vec<Scalar>,
    params: groth16::Parameters<Bls12>
) {
    let rng = &mut thread_rng();
    let pvk = groth16::prepare_verifying_key(&params.vk);

    c.bench_function("setup", |b| {
        b.iter(|| { 
            let _ = groth16::generate_random_parameters::<Bls12, _, _>(circuit.clone(), rng).unwrap();
        })
    });

    c.bench_function("prove", |b| {
        b.iter(|| { 
            let _ = groth16::create_random_proof(circuit.clone(), &params, rng); 
        })
    });

    let proof = groth16::create_random_proof(circuit.clone(), &params, rng).unwrap(); 

    c.bench_function("verify", |b| {
        b.iter(|| {
            let _ = groth16::verify_proof(&pvk, &proof, &public_inputs);        
        })
    });
}

// Benchmark for SHA-256
fn bench_sha256(c: &mut Criterion, input_str: String) {
    let mut group = c.benchmark_group("sha256");
    let (preimage_length, preimage) = sha256::get_sha256_data(input_str);

    // Pre-Compute public inputs
    let hash = Sha256::digest(&Sha256::digest(&preimage));
    let hash_bits = multipack::bytes_to_bits_le(&hash);
    let inputs = multipack::compute_multipacking(&hash_bits);

    // Define the circuit
    let circuit = sha256::Sha256Circuit {
        preimage: Some(preimage.clone()),
        preimage_length: preimage_length,
    };

    // Generate Parameters
    let params = groth16::generate_random_parameters::<Bls12, _, _>(circuit.clone(), &mut thread_rng()).unwrap();
    
    bench_circuit(&mut group, circuit, inputs, params);
}

// Benchmark for Exponentiation
fn bench_exponentiate(c: &mut Criterion, input_str: String) {
    let mut group = c.benchmark_group("exponentiate");

    // Get data from config
    let (x, e, y) = exponentiate::get_exponentiate_data(input_str);
    
    // Public inputs are x and y
    let x_bits = multipack::bytes_to_bits_le(&x.to_repr().as_ref());
    let y_bits = multipack::bytes_to_bits_le(&y.to_repr().as_ref());
    let inputs = [multipack::compute_multipacking(&x_bits), multipack::compute_multipacking(&y_bits)].concat();

    // Define the circuit
    let circuit = exponentiate::ExponentiationCircuit {
        x: Some(x),
        e: e,
        y: Some(y),
    };

    // Generate Parameters
    let params = groth16::generate_random_parameters::<Bls12, _, _>(circuit.clone(), &mut thread_rng()).unwrap();

    // Create a mock constraint system
    let mut cs = TestConstraintSystem::<Scalar>::new();
    // Synthesize the circuit with our mock constraint system
    circuit.clone().synthesize(&mut cs).unwrap();
    println!("Number of constraints: {}", cs.num_constraints());

    bench_circuit(&mut group, circuit, inputs, params);
}

fn main() {
    let mut criterion = Criterion::default()
        .configure_from_args()
        .sample_size(10);

    let input_file_str = read_file_from_env_var("INPUT_FILE".to_string());

    let circuit_str = read_env_variable("CIRCUIT".to_string());

    match circuit_str.as_str() {
        "sha256" => bench_sha256(&mut criterion, input_file_str),
        "exponentiate" => bench_exponentiate(&mut criterion, input_file_str),
        _ => println!("Unsupported circuit"),
    }
    criterion.final_summary();
}
