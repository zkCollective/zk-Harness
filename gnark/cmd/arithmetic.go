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

	bls12377fr "github.com/consensys/gnark-crypto/ecc/bls12-377/fr"
	bls12381fr "github.com/consensys/gnark-crypto/ecc/bls12-381/fr"
	bls24315fr "github.com/consensys/gnark-crypto/ecc/bls24-315/fr"
	bn254fr "github.com/consensys/gnark-crypto/ecc/bn254/fr"
	bw6633fr "github.com/consensys/gnark-crypto/ecc/bw6-633/fr"
	bw6761fr "github.com/consensys/gnark-crypto/ecc/bw6-761/fr"
	"github.com/pkg/profile"

	"github.com/consensys/gnark-crypto/ecc"
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

var (
	start time.Time
	took  time.Duration
	prof  interface{ Stop() }
)

func startProfile() {
	start = time.Now()
	if p != nil {
		prof = profile.Start(p, profile.ProfilePath("."), profile.NoShutdownHook)
	}
}

func stopProfile() {
	took = time.Since(start)
	if p != nil {
		prof.Stop()
	}
	took /= time.Duration(*fCount)
}

func ExecuteOperation254(operation string) (time.Duration, *big.Int) {

	switch operation {
	case "add":
		var x, y bn254fr.Element
		x.SetRandom()
		y.SetRandom()
		startProfile()
		for i := 0; i < *fCount; i++ {
			x.Add(&x, &y)
		}
		stopProfile()
		order := bn254fr.Modulus()
		return took, order
	case "sub":
		var x, y bn254fr.Element
		x.SetRandom()
		y.SetRandom()
		startProfile()
		for i := 0; i < *fCount; i++ {
			x.Sub(&x, &y)
		}
		stopProfile()
		order := bn254fr.Modulus()
		return took, order
	case "mul":
		var x, y bn254fr.Element
		x.SetRandom()
		y.SetRandom()
		startProfile()
		for i := 0; i < *fCount; i++ {
			x.Mul(&x, &y)
		}
		stopProfile()
		order := bn254fr.Modulus()
		return took, order
	case "div":
		var x, y bn254fr.Element
		x.SetRandom()
		y.SetRandom()
		startProfile()
		for i := 0; i < *fCount; i++ {
			x.Div(&x, &y)
		}
		stopProfile()
		order := bn254fr.Modulus()
		return took, order
	case "exp":
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
		return took, order
	default:
		panic("arithmetic operation not implemented")
	}
}

func ExecuteOperationBLS12381(operation string) (time.Duration, *big.Int) {

	switch operation {
	case "add":
		var x, y bls12381fr.Element
		x.SetRandom()
		y.SetRandom()
		startProfile()
		for i := 0; i < *fCount; i++ {
			x.Add(&x, &y)
		}
		stopProfile()
		order := bls12381fr.Modulus()
		return took, order
	case "sub":
		var x, y bls12381fr.Element
		x.SetRandom()
		y.SetRandom()
		startProfile()
		for i := 0; i < *fCount; i++ {
			x.Sub(&x, &y)
		}
		stopProfile()
		order := bls12381fr.Modulus()
		return took, order
	case "mul":
		var x, y bls12381fr.Element
		x.SetRandom()
		y.SetRandom()
		startProfile()
		for i := 0; i < *fCount; i++ {
			x.Mul(&x, &y)
		}
		stopProfile()
		order := bls12381fr.Modulus()
		return took, order
	case "div":
		var x, y bls12381fr.Element
		x.SetRandom()
		y.SetRandom()
		startProfile()
		for i := 0; i < *fCount; i++ {
			x.Div(&x, &y)
		}
		stopProfile()
		order := bls12381fr.Modulus()
		return took, order
	case "exp":
		var x bls12381fr.Element
		x.SetRandom()
		max := big.NewInt(1000000)
		k, err := rand.Int(rand.Reader, max)
		if err != nil {
			panic(err)
		}
		startProfile()
		x.Exp(x, k)
		stopProfile()
		order := bls12381fr.Modulus()
		return took, order
	default:
		panic("arithmetic operation not implemented")
	}
}

func ExecuteOperationBLS12377(operation string) (time.Duration, *big.Int) {

	switch operation {
	case "add":
		var x, y bls12377fr.Element
		x.SetRandom()
		y.SetRandom()
		startProfile()
		for i := 0; i < *fCount; i++ {
			x.Add(&x, &y)
		}
		stopProfile()
		order := bls12377fr.Modulus()
		return took, order
	case "sub":
		var x, y bls12377fr.Element
		x.SetRandom()
		y.SetRandom()
		startProfile()
		for i := 0; i < *fCount; i++ {
			x.Sub(&x, &y)
		}
		stopProfile()
		order := bls12377fr.Modulus()
		return took, order
	case "mul":
		var x, y bls12377fr.Element
		x.SetRandom()
		y.SetRandom()
		startProfile()
		for i := 0; i < *fCount; i++ {
			x.Mul(&x, &y)
		}
		stopProfile()
		order := bls12377fr.Modulus()
		return took, order
	case "div":
		var x, y bls12377fr.Element
		x.SetRandom()
		y.SetRandom()
		startProfile()
		for i := 0; i < *fCount; i++ {
			x.Div(&x, &y)
		}
		stopProfile()
		order := bls12377fr.Modulus()
		return took, order
	case "exp":
		var x bls12377fr.Element
		x.SetRandom()
		max := big.NewInt(1000000)
		k, err := rand.Int(rand.Reader, max)
		if err != nil {
			panic(err)
		}
		startProfile()
		x.Exp(x, k)
		stopProfile()
		order := bls12377fr.Modulus()
		return took, order
	default:
		panic("arithmetic operation not implemented")
	}
}

func ExecuteOperationBLS24315(operation string) (time.Duration, *big.Int) {

	switch operation {
	case "add":
		var x, y bls24315fr.Element
		x.SetRandom()
		y.SetRandom()
		startProfile()
		for i := 0; i < *fCount; i++ {
			x.Add(&x, &y)
		}
		stopProfile()
		order := bls24315fr.Modulus()
		return took, order
	case "sub":
		var x, y bls24315fr.Element
		x.SetRandom()
		y.SetRandom()
		startProfile()
		for i := 0; i < *fCount; i++ {
			x.Sub(&x, &y)
		}
		stopProfile()
		order := bls24315fr.Modulus()
		return took, order
	case "mul":
		var x, y bls24315fr.Element
		x.SetRandom()
		y.SetRandom()
		startProfile()
		for i := 0; i < *fCount; i++ {
			x.Mul(&x, &y)
		}
		stopProfile()
		order := bls24315fr.Modulus()
		return took, order
	case "div":
		var x, y bls24315fr.Element
		x.SetRandom()
		y.SetRandom()
		startProfile()
		for i := 0; i < *fCount; i++ {
			x.Div(&x, &y)
		}
		stopProfile()
		order := bls24315fr.Modulus()
		return took, order
	case "exp":
		var x bls24315fr.Element
		x.SetRandom()
		max := big.NewInt(1000000)
		k, err := rand.Int(rand.Reader, max)
		if err != nil {
			panic(err)
		}
		startProfile()
		x.Exp(x, k)
		stopProfile()
		order := bls24315fr.Modulus()
		return took, order
	default:
		panic("arithmetic operation not implemented")
	}
}

func ExecuteOperationBW6633(operation string) (time.Duration, *big.Int) {

	switch operation {
	case "add":
		var x, y bw6633fr.Element
		x.SetRandom()
		y.SetRandom()
		startProfile()
		for i := 0; i < *fCount; i++ {
			x.Add(&x, &y)
		}
		stopProfile()
		order := bw6633fr.Modulus()
		return took, order
	case "sub":
		var x, y bw6633fr.Element
		x.SetRandom()
		y.SetRandom()
		startProfile()
		for i := 0; i < *fCount; i++ {
			x.Sub(&x, &y)
		}
		stopProfile()
		order := bw6633fr.Modulus()
		return took, order
	case "mul":
		var x, y bw6633fr.Element
		x.SetRandom()
		y.SetRandom()
		startProfile()
		for i := 0; i < *fCount; i++ {
			x.Mul(&x, &y)
		}
		stopProfile()
		order := bw6633fr.Modulus()
		return took, order
	case "div":
		var x, y bw6633fr.Element
		x.SetRandom()
		y.SetRandom()
		startProfile()
		for i := 0; i < *fCount; i++ {
			x.Div(&x, &y)
		}
		stopProfile()
		order := bw6633fr.Modulus()
		return took, order
	case "exp":
		var x bw6633fr.Element
		x.SetRandom()
		max := big.NewInt(1000000)
		k, err := rand.Int(rand.Reader, max)
		if err != nil {
			panic(err)
		}
		startProfile()
		x.Exp(x, k)
		stopProfile()
		order := bw6633fr.Modulus()
		return took, order
	default:
		panic("arithmetic operation not implemented")
	}
}

func ExecuteOperationBW6761(operation string) (time.Duration, *big.Int) {

	switch operation {
	case "add":
		var x, y bw6761fr.Element
		x.SetRandom()
		y.SetRandom()
		startProfile()
		for i := 0; i < *fCount; i++ {
			x.Add(&x, &y)
		}
		stopProfile()
		order := bw6761fr.Modulus()
		return took, order
	case "sub":
		var x, y bw6761fr.Element
		x.SetRandom()
		y.SetRandom()
		startProfile()
		for i := 0; i < *fCount; i++ {
			x.Sub(&x, &y)
		}
		stopProfile()
		order := bw6761fr.Modulus()
		return took, order
	case "mul":
		var x, y bw6761fr.Element
		x.SetRandom()
		y.SetRandom()
		startProfile()
		for i := 0; i < *fCount; i++ {
			x.Mul(&x, &y)
		}
		stopProfile()
		order := bw6761fr.Modulus()
		return took, order
	case "div":
		var x, y bw6761fr.Element
		x.SetRandom()
		y.SetRandom()
		startProfile()
		for i := 0; i < *fCount; i++ {
			x.Div(&x, &y)
		}
		stopProfile()
		order := bw6761fr.Modulus()
		return took, order
	case "exp":
		var x bw6761fr.Element
		x.SetRandom()
		max := big.NewInt(1000000)
		k, err := rand.Int(rand.Reader, max)
		if err != nil {
			panic(err)
		}
		startProfile()
		x.Exp(x, k)
		stopProfile()
		order := bw6761fr.Modulus()
		return took, order
	default:
		panic("arithmetic operation not implemented")
	}
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
	}

	switch *fOrder {
	case "bn254":
		took, order := ExecuteOperation254(*fOperation)
		writeResults(took, *order)
	case "bls12381":
		took, order := ExecuteOperationBLS12381(*fOperation)
		writeResults(took, *order)
	case "bls12377":
		took, order := ExecuteOperationBLS12377(*fOperation)
		writeResults(took, *order)
	case "bls24315":
		took, order := ExecuteOperationBLS24315(*fOperation)
		writeResults(took, *order)
	case "bw6633":
		took, order := ExecuteOperationBW6633(*fOperation)
		writeResults(took, *order)
	case "bw6761":
		took, order := ExecuteOperationBW6761(*fOperation)
		writeResults(took, *order)
	default:
		panic("field order not implemented")
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
	fOrder = arithmeticCmd.Flags().String("field", "bn254", "operation to benchmark")
	fCount = arithmeticCmd.Flags().Int("count", 1, "bench count (time is averaged on number of executions)")

	rootCmd.AddCommand(arithmeticCmd)
}
