package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/backend/witness"
	"github.com/spf13/cobra"
	"github.com/zkCollective/zk-Harness/gnark/parser"
)

// groth16Cmd represents the groth16 command
var plonkMemoryVerifyCmd = &cobra.Command{
	Use:   "plonkMemoryVerify",
	Short: "runs benchmarks and profiles using Groth16 proof system",
	Run:   runPlonkMemoryVerify,
}

var cfg = parser.NewConfig()

func runPlonkMemoryVerify(cmd *cobra.Command, args []string) {

	if err := parser.ParseFlagsMemory(cfg); err != nil {
		fmt.Println("error: ", err.Error())
		cmd.Help()
		os.Exit(-1)
	}

	// Initialize variables
	proof := groth16.NewProof(parser.CurveID)
	vk := groth16.NewVerifyingKey(parser.CurveID)
	publicWitness, err := witness.New(parser.CurveID.ScalarField())

	// Read vk
	f, err := os.Open("tmp/vk.dat")
	if err != nil {
		panic("Failed to open file: " + err.Error())
	}

	_, err = vk.ReadFrom(f)
	if err != nil {
		panic("Failed to read from file: " + err.Error())
	}

	// Read Public Witness
	f, err = os.Open("tmp/publicWitness.dat")
	if err != nil {
		panic("Failed to open file: " + err.Error())
	}

	_, err = publicWitness.ReadFrom(f)
	if err != nil {
		panic("Failed to read from file: " + err.Error())
	}

	f.Close()

	// Read proof
	f, err = os.Open("tmp/proof.dat")
	if err != nil {
		panic("Failed to open file: " + err.Error())
	}

	_, err = proof.ReadFrom(f)
	if err != nil {
		panic("Failed to read from file: " + err.Error())
	}

	f.Close()

	// Proof Verification
	err = groth16.Verify(proof, vk, publicWitness)
	if err != nil {
		panic("Failed Verification!")
	}

	return
}

func Execute() {
	if err := plonkMemoryVerifyCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {

	cobra.OnInitialize()

	_curves := ecc.Implemented()
	curves := make([]string, len(_curves))
	for i := 0; i < len(_curves); i++ {
		curves[i] = strings.ToLower(_curves[i].String())
	}

	cfg.InputPath = plonkMemoryVerifyCmd.PersistentFlags().String("input", "none", "input path to the dedicated input")
	plonkMemoryVerifyCmd.MarkPersistentFlagRequired("input")
	cfg.Circuit = plonkMemoryVerifyCmd.PersistentFlags().String("circuit", "expo", "name of the circuit to use")
	cfg.CircuitSize = plonkMemoryVerifyCmd.PersistentFlags().Int("size", 10000, "size of the circuit, parameter to circuit constructor")
	cfg.Count = plonkMemoryVerifyCmd.PersistentFlags().Int("count", 2, "bench count (time is averaged on number of executions)")
	cfg.Curve = plonkMemoryVerifyCmd.PersistentFlags().String("curve", "bn254", "curve name. must be "+fmt.Sprint(curves))
	cfg.FileType = plonkMemoryVerifyCmd.PersistentFlags().String("filetype", "csv", "Type of file to output for benchmarks")
}
