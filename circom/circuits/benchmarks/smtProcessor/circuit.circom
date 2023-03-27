pragma circom 2.0.0;

include "../../circomlib/smt/smtprocessor.circom";

component main {public[fnc,oldRoot,siblings,oldKey,oldValue,isOld0,newKey,newValue]} = SMTProcessor(10);

//./scripts/run_circuit.sh circuits/benchmarks/smtProcessor/circuit.circom smtProcessor inputs/benchmarks/smtProcessor/input.json phase1/powersOfTau28_hez_final_16.ptau res.csv smtProcessor_output
       