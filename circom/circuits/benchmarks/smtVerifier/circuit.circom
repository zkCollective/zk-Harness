pragma circom 2.0.0;

include "../../circomlib/smt/smtverifier.circom";

component main {public[enabled,fnc,root,siblings,oldKey,oldValue,isOld0,key,value]} = SMTVerifier(10);

//./scripts/run_circuit.sh circuits/benchmarks/smtVerifier/circuit.circom smtVerifier inputs/benchmarks/smtVerifier/input.json phase1/powersOfTau28_hez_final_16.ptau res.csv smtVerifier_output
       