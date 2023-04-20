package cmd

import (
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/constraint"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
	"github.com/consensys/gnark/logger"
	"github.com/spf13/cobra"
	"github.com/zkCollective/zk-Harness/gnark/circuits"
	"github.com/zkCollective/zk-Harness/gnark/util"
)

var recursionCmd = &cobra.Command{
	Use:   "recursion",
	Short: "runs benchmarks for recursion",
	Run:   runOneStep,
}

var witness frontend.Variable

var recursiveCircuit string

func computeInnerProofG16(fcircuitSize int, fcircuit string, finputPath string, innerCurveID ecc.ID) (groth16.VerifyingKey, groth16.Proof, constraint.ConstraintSystem) {
	fmt.Println("COMPUTING INNER PROOF")
	circuit := c.Circuit(fcircuitSize, fcircuit, circuits.WithInputCircuit(finputPath))
	ccs, err := frontend.Compile(innerCurveID.ScalarField(), r1cs.NewBuilder, circuit, frontend.WithCapacity(fcircuitSize))
	witness := c.Witness(fcircuitSize, innerCurveID, fcircuit, circuits.WithInputWitness(finputPath))
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
	return vk, proof, ccs
}

func runOneStep(cmd *cobra.Command, args []string) {
	log := logger.Logger()
	log.Info().Msg("Benchmarking " + *fCircuit + " - gnark, recursion: " + *fAlgo + " " + *fCurve + " " + *fInputPath)

	var filename = "../benchmarks/gnark/gnark_" +
		"recursion" + "_" +
		*fCircuit + "_" +
		*fCurve + "." +
		*fFileType

	if err := parseFlags(); err != nil {
		fmt.Println("error: ", err.Error())
		cmd.Help()
		os.Exit(-1)
	}

	var data map[string]interface{}
	if *fInputPath != "none" {
		var err error
		data, err = util.ReadFromInputPath(*fInputPath)
		if err != nil {
			panic(err)
		}
	}

	// Set inner curve based on outer curve
	switch *fCurve {
	case "bw6_761":
		innerCurveID = ecc.BLS12_377
		recursiveCircuit = "groth16_bls12377"
	case "bw6_633":
		innerCurveID = ecc.BLS24_315
		recursiveCircuit = "groth16_bls24315"
	}

	// pre-compute the inner G16 proof, return innerCCS to get num Constraints inner
	innerVk, innerProof, innerCCS := computeInnerProofG16(*fCircuitSize, *fCircuit, *fInputPath, innerCurveID)

	writeResults := func(took time.Duration, ccs constraint.ConstraintSystem, proof_size int) {

		// check memory usage, max ram requested from OS
		var m runtime.MemStats
		runtime.ReadMemStats(&m)

		_, secret, public := ccs.GetNbVariables()

		// FIXME - hardcoded inner backend
		bData := util.BenchDataRecursion{
			Framework:          "gnark",
			Category:           "circuit",
			InnerBackend:       "groth16",
			InnerCurve:         innerCurveID.String(),
			OuterBackend:       *fOuterBackend,
			OuterCurve:         curveID.String(),
			Circuit:            *fCircuit,
			Input:              *fInputPath,
			Operation:          *fAlgo,
			InnerNbConstraints: innerCCS.GetNbConstraints(),
			NbConstraints:      ccs.GetNbConstraints(),
			NbSecretVariables:  secret,
			NbPublicVariables:  public,
			ProofSize:          proof_size,
			MaxRAM:             (m.Sys / 1024 / 1024),
			Count:              *fCount,
			RunTime:            took.Milliseconds(),
		}

		if err := util.WriteData("csv", bData, filename); err != nil {
			panic(err)
		}
	}

	switch *fOuterBackend {
	case "groth16":
		switch *fCircuit {
		case "mimc":
			witness = util.PreCalcMIMC(innerCurveID, (data["PreImage"].(string)))
		case "cubic":
			witness = (data["Y"].(string))
		case "bench":
			witness = util.PreCalcBench(*fCircuitSize, innerCurveID)
		case "expo":
			witness = util.PreCalcBench(*fCircuitSize, innerCurveID)
		default:
			panic("Circuit not implemented for Groth16 recursion!")
		}
		benchGroth16(
			writeResults,
			*fAlgo,
			*fCount,
			*fCircuitSize,
			recursiveCircuit,
			util.WithVK(innerVk),
			util.WithProof(innerProof),
			util.WithWitness(witness))
	case "plonk":
		switch *fCircuit {
		case "mimc":
			witness = util.PreCalcMIMC(innerCurveID, (data["PreImage"].(string)))
		case "cubic":
			witness = (data["Y"].(string))
		case "bench":
			witness = util.PreCalcBench(*fCircuitSize, innerCurveID)
		case "expo":
			witness = util.PreCalcBench(*fCircuitSize, innerCurveID)
		default:
			panic("Circuit not implemented for Plonk recursion!")
		}
		benchPlonk(
			writeResults,
			*fAlgo,
			*fCount,
			*fCircuitSize,
			recursiveCircuit,
			util.WithVK(innerVk),
			util.WithProof(innerProof),
			util.WithWitness(witness))
	case "plonkFRI":
		switch *fCircuit {
		case "mimc":
			witness = util.PreCalcMIMC(innerCurveID, (data["PreImage"].(string)))
		case "cubic":
			witness = (data["Y"].(string))
		case "bench":
			witness = util.PreCalcBench(*fCircuitSize, innerCurveID)
		case "expo":
			witness = util.PreCalcBench(*fCircuitSize, innerCurveID)
		default:
			panic("Circuit not implemented for PlonkFRI recursion!")
		}
		benchPlonkFRI(
			writeResults,
			*fAlgo,
			*fCount,
			*fCircuitSize,
			recursiveCircuit,
			util.WithVK(innerVk),
			util.WithProof(innerProof),
			util.WithWitness(witness))
	default:
		panic("Outer backend not supported!")
	}

}

func init() {
	rootCmd.AddCommand(recursionCmd)
}
