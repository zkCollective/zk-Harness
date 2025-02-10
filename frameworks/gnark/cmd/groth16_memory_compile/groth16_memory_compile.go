package main

import (
	"fmt"
	"os"
	"strings"
	"bytes"
	"log"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
	"github.com/spf13/cobra"
	"github.com/zkCollective/zk-Harness/frameworks/gnark/circuits"
	"github.com/zkCollective/zk-Harness/frameworks/gnark/parser"
)

var groth16MemoryCompileCmd = &cobra.Command{
	Use:   "groth16MemoryCompile",
	Short: "runs benchmarks and profiles using Groth16 proof system",
	Run:   runGroth16MemoryCompile,
}

var cfg = parser.NewConfig()

func runGroth16MemoryCompile(cmd *cobra.Command, args []string) {

	if err := parser.ParseFlagsMemory(cfg); err != nil {
		fmt.Println("error: ", err.Error())
		cmd.Help()
		os.Exit(-1)
	}

	ccs, err := frontend.Compile(parser.CurveID.ScalarField(),
		r1cs.NewBuilder,
		parser.C.Circuit(
			*cfg.CircuitSize,
			*cfg.Circuit,
			circuits.WithInputCircuit(*cfg.InputPath)),
		frontend.WithCapacity(*cfg.CircuitSize))
	parser.AssertNoError(err)

	// SERIALIZE - WRITE THE CCS
	var buf bytes.Buffer
	_, _ = ccs.WriteTo(&buf)

	err = os.WriteFile("tmp/ccs.dat", buf.Bytes(), 0644)
	if err != nil {
		log.Fatal(err)
	}
	return
}

func Execute() {
	if err := groth16MemoryCompileCmd.Execute(); err != nil {
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

	cfg.InputPath = groth16MemoryCompileCmd.PersistentFlags().String("input", "none", "input path to the dedicated input")
	groth16MemoryCompileCmd.MarkPersistentFlagRequired("input")
	cfg.Circuit = groth16MemoryCompileCmd.PersistentFlags().String("circuit", "expo", "name of the circuit to use")
	cfg.CircuitSize = groth16MemoryCompileCmd.PersistentFlags().Int("size", 10000, "size of the circuit, parameter to circuit constructor")
	cfg.Count = groth16MemoryCompileCmd.PersistentFlags().Int("count", 2, "bench count (time is averaged on number of executions)")
	cfg.Curve = groth16MemoryCompileCmd.PersistentFlags().String("curve", "bn254", "curve name. must be "+fmt.Sprint(curves))
}
