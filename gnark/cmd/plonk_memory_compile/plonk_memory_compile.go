package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/scs"
	"github.com/spf13/cobra"
	"github.com/zkCollective/zk-Harness/gnark/parser"
)

var plonkMemoryCompileCmd = &cobra.Command{
	Use:   "plonkMemoryCompile",
	Short: "runs memory benchmarks using Plonk proof system",
	Run:   runPlonkMemoryCompile,
}

var cfg = parser.NewConfig()

func runPlonkMemoryCompile(cmd *cobra.Command, args []string) {

	if err := parser.ParseFlagsMemory(cfg); err != nil {
		fmt.Println("error: ", err.Error())
		cmd.Help()
		os.Exit(-1)
	}

	ccs, err := frontend.Compile(parser.CurveID.ScalarField(),
		scs.NewBuilder,
		parser.C.Circuit(
			*cfg.CircuitSize,
			*cfg.Circuit,
			*cfg.InputPath),
		frontend.WithCapacity(*cfg.CircuitSize))
	parser.AssertNoError(err)

	f, err := os.Create("tmp/ccs.dat")
	if err != nil {
		panic("Failed to create file: " + err.Error())
	}
	defer f.Close()

	_, err = ccs.WriteTo(f)
	if err != nil {
		panic("Failed to write to file: " + err.Error())
	}
	return
}

func Execute() {
	if err := plonkMemoryCompileCmd.Execute(); err != nil {
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

	cfg.InputPath = plonkMemoryCompileCmd.PersistentFlags().String("input", "none", "input path to the dedicated input")
	plonkMemoryCompileCmd.MarkPersistentFlagRequired("input")
	cfg.Circuit = plonkMemoryCompileCmd.PersistentFlags().String("circuit", "expo", "name of the circuit to use")
	cfg.CircuitSize = plonkMemoryCompileCmd.PersistentFlags().Int("size", 10000, "size of the circuit, parameter to circuit constructor")
	cfg.Count = plonkMemoryCompileCmd.PersistentFlags().Int("count", 2, "bench count (time is averaged on number of executions)")
	cfg.Curve = plonkMemoryCompileCmd.PersistentFlags().String("curve", "bn254", "curve name. must be "+fmt.Sprint(curves))
	cfg.FileType = plonkMemoryCompileCmd.PersistentFlags().String("filetype", "csv", "Type of file to output for benchmarks")
}
