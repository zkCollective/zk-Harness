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
	"github.com/zkCollective/zk-Harness/frameworks/gnark/circuits"
	"github.com/zkCollective/zk-Harness/frameworks/gnark/parser"
	"github.com/zkCollective/zk-Harness/frameworks/gnark/util"
)

var recursionCmd = &cobra.Command{
	Use:   "recursion",
	Short: "runs benchmarks for recursion",
	Run:   runOneStep,
}

var witness frontend.Variable

var recursiveCircuit string

func computeInnerProofG16(fcircuitSize int, fcircuit string, finputPath string, finnerCurveID ecc.ID) (groth16.VerifyingKey, groth16.Proof, constraint.ConstraintSystem) {
	circuit := parser.C.Circuit(fcircuitSize, fcircuit, circuits.WithInputCircuit(finputPath))
	ccs, err := frontend.Compile(finnerCurveID.ScalarField(), r1cs.NewBuilder, circuit, frontend.WithCapacity(fcircuitSize))
	witness := parser.C.Witness(fcircuitSize, finnerCurveID, fcircuit, circuits.WithInputWitness(finputPath))
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
	log.Info().Msg("Benchmarking " + *cfg.Circuit + " - gnark, recursion: " + *cfg.Algo + " " + *cfg.Curve + " " + *cfg.InputPath)

	var filename = *cfg.OutputPath

	if err := parser.ParseFlags(cfg); err != nil {
		fmt.Println("error: ", err.Error())
		cmd.Help()
		os.Exit(-1)
	}

	var data map[string]interface{}
	if *cfg.InputPath != "none" {
		var err error
		data, err = util.ReadFromInputPath(*cfg.InputPath)
		if err != nil {
			panic(err)
		}
	}

	// Set inner curve based on outer curve
	switch *cfg.Curve {
	case "bw6_761":
		parser.InnerCurveID = ecc.BLS12_377
		recursiveCircuit = "groth16_bls12377"
	case "bw6_633":
		parser.InnerCurveID = ecc.BLS24_315
		recursiveCircuit = "groth16_bls24315"
	default:
		panic("Chosen Curve not implemented for 2-Chain recursion! Must be bw6_761 or bw6_633")
	}

	// pre-compute the inner G16 proof, return innerCCS to get num Constraints inner
	innerVk, innerProof, innerCCS := computeInnerProofG16(*cfg.CircuitSize, *cfg.Circuit, *cfg.InputPath, parser.InnerCurveID)

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
			InnerCurve:         parser.InnerCurveID.String(),
			OuterBackend:       *cfg.OuterBackend,
			OuterCurve:         parser.CurveID.String(),
			Circuit:            *cfg.Circuit,
			Input:              *cfg.InputPath,
			Operation:          *cfg.Algo,
			InnerNbConstraints: innerCCS.GetNbConstraints(),
			NbConstraints:      ccs.GetNbConstraints(),
			NbSecretVariables:  secret,
			NbPublicVariables:  public,
			ProofSize:          proof_size,
			MaxRAM:             (m.Sys / 1024 / 1024),
			Count:              *cfg.Count,
			RunTime:            took.Milliseconds(),
		}

		if err := util.WriteData("csv", bData, filename); err != nil {
			panic(err)
		}
	}

	switch *cfg.OuterBackend {
	case "groth16":
		switch *cfg.Circuit {
		case "mimc":
			witness = util.PreCalcMIMC(parser.InnerCurveID, (data["PreImage"].(string)))
		case "cubic":
			witness = (data["Y"].(string))
		case "bench":
			witness = (data["Y"].(string))
		default:
			panic("Chosen Circuit not implemented for Groth16 recursion!")
		}
		benchGroth16(
			writeResults,
			*cfg.Algo,
			*cfg.Count,
			*cfg.CircuitSize,
			recursiveCircuit,
			util.WithVK(innerVk),
			util.WithProof(innerProof),
			util.WithWitness(witness))
	case "plonk":
		switch *cfg.Circuit {
		case "mimc":
			witness = util.PreCalcMIMC(parser.InnerCurveID, (data["PreImage"].(string)))
		case "cubic":
			witness = (data["Y"].(string))
		case "bench":
			witness = (data["Y"].(string))
		default:
			panic("Circuit not implemented for Plonk recursion!")
		}
		benchPlonk(
			writeResults,
			*cfg.Algo,
			*cfg.Count,
			*cfg.CircuitSize,
			recursiveCircuit,
			util.WithVK(innerVk),
			util.WithProof(innerProof),
			util.WithWitness(witness))
	case "plonkFRI":
		switch *cfg.Circuit {
		case "mimc":
			witness = util.PreCalcMIMC(parser.InnerCurveID, (data["PreImage"].(string)))
		case "cubic":
			witness = (data["Y"].(string))
		case "bench":
			witness = (data["Y"].(string))
		default:
			panic("Circuit not implemented for PlonkFRI recursion!")
		}
		benchPlonkFRI(
			writeResults,
			*cfg.Algo,
			*cfg.Count,
			*cfg.CircuitSize,
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
