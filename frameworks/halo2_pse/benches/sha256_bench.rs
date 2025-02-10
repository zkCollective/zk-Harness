use halo2_proofs::{
    halo2curves::bn256::Fr
};
use criterion::{
    Criterion
};
use halo2_pse_circuits::circuits::sha256::{Sha256Circuit, get_sha256_data};

use utilities::{bench_circuit, read_file_from_env_var};


fn bench_sha256(c: &mut Criterion, input_file_str: String) {
    let mut group = c.benchmark_group("sha256");
    let (k, sha_data) = get_sha256_data(input_file_str);
    let circuit = Sha256Circuit {
        sha_data: sha_data, 
    };
    let public_input: &[&[Fr]] = &[];
    bench_circuit(&mut group, k, circuit, public_input);
}


fn main() {
    let mut criterion = Criterion::default().configure_from_args().sample_size(10);

    let input_file_str = read_file_from_env_var("INPUT_FILE".to_string());

    bench_sha256(&mut criterion, input_file_str);

    criterion.final_summary();
}