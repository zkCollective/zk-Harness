pragma circom 2.0.0;

// Source: https://github.com/iden3/circomlib/blob/master/circuits/sha256/sha256.circom
include "../circomlib/sha256/sha256.circom";
include "../circomlib/bitify.circom";

template Sha256Bench(){
    signal input PreImage;
    component n2b = Num2Bits_strict();
    n2b.in <== PreImage;
    component sha256 = Sha256(254);
    sha256.in <== n2b.out;
}
component main {public[PreImage]} = Sha256Bench();

//./scripts/run_circuit.sh circuits/benchmarks/pedersen/circuit.circom sha256 inputs/benchmarks/pedersen/input.json phase1/powersOfTau28_hez_final_16.ptau res.csv sha256_output

