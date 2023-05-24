#!/usr/bin/env bash
#
# Compile, setup, prove, and verify a proof for a circom circuit.
# If the .csv files already exist, then just append the results
# TODO add C++ witness generation support
# TODO try add PLONK and FFLONKsupport

if [ $# -lt 5 ]; then
    echo $0: usage: run_circuit.sh circuit.circom circuit_name input.json powersOfTau.ptau results.csv tmp template_vars
    exit 1
fi

if [ ! -z "$IN_NIX_SHELL" ]; then
    TIMEBIN="/usr/bin/env time"
else
    TIMEBIN="/usr/bin/time"
fi

SCRIPT_PATH=$(dirname "$(readlink -f "$0")")
NODE_MODULES="${SCRIPT_PATH}/../../node_modules/"
CIRCUIT=$1
CIRCUIT_NAME=$2
CIRCUIT_NAME_INT=${CIRCUIT##*/}
CIRCUIT_NAME_INT=${CIRCUIT_NAME_INT%.circom}
INPUT=$3
TAU=$4
RES=$5
if [ ! -z "$6" ]; then
    TMP=$6
else
    TMP=tmp
fi
if [ ! -z "$7" ]; then
    TEMPLATE_VARS=$7
else
    TEMPLATE_VARS=
fi
echo ">>> Running with: $CIRCUIT, $CIRCUIT_NAME, $INPUT, $TAU, $RES, $TMP, $TEMPLATE_VARS"

if [[ $(uname) == "Linux" ]]; then
    TIMECMD="$TIMEBIN -f \"Real time (seconds): %e\nMaximum resident set size (bytes): %M\" -o"
    STATCMD='stat --printf="%s" '
    OS="Linux"
elif [[ $(uname) == "Darwin" ]]; then
    if [ ! -z "$IN_NIX_SHELL" ]; then
        TIMECMD="$TIMEBIN -f \"Real time (seconds): %e\nMaximum resident set size (bytes): %M\" -o"
        STATCMD='stat --printf="%s" '
    else
        TIMECMD="$TIMEBIN -h -l -o"
        STATCMD='stat -f%z '
    fi
    OS="Darwin"
else
    echo "Unsupported operating system."
    exit 1
fi

INPUT_FILENAME=$(basename $INPUT)
NEW_INPUT=${TMP}/${INPUT_FILENAME}
ORIGINAL_CIRCOM_CONTENTS=$(cat "$CIRCUIT")

handle_template_vars() {
    if [ ! -z "$TEMPLATE_VARS" ]; then
        json=$(cat "$NEW_INPUT")
        # Initialize an empty array to store the keys to remove
        keys_to_remove=()
        # Initialize an empty array to store values to replace TEMPLATE_VARS
        values_to_replace=()
        # Split the input values by commas
        IFS=',' read -ra values_array <<< "$TEMPLATE_VARS"
        # Loop through the values
        for value in "${values_array[@]}"; do
            # Retrieve the value from the JSON
            retrieved_value=$(echo "$json" | jq -r ".$value")
            # Add the value to the array
            values_to_replace+=("$retrieved_value")
            # Remove the corresponding key from the JSON 
            # and add it to the keys to remove
            json=$(echo "$json" | jq "del(.${value})")
            keys_to_remove+=("$value")
        done
        # Remove the keys from the JSON
        for key in "${keys_to_remove[@]}"; do
            json=$(echo "$json" | jq "del(.${key})")
        done
        # Replace {TEMPLATE_VARS} with the passed values
        new_contents=$(echo "$ORIGINAL_CIRCOM_CONTENTS" | sed "s/{TEMPLATE_VARS}/$(echo "${values_to_replace[*]}" | tr ' ' ',')/g")
        # Save the modified JSON back to input.json
        echo "$json" > "$NEW_INPUT"
        # Save the modified circom files back to circom
        echo "$new_contents" > "$CIRCUIT"
    fi
}


### EXECUTION ###
echo ">>>Step 0: cleaning and creating ${TMP}" && \
rm -rf ${TMP} && mkdir ${TMP} && \
cp ${INPUT} ${NEW_INPUT} && \
handle_template_vars && \
echo ">>>Step 1: compiling the circuit" && \
eval """
$TIMECMD ${TMP}/compiler_times.txt circom ${CIRCUIT} --r1cs --wasm --sym --c --output ${TMP} | tee ${TMP}/circom_output 
""" && \
# Revert the circom file contents
echo "$ORIGINAL_CIRCOM_CONTENTS" > ${CIRCUIT} && \
echo ">>>Step 2: generating the witness JS" && \
eval """
$TIMECMD ${TMP}/witness_times.txt node --max_old_space_size=50000 ${TMP}/${CIRCUIT_NAME_INT}_js/generate_witness.js ${TMP}/${CIRCUIT_NAME_INT}_js/${CIRCUIT_NAME_INT}.wasm ${NEW_INPUT} ${TMP}/witness.wtns
""" && \
# We only care about phase 2 which is circuit-specific
# .zkey file that will contain the proving and verification keys together with 
# all phase 2 contributions.
echo ">>>Step 3: Setup" && \
eval """
$TIMECMD ${TMP}/setup_times.txt node --max_old_space_size=50000 ${NODE_MODULES}/snarkjs/cli.js groth16 setup ${TMP}/${CIRCUIT_NAME_INT}.r1cs ${TAU} ${TMP}/${CIRCUIT_NAME_INT}_0.zkey
""" && \
# TODO Should we contribute here?
# We could contribute here using: snarkjs zkey contribute ${TMP}/${CIRCUIT_NAME}_0.zkey ${TMP}/${CIRCUIT_NAME}_1.zkey --name="1st Contributor Name" -v
echo ">>>Step 4: Export verification key" && \
eval """
$TIMECMD ${TMP}/export_times.txt node --max_old_space_size=50000 ${NODE_MODULES}/snarkjs/cli.js zkey export verificationkey ${TMP}/${CIRCUIT_NAME_INT}_0.zkey ${TMP}/verification_key.json
""" && \
echo ">>>Step 5: Prove" && \
eval """
$TIMECMD ${TMP}/prove_times.txt node --max_old_space_size=50000 ${NODE_MODULES}/snarkjs/cli.js groth16 prove ${TMP}/${CIRCUIT_NAME_INT}_0.zkey ${TMP}/witness.wtns ${TMP}/proof.json ${TMP}/public.json
""" && \
echo ">>>Step 6: Verify" && \
eval """
$TIMECMD ${TMP}/verify_times.txt node --max_old_space_size=50000 ${NODE_MODULES}/snarkjs/cli.js groth16 verify ${TMP}/verification_key.json ${TMP}/public.json ${TMP}/proof.json
"""

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

    if [[ "$OS" == "Linux" ]] || [ ! -z "$IN_NIX_SHELL" ]; then
        ram=$(grep Maximum ${timeRes} | cut -d ":" -f2 | xargs)
        realTime=$(grep Real ${timeRes} | cut -d ":" -f2 | xargs)
        # RAM here is in kbytes
        ramMb=$(echo ${ram}/1024 | bc)
    elif [[ "$OS" == "Darwin" ]]; then
        ram=$(grep maximum ${timeRes} | xargs | cut -d " " -f1) 
        realTime=$(grep real ${timeRes} | xargs | cut -d " " -f1)
        ramMb=$(echo ${ram}/1024/1024 | bc)
    fi
    # NOTE: if real contains minutes in Mac it won't work
    realTime=$(echo "$realTime" | sed 's/s//')
    millisecs=$(echo "${realTime} * 1000" | bc)
    millisecs_without_dec=${millisecs%.*}
    echo "$ramMb,$millisecs_without_dec"
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
        proofSize=$(eval "$STATCMD ${TMP}/proof.json")
    else
        proofSize=""
    fi
    # If phaseTimeFileToMerge is not empty then merge its results with phaseTimeFile 
    if [ ! -z "$phaseTimeFileToMerge" ]; then
        ramtimeToMerge="$(get_time_results $phaseTimeFileToMerge)"
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
    # We don't want to print the whole input file but only the part that it is
    # Count is always 1
    echo "circom/snarkjs,circuit,groth16,bn128,$CIRCUIT_NAME,$INPUT,$phase,$nbConstraints,$nbPrivateInputSignals,$nbPublicInputSignals,$ramtimeFinal,$proofSize,$physical,$virtual,1,$PROC"

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

    # Check if RES file already exist.
    if [ ! -f "$RES" ]; then
        echo "framework,category,backend,curve,circuit,input,operation,nbConstraints,nbSecret,nbPublic,ram(mb),time(ms),proofSize,nbPhysicalCores,nbLogicalCores,count,cpu" > ${RES}
    fi
    for (( i=0; i<${arraylength}; i++ ));
    do
      echo "$(get_phase_stats ${stages[$i]} ${times[$i]})" >> ${RES}
    done
fi
