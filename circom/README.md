# Circom Compiler

[Circom](https://github.com/iden3/circom) is a comiler (Hardware Description Language / HDL) written in Rust.
The compiler outputs the representation of the circuit as R1CS. Successively, one can apply the respective proof system.

## Plain Setup

### Installation

Installation and setup descriptions can be found [here](https://docs.circom.io/getting-started/installation/#installing-dependencies)

In short, you need to clone the Circom repository, run ``` cargo build --release ``` and then install circom with ``` cargo install --path circom ```.

As Circom is only the compiler to compile from HDL circuit description to R1CS, you need to additionally install [snarkjs](https://github.com/iden3/snarkjs) to create and verify proofs with e.g. [Groth16](https://github.com/iden3/snarkjs/blob/master/src/groth16_prove.js).

You can install snarkjs with ``` npm install -g snarkjs ``` using [npm](https://docs.npmjs.com/downloading-and-installing-node-js-and-npm).

### Compilation & Proof - Toy Examples

To compile a circuit run e.g. ``` circom circuits/toy/multiply.circom ``` in this folder.

The Makefile creates additional files in the folder ``` tmp ```, with the following phases:

- ``` make ```
  - multiply.r1cs - The binary format file of the R1CS description of the circuit
  - multiply_js/multiply.wasm - Contains the wasm code 
  - multiply_js/witness_calculator.js and generate_witness.js - files for witness generation

- ``` make setup ```
  - ``` sudok.wtns ``` - witness, contains all the computed signals given the input file (json)
  - ``` sudok.ptau ``` - File for the ["Powers of Tau"](https://eprint.iacr.org/2022/1592.pdf), in this Makefile we simply run it without any contribution
  - ``` sudoku.pk ``` and ``` sudoku.vk ``` - prover key and verifier key as output of the trusted setup
  
- ``` make prove ```
  - ``` sudoku.inst.json ``` - Contains the public inputs and outputs
  - ``` sudoku.pf.json ``` - The proof of the to be proven relation

- ``` make verify ```
  - Outputs ``` OK ``` if the proof verifies.

- ``` make clean ```
  - Cleans up the ``` tmp ``` folder

## Docker Setup