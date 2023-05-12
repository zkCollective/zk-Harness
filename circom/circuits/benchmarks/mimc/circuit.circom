pragma circom 2.0.0;

// Source: https://github.com/iden3/circomlib/blob/master/circuits/mimc.circom
include "../circomlib/mimc.circom";

template OurMimc() {
    signal input X;
    signal input K;
    component mimc = MiMC7(91);
    mimc.x_in <== X;
    mimc.k <== K;
}

component main {public[X,K]} = OurMimc();


//./scripts/run_circuit.sh circuits/benchmarks/mimc/circuit.circom mimc inputs/benchmarks/mimc/input.json phase1/powersOfTau28_hez_final_16.ptau res.csv mimc_output
