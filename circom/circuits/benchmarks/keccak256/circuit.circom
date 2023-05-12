pragma circom 2.0.0;

// Source: https://github.com/vocdoni/keccak256-circom
include "./keccak.circom";
include "../circomlib/bitify.circom";

template KeccakBench(){
    signal input PreImage;
    component n2b = Num2Bits_strict();
    n2b.in <== PreImage;
    component keccak = Keccak(254,254);
    keccak.in <== n2b.out;
}
component main {public[PreImage]} = KeccakBench();


//./scripts/run_circuit.sh circuits/benchmarks/keccak256/circuit.circom keccak256 inputs/benchmarks/keccak256/input.json phase1/powersOfTau28_hez_final_20.ptau res.csv keccak256_output