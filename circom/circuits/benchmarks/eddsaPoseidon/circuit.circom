pragma circom 2.0.0;

// Source: https://github.com/iden3/circomlib/blob/master/circuits/eddsaposeidon.circom
include "../circomlib/eddsaposeidon.circom";

component main {public [enabled, Ax, Ay, R8x, R8y, S, M]} = EdDSAPoseidonVerifier();

//./scripts/run_circuit.sh circuits/benchmarks/eddsaPoseidon/circuit.circom eddsaPoseidon inputs/benchmarks/eddsaPoseidon/input.json phase1/powersOfTau28_hez_final_16.ptau res.csv eddsaPoseidon_output
