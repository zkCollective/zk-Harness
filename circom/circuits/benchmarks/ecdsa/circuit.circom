pragma circom 2.0.0;
// Source: https://github.com/0xPARC/circom-ecdsa
include "../0xParcECDSA/ecdsa.circom";

// r, s, msghash, and pubkey have coordinates
// encoded with k registers of n bits each
// signature is (r, s)
// Does not check that pubkey is valid
// Requires 56GB of RAM to run
component main {public [r, s, msghash, pubkey]} = ECDSAVerifyNoPubkeyCheck(64, 4);

//./scripts/run_circuit.sh circuits/benchmarks/ecdsa/circuit.circom ecdsa inputs/benchmarks/ecdsa/input.json phase1/powersOfTau28_hez_final_21.ptau res.csv ecdsa_output