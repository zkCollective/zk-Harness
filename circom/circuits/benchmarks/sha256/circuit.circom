pragma circom 2.0.0;

include "../../circomlib/sha256/sha256.circom";

component main {public[in]} = Sha256(32);

//./scripts/run_circuit.sh circuits/benchmarks/pedersen/circuit.circom sha256 inputs/benchmarks/pedersen/input.json phase1/powersOfTau28_hez_final_16.ptau res.csv sha256_output

