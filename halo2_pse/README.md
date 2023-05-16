# Benchmarking halo2-pse Library

This directory benchmarks the implementation of halo2's fork form PSE that
can be found [here](https://github.com/privacy-scaling-explorations/halo2/).
The main difference of PSE's halo2 fork is that it uses KZG instead of IPA.
You can find more details about halo2 in the [halo2 book](https://zcash.github.io/halo2/index.html).

Currently the following circuits have been implemented:

* Exponentiation (using a custom implementation)
* Sha256 (using the sha256 gadget from [here](https://github.com/privacy-scaling-explorations/halo2/tree/main/halo2_gadgets/src/sha256)

Note: currently we use Celar network 
[fork](https://github.com/celer-network/halo2/tree/KZG-bench-sha256)
because there is no implementation of sha256 to work with KZG in 
[PSE's fork](https://github.com/privacy-scaling-explorations/halo2/issues/182).
We should update the dependencies when there would be support in PSE's fork.

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
INPUT_FILE=../_input/circuit/exponentiate/input_1.json cargo criterion --message-format=json \
    --bench exponentiation_bench 1> ../benchmarks/halo2_pse/jsons/exponentiate_input_1_bench.json
```

* Run benchmarks for measuring memory consuption and proof size

```
cargo run --bin exponentiation --release -- \
    --input ../_input/circuit/exponentiate/input_1.json \
    --output ../benchmarks/halo2_pse/jsons/exponentiate_input_1.json
```

## Adding new circuits

See `TUTORIAL.md`
