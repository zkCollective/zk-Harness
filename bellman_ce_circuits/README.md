# Benchmarking Bellman

This directory contains a benchmarking framework for the [Bellman Community Edition](https://github.com/matter-labs/bellman).

Currently the following circuits have been implemented:

* MIMC (using [Matter Labs implementation](https://github.com/matter-labs/bellman/blob/dev/tests/mimc.rs))

With the following proof system(s):

* Groth16

By default, the number of rounds for MIMC is set to 1000000 in their implementation.

Proof systems that are not yet supported for benchmarks, and are to be integrated in future work:

* Marlin
* Plonk
* Sonic
* gm17

Primitives that should be added for comparison with other libraries:

* Dummy Exponentiation for constraint simulation
* SHA-256

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
RUSTFLAGS=-Awarnings INPUT_FILE=../_input/circuit/mimc/input_1.json CIRCUIT=mimc cargo criterion --message-format=json --bench benchmark_circuit 1> ../benchmarks/bellman_ce/jsons/mimc_input_1_bench.json
```

* Run benchmarks for measuring memory consuption and proof size

```ToDo```

## Adding new circuits

```ToDo```
