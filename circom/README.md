# Circom Compiler / snarkjs

[Circom](https://github.com/iden3/circom) is a comiler (Hardware Description Language / HDL) written in Rust.
The compiler outputs the representation of the circuit as R1CS. Successively, one can apply the respective proof system.
Circom's output can be used via a backend (typically using [snarkjs](https://github.com/iden3/snarkjs)).

## Installation

Installation and setup descriptions can be found [here](https://docs.circom.io/getting-started/installation/#installing-dependencies)

In short, you need to clone the Circom repository, run `cargo build --release` and then install circom with `cargo install --path circom`.

As Circom is only the compiler to compile from HDL circuit description to R1CS, you need to additionally install [snarkjs](https://github.com/iden3/snarkjs) to create and verify proofs with e.g. [Groth16](https://github.com/iden3/snarkjs/blob/master/src/groth16_prove.js).

You can install snarkjs with `npm install -g snarkjs` using [npm](https://docs.npmjs.com/downloading-and-installing-node-js-and-npm).

### Using rapidsnark

It is possible to use rapidsnark instead of snarkjs for the proving phase.
To install rapidsnark you should run the following code:

```
git submodule init && git submodule update
cd rapidsnark
```

Then you need to follow the commands described in the README file.

__NOTE:__ You can rapidsnark only in Intel64 machines.

## Compilation & Proof - Toy Examples

To produce and verify a circom circuit you need to perform 6 steps:

1. Compile the circuit
2. Generating the witness for the circuit
3. Setup ceremony for proving and verifying the witness
4. Prove
5. Verify the proof

To automate that process we have wrap everything in a single script (`scripts/run_circuit.sh`).

For example, to prove and verify `circuits/toy/sudoku.circom` using as input `inputs/toy/sudoku.input.json`,
you can run the following command:

```
./scripts/run_circuit.sh circuits/toy/sudoku.circom sudoku inputs/toy/sudoku.input.json phase1/powersOfTau28_hez_final_16.ptau res.csv sudoku_output
```

This script will produce the following files:

* Compile
  - `sudoku_output/sudoku.r1cs` -- the binary format file of the R1CS description of the circuit.
  - `sudoku_output/sudoku_js/sudoku.wasm` -- contains the wasm code describing the circuit.
  - `sudoku_output/sudoku_js/generate_witness.js` and `sudoku_output/sudoku_js/witness_calculator.js` -- files for witness generation.
* Witness Generation
  - `sudoku_output/witness.wtns` -- witness, contains all the computed signals given the input file.
* Setup
  - `sudoku_output/sudoku_0.zkey` -- prover key
  - `sudoku_output/verification_key.json` -- verifier key
* Prove
  - `sudoku_output/public.json` -- contain the public inputs and outputs.
  - `sudoku_output/proof.json` -- contains the proof.
* Verify 
  - Prints the message `[INFO]  snarkJS: OK!` if the verification succeeds.

Finally, `res.csv` will contain statistics about the execution of each step.

__Note__: We currently using the precomputed ceremony from `phase1/powersOfTau28_hez_final_16.ptau`, but in order to safely prove a circuit using Circom you need to safely run a setup ceremony. 
Furthermore, to execute larger circuits you might need a larger powers of tau.

## Running Example with parameterized template variables

```
./scripts/run_circuit.sh circuits/benchmarks/exponentiate/circuit.circom exponentiate \
    inputs/input_1.json phase1/powersOfTau28_hez_final_16.ptau res.csv expo_out E
```

## Adding new circuits

See `TUTORIAL.md`
