pragma circom 2.0.0;

include "./keccak.circom";

component main {public [in]} = Keccak(32*8, 32*8);

//./scripts/run_circuit.sh circuits/benchmarks/keccak256/circuit.circom keccak256 inputs/benchmarks/keccak256/input.json phase1/powersOfTau28_hez_final_16.ptau res.csv keccak256_output

