# Benchmarking Bellman

This directory contains a benchmarking framework for Bellman.

Currently the following circuits have been implemented:

* Exponentiation (using a custom implementation)
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
RUSTFLAGS=-Awarnings INPUT_FILE=../_input/circuit/exponentiate/input_10.json CIRCUIT=exponentiate cargo criterion --message-format=json --bench benchmark_circuit 1> ../benchmarks/bellman/jsons/exponentiate_input_10_bench.json
```

```
RUSTFLAGS=-Awarnings INPUT_FILE=../_input/circuit/sha256/input_1.json CIRCUIT=sha256 cargo criterion --message-format=json --bench benchmark_circuit 1> ../benchmarks/bellman/jsons/sha256_input_1_bench.json
```

* Run benchmarks for measuring memory consuption and proof size

```
cargo run --bin sha256 --release -- \
    --input ../_input/circuit/sha256/input_1.json \
    --output ../benchmarks/bellman/jsons/sha256_input_1.json
```

## Adding new circuits

ToDo
