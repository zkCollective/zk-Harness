#!/bin/bash
#
# Compile, setup, prove, and verify a proof for a circom circuit.
# TODO add C++ witness generation support
# TODO try add PLONK and FFLONKsupport

if [ $# -lt 3 ]; then
    echo $0: usage: run_circuit.sh circuit.circom input.json powersOfTau.ptau results.csv tmp
    exit 1
fi

CIRCUIT=$1
CIRCUIT_NAME=${CIRCUIT##*/}
CIRCUIT_NAME=${CIRCUIT_NAME%.circom}
INPUT=$2
TAU=$3
RES=$4
if [ ! -z "$5" ]; then
    TMP=$5
else
    TMP=tmp
fi

### EXECUTION ###
echo ">>>Step 0: cleaning and creating ${TMP}" && \
rm -rf ${TMP} && mkdir ${TMP} && \
echo ">>>Step 1: compiling the circuit" && \
/usr/bin/time -h -l -o ${TMP}/compiler_times.txt circom ${CIRCUIT} --r1cs --wasm --sym --c --output ${TMP} | tee ${TMP}/circom_output && \
echo ">>>Step 2: generating the witness JS" && \
/usr/bin/time -h -l -o ${TMP}/witness_times.txt node  ${TMP}/${CIRCUIT_NAME}_js/generate_witness.js ${TMP}/${CIRCUIT_NAME}_js/${CIRCUIT_NAME}.wasm ${INPUT} ${TMP}/witness.wtns && \
# We only care about phase 2 which is circuit-specific
# .zkey file that will contain the proving and verification keys together with 
# all phase 2 contributions.
echo ">>>Step 3: Setup" && \
/usr/bin/time -h -l -o ${TMP}/setup_times.txt snarkjs groth16 setup ${TMP}/${CIRCUIT_NAME}.r1cs ${TAU} ${TMP}/${CIRCUIT_NAME}_0.zkey && \
# TODO Should we contribute here?
# We could contribute here using: snarkjs zkey contribute ${TMP}/${CIRCUIT_NAME}_0.zkey ${TMP}/${CIRCUIT_NAME}_1.zkey --name="1st Contributor Name" -v
echo ">>>Step 4: Export verification key" && \
/usr/bin/time -h -l -o ${TMP}/export_times.txt snarkjs zkey export verificationkey ${TMP}/${CIRCUIT_NAME}_0.zkey ${TMP}/verification_key.json && \
echo ">>>Step 5: Prove" && \
/usr/bin/time -h -l -o ${TMP}/prove_times.txt snarkjs groth16 prove ${TMP}/${CIRCUIT_NAME}_0.zkey ${TMP}/witness.wtns ${TMP}/proof.json ${TMP}/public.json && \
echo ">>>Step 6: Verify" && \
/usr/bin/time -h -l -o ${TMP}/verify_times.txt snarkjs groth16 verify ${TMP}/verification_key.json ${TMP}/public.json ${TMP}/proof.json

portable_proc() {
    OS="$(uname -s)"
    if [ "$OS" = "Linux" ]; then
        PROC="$(lscpu | grep 'Model name:' | cut -d ':' -f2 | xargs)"
    elif [ "$OS" = "Darwin" ] || \
         [ "$(echo "$OS" | grep -q BSD)" = "BSD" ]; then
        PROC="$(sysctl -a | grep machdep.cpu.brand_string | cut -d ':' -f2 | xargs)"
    else
        PROC=""  
    fi
    echo "$PROC"
}

get_time_results() {
    timeRes=$1
    ram=$(grep maximum ${timeRes} | xargs | cut -d " " -f1) 
    ramMb=$(echo ${ram}/1024/1024 | bc)
    realTime=$(grep real ${timeRes} | xargs | cut -d " " -f1)
    realTime=${realTime::${#realTime}-1}
    milisecs=$(echo "$realTime * 1000" | bc)
    milisecs=${milisecs::${#milisecs}-3}
    echo "$ramMb,$milisecs"
}

get_phase_stats() {
    phase=$1
    phaseTimeFile=$2
    phaseTimeFileToMerge=$3

    ramtime="$(get_time_results $phaseTimeFile)"
    # TODO Node uses 1 single thread in one core. Nevertheless, snarkjs
    # (and the underlying library ffjavascript)
    # use workers to perfrom operations, hence it isn't actually single threaded.
    # We could instrument snarkjs so it uses only a single worker.
    # For circom compiler, again we could potentially enforce single core and
    # thread execution if we instrument it.
    # Finally, we might need to do the same for the witness generator.
    physical=1
    virtual=1
    # Proof size in bytes
    if [ $phase == "prove" ]; then 
        proofSize=$(stat -f%z ${TMP}/proof.json)
    else
        proofSize=""
    fi
    # If phaseTimeFileToMerge is not empty then merge its results with phaseTimeFile 
    if [ ! -z "$phaseTimeFileToMerge" ]; then
        ramtimeToMerge="$(get_time_results $phaseTimeFileToMerge)"
        echo $ramtimeToMerge > fff
        ramInitial=$(echo $ramtime | cut -d ',' -f1)
        ramToMerge=$(echo $ramtimeToMerge | cut -d ',' -f1)
        timeInitial=$(echo $ramtime | cut -d ',' -f2)
        timeToMerge=$(echo $ramtimeToMerge | cut -d ',' -f2)
        if [ "$ramInitial" -gt "$ramToMerge" ]; then
            ramFinal=$ramInitial
        else
            ramFinal=$ramToMerge
        fi
        timeFinal=$(($timeInitial + $timeToMerge))
        ramtimeFinal="${ramFinal},${timeFinal}"
    else
        ramtimeFinal=$ramtime
    fi
    # TODO add proof size before phyical cores
    echo "circom,circuit,groth16,bn128,$CIRCUIT_NAME,$INPUT,$phase,$nbConstraints,$nbPrivateInputSignals,$nbPublicInputSignals,$ramtimeFinal,$physical,$virtual,$PROC"

}

if [ ! -z "$RES" ]; then
    PROC=$(portable_proc)
    nbLinearConstraints=$(grep "^linear constraints" ${TMP}/circom_output | cut -d ":" -f2 | xargs)
    nbNonLinearConstraints=$(grep "^non-linear constraints" ${TMP}/circom_output | cut -d ":" -f2 | xargs)
    nbConstraints=$(($nbLinearConstraints+$nbNonLinearConstraints))
    nbPrivateInputSignals=$(grep "^private inputs" ${TMP}/circom_output | cut -d ":" -f2 | xargs)
    nbPublicInputSignals=$(grep "^public inputs" ${TMP}/circom_output | cut -d ":" -f2 | xargs)
    nbIntermediateSignals=$(grep "signal" ${CIRCUIT} | wc -l | xargs)
    nbIntermediateSignals=$(($nbIntermediateSignals-$nbPrivateInputSignals-$nbPublicInputSignals))
    nbPrivateOuputSignals=$(grep "^private outputs" ${TMP}/circom_output | cut -d ":" -f2 | xargs)
    nbPublicOutputSignals=$(grep "^public outputs" ${TMP}/circom_output | cut -d ":" -f2 | xargs)
    nbWires=$(grep "^wires" ${TMP}/circom_output | cut -d ":" -f2 | xargs)
    nbLabels=$(grep "^labels" ${TMP}/circom_output | cut -d ":" -f2 | xargs)
    declare -a stages=("compile" 
                       "witness" 
                       "setup"
                       "prove"
                       "verify"
                      )
    declare -a times=("${TMP}/compiler_times.txt" 
                      "${TMP}/witness_times.txt" 
                      "${TMP}/setup_times.txt ${TMP}/export_times.txt" 
                      "${TMP}/prove_times.txt" 
                      "${TMP}/verify_times.txt" 
                     )
    arraylength=${#stages[@]}

    # TODO add proof size before phyical cores
    echo "framework,category,backend,curve,circuit,input,operation,nbConstraints,nbSecret,nbPublic,ram(mb),time(ms),nbPhysicalCores,nbLogicalCores,cpu" > ${RES}
    for (( i=0; i<${arraylength}; i++ ));
    do
      echo "$(get_phase_stats ${stages[$i]} ${times[$i]})" >> ${RES}
    done
fi

