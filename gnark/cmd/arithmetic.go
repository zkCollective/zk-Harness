package cmd

import (
	"crypto/rand"
	"encoding/csv"
	"fmt"
	"math/big"
	"os"
	"runtime"
	"strings"
	"time"

	bn254fr "github.com/consensys/gnark-crypto/ecc/bn254/fr"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/pkg/profile"
	"github.com/spf13/cobra"
	"github.com/tumberger/zk-compilers/gnark/util"
)

// Fields:
// Goldilocks

// plonkCmd represents the plonk command
var arithmeticCmd = &cobra.Command{
	Use:   "arithmetic",
	Short: "runs benchmarks and profiles for the gnark arithmetic operations",
	Run:   benchArithmetic,
}

func benchArithmetic(cmd *cobra.Command, args []string) {

	var filename = "../benchmarks/gnark/gnark_" +
		"arithmetic" +
		"." +
		*fFileType

	if err := parseFlags(); err != nil {
		fmt.Println("error: ", err.Error())
		cmd.Help()
		os.Exit(-1)
	}

	// write to stdout
	w := csv.NewWriter(os.Stdout)
	if err := w.Write(util.BenchDataArithmetic{}.Headers()); err != nil {
		fmt.Println("error: ", err.Error())
		os.Exit(-1)
	}

	writeResults := func(took time.Duration, p big.Int) {
		// check memory usage, max ram requested from OS
		var m runtime.MemStats
		runtime.ReadMemStats(&m)

		bDataArith := util.BenchDataArithmetic{
			Framework: "gnark",
			Category:  "arithmetic",
			Field:     "native",
			Order:     p.BitLen(),
			Operation: *fOperation,
			Input:     *fInputPath,
			MaxRAM:    (m.Sys / 1024 / 1024),
			RunTime:   took.Nanoseconds(),
		}

		if err := util.WriteData("csv", bDataArith, filename); err != nil {
			panic(err)
		}
		// if err := w.Write(bDataArith.Values()); err != nil {
		// 	panic(err)
		// }
		// w.Flush()
	}

	var (
		start time.Time
		took  time.Duration
		prof  interface{ Stop() }
	)

	startProfile := func() {
		start = time.Now()
		if p != nil {
			prof = profile.Start(p, profile.ProfilePath("."), profile.NoShutdownHook)
		}
	}

	stopProfile := func() {
		took = time.Since(start)
		if p != nil {
			prof.Stop()
		}
		took /= time.Duration(*fCount)
	}

	if *fOperation == "add" {
		var x, y bn254fr.Element
		x.SetRandom()
		y.SetRandom()
		startProfile()
		for i := 0; i < *fCount; i++ {
			x.Add(&x, &y)
		}
		stopProfile()
		order := bn254fr.Modulus()
		writeResults(took, *order)
		return
	}

	if *fOperation == "sub" {
		var x, y bn254fr.Element
		x.SetRandom()
		y.SetRandom()
		startProfile()
		for i := 0; i < *fCount; i++ {
			x.Sub(&x, &y)
		}
		stopProfile()
		order := bn254fr.Modulus()
		writeResults(took, *order)
		return
	}

	if *fOperation == "mul" {
		var x, y bn254fr.Element
		x.SetRandom()
		y.SetRandom()
		startProfile()
		for i := 0; i < *fCount; i++ {
			x.Mul(&x, &y)
		}
		stopProfile()
		order := bn254fr.Modulus()
		writeResults(took, *order)
		return
	}

	if *fOperation == "div" {
		var x, y bn254fr.Element
		x.SetRandom()
		y.SetRandom()
		startProfile()
		for i := 0; i < *fCount; i++ {
			x.Div(&x, &y)
		}
		stopProfile()
		order := bn254fr.Modulus()
		writeResults(took, *order)
		return
	}

	if *fOperation == "exp" {
		var x bn254fr.Element
		x.SetRandom()
		max := big.NewInt(1000000)
		k, err := rand.Int(rand.Reader, max)
		if err != nil {
			panic(err)
		}
		startProfile()
		x.Exp(x, k)
		stopProfile()
		order := bn254fr.Modulus()
		writeResults(took, *order)
		return
	}

}

func init() {
	// Here the commands for the "arithmetic" category are defined

	_curves := ecc.Implemented()
	curves := make([]string, len(_curves))
	for i := 0; i < len(_curves); i++ {
		curves[i] = strings.ToLower(_curves[i].String())
	}

	// Possible Operations: add, sub, mul, div, exp
	fOperation = arithmeticCmd.Flags().String("operation", "add", "operation to benchmark")
	fCount = arithmeticCmd.Flags().Int("count", 2, "bench count (time is averaged on number of executions)")

	rootCmd.AddCommand(arithmeticCmd)
}
