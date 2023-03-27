pragma circom 2.0.0;

include "../../circomlib/mimc.circom";

component main {public[x_in,k]} = MiMC7(91);


//./scripts/run_circuit.sh circuits/benchmarks/mimc/circuit.circom mimc inputs/benchmarks/mimc/input.json phase1/powersOfTau28_hez_final_16.ptau res.csv mimc_output
