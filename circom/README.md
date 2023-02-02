# Circom Compiler

[Circom](https://github.com/iden3/circom) is a comiler (Hardware Description Language / HDL) written in Rust.
The compiler outputs the representation of the circuit as R1CS. Successively, one can apply the respective proof system.

## Plain Setup

### Installation

Installation and setup descriptions can be found [here](https://docs.circom.io/getting-started/installation/#installing-dependencies)

In short, you need to clone the Circom repository, run ``` cargo build --release ``` and then install circom with ``` cargo install --path circom ```.

As Circom is only the compiler to compile from HDL circuit description to R1CS, you need to additionally install [snarkjs](https://github.com/iden3/snarkjs) to create and verify proofs with e.g. [Groth16](https://github.com/iden3/snarkjs/blob/master/src/groth16_prove.js).

You can install snarkjs with ``` npm install -g snarkjs ``` using [npm](https://docs.npmjs.com/downloading-and-installing-node-js-and-npm).

### Compilation

To compile a circuit run e.g. ``` circom circuits/toy/multiply.circom ``` in this folder.

The Makefile creates additional files in the folder ``` multiply_js ```, such as:

    - multiply.r1cs - The binary format file of the R1CS description of the circuit
    - multiply.wasm - Contains the wasm code 
    - witness_calculator.js and generate_witness.js - files for witness generation

### Proving the circuit


## Docker Setup