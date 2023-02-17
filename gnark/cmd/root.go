package cmd

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/pkg/profile"
	"github.com/spf13/cobra"
	"github.com/tumberger/zk-compilers/gnark/circuits"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gnark-toy-bench",
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
	fFileType    *string
	fInputPath   *string

	// Variables Arithmetic
	fOperation *string
	fOrder     *string
)

var (
	curveID ecc.ID
	p       func(p *profile.Profile)
	c       circuits.BenchCircuit
)

func init() {

	fInputPath = rootCmd.PersistentFlags().String("input", "none", "input path to the dedicated input")
	rootCmd.MarkPersistentFlagRequired("input")

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// plonkCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// plonkCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func parseFlags() error {
	if *fCircuitSize <= 0 {
		return errors.New("circuit size must be >= 0")
	}
	if *fCount <= 0 {
		return errors.New("bench count must be >= 0")
	}

	switch *fAlgo {
	case "compile", "setup", "prove", "verify":
	default:
		return errors.New("invalid algo")
	}

	switch *fProfile {
	case "none":
	case "trace":
		p = profile.TraceProfile
	case "cpu":
		p = profile.CPUProfile
	case "mem":
		p = profile.MemProfile
	default:
		return errors.New("invalid profile")
	}

	curves := ecc.Implemented()
	for _, id := range curves {
		if *fCurve == strings.ToLower(id.String()) {
			curveID = id
		}
	}
	if curveID == ecc.UNKNOWN {
		return errors.New("invalid curve")
	}

	if *fFileType != "csv" {
		return errors.New("invalid file type for log")
	}

	var ok bool
	c, ok = circuits.BenchCircuits[*fCircuit]
	if !ok {
		return errors.New("unknown circuit")
	}

	return nil
}
