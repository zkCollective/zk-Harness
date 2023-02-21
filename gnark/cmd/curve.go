/*
Benchmarking Elliptic Curve Operations over a variety of curves in gnark
*/
package cmd

import (
	"fmt"
	"math/big"
	"os"
	"runtime"
	"time"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/spf13/cobra"
	"github.com/tumberger/zk-compilers/gnark/util"

	bn254 "github.com/consensys/gnark-crypto/ecc/bn254"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	bn254fr "github.com/consensys/gnark-crypto/ecc/bn254/fr"
)

var curveCmd = &cobra.Command{
	Use:   "ec",
	Short: "runs benchmarks and profiles for the gnark arithmetic operations",
	Run:   benchCurveOperations,
}

// Operations for Curve BN254
func CurveOperation254(operation string) time.Duration {
	switch operation {

	case "scalar-multiplication":
		// Scalar Multiplication in Jacobian coordinates as it's more efficient than affine coordinates in gnark
		if *fGroup == "g1" {
			var g1Jac bn254.G1Jac
			var a bn254fr.Element
			var b big.Int
			startProfile()
			for i := 0; i < *fCount; i++ {
				g1Jac.ScalarMultiplication(&g1Jac, a.BigInt(&b))
			}
			stopProfile()
			return took
		} else if *fGroup == "g2" {
			var g2Jac bn254.G2Jac
			var test bn254fr.Element
			var b big.Int
			startProfile()
			for i := 0; i < *fCount; i++ {
				g2Jac.ScalarMultiplication(&g2Jac, test.BigInt(&b))
			}
			stopProfile()
			return took
		} else {
			panic("group not defined for this operation")
		}
	case "multi-scalar-multiplication":
		if *fGroup == "g1" {
			// size of the multiExp
			const nbSamples = 73
			// multi exp points
			var samplePoints [nbSamples]bn254.G1Affine
			var sampleScalars [nbSamples]fr.Element
			fillBenchScalars(sampleScalars[:])
			var g bn254.G1Jac
			g.Set(&g)
			for i := 1; i <= nbSamples; i++ {
				samplePoints[i-1].FromJacobian(&g)
				g.AddAssign(&g)
			}
			n := runtime.NumCPU()
			startProfile()
			for i := 0; i < *fCount; i++ {
				_, err := g.MultiExp(samplePoints[:], sampleScalars[:], ecc.MultiExpConfig{NbTasks: n / 2})
				if err != nil {
					panic(err)
				}
			}
			stopProfile()
			return took
		} else if *fGroup == "g2" {
			// size of the multiExp
			const nbSamples = 73
			// multi exp points
			var samplePoints [nbSamples]bn254.G2Affine
			var sampleScalars [nbSamples]fr.Element
			fillBenchScalars(sampleScalars[:])
			var g bn254.G2Jac
			g.Set(&g)
			for i := 1; i <= nbSamples; i++ {
				samplePoints[i-1].FromJacobian(&g)
				g.AddAssign(&g)
			}
			n := runtime.NumCPU()
			startProfile()
			for i := 0; i < *fCount; i++ {
				_, err := g.MultiExp(samplePoints[:], sampleScalars[:], ecc.MultiExpConfig{NbTasks: n / 2})
				if err != nil {
					panic(err)
				}
			}
			stopProfile()
			return took
		} else {
			panic("group not defined for this operation")
		}
	case "pairing":
		_, _, g1GenAff, g2GenAff := bn254.Generators()
		var ag1 bn254.G1Affine
		var bg2 bn254.G2Affine
		var abigint, bbigint big.Int
		// Get Points on EC
		ag1.ScalarMultiplication(&g1GenAff, &abigint)
		bg2.ScalarMultiplication(&g2GenAff, &bbigint)
		startProfile()
		for i := 0; i < *fCount; i++ {
			_, err := bn254.Pair([]bn254.G1Affine{ag1}, []bn254.G2Affine{bg2})
			// Pair(api, []G1Affine{circuit.P}, []G2Affine{circuit.Q})
			// _, err := bn254.Pair([]bn254.G1Affine{*b)
			if err != nil {
				panic(err)
			}
		}
		stopProfile()
		return took
	default:
		panic("arithmetic operation not implemented")
	}
}

func benchCurveOperations(cmd *cobra.Command, args []string) {
	var filename = "../benchmarks/gnark/gnark_" +
		"curve" +
		"." +
		*fFileType

	if err := parseFlags(); err != nil {
		fmt.Println("error: ", err.Error())
		cmd.Help()
		os.Exit(-1)
	}

	writeResults := func(took time.Duration) {
		// check memory usage, max ram requested from OS
		var operationString string
		var m runtime.MemStats
		runtime.ReadMemStats(&m)

		if *fGroup != "None" {
			operationString = string(*fGroup) + "-" + string(*fOperation)
		} else {
			operationString = string(*fOperation)
		}

		bDataArith := util.BenchDataCurve{
			Framework: "gnark",
			Category:  "arithmetic",
			Curve:     curveID.String(),
			Operation: operationString,
			Input:     *fInputPath,
			MaxRAM:    (m.Sys / 1024 / 1024),
			RunTime:   took.Nanoseconds(),
		}

		if err := util.WriteData("csv", bDataArith, filename); err != nil {
			panic(err)
		}
	}

	// Read input data given the input path
	// data, err := util.ReadFromInputPath(*fInputPath)
	// if err != nil {
	// 	panic(err)
	// }

	// func (p *G1Affine) MultiExp(points []G1Affine, scalars []fr.Element, config ecc.MultiExpConfig) (*G1Affine, error)
	// func (p *G1Jac) MultiExp(points []G1Affine, scalars []fr.Element, config ecc.MultiExpConfig) (*G1Jac, error)
	switch curveID {
	case ecc.BN254:
		took := CurveOperation254(*fOperation)
		writeResults(took)
	// case ecc.BLS12_381:
	// 	took, order := ExecuteOperationBLS12381(*fOperation, data["x"].(float64), data["y"].(float64))
	// 	writeResults(took, *order)
	// case ecc.BLS12_377:
	// 	took, order := ExecuteOperationBLS12377(*fOperation, data["x"].(float64), data["y"].(float64))
	// 	writeResults(took, *order)
	// case ecc.BLS24_315:
	// 	took, order := ExecuteOperationBLS24315(*fOperation, data["x"].(float64), data["y"].(float64))
	// 	writeResults(took, *order)
	// case ecc.BW6_633:
	// 	took, order := ExecuteOperationBW6633(*fOperation, data["x"].(float64), data["y"].(float64))
	// 	writeResults(took, *order)
	// case ecc.BW6_761:
	// 	took, order := ExecuteOperationBW6761(*fOperation, data["x"].(float64), data["y"].(float64))
	// 	writeResults(took, *order)
	default:
		panic("field order not implemented")
	}
}

// Helper function for multi-exponentiation - fill
func fillBenchScalars(sampleScalars []bn254fr.Element) {
	// ensure every words of the scalars are filled
	for i := 0; i < len(sampleScalars); i++ {
		sampleScalars[i].SetRandom()
	}
}

func init() {

	// Possible Operations: add, sub, mul, div, exp
	fGroup = curveCmd.Flags().String("group", "None", "group to benchmark")
	fOperation = curveCmd.Flags().String("operation", "scalar-multiplication", "operation to benchmark")

	rootCmd.AddCommand(curveCmd)
}
