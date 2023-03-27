pragma circom 2.0.0;

include "../../circomlib/poseidon.circom";

component main {public[inputs]} = Poseidon(2);

//./scripts/run_circuit.sh circuits/benchmarks/poseidon/circuit.circom poseidon inputs/benchmarks/poseidon/input.json phase1/powersOfTau28_hez_final_16.ptau res.csv poseidon_output
