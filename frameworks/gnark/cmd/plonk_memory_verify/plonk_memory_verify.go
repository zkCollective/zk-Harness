package main

import (
	"fmt"
	"os"
	"strings"
	"bytes"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/plonk"
	"github.com/consensys/gnark/backend/witness"
	"github.com/spf13/cobra"
	"github.com/zkCollective/zk-Harness/frameworks/gnark/parser"
)

var plonkMemoryVerifyCmd = &cobra.Command{
	Use:   "plonkMemoryVerify",
	Short: "benchmarking memory consumption of Plonk verification",
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
	reconstructedProof := plonk.NewProof(parser.CurveID)
	reconstructedVK := plonk.NewVerifyingKey(parser.CurveID)
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

	// VERIFY
	err = plonk.Verify(reconstructedProof, reconstructedVK, reconstructedPublicWitness)
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
}
