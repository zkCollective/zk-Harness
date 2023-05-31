package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/spf13/cobra"
	"github.com/zkCollective/zk-Harness/gnark/parser"
)

var groth16MemoryWitnessCmd = &cobra.Command{
	Use:   "groth16MemorySetup",
	Short: "runs benchmarks and profiles using Groth16 proof system",
	Run:   runGroth16MemoryWitness,
}

var cfg = parser.NewConfig()

func runGroth16MemoryWitness(cmd *cobra.Command, args []string) {

	if err := parser.ParseFlagsMemory(cfg); err != nil {
		fmt.Println("error: ", err.Error())
		cmd.Help()
		os.Exit(-1)
	}

	witness := parser.C.Witness(*cfg.CircuitSize, parser.CurveID, *cfg.Circuit, *cfg.InputPath)

	f, err := os.Create("tmp/witness.dat")
	if err != nil {
		panic("Failed to create file: " + err.Error())
	}
	defer f.Close()

	_, err = witness.WriteTo(f)
	if err != nil {
		panic("Failed to write to file: " + err.Error())
	}
	return
}

func Execute() {
	if err := groth16MemoryWitnessCmd.Execute(); err != nil {
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

	cfg.InputPath = groth16MemoryWitnessCmd.PersistentFlags().String("input", "none", "input path to the dedicated input")
	groth16MemoryWitnessCmd.MarkPersistentFlagRequired("input")
	cfg.Circuit = groth16MemoryWitnessCmd.PersistentFlags().String("circuit", "expo", "name of the circuit to use")
	cfg.CircuitSize = groth16MemoryWitnessCmd.PersistentFlags().Int("size", 10000, "size of the circuit, parameter to circuit constructor")
	cfg.Count = groth16MemoryWitnessCmd.PersistentFlags().Int("count", 2, "bench count (time is averaged on number of executions)")
	cfg.Curve = groth16MemoryWitnessCmd.PersistentFlags().String("curve", "bn254", "curve name. must be "+fmt.Sprint(curves))
	cfg.FileType = groth16MemoryWitnessCmd.PersistentFlags().String("filetype", "csv", "Type of file to output for benchmarks")
}
