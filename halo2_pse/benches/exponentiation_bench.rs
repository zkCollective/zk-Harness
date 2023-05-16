use halo2_proofs::{
    halo2curves::bn256::Fr
};
use criterion::{
    Criterion
};
use halo2_pse_circuits::circuits::exponentiation::{ExponentiationCircuit, get_exponentiation_data};

extern crate utilities;
use utilities::{bench_circuit, read_file_from_env_var};


fn bench_exponentiation(c: &mut Criterion, input_file_str: String) {
    let mut group = c.benchmark_group("exponentiation");
    let (k, e_value, x, e, y) = get_exponentiation_data(input_file_str);
    let circuit = ExponentiationCircuit {
        row: e_value,
    };
    let public_input: &[&[Fr]] = &[&[x, e, y]];
    bench_circuit(&mut group, k, circuit, public_input);
    group.finish();
}

fn main() {
    let mut criterion = Criterion::default().configure_from_args().sample_size(10);

    let input_file_str = read_file_from_env_var("INPUT_FILE".to_string());

    bench_exponentiation(&mut criterion, input_file_str);

    criterion.final_summary();
}