package main

import (
	"fmt"
	"os"
	"strings"
	"bytes"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/backend/witness"
	"github.com/spf13/cobra"
	"github.com/zkCollective/zk-Harness/frameworks/gnark/parser"
)

// groth16Cmd represents the groth16 command
var groth16MemoryVerifyCmd = &cobra.Command{
	Use:   "groth16MemoryVerify",
	Short: "runs benchmarks and profiles using Groth16 proof system",
	Run:   runGroth16MemoryVerify,
}

var cfg = parser.NewConfig()

func runGroth16MemoryVerify(cmd *cobra.Command, args []string) {

	if err := parser.ParseFlagsMemory(cfg); err != nil {
		fmt.Println("error: ", err.Error())
		cmd.Help()
		os.Exit(-1)
	}

	// Initialize variables
	reconstructedProof := groth16.NewProof(parser.CurveID)
	reconstructedVK := groth16.NewVerifyingKey(parser.CurveID)
	newWitness, err := witness.New(parser.CurveID.ScalarField())
	reconstructedPublicWitness, _ := newWitness.Public()

	// Read vk
	_vk, err := os.ReadFile("tmp/vk.dat")
	if err != nil {
		panic("Failed to open file: " + err.Error())
	}

	_buf := *bytes.NewBuffer(_vk)
	_, err = reconstructedVK.ReadFrom(&_buf)
	if err != nil {
		panic("Failed to read verifier key: " + err.Error())
	}

	// Read Public Witness
	_pubWit, err := os.ReadFile("tmp/publicWitness.dat")
	if err != nil {
		panic("Failed to read from file: " + err.Error())
	}

	// Binary marshalling
	reconstructedPublicWitness.UnmarshalBinary(_pubWit)

	// Read proof
	_proof, err := os.ReadFile("tmp/proof.dat")
	if err != nil {
		panic("Failed to read from file: " + err.Error())
	}
	
	_, err = reconstructedProof.ReadFrom(bytes.NewReader(_proof))

	// Proof Verification
	err = groth16.Verify(reconstructedProof, reconstructedVK, reconstructedPublicWitness)
	if err != nil {
		panic("Failed Verification!")
	}
	return
}

func Execute() {
	if err := groth16MemoryVerifyCmd.Execute(); err != nil {
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

	cfg.InputPath = groth16MemoryVerifyCmd.PersistentFlags().String("input", "none", "input path to the dedicated input")
	groth16MemoryVerifyCmd.MarkPersistentFlagRequired("input")
	cfg.Circuit = groth16MemoryVerifyCmd.PersistentFlags().String("circuit", "expo", "name of the circuit to use")
	cfg.CircuitSize = groth16MemoryVerifyCmd.PersistentFlags().Int("size", 10000, "size of the circuit, parameter to circuit constructor")
	cfg.Count = groth16MemoryVerifyCmd.PersistentFlags().Int("count", 2, "bench count (time is averaged on number of executions)")
	cfg.Curve = groth16MemoryVerifyCmd.PersistentFlags().String("curve", "bn254", "curve name. must be "+fmt.Sprint(curves))
}
