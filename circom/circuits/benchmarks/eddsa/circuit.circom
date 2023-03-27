pragma circom 2.0.0;

include "../../circomlib/eddsa.circom";

component main {public [A, R8, S, msg]} = EdDSAVerifier(80);

//./scripts/run_circuit.sh circuits/benchmarks/eddsa/circuit.circom eddsa inputs/benchmarks/eddsa/input.json phase1/powersOfTau28_hez_final_16.ptau res.csv eddsa_output

