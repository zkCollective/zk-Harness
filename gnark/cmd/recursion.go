package cmd

import (
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/constraint"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
	"github.com/consensys/gnark/logger"
	"github.com/spf13/cobra"
	"github.com/zkCollective/zk-Harness/gnark/circuits"
	"github.com/zkCollective/zk-Harness/gnark/util"
)

// groth16Cmd represents the groth16 command
var recursionCmd = &cobra.Command{
	Use:   "recursion",
	Short: "runs benchmarks for recursion",
	Run:   runOneStep,
}

// TODO - Remove hardcoded MIMC values
const (
	preImage   = "4992816046196248432836492760315135318126925090839638585255611512962528270024"
	publicHash = "7831393781387060555412927989411398077996792073838215843928284475008119358174"
)

func computeInnerProofG16(fcircuitSize int, fcircuit string, finputPath string) (groth16.VerifyingKey, groth16.Proof) {
	ccs, err := frontend.Compile(curveID.ScalarField(), r1cs.NewBuilder, c.Circuit(fcircuitSize, fcircuit, circuits.WithInputCircuit(finputPath)), frontend.WithCapacity(fcircuitSize))
	witness := c.Witness(fcircuitSize, curveID, fcircuit, circuits.WithInputWitness(finputPath))
	pk, vk, err := groth16.Setup(ccs)
	assertNoError(err)
	proof, err := groth16.Prove(ccs, pk, witness)
	assertNoError(err)
	publicWitness, err := witness.Public()
	assertNoError(err)
	// Check that proof verifies before continuing
	if err := groth16.Verify(proof, vk, publicWitness); err != nil {
		panic(err)
	}
	return vk, proof
}

func runOneStep(cmd *cobra.Command, args []string) {
	log := logger.Logger()
	log.Info().Msg("Benchmarking " + *fCircuit + " - gnark, recursion: " + *fAlgo + " " + *fCurve + " " + *fInputPath)

	var filename = "../benchmarks/gnark/gnark_" +
		"recursion" + "_" +
		*fCircuit + "." +
		*fFileType

	if err := parseFlags(); err != nil {
		fmt.Println("error: ", err.Error())
		cmd.Help()
		os.Exit(-1)
	}

	writeResults := func(took time.Duration, ccs constraint.ConstraintSystem, proof_size int) {

		// check memory usage, max ram requested from OS
		var m runtime.MemStats
		runtime.ReadMemStats(&m)

		_, secret, public := ccs.GetNbVariables()
		bData := util.BenchDataCircuit{
			Framework:         "gnark",
			Category:          "circuit",
			Backend:           "groth16",
			Curve:             curveID.String(),
			Circuit:           *fCircuit,
			Input:             *fInputPath,
			Operation:         *fAlgo,
			NbConstraints:     ccs.GetNbConstraints(),
			NbSecretVariables: secret,
			NbPublicVariables: public,
			ProofSize:         proof_size,
			MaxRAM:            (m.Sys / 1024 / 1024),
			Count:             *fCount,
			RunTime:           took.Milliseconds(),
		}

		if err := util.WriteData("csv", bData, filename); err != nil {
			panic(err)
		}
	}

	// pre-compute the plain result
	// TODO - replace with switch case
	hash := util.PreCalcMIMC(curveID, preImage)

	// pre-compute the inner proof
	innerVk, innerProof := computeInnerProofG16(*fCircuitSize, *fCircuit, *fInputPath)

	// compute the outer proof
	// TODO - replace hardcoded value
	*fCircuit = "groth16_bls12377"
	// TODO - remove hardcode
	*fAlgo = "recursion"

	// Run Benchmarks for Groth16 for the outer proof
	benchGroth16(
		writeResults,
		*fAlgo,
		*fCount,
		*fCircuitSize,
		*fCircuit,
		WithVK(innerVk),
		WithProof(innerProof),
		WithWitness(hash))
}

func init() {
	rootCmd.AddCommand(groth16Cmd)
	// groth16Cmd.Flags().StringVar(&inputPath, "input", "none", "input path to the dedicated input")
	// print(inputPath)
}
