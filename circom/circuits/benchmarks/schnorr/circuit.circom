pragma circom 2.0.0;

include "./schnorr.circom";

component main {public[enabled, M, yx, yy, S, e]} = SchnorrPosedion(
    2736030358979909402780800718157159386076813972158567259200215660948447373041,5299619240641551281634865583518297030282874472190772894086521144482721001553,16950150798460657717958625567821834550301663161624707787222815936182638968203);

//./scripts/run_circuit.sh circuits/benchmarks/schnorr/circuit.circom schnorr inputs/benchmarks/schnorr/input.json phase1/powersOfTau28_hez_final_16.ptau res.csv schnorr_output