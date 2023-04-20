use anyhow::Result;
use ark_bw6_761::Fr as Fr761;
use ark_ff::PrimeField;
use clap::Parser;
use jf_plonk::errors::PlonkError;
use jf_relation::{Circuit, PlonkCircuit};
use serde::{Deserialize, Serialize};
use std::{fs::OpenOptions, path::PathBuf, time::Instant};
use sysinfo::{CpuExt, SystemExt};

pub const FRAMEWORK: &str = "jellyfish";
const CSV_FILENAME: &str = "../benchmarks/jellyfish/jellyfish_plonk_cubic.csv";

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

/// A macro that sets up a cubic circuit to benchmark the `Compile` operation.
macro_rules! plonk_compile_bench {
    ($curve:ty, $x:expr, $y:expr, $count:expr) => {{
        for _ in 0..$count - 1 {
            let _: PlonkCircuit<$curve> = cubic_circuit($x, $y).unwrap();
        }

        let cs = cubic_circuit::<$curve>($x, $y).unwrap();
        cs
    }};
}

/// Defines a simple circuit: x**3 + x + 5 == y
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

/// A single row of benchmark results within the CSV log.
#[derive(Debug, Serialize)]
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
    /// Time (in milliseconds) for the operation to finish.
    #[serde(rename(serialize = "time(ms)"))]
    time: String,
    proof_size: String,
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

/// Kind of operation to benchmark.
#[derive(Debug, Serialize, Clone, clap::ValueEnum)]
#[serde(rename_all = "camelCase")]
enum Operation {
    Compile,
    Setup,
    Witness,
    Prove,
    Verify,
}

#[derive(Debug, Serialize)]
#[serde(rename_all = "camelCase")]
enum Category {
    Circuit,
}

fn main() -> Result<()> {
    let args = Args::parse();

    let bench_dir = PathBuf::from(env!("CARGO_MANIFEST_DIR"))
        .parent()
        .unwrap()
        .join("benchmarks/jellyfish/jellyfish_plonk_cubic.csv");

    let mut writer = if !bench_dir.exists() {
        let file = std::fs::File::create(CSV_FILENAME)?;

        csv::WriterBuilder::new().from_writer(file)
    } else {
        csv::WriterBuilder::new()
            .has_headers(false)
            .from_writer(OpenOptions::new().append(true).open(CSV_FILENAME)?)
    };

    println!(
        "Benching:\n framework: {}\n backend: {}\n circuit: {}\n curve: {}\n count: {}\n input: {}\n",
        FRAMEWORK, args.backend, args.circuit, args.curve, args.count, args.input
    );

    let input = std::fs::read_to_string(format!("../{}", args.input)).unwrap();
    let input: CubicInput = serde_json::de::from_str(&input).unwrap();

    let x = input.X.parse().ok().unwrap();
    let y = input.Y.parse().ok().unwrap();

    // TODO: we allow a single match for now since we're only benchmarking the circuit compilation.
    // This should be removed once we add more operations.
    #[allow(clippy::single_match)]
    match args.op {
        Operation::Compile => {
            let cs = plonk_compile_bench!(Fr761, x, y, args.count);
            let start = Instant::now();
            let end = start.elapsed().as_secs_f32() * 1000.0;
            let time = if end < 1.0 {
                "1".to_string()
            } else {
                end.to_string()
            };

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
                ram: Default::default(), // TODO: add memory usage
                time,
                proof_size: 0.to_string(),
                nb_physical_cores: num_cpus::get_physical(),
                nb_logical_cores: num_cpus::get(),
                count: args.count,
                cpu: system.global_cpu_info().brand().to_string(),
            };
            writer.serialize(record)?;
        }
        _ => {}
    }

    writer.flush()?;
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

    /// Test driver using different curves for our cubic_circuit() circuit construction.
    #[test]
    fn test_cubic() -> Result<(), CircuitError> {
        test_cubic_helper::<FqEd254>()?;
        test_cubic_helper::<FqEd377>()?;
        test_cubic_helper::<FqEd381>()?;
        test_cubic_helper::<Fq377>()
    }

    /// Sanity check for our cubic_circuit() circuit construction.
    fn test_cubic_helper<F: PrimeField>() -> Result<(), CircuitError> {
        let circuit: PlonkCircuit<F> = cubic_circuit(2u32, 15u32).unwrap();
        // 2 mul gates, 2 additional, 2 constant gates
        assert_eq!(circuit.num_gates(), 7);
        // Check the number of public inputs:
        assert_eq!(circuit.num_inputs(), 1);

        // Check circuit satisfiability
        let pub_input = &[F::from(15u32)];
        let verify = circuit.check_circuit_satisfiability(pub_input);
        assert!(verify.is_ok(), "{:?}", verify.unwrap_err());

        Ok(())
    }
}
