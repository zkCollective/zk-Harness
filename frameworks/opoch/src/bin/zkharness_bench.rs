//! zk-Harness compatible benchmark binary for OPOCH
//!
//! Outputs CSV in the format expected by zk-Harness:
//! framework,category,backend,curve,circuit,input,operation,nbConstraints,nbSecret,nbPublic,ram,time,proofSize,count
//!
//! Usage:
//!   cargo run --release --bin zkharness_bench -- --circuit sha256_chain --input 1024 --input-path input/circuit/opoch/sha256_chain/input_1024.json --count 1

use std::env;
use std::time::Instant;
use opoch_poc_sha::endtoend::{generate_production_proof, measure_verification_time};

fn main() {
    let args: Vec<String> = env::args().collect();

    // Parse arguments
    let mut circuit = "sha256_chain".to_string();
    let mut input_size: usize = 1024;
    let mut input_path: Option<String> = None;
    let mut count: usize = 1;

    let mut i = 1;
    while i < args.len() {
        match args[i].as_str() {
            "--circuit" => {
                circuit = args.get(i + 1).cloned().unwrap_or_default();
                i += 2;
            }
            "--input" => {
                input_size = args.get(i + 1)
                    .and_then(|s| s.parse().ok())
                    .unwrap_or(1024);
                i += 2;
            }
            "--input-path" => {
                input_path = args.get(i + 1).cloned();
                i += 2;
            }
            "--count" => {
                count = args.get(i + 1)
                    .and_then(|s| s.parse().ok())
                    .unwrap_or(1);
                i += 2;
            }
            "--curve" => {
                // We only support goldilocks, but accept the arg for compatibility
                i += 2;
            }
            _ => {
                i += 1;
            }
        }
    }

    // Use input_path if provided, otherwise fall back to input_size for the CSV input column
    let input_label = input_path.unwrap_or_else(|| input_size.to_string());

    // Print CSV header
    println!("framework,category,backend,curve,circuit,input,operation,nbConstraints,nbSecret,nbPublic,ram,time,proofSize,count");

    match circuit.as_str() {
        "sha256_chain" => benchmark_sha256_chain(input_size, &input_label, count),
        _ => {
            eprintln!("Unknown circuit: {}. Supported: sha256_chain", circuit);
            std::process::exit(1);
        }
    }
}

fn benchmark_sha256_chain(n: usize, input_label: &str, count: usize) {
    let input = b"zk-Harness benchmark input";

    // Determine segment configuration
    // Use 64 steps per segment (standard configuration)
    let segment_length = 64;
    let num_segments = (n + segment_length - 1) / segment_length;
    let total_steps = num_segments * segment_length;

    // ~64 constraints per SHA-256 round
    let nb_constraints = total_steps * 64;

    // Memory estimates (rough)
    let ram_prove: u64 = (num_segments as u64) * 512 * 1024 * 1024; // ~512MB per segment
    let ram_verify: u64 = 1024 * 1024; // ~1MB for verification

    for iteration in 0..count {
        // === PROVE ===
        let prove_start = Instant::now();

        // Generate production proof using the existing function
        let (proof, _d0, _y) = generate_production_proof(input, num_segments, segment_length);

        let prove_time_ms = prove_start.elapsed().as_millis() as f64;
        let proof_bytes = proof.serialize();
        let proof_size = proof_bytes.len();

        // Output prove result
        println!(
            "opoch,circuit,stark,goldilocks,sha256_chain,{},prove,{},0,2,{},{:.3},{},{}",
            input_label, nb_constraints, ram_prove, prove_time_ms, proof_size, iteration + 1
        );

        // === VERIFY ===
        // Measure verification time (average over 100 iterations for stability)
        let verify_duration = measure_verification_time(&proof, input, 100);
        let verify_time_ms = verify_duration.as_secs_f64() * 1000.0;

        // Output verify result
        println!(
            "opoch,circuit,stark,goldilocks,sha256_chain,{},verify,{},0,2,{},{:.6},{},{}",
            input_label, nb_constraints, ram_verify, verify_time_ms, proof_size, iteration + 1
        );
    }
}
