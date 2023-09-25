package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/spf13/cobra"
	"github.com/zkCollective/zk-Harness/frameworks/gnark/parser"
)

var cfg = parser.NewConfig()

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gnark-harness",
	Short: "runs benchmarks and profiles using gnark",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var (
	// Variables Circuit
	fCircuit     *string
	fCircuitSize *int
	fAlgo        *string
	fProfile     *string
	fCount       *int
	fCurve       *string
	fInputPath   *string

	// Variables Arithmetic / Curve
	fOperation *string
	fField     *string
	fGroup     *string

	// Variables Recursion
	fOuterBackend *string

	// Machine Variable
	fOutputPath *string
)

func init() {

	cobra.OnInitialize()

	_curves := ecc.Implemented()
	curves := make([]string, len(_curves))
	for i := 0; i < len(_curves); i++ {
		curves[i] = strings.ToLower(_curves[i].String())
	}

	cfg.InputPath = rootCmd.PersistentFlags().String("input", "none", "input path to the dedicated input")
	rootCmd.MarkPersistentFlagRequired("input")
	cfg.Circuit = rootCmd.PersistentFlags().String("circuit", "expo", "name of the circuit to use")
	cfg.CircuitSize = rootCmd.PersistentFlags().Int("size", 10000, "size of the circuit, parameter to circuit constructor")
	cfg.Count = rootCmd.PersistentFlags().Int("count", 2, "bench count (time is averaged on number of executions)")
	cfg.Curve = rootCmd.PersistentFlags().String("curve", "bn254", "curve name. must be "+fmt.Sprint(curves))
	cfg.Algo = rootCmd.PersistentFlags().String("algo", "prove", "name of the algorithm to benchmark. must be compile, setup, prove or verify")
	cfg.Operation = rootCmd.PersistentFlags().String("operation", "None", "operation to benchmark")
	cfg.Profile = rootCmd.PersistentFlags().String("profile", "none", "type of profile. must be none, trace, cpu or mem")

	cfg.OuterBackend = rootCmd.PersistentFlags().String("outerBackend", "groth16", "Backend for the outer circuit")

	cfg.OutputPath = rootCmd.PersistentFlags().String("outputPath", "None", "The output path for the log file")

	rootCmd.AddCommand(groth16Cmd)
	rootCmd.AddCommand(plonkCmd)
}
