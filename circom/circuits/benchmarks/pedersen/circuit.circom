pragma circom 2.0.0;

// Source: https://github.com/iden3/circomlib/blob/master/circuits/pedersen.circom
include "../circomlib/pedersen.circom";
include "../circomlib/bitify.circom";

template PedersenBench(){
    signal input PreImage;
    component n2b = Num2Bits_strict();
    n2b.in <== PreImage;
    component pedersen = Pedersen(254);
    pedersen.in <== n2b.out;
}
component main {public[PreImage]} = PedersenBench();

//./scripts/run_circuit.sh circuits/benchmarks/pedersen/circuit.circom pedersen inputs/benchmarks/pedersen/input.json phase1/powersOfTau28_hez_final_16.ptau res.csv pedersen_output
