pragma circom 2.0.0;

// Source: https://github.com/iden3/circomlib/blob/master/circuits/smt/smtverifier.circom
include "../circomlib/smt/smtverifier.circom";

component main {public[enabled,fnc,root,siblings,oldKey,oldValue,isOld0,key,value]} = SMTVerifier(10);

//./scripts/run_circuit.sh circuits/benchmarks/smt/circuit.circom smt inputs/benchmarks/smt/input.json phase1/powersOfTau28_hez_final_16.ptau res.csv smt_output
       