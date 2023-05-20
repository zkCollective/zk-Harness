extern crate rand;
extern crate bellman_ce;
extern crate criterion;

use bellman_ce_circuits::circuits::{mimc};
use std::rc::Rc;
use rand::{thread_rng, Rng};
use bellman_ce::{Circuit};
use bellman_ce::pairing::bn256::{Bn256};
use bellman_ce::groth16::{Parameters, generate_random_parameters, prepare_verifying_key, create_random_proof, verify_proof};
use criterion::{Criterion, BenchmarkGroup};
use criterion::measurement::Measurement;
use rust_utils::{read_file_from_env_var, read_env_variable};

// Modify the bench_mimc_proof function to accept the circuit as an argument
fn bench_circuit<M: Measurement, C: Circuit<Bn256> + Clone + 'static>(
    c: &mut BenchmarkGroup<'_, M>,
    circuit: C,
    params: Parameters<Bn256>
) {
    let rng = &mut thread_rng();

    let pvk = prepare_verifying_key(&params.vk);

    c.bench_function("setup", |b| {
        b.iter(|| { 
            let _ = generate_random_parameters(circuit.clone(), rng);
        })
    });

    c.bench_function("prove", |b| {
        b.iter(|| {
            let c = circuit.clone();     
            let _ = create_random_proof(c, &params, rng); 
        })
    });

    c.bench_function("verify", |b| {
        let c = circuit.clone();
        let proof = create_random_proof(c, &params, rng).unwrap();  
        b.iter(|| {
            let _ = verify_proof(&pvk, &proof, &[]);        
        })
    });
}

fn bench_mimc(c: &mut Criterion, input_str: String) {
    let mut group = c.benchmark_group("mimc");

    // TODO - Add parsing from input according to MIMC spec, currently input is fixed / random
    // Generate the MiMC round constants
    let rng = &mut thread_rng();
    let constants = mimc::generate_constants(rng);

    let xl = rng.gen();
    let xr = rng.gen();  

    // Create the circuit
    let circuit = mimc::MiMCDemo::<Bn256> {
        xl: Some(xl),
        xr: Some(xr),
        constants: Rc::clone(&constants),
    };

    // Generate the parameters for the circuit
    let params = mimc::generate_circuit_parameters(circuit.clone(), rng).unwrap();

    // Run the benchmark with the given circuit and parameters
    bench_circuit(&mut group, circuit, params);
}


fn main() {
    let mut criterion = Criterion::default()
        .configure_from_args()
        .sample_size(10);

    let input_file_str = read_file_from_env_var("INPUT_FILE".to_string());

    let circuit_str = read_env_variable("CIRCUIT".to_string());

    match circuit_str.as_str() {
        "mimc" => bench_mimc(&mut criterion, input_file_str),
        _ => println!("Unsupported circuit"),
    }
    criterion.final_summary();
}
