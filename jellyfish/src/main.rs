use anyhow::Result;
use ark_bw6_761::Fr as Fr761;
use ark_ff::PrimeField;
use clap::Parser;
use jf_plonk::errors::PlonkError;
use jf_relation::{Circuit, PlonkCircuit};
use serde::{Deserialize, Serialize};
use std::{path::PathBuf, time::Instant};
use sysinfo::{CpuExt, SystemExt};

pub const FRAMEWORK: &str = "jellyfish";

#[derive(Parser, Debug)]
#[command(author, version, about, long_about = None)]
struct Args {
    /// Backend to use.
    ///
    /// This input should usually be provided by process_circuit.py, under _scripts.
    #[arg(short, long)]
    backend: String,
    /// Type of circuit to benchmark against.
    ///
    /// This input should usually be provided by process_circuit.py, under _scripts.
    #[arg(long)]
    circuit: String,
    /// Curve to use.
    ///
    /// This input should usually be provided by process_circuit.py, under _scripts.
    #[arg(long)]
    curve: String,
    /// Path(s) to the JSON inputs used.
    ///
    /// This input should usually be provided by process_circuit.py, under _scripts.
    #[arg(long)]
    input: String,
    /// Kind of operation to run.
    ///
    /// This input should usually be provided by process_circuit.py, under _scripts.
    #[arg(long)]
    op: Operation,
    /// Number of times to run.
    ///
    /// This input should usually be provided by process_circuit.py, under _scripts.
    #[arg(long)]
    count: usize,
}

fn compile_bench<F: PrimeField>(
    x: u32,
    y: u32,
    count: usize,
) -> Result<PlonkCircuit<F>, PlonkError> {
    for _ in 0..count - 1 {
        let _: PlonkCircuit<F> = cubic_circuit(x, y).unwrap();
    }

    cubic_circuit::<F>(x, y)
}

fn cubic_circuit<F: PrimeField>(x: u32, y: u32) -> Result<PlonkCircuit<F>, PlonkError> {
    let mut circuit: PlonkCircuit<F> = PlonkCircuit::new_turbo_plonk();

    // Setup required variables
    let a = circuit.create_variable(F::from(x)).unwrap();
    let _y = circuit.create_public_variable(F::from(y)).unwrap();
    let five = circuit.create_variable(F::from(5u32)).unwrap();
    let b = circuit.create_variable(F::from(x)).unwrap();

    // Starting here, we build: x^3 + x + 5 == y
    // x = x * x * x
    let x_to_two = circuit.mul(a, b).unwrap();
    let x_to_three = circuit.mul(a, x_to_two).unwrap();

    // e = x^3 + x
    let e = circuit.add(a, x_to_three).unwrap();

    // f = x^3 + x + 5
    let _f = circuit.add(e, five).unwrap();

    Ok(circuit)
}

/// A row of benchmark results within the CSV log.
#[derive(Serialize)]
#[serde(rename_all = "camelCase")]
struct Record {
    framework: String,
    category: Category,
    backend: String,
    curve: String,
    circuit: String,
    input: PathBuf,
    operation: Operation,
    nb_constraints: String,
    nb_secret: String,
    nb_public: String,
    ram: usize,
    proof_size: usize,
    /// Time (in milliseconds) for the operation to finish.
    #[serde(rename(serialize = "time(ms)"))]
    time: String,
    nb_physical_cores: usize,
    nb_logical_cores: usize,
    count: usize,
    cpu: String,
}

/// Input for cubic curve.
#[allow(non_snake_case)]
#[derive(Debug, Deserialize)]
struct CubicInput {
    X: String,
    Y: String,
}

#[derive(Debug, Serialize, Clone, clap::ValueEnum)]
#[serde(rename_all = "camelCase")]
enum Operation {
    Compile,
    Setup,
    Witness,
    Prove,
    Verify,
}

#[derive(Serialize)]
#[serde(rename_all = "camelCase")]
enum Category {
    Circuit,
}

fn main() -> Result<()> {
    let args = Args::parse();

    println!(
        "Running:\n backend: {}\n circuit: {}\n curve: {}\n count: {} input: {}\n",
        args.backend, args.circuit, args.curve, args.count, args.input
    );

    let mut writer =
        csv::Writer::from_path("../benchmarks/jellyfish/jellyfish_plonk_cubic.csv").unwrap();

    let file = format!("../{}", args.input);
    let input = std::fs::read_to_string(file).unwrap();
    let input: CubicInput = serde_json::de::from_str(&input).unwrap();

    let x = input.X.parse().ok().unwrap();
    let y = input.Y.parse().ok().unwrap();

    match args.op {
        Operation::Compile => {
            let start = Instant::now();
            let cs: PlonkCircuit<Fr761> = compile_bench(x, y, args.count).unwrap();
            let end = start.elapsed().as_secs_f32() * 1000.0;
            println!("end: {}", end);

            let system = sysinfo::System::new_all();
            let record = Record {
                framework: FRAMEWORK.to_string(),
                category: Category::Circuit,
                backend: args.backend.to_string(),
                curve: args.curve.to_string(),
                circuit: args.circuit.to_string(),
                input: args.input.into(),
                operation: args.op,
                nb_constraints: cs.num_gates().to_string(),
                nb_secret: cs.num_vars().to_string(),
                nb_public: cs.num_inputs().to_string(),
                ram: Default::default(),        // TODO: add memory usage
                proof_size: Default::default(), // TODO: add proof_size
                time: end.to_string(),
                nb_physical_cores: num_cpus::get_physical(),
                nb_logical_cores: num_cpus::get(),
                count: args.count,
                cpu: system.global_cpu_info().brand().to_string(),
            };

            writer.serialize(record).unwrap();
            writer.flush().unwrap();
        }
        _ => {}
    }
    Ok(())
}

#[cfg(test)]
mod tests {
    use ark_bls12_377::Fq as Fq377;
    use ark_ed_on_bls12_377::Fq as FqEd377;
    use ark_ed_on_bls12_381::Fq as FqEd381;
    use ark_ed_on_bn254::Fq as FqEd254;
    use jf_relation::errors::CircuitError;

    use super::*;

    #[test]
    fn test_cubic() -> Result<(), CircuitError> {
        test_cubic_helper::<FqEd254>()?;
        test_cubic_helper::<FqEd377>()?;
        test_cubic_helper::<FqEd381>()?;
        test_cubic_helper::<Fq377>()
    }

    fn test_cubic_helper<F: PrimeField>() -> Result<(), CircuitError> {
        let circuit: PlonkCircuit<F> = cubic_circuit(1u32, 1u32).unwrap();
        // 2 mul gates, 2 additional, 2 constant gates
        assert_eq!(circuit.num_gates(), 7);
        // Check the number of public inputs:
        assert_eq!(circuit.num_inputs(), 1);
        // circuit.enforce_equal(result, f).unwrap();

        Ok(())
    }
}
