# Tutorial for Adding Starky Circuit in zk-Harness

To add a new circuit please follow the next steps.

## Adding a new circuit

1. Add the circuit source code. You should implement the circuit in a file named `src/circuits/<circuit_name>.rs`. The circuit added should implement the `Stark` trait (example - see `src/circuits/exponentiate.rs`). As such, you should implement the functions `eval_packed_generic` and `eval_ext_circuit`. `eval_ext_circuit` is only relevant for recursion, which is currently not supported for starky in zk-Harness.

2. In the same file you should add a struct called `<CircuitName>Input` that will derive `#[derive(Debug, Deserialize, Serialize)]`
to serialize/desiarilize input files from zk-harness.

3. In the same file you should also implement a function called `get_<circuit_name>_data`
the will get a parameter called `input_file_str` of String type, 
it will read that file, and it will return the deserialized inputs.

4.  Add a test in the same file to test the correctness of the circuit.

5. In `benches/benchmark_circuit.rs` add a function `bench_<circuit_name>` that will get
the file path for the input file. It will read the data using `get_<circuit_name>_data`.
This benchmark will be responsible to measure the performance of the circuit.

6. Adding `src/bin/<circuit_name>_<operation>.rs` that will use `clap` to parse two cli arguments: 
`input` and `proof` (the location of the proof file).

7. You should not forget to update `Cargo.toml` accordingly.

## Add Input File (Optionally)

If you add a new circuit that is not implemented in another framework please include some input files in `../_input/circuit/circuit_name/`.
For more information about the input files please check `../documentation/`.