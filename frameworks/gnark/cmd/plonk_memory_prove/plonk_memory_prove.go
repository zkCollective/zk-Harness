package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/plonk"
	"github.com/spf13/cobra"
	"github.com/zkCollective/zk-Harness/gnark/circuits"
	"github.com/zkCollective/zk-Harness/gnark/parser"
)

var plonkMemoryProveCmd = &cobra.Command{
	Use:   "plonkMemoryProve",
	Short: "runs benchmarks and profiles using Groth16 proof system",
	Run:   runPlonkMemoryProve,
}

var cfg = parser.NewConfig()

func runPlonkMemoryProve(cmd *cobra.Command, args []string) {

	if err := parser.ParseFlagsMemory(cfg); err != nil {
		fmt.Println("error: ", err.Error())
		cmd.Help()
		os.Exit(-1)
	}

	// Initialize variables
	ccs := plonk.NewCS(parser.CurveID)
	pk := plonk.NewProvingKey(parser.CurveID)

	// Read CCS
	f, err := os.Open("tmp/ccs.dat")
	if err != nil {
		panic("Failed to open file: " + err.Error())
	}

	_, err = ccs.ReadFrom(f)
	if err != nil {
		panic("Failed to read from file: " + err.Error())
	}

	f.Close()

	// Read PK
	f, err = os.Open("tmp/pk.dat")
	if err != nil {
		panic("Failed to open file: " + err.Error())
	}

	_, err = pk.ReadFrom(f)
	if err != nil {
		panic("Failed to read from file: " + err.Error())
	}

	f.Close()

	// Witness creation is included in Prover Memory benchmarks
	witness := parser.C.Witness(*cfg.CircuitSize, parser.CurveID, *cfg.Circuit, circuits.WithInputWitness(*cfg.InputPath))

	proof, err := plonk.Prove(ccs, pk, witness)
	if err != nil {
		panic("Error when computing proof. Ensure that Constraint System and pk/vk are generated for the same parameters.")
	}

	// Extract the public part only
	publicWitness, _ := witness.Public()

	// Serialize the publicWitness
	f, err = os.Create("tmp/publicWitness.dat")
	if err != nil {
		panic("Failed to create file: " + err.Error())
	}
	defer f.Close()

	_, err = publicWitness.WriteTo(f)
	if err != nil {
		panic("Failed to write to file: " + err.Error())
	}

	// Serialize the Proof
	f, err = os.Create("tmp/proof.dat")
	if err != nil {
		panic("Failed to create file: " + err.Error())
	}
	defer f.Close()

	_, err = proof.WriteTo(f)
	if err != nil {
		panic("Failed to write to file: " + err.Error())
	}

	return
}

func Execute() {
	if err := plonkMemoryProveCmd.Execute(); err != nil {
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

	cfg.InputPath = plonkMemoryProveCmd.PersistentFlags().String("input", "none", "input path to the dedicated input")
	plonkMemoryProveCmd.MarkPersistentFlagRequired("input")
	cfg.Circuit = plonkMemoryProveCmd.PersistentFlags().String("circuit", "expo", "name of the circuit to use")
	cfg.CircuitSize = plonkMemoryProveCmd.PersistentFlags().Int("size", 10000, "size of the circuit, parameter to circuit constructor")
	cfg.Count = plonkMemoryProveCmd.PersistentFlags().Int("count", 2, "bench count (time is averaged on number of executions)")
	cfg.Curve = plonkMemoryProveCmd.PersistentFlags().String("curve", "bn254", "curve name. must be "+fmt.Sprint(curves))
	cfg.FileType = plonkMemoryProveCmd.PersistentFlags().String("filetype", "csv", "Type of file to output for benchmarks")
}
