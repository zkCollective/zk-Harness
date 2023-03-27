pragma circom 2.0.0;

include "../../circomlib/pedersen.circom";

component main {public[in]} = Pedersen(32);

//./scripts/run_circuit.sh circuits/benchmarks/pedersen/circuit.circom pedersen inputs/benchmarks/pedersen/input.json phase1/powersOfTau28_hez_final_16.ptau res.csv pedersen_output
