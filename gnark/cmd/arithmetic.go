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

	bls12377fp "github.com/consensys/gnark-crypto/ecc/bls12-377/fp"
	bls12377fr "github.com/consensys/gnark-crypto/ecc/bls12-377/fr"
	bls12381fp "github.com/consensys/gnark-crypto/ecc/bls12-381/fp"
	bls12381fr "github.com/consensys/gnark-crypto/ecc/bls12-381/fr"
	bls24315fp "github.com/consensys/gnark-crypto/ecc/bls24-315/fp"
	bls24315fr "github.com/consensys/gnark-crypto/ecc/bls24-315/fr"
	bn254fp "github.com/consensys/gnark-crypto/ecc/bn254/fp"
	bn254fr "github.com/consensys/gnark-crypto/ecc/bn254/fr"
	bw6633fp "github.com/consensys/gnark-crypto/ecc/bw6-633/fp"
	bw6633fr "github.com/consensys/gnark-crypto/ecc/bw6-633/fr"
	bw6761fp "github.com/consensys/gnark-crypto/ecc/bw6-761/fp"
	bw6761fr "github.com/consensys/gnark-crypto/ecc/bw6-761/fr"
	"github.com/consensys/gnark/logger"
	"github.com/pkg/profile"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/spf13/cobra"
	"github.com/tumberger/zk-compilers/gnark/util"
)

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

// TODO
// Currently solved very ugly, need to refactor
// Similar to below, requires handling reflection
// func newElement(curveID ecc.ID, field string) (reflect.Type, reflect.Type, error) {
//     if field == “scalar” {
//         switch curveID {
//         case ecc.BN254:
//             return reflect.TypeOf(bn254fr.Element{}), reflect.TypeOf(bn254fr.Element{}), nil
//         case ecc.BLS12_377:
//             return reflect.TypeOf(bls12377fr.Element{}), reflect.TypeOf(bls12377fr.Element{}), nil
//         case ecc.BLS12_381:
//             return reflect.TypeOf(bls12381fr.Element{}), reflect.TypeOf(bls12381fr.Element{}), nil
//         case ecc.BW6_761:
//             return reflect.TypeOf(bw6761fr.Element{}), reflect.TypeOf(bw6761fr.Element{}), nil
//         case ecc.BLS24_315:
//             return reflect.TypeOf(bls24315fr.Element{}), reflect.TypeOf(bls24315fr.Element{}), nil
//         case ecc.BW6_633:
//             return reflect.TypeOf(bw6633fr.Element{}), reflect.TypeOf(bw6633fr.Element{}), nil
//         default:
//             return nil, nil, errors.New(“unsupported curve)
//         }
//     } else {
//         return nil, nil, errors.New(“unsupported field”)
//     }
// }

func ExecuteOperation254(operation string, x float64, y float64) (time.Duration, *big.Int) {
	if *fField == "scalar" {
		switch operation {
		case "add":
			var x, y bn254fr.Element
			x.SetInterface(x)
			y.SetInterface(y)
			startProfile()
			for i := 0; i < *fCount; i++ {
				x.Add(&x, &y)
			}
			stopProfile()
			order := bn254fr.Modulus()
			return took, order
		case "sub":
			var x, y bn254fr.Element
			x.SetInterface(x)
			y.SetInterface(y)
			startProfile()
			for i := 0; i < *fCount; i++ {
				x.Sub(&x, &y)
			}
			stopProfile()
			order := bn254fr.Modulus()
			return took, order
		case "mul":
			var x, y bn254fr.Element
			x.SetInterface(x)
			y.SetInterface(y)
			startProfile()
			for i := 0; i < *fCount; i++ {
				x.Mul(&x, &y)
			}
			stopProfile()
			order := bn254fr.Modulus()
			return took, order
		case "div":
			var x, y bn254fr.Element
			x.SetInterface(x)
			y.SetInterface(y)
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
	} else if *fField == "base" {
		switch operation {
		case "add":
			var x, y bn254fp.Element
			x.SetInterface(x)
			y.SetInterface(y)
			startProfile()
			for i := 0; i < *fCount; i++ {
				x.Add(&x, &y)
			}
			stopProfile()
			order := bn254fp.Modulus()
			return took, order
		case "sub":
			var x, y bn254fp.Element
			x.SetInterface(x)
			y.SetInterface(y)
			startProfile()
			for i := 0; i < *fCount; i++ {
				x.Sub(&x, &y)
			}
			stopProfile()
			order := bn254fp.Modulus()
			return took, order
		case "mul":
			var x, y bn254fp.Element
			x.SetInterface(x)
			y.SetInterface(y)
			startProfile()
			for i := 0; i < *fCount; i++ {
				x.Mul(&x, &y)
			}
			stopProfile()
			order := bn254fp.Modulus()
			return took, order
		case "div":
			var x, y bn254fp.Element
			x.SetInterface(x)
			y.SetInterface(y)
			startProfile()
			for i := 0; i < *fCount; i++ {
				x.Div(&x, &y)
			}
			stopProfile()
			order := bn254fp.Modulus()
			return took, order
		case "exp":
			var x bn254fp.Element
			x.SetRandom()
			max := big.NewInt(1000000)
			k, err := rand.Int(rand.Reader, max)
			if err != nil {
				panic(err)
			}
			startProfile()
			x.Exp(x, k)
			stopProfile()
			order := bn254fp.Modulus()
			return took, order
		default:
			panic("arithmetic operation not implemented")
		}
	} else {
		panic("field not valid")
	}
}

func ExecuteOperationBLS12381(operation string, x float64, y float64) (time.Duration, *big.Int) {

	if *fField == "scalar" {
		switch operation {
		case "add":
			var x, y bls12381fr.Element
			x.SetInterface(x)
			y.SetInterface(y)
			startProfile()
			for i := 0; i < *fCount; i++ {
				x.Add(&x, &y)
			}
			stopProfile()
			order := bls12381fr.Modulus()
			return took, order
		case "sub":
			var x, y bls12381fr.Element
			x.SetInterface(x)
			y.SetInterface(y)
			startProfile()
			for i := 0; i < *fCount; i++ {
				x.Sub(&x, &y)
			}
			stopProfile()
			order := bls12381fr.Modulus()
			return took, order
		case "mul":
			var x, y bls12381fr.Element
			x.SetInterface(x)
			y.SetInterface(y)
			startProfile()
			for i := 0; i < *fCount; i++ {
				x.Mul(&x, &y)
			}
			stopProfile()
			order := bls12381fr.Modulus()
			return took, order
		case "div":
			var x, y bls12381fr.Element
			x.SetInterface(x)
			y.SetInterface(y)
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
	} else if *fField == "base" {
		switch operation {
		case "add":
			var x, y bls12381fp.Element
			x.SetInterface(x)
			y.SetInterface(y)
			startProfile()
			for i := 0; i < *fCount; i++ {
				x.Add(&x, &y)
			}
			stopProfile()
			order := bls12381fp.Modulus()
			return took, order
		case "sub":
			var x, y bls12381fp.Element
			x.SetInterface(x)
			y.SetInterface(y)
			startProfile()
			for i := 0; i < *fCount; i++ {
				x.Sub(&x, &y)
			}
			stopProfile()
			order := bls12381fp.Modulus()
			return took, order
		case "mul":
			var x, y bls12381fp.Element
			x.SetInterface(x)
			y.SetInterface(y)
			startProfile()
			for i := 0; i < *fCount; i++ {
				x.Mul(&x, &y)
			}
			stopProfile()
			order := bls12381fp.Modulus()
			return took, order
		case "div":
			var x, y bls12381fp.Element
			x.SetInterface(x)
			y.SetInterface(y)
			startProfile()
			for i := 0; i < *fCount; i++ {
				x.Div(&x, &y)
			}
			stopProfile()
			order := bls12381fp.Modulus()
			return took, order
		case "exp":
			var x bls12381fp.Element
			x.SetRandom()
			max := big.NewInt(1000000)
			k, err := rand.Int(rand.Reader, max)
			if err != nil {
				panic(err)
			}
			startProfile()
			x.Exp(x, k)
			stopProfile()
			order := bls12381fp.Modulus()
			return took, order
		default:
			panic("arithmetic operation not implemented")
		}
	} else {
		panic("field not valid")
	}
}

func ExecuteOperationBLS12377(operation string, x float64, y float64) (time.Duration, *big.Int) {

	if *fField == "scalar" {
		switch operation {
		case "add":
			var x, y bls12377fr.Element
			x.SetInterface(x)
			y.SetInterface(y)
			startProfile()
			for i := 0; i < *fCount; i++ {
				x.Add(&x, &y)
			}
			stopProfile()
			order := bls12377fr.Modulus()
			return took, order
		case "sub":
			var x, y bls12377fr.Element
			x.SetInterface(x)
			y.SetInterface(y)
			startProfile()
			for i := 0; i < *fCount; i++ {
				x.Sub(&x, &y)
			}
			stopProfile()
			order := bls12377fr.Modulus()
			return took, order
		case "mul":
			var x, y bls12377fr.Element
			x.SetInterface(x)
			y.SetInterface(y)
			startProfile()
			for i := 0; i < *fCount; i++ {
				x.Mul(&x, &y)
			}
			stopProfile()
			order := bls12377fr.Modulus()
			return took, order
		case "div":
			var x, y bls12377fr.Element
			x.SetInterface(x)
			y.SetInterface(y)
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
	} else if *fField == "base" {
		switch operation {
		case "add":
			var x, y bls12377fp.Element
			x.SetInterface(x)
			y.SetInterface(y)
			startProfile()
			for i := 0; i < *fCount; i++ {
				x.Add(&x, &y)
			}
			stopProfile()
			order := bls12377fp.Modulus()
			return took, order
		case "sub":
			var x, y bls12377fp.Element
			x.SetInterface(x)
			y.SetInterface(y)
			startProfile()
			for i := 0; i < *fCount; i++ {
				x.Sub(&x, &y)
			}
			stopProfile()
			order := bls12377fp.Modulus()
			return took, order
		case "mul":
			var x, y bls12377fp.Element
			x.SetInterface(x)
			y.SetInterface(y)
			startProfile()
			for i := 0; i < *fCount; i++ {
				x.Mul(&x, &y)
			}
			stopProfile()
			order := bls12377fp.Modulus()
			return took, order
		case "div":
			var x, y bls12377fp.Element
			x.SetInterface(x)
			y.SetInterface(y)
			startProfile()
			for i := 0; i < *fCount; i++ {
				x.Div(&x, &y)
			}
			stopProfile()
			order := bls12377fp.Modulus()
			return took, order
		case "exp":
			var x bls12377fp.Element
			x.SetRandom()
			max := big.NewInt(1000000)
			k, err := rand.Int(rand.Reader, max)
			if err != nil {
				panic(err)
			}
			startProfile()
			x.Exp(x, k)
			stopProfile()
			order := bls12377fp.Modulus()
			return took, order
		default:
			panic("arithmetic operation not implemented")
		}
	} else {
		panic("field not valid")
	}
}

func ExecuteOperationBLS24315(operation string, x float64, y float64) (time.Duration, *big.Int) {

	if *fField == "scalar" {
		switch operation {
		case "add":
			var x, y bls24315fr.Element
			x.SetInterface(x)
			y.SetInterface(y)
			startProfile()
			for i := 0; i < *fCount; i++ {
				x.Add(&x, &y)
			}
			stopProfile()
			order := bls24315fr.Modulus()
			return took, order
		case "sub":
			var x, y bls24315fr.Element
			x.SetInterface(x)
			y.SetInterface(y)
			startProfile()
			for i := 0; i < *fCount; i++ {
				x.Sub(&x, &y)
			}
			stopProfile()
			order := bls24315fr.Modulus()
			return took, order
		case "mul":
			var x, y bls24315fr.Element
			x.SetInterface(x)
			y.SetInterface(y)
			startProfile()
			for i := 0; i < *fCount; i++ {
				x.Mul(&x, &y)
			}
			stopProfile()
			order := bls24315fr.Modulus()
			return took, order
		case "div":
			var x, y bls24315fr.Element
			x.SetInterface(x)
			y.SetInterface(y)
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
	} else if *fField == "base" {
		switch operation {
		case "add":
			var x, y bls24315fp.Element
			x.SetInterface(x)
			y.SetInterface(y)
			startProfile()
			for i := 0; i < *fCount; i++ {
				x.Add(&x, &y)
			}
			stopProfile()
			order := bls24315fp.Modulus()
			return took, order
		case "sub":
			var x, y bls24315fp.Element
			x.SetInterface(x)
			y.SetInterface(y)
			startProfile()
			for i := 0; i < *fCount; i++ {
				x.Sub(&x, &y)
			}
			stopProfile()
			order := bls24315fp.Modulus()
			return took, order
		case "mul":
			var x, y bls24315fp.Element
			x.SetInterface(x)
			y.SetInterface(y)
			startProfile()
			for i := 0; i < *fCount; i++ {
				x.Mul(&x, &y)
			}
			stopProfile()
			order := bls24315fp.Modulus()
			return took, order
		case "div":
			var x, y bls24315fp.Element
			x.SetInterface(x)
			y.SetInterface(y)
			startProfile()
			for i := 0; i < *fCount; i++ {
				x.Div(&x, &y)
			}
			stopProfile()
			order := bls24315fp.Modulus()
			return took, order
		case "exp":
			var x bls24315fp.Element
			x.SetRandom()
			max := big.NewInt(1000000)
			k, err := rand.Int(rand.Reader, max)
			if err != nil {
				panic(err)
			}
			startProfile()
			x.Exp(x, k)
			stopProfile()
			order := bls24315fp.Modulus()
			return took, order
		default:
			panic("arithmetic operation not implemented")
		}
	} else {
		panic("field not valid")
	}
}

func ExecuteOperationBW6633(operation string, x float64, y float64) (time.Duration, *big.Int) {

	if *fField == "scalar" {
		switch operation {
		case "add":
			var x, y bw6633fr.Element
			x.SetInterface(x)
			y.SetInterface(y)
			startProfile()
			for i := 0; i < *fCount; i++ {
				x.Add(&x, &y)
			}
			stopProfile()
			order := bw6633fr.Modulus()
			return took, order
		case "sub":
			var x, y bw6633fr.Element
			x.SetInterface(x)
			y.SetInterface(y)
			startProfile()
			for i := 0; i < *fCount; i++ {
				x.Sub(&x, &y)
			}
			stopProfile()
			order := bw6633fr.Modulus()
			return took, order
		case "mul":
			var x, y bw6633fr.Element
			x.SetInterface(x)
			y.SetInterface(y)
			startProfile()
			for i := 0; i < *fCount; i++ {
				x.Mul(&x, &y)
			}
			stopProfile()
			order := bw6633fr.Modulus()
			return took, order
		case "div":
			var x, y bw6633fr.Element
			x.SetInterface(x)
			y.SetInterface(y)
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
	} else if *fField == "base" {
		switch operation {
		case "add":
			var x, y bw6633fp.Element
			x.SetInterface(x)
			y.SetInterface(y)
			startProfile()
			for i := 0; i < *fCount; i++ {
				x.Add(&x, &y)
			}
			stopProfile()
			order := bw6633fp.Modulus()
			return took, order
		case "sub":
			var x, y bw6633fp.Element
			x.SetInterface(x)
			y.SetInterface(y)
			startProfile()
			for i := 0; i < *fCount; i++ {
				x.Sub(&x, &y)
			}
			stopProfile()
			order := bw6633fp.Modulus()
			return took, order
		case "mul":
			var x, y bw6633fp.Element
			x.SetInterface(x)
			y.SetInterface(y)
			startProfile()
			for i := 0; i < *fCount; i++ {
				x.Mul(&x, &y)
			}
			stopProfile()
			order := bw6633fp.Modulus()
			return took, order
		case "div":
			var x, y bw6633fp.Element
			x.SetInterface(x)
			y.SetInterface(y)
			startProfile()
			for i := 0; i < *fCount; i++ {
				x.Div(&x, &y)
			}
			stopProfile()
			order := bw6633fp.Modulus()
			return took, order
		case "exp":
			var x bw6633fp.Element
			x.SetRandom()
			max := big.NewInt(1000000)
			k, err := rand.Int(rand.Reader, max)
			if err != nil {
				panic(err)
			}
			startProfile()
			x.Exp(x, k)
			stopProfile()
			order := bw6633fp.Modulus()
			return took, order
		default:
			panic("arithmetic operation not implemented")
		}
	} else {
		panic("field not valid")
	}
}

func ExecuteOperationBW6761(operation string, x float64, y float64) (time.Duration, *big.Int) {

	if *fField == "scalar" {
		switch operation {
		case "add":
			var x, y bw6761fr.Element
			x.SetInterface(x)
			y.SetInterface(y)
			startProfile()
			for i := 0; i < *fCount; i++ {
				x.Add(&x, &y)
			}
			stopProfile()
			order := bw6761fr.Modulus()
			return took, order
		case "sub":
			var x, y bw6761fr.Element
			x.SetInterface(x)
			y.SetInterface(y)
			startProfile()
			for i := 0; i < *fCount; i++ {
				x.Sub(&x, &y)
			}
			stopProfile()
			order := bw6761fr.Modulus()
			return took, order
		case "mul":
			var x, y bw6761fr.Element
			x.SetInterface(x)
			y.SetInterface(y)
			startProfile()
			for i := 0; i < *fCount; i++ {
				x.Mul(&x, &y)
			}
			stopProfile()
			order := bw6761fr.Modulus()
			return took, order
		case "div":
			var x, y bw6761fr.Element
			x.SetInterface(x)
			y.SetInterface(y)
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
	} else if *fField == "base" {
		switch operation {
		case "add":
			var x, y bw6761fp.Element
			x.SetInterface(x)
			y.SetInterface(y)
			startProfile()
			for i := 0; i < *fCount; i++ {
				x.Add(&x, &y)
			}
			stopProfile()
			order := bw6761fp.Modulus()
			return took, order
		case "sub":
			var x, y bw6761fp.Element
			x.SetInterface(x)
			y.SetInterface(y)
			startProfile()
			for i := 0; i < *fCount; i++ {
				x.Sub(&x, &y)
			}
			stopProfile()
			order := bw6761fp.Modulus()
			return took, order
		case "mul":
			var x, y bw6761fp.Element
			x.SetInterface(x)
			y.SetInterface(y)
			startProfile()
			for i := 0; i < *fCount; i++ {
				x.Mul(&x, &y)
			}
			stopProfile()
			order := bw6761fp.Modulus()
			return took, order
		case "div":
			var x, y bw6761fp.Element
			x.SetInterface(x)
			y.SetInterface(y)
			startProfile()
			for i := 0; i < *fCount; i++ {
				x.Div(&x, &y)
			}
			stopProfile()
			order := bw6761fp.Modulus()
			return took, order
		case "exp":
			var x bw6761fp.Element
			x.SetRandom()
			max := big.NewInt(1000000)
			k, err := rand.Int(rand.Reader, max)
			if err != nil {
				panic(err)
			}
			startProfile()
			x.Exp(x, k)
			stopProfile()
			order := bw6761fp.Modulus()
			return took, order
		default:
			panic("arithmetic operation not implemented")
		}
	} else {
		panic("field not valid")
	}
}

func benchArithmetic(cmd *cobra.Command, args []string) {

	log := logger.Logger()
	log.Info().Msg("Benchmarking arithmetics - gnark: " + *fCurve + " " + *fField + " " + *fOperation + " " + *fInputPath)

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
			Curve:     curveID.String(),
			Field:     *fField,
			Operation: *fOperation,
			Input:     *fInputPath,
			MaxRAM:    (m.Sys),
            Count:     *fCount,
			RunTime:   took.Nanoseconds(),
		}

		if err := util.WriteData("csv", bDataArith, filename); err != nil {
			panic(err)
		}
	}

	data, err := util.ReadFromInputPath(*fInputPath)
	if err != nil {
		panic(err)
	}

	switch curveID {
	case ecc.BN254:
		took, order := ExecuteOperation254(*fOperation, data["x"].(float64), data["y"].(float64))
		writeResults(took, *order)
	case ecc.BLS12_381:
		took, order := ExecuteOperationBLS12381(*fOperation, data["x"].(float64), data["y"].(float64))
		writeResults(took, *order)
	case ecc.BLS12_377:
		took, order := ExecuteOperationBLS12377(*fOperation, data["x"].(float64), data["y"].(float64))
		writeResults(took, *order)
	case ecc.BLS24_315:
		took, order := ExecuteOperationBLS24315(*fOperation, data["x"].(float64), data["y"].(float64))
		writeResults(took, *order)
	case ecc.BW6_633:
		took, order := ExecuteOperationBW6633(*fOperation, data["x"].(float64), data["y"].(float64))
		writeResults(took, *order)
	case ecc.BW6_761:
		took, order := ExecuteOperationBW6761(*fOperation, data["x"].(float64), data["y"].(float64))
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

	// Possible Operations: add, sub, mul, exp
	fField = arithmeticCmd.Flags().String("field", "scalar", "field to benchmark over")

	rootCmd.AddCommand(arithmeticCmd)
}
