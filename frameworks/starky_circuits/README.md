# Benchmarking Bellman

This directory contains a benchmarking framework for starky.

Currently the following circuits have been implemented:

* Exponentiation (using a custom implementation)
* Fibonacci Sequence (using the implementation [here](https://github.com/tumberger/plonky2/tree/sha256-starky))
* Sha256 (using the sha256 gadget from [here](https://docs.rs/bellman/latest/bellman/)

## Plain Setup

### Installation

To run these benchmarks you need to install `rust`, `cargo`, and `cargo-criterion`.

```bash
# install rust and cargo: https://www.rust-lang.org/tools/install
curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh
# install cargo-criterion
cargo install cargo-criterion
```

### Run the benchmarks

Running a benchmark can be facilitated through the following commands:

* Run benchmarks for measuring performance
```
RUSTFLAGS=-Awarnings INPUT_FILE=../_input/circuit/exponentiate/input_10.json CIRCUIT=exponentiate cargo criterion --message-format=json --bench benchmark_circuit 1> ../benchmarks/starky/jsons/exponentiate_input_10_bench.json
```

* Run benchmarks for measuring memory consuption and proof size

```
cargo run --bin exponentiate_prove --release -- \
    --input ../_input/circuit/exponentiate/input_10.json \
    --output ../benchmarks/starky/jsons/exponentiate_input_10.json
```

## OpenTO

## Adding new circuits

ToDo
