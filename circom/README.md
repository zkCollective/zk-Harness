# Circom Compiler / snarkjs

[Circom](https://github.com/iden3/circom) is a comiler (Hardware Description Language / HDL) written in Rust.
The compiler outputs the representation of the circuit as R1CS. Successively, one can apply the respective proof system.
Circom's output can be used via a backend (typically using [snarkjs](https://github.com/iden3/snarkjs)).

## Installation

Installation and setup descriptions can be found [here](https://docs.circom.io/getting-started/installation/#installing-dependencies)

In short, you need to clone the Circom repository, run `cargo build --release` and then install circom with `cargo install --path circom`.

As Circom is only the compiler to compile from HDL circuit description to R1CS, you need to additionally install [snarkjs](https://github.com/iden3/snarkjs) to create and verify proofs with e.g. [Groth16](https://github.com/iden3/snarkjs/blob/master/src/groth16_prove.js).

You can install snarkjs with `npm install -g snarkjs` using [npm](https://docs.npmjs.com/downloading-and-installing-node-js-and-npm).

#### Arithmetics and elliptic curves benchmarks

If you want to run the benchmarks for Arithmetics and EC then you need to 
execute the following command inside `circom` directory to download the
required library, ffjavascript. This will also install the libraries needed to run the test scripts.

```
yarn install
```

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

A list of available scripts is at the bottom: [Compilation & Proof Scripts](#compilation--proof-scripts)

## Adding new circuits

See `TUTORIAL.md`

### Compilation & Proof Scripts
- Cubic
  ```
  ./scripts/run_circuit.sh circuits/benchmarks/cubic/circuit.circom cubic inputs/benchmarks/cubic/input.json phase1/powersOfTau28_hez_final_16.ptau res.csv cubic_output
  ```
- Poseidon 
  ```
  ./scripts/run_circuit.sh circuits/benchmarks/poseidon/circuit.circom poseidon inputs/benchmarks/poseidon/input.json phase1/powersOfTau28_hez_final_16.ptau res.csv poseidon_output
  ```
- Pedersen 
  ```
  ./scripts/run_circuit.sh circuits/benchmarks/pedersen/circuit.circom pedersen inputs/benchmarks/pedersen/input.json phase1/powersOfTau28_hez_final_16.ptau res.csv pedersen_output
  ```
- MiMC 
  ```
  ./scripts/run_circuit.sh circuits/benchmarks/mimc/circuit.circom mimc inputs/benchmarks/mimc/input.json phase1/powersOfTau28_hez_final_16.ptau res.csv mimc_output
  ```
- SHA-256 
  ```
  ./scripts/run_circuit.sh circuits/benchmarks/pedersen/circuit.circom sha256 inputs/benchmarks/pedersen/input.json phase1/powersOfTau28_hez_final_16.ptau res.csv sha256_output
  ```
- ECDSA (verifier) (this required 56G of RAM)
  ```
  ./scripts/run_circuit.sh circuits/benchmarks/ecdsa/circuit.circom ecdsa inputs/benchmarks/ecdsa/input.json phase1/powersOfTau28_hez_final_21.ptau res.csv ecdsa_output
  ```
- EdDSA Poseidon
  ```
  ./scripts/run_circuit.sh circuits/benchmarks/eddsaPoseidon/circuit.circom eddsaPoseidon inputs/benchmarks/eddsaPoseidon/input.json phase1/powersOfTau28_hez_final_16.ptau res.csv eddsaPoseidon_output
  ```
- SMT
  ```
  ./scripts/run_circuit.sh circuits/benchmarks/smt/circuit.circom smt inputs/benchmarks/smt/input.json phase1/powersOfTau28_hez_final_16.ptau res.csv smt_output
  ```
- Keccak-256
  ```
  ./scripts/run_circuit.sh circuits/benchmarks/keccak256/circuit.circom keccak256 inputs/benchmarks/keccak256/input.json phase1/powersOfTau28_hez_final_20.ptau res.csv keccak256_output
  ```
- Schnorr
  ```
  ./scripts/run_circuit.sh circuits/benchmarks/schnorr/circuit.circom schnorr inputs/benchmarks/schnorr/input.json phase1/powersOfTau28_hez_final_16.ptau res.csv schnorr_output
  ```



