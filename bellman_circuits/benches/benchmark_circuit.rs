extern crate rand;
extern crate criterion;

use bellman_circuits::circuits::{sha256};
use rand::{thread_rng};
use bellman::{Circuit, groth16};
use criterion::{Criterion};
use sha2::{Digest, Sha256};
use bellman::gadgets::multipack;
use bls12_381::Bls12;
use bls12_381::Scalar;

fn bench_proof<C: Circuit<Scalar> + Clone + 'static>(
    c: &mut Criterion,
    circuit: C,
    params: groth16::Parameters<Bls12>
) {
    let rng = &mut thread_rng();

    let pvk = groth16::prepare_verifying_key(&params.vk);

    c.bench_function("benching_sha256_setup_time", |b| {
        b.iter(|| { 
            let _ = groth16::generate_random_parameters::<Bls12, _, _>(circuit.clone(), rng).unwrap();
        })
    });

    c.bench_function("benching_sha256_prover_time", |b| {
        b.iter(|| { 
            let _ = groth16::create_random_proof(circuit.clone(), &params, rng); 
        })
    });

    c.bench_function("benching_sha256_verifier_time", |b| {
        let proof = groth16::create_random_proof(circuit.clone(), &params, rng).unwrap(); 
        let hash = Sha256::digest(&Sha256::digest(&[42; 80]));
        let hash_bits = multipack::bytes_to_bits_le(&hash);
        let inputs = multipack::compute_multipacking(&hash_bits); 
        b.iter(|| {
            let _ = groth16::verify_proof(&pvk, &proof, &inputs);        
        })
    });
}

fn main() {
    let mut criterion = Criterion::default().configure_from_args();

    // Pick a preimage and create an instance of our circuit (with the preimage as a witness).
    let hex_value = "68656c6c6f20776f726c64";
    let preimage = hex::decode(hex_value).unwrap();
    let circuit = sha256::Sha256Circuit {
        preimage: Some(preimage.clone()),
        preimage_length: preimage.len(),
    };

    // Generate the parameters for the circuit
    let params = groth16::generate_random_parameters::<Bls12, _, _>(circuit.clone(), &mut thread_rng()).unwrap();

    // Run the benchmark with the given circuit and parameters
    bench_proof(&mut criterion, circuit, params);

    criterion.final_summary();
}

