package main

import (
	"fmt"
	"os"
	"strings"
	"bytes"
	"log"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/spf13/cobra"
	"github.com/zkCollective/zk-Harness/frameworks/gnark/circuits"
	"github.com/zkCollective/zk-Harness/frameworks/gnark/parser"
)

var groth16MemoryProveCmd = &cobra.Command{
	Use:   "groth16MemoryProve",
	Short: "runs benchmarks and profiles using Groth16 proof system",
	Run:   runGroth16MemoryProve,
}

var cfg = parser.NewConfig()

func runGroth16MemoryProve(cmd *cobra.Command, args []string) {

	if err := parser.ParseFlagsMemory(cfg); err != nil {
		fmt.Println("error: ", err.Error())
		cmd.Help()
		os.Exit(-1)
	}

	// Initialize variables
	reconstructedCCS := groth16.NewCS(parser.CurveID)
	reconstructedPK := groth16.NewProvingKey(parser.CurveID)

	// Read CCS
	_ccs, err := os.ReadFile("tmp/ccs.dat")
	if err != nil {
		log.Fatal(err)
	}
	_buf := *bytes.NewBuffer(_ccs)
	_, _ = reconstructedCCS.ReadFrom(&_buf)

	// Read PK
	_pk, err := os.ReadFile("tmp/pk.dat")
	if err != nil {
		panic("Failed to read from file: " + err.Error())
	}
	_buf = *bytes.NewBuffer(_pk)
	_, err = reconstructedPK.ReadFrom(&_buf)
	if err != nil {
		panic("Failed to read prover key: " + err.Error())
	}

	// Witness creation is included in Prover Memory benchmarks
	witness := parser.C.Witness(*cfg.CircuitSize,
		parser.CurveID,
		*cfg.Circuit,
		circuits.WithInputWitness(*cfg.InputPath))

	proof, err := groth16.Prove(reconstructedCCS, reconstructedPK, witness)
	if err != nil {
		panic("Error when computing proof. Ensure that Constraint System and pk/vk are generated for the same parameters.")
	}

	// Extract the public part only
	publicWitness, _ := witness.Public()

	// Serialize the publicWitness
	data, err := publicWitness.MarshalBinary()
	if err != nil {
		panic("Failed to marshal binary: " + err.Error())
	}
	err = os.WriteFile("tmp/publicWitness.dat", data, 0644)
	if err != nil {
		panic("Failed to write to file: " + err.Error())
	}

	// Serialize the Proof
	var bufProof bytes.Buffer
	_, _ = proof.WriteTo(&bufProof)
	err = os.WriteFile("tmp/proof.dat", bufProof.Bytes(), 0644)
	if err != nil {
		log.Fatal(err)
	}

	return
}

func Execute() {
	if err := groth16MemoryProveCmd.Execute(); err != nil {
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

	cfg.InputPath = groth16MemoryProveCmd.PersistentFlags().String("input", "none", "input path to the dedicated input")
	groth16MemoryProveCmd.MarkPersistentFlagRequired("input")
	cfg.Circuit = groth16MemoryProveCmd.PersistentFlags().String("circuit", "expo", "name of the circuit to use")
	cfg.CircuitSize = groth16MemoryProveCmd.PersistentFlags().Int("size", 10000, "size of the circuit, parameter to circuit constructor")
	cfg.Count = groth16MemoryProveCmd.PersistentFlags().Int("count", 2, "bench count (time is averaged on number of executions)")
	cfg.Curve = groth16MemoryProveCmd.PersistentFlags().String("curve", "bn254", "curve name. must be "+fmt.Sprint(curves))
}
