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
	"github.com/consensys/gnark/logger"
	"github.com/spf13/cobra"
	"github.com/tumberger/zk-compilers/gnark/util"

	bls12377 "github.com/consensys/gnark-crypto/ecc/bls12-377"
	bls12377fr "github.com/consensys/gnark-crypto/ecc/bls12-377/fr"
	bls12381 "github.com/consensys/gnark-crypto/ecc/bls12-381"
	bls12381fr "github.com/consensys/gnark-crypto/ecc/bls12-381/fr"
	bls24315 "github.com/consensys/gnark-crypto/ecc/bls24-315"
	bls24315fr "github.com/consensys/gnark-crypto/ecc/bls24-315/fr"
	bn254 "github.com/consensys/gnark-crypto/ecc/bn254"
	bn254fr "github.com/consensys/gnark-crypto/ecc/bn254/fr"
	bw6633 "github.com/consensys/gnark-crypto/ecc/bw6-633"
	bw6633fr "github.com/consensys/gnark-crypto/ecc/bw6-633/fr"
	bw6761 "github.com/consensys/gnark-crypto/ecc/bw6-761"
	bw6761fr "github.com/consensys/gnark-crypto/ecc/bw6-761/fr"
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
			var sampleScalars [nbSamples]bn254fr.Element
			for i := 0; i < len(sampleScalars); i++ {
				sampleScalars[i].SetRandom()
			}
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
			var sampleScalars [nbSamples]bn254fr.Element
			for i := 0; i < len(sampleScalars); i++ {
				sampleScalars[i].SetRandom()
			}
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

func CurveOperationBLS12377(operation string) time.Duration {
	switch operation {
	case "scalar-multiplication":
		// Scalar Multiplication in Jacobian coordinates as it's more efficient than affine coordinates in gnark
		if *fGroup == "g1" {
			var g1Jac bls12377.G1Jac
			var a bls12377fr.Element
			var b big.Int
			startProfile()
			for i := 0; i < *fCount; i++ {
				g1Jac.ScalarMultiplication(&g1Jac, a.BigInt(&b))
			}
			stopProfile()
			return took
		} else if *fGroup == "g2" {
			var g2Jac bls12377.G2Jac
			var test bls12377fr.Element
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
			var samplePoints [nbSamples]bls12377.G1Affine
			var sampleScalars [nbSamples]bls12377fr.Element
			for i := 0; i < len(sampleScalars); i++ {
				sampleScalars[i].SetRandom()
			}
			var g bls12377.G1Jac
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
			var samplePoints [nbSamples]bls12377.G2Affine
			var sampleScalars [nbSamples]bls12377fr.Element
			for i := 0; i < len(sampleScalars); i++ {
				sampleScalars[i].SetRandom()
			}
			var g bls12377.G2Jac
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
		_, _, g1GenAff, g2GenAff := bls12377.Generators()
		var ag1 bls12377.G1Affine
		var bg2 bls12377.G2Affine
		var abigint, bbigint big.Int
		// Get Points on EC
		ag1.ScalarMultiplication(&g1GenAff, &abigint)
		bg2.ScalarMultiplication(&g2GenAff, &bbigint)
		startProfile()
		for i := 0; i < *fCount; i++ {
			_, err := bls12377.Pair([]bls12377.G1Affine{ag1}, []bls12377.G2Affine{bg2})
			// Pair(api, []G1Affine{circuit.P}, []G2Affine{circuit.Q})
			// _, err := bls12377.Pair([]bls12377.G1Affine{*b)
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

func CurveOperationBLS12381(operation string) time.Duration {
	switch operation {
	case "scalar-multiplication":
		// Scalar Multiplication in Jacobian coordinates as it's more efficient than affine coordinates in gnark
		if *fGroup == "g1" {
			var g1Jac bls12381.G1Jac
			var a bls12381fr.Element
			var b big.Int
			startProfile()
			for i := 0; i < *fCount; i++ {
				g1Jac.ScalarMultiplication(&g1Jac, a.BigInt(&b))
			}
			stopProfile()
			return took
		} else if *fGroup == "g2" {
			var g2Jac bls12381.G2Jac
			var test bls12381fr.Element
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
			var samplePoints [nbSamples]bls12381.G1Affine
			var sampleScalars [nbSamples]bls12381fr.Element
			for i := 0; i < len(sampleScalars); i++ {
				sampleScalars[i].SetRandom()
			}
			var g bls12381.G1Jac
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
			var samplePoints [nbSamples]bls12381.G2Affine
			var sampleScalars [nbSamples]bls12381fr.Element
			for i := 0; i < len(sampleScalars); i++ {
				sampleScalars[i].SetRandom()
			}
			var g bls12381.G2Jac
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
		_, _, g1GenAff, g2GenAff := bls12381.Generators()
		var ag1 bls12381.G1Affine
		var bg2 bls12381.G2Affine
		var abigint, bbigint big.Int
		// Get Points on EC
		ag1.ScalarMultiplication(&g1GenAff, &abigint)
		bg2.ScalarMultiplication(&g2GenAff, &bbigint)
		startProfile()
		for i := 0; i < *fCount; i++ {
			_, err := bls12381.Pair([]bls12381.G1Affine{ag1}, []bls12381.G2Affine{bg2})
			// Pair(api, []G1Affine{circuit.P}, []G2Affine{circuit.Q})
			// _, err := bls12381.Pair([]bls12381.G1Affine{*b)
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

func CurveOperationBLS24315(operation string) time.Duration {
	switch operation {
	case "scalar-multiplication":
		// Scalar Multiplication in Jacobian coordinates as it's more efficient than affine coordinates in gnark
		if *fGroup == "g1" {
			var g1Jac bls24315.G1Jac
			var a bls24315fr.Element
			var b big.Int
			startProfile()
			for i := 0; i < *fCount; i++ {
				g1Jac.ScalarMultiplication(&g1Jac, a.BigInt(&b))
			}
			stopProfile()
			return took
		} else if *fGroup == "g2" {
			var g2Jac bls24315.G2Jac
			var test bls24315fr.Element
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
			var samplePoints [nbSamples]bls24315.G1Affine
			var sampleScalars [nbSamples]bls24315fr.Element
			for i := 0; i < len(sampleScalars); i++ {
				sampleScalars[i].SetRandom()
			}
			var g bls24315.G1Jac
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
			var samplePoints [nbSamples]bls24315.G2Affine
			var sampleScalars [nbSamples]bls24315fr.Element
			for i := 0; i < len(sampleScalars); i++ {
				sampleScalars[i].SetRandom()
			}
			var g bls24315.G2Jac
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
		_, _, g1GenAff, g2GenAff := bls24315.Generators()
		var ag1 bls24315.G1Affine
		var bg2 bls24315.G2Affine
		var abigint, bbigint big.Int
		// Get Points on EC
		ag1.ScalarMultiplication(&g1GenAff, &abigint)
		bg2.ScalarMultiplication(&g2GenAff, &bbigint)
		startProfile()
		for i := 0; i < *fCount; i++ {
			_, err := bls24315.Pair([]bls24315.G1Affine{ag1}, []bls24315.G2Affine{bg2})
			// Pair(api, []G1Affine{circuit.P}, []G2Affine{circuit.Q})
			// _, err := bls24315.Pair([]bls24315.G1Affine{*b)
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

func CurveOperationBW6633(operation string) time.Duration {
	switch operation {
	case "scalar-multiplication":
		// Scalar Multiplication in Jacobian coordinates as it's more efficient than affine coordinates in gnark
		if *fGroup == "g1" {
			var g1Jac bw6633.G1Jac
			var a bw6633fr.Element
			var b big.Int
			startProfile()
			for i := 0; i < *fCount; i++ {
				g1Jac.ScalarMultiplication(&g1Jac, a.BigInt(&b))
			}
			stopProfile()
			return took
		} else if *fGroup == "g2" {
			var g2Jac bw6633.G2Jac
			var test bw6633fr.Element
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
			var samplePoints [nbSamples]bw6633.G1Affine
			var sampleScalars [nbSamples]bw6633fr.Element
			for i := 0; i < len(sampleScalars); i++ {
				sampleScalars[i].SetRandom()
			}
			var g bw6633.G1Jac
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
			var samplePoints [nbSamples]bw6633.G2Affine
			var sampleScalars [nbSamples]bw6633fr.Element
			for i := 0; i < len(sampleScalars); i++ {
				sampleScalars[i].SetRandom()
			}
			var g bw6633.G2Jac
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
		_, _, g1GenAff, g2GenAff := bw6633.Generators()
		var ag1 bw6633.G1Affine
		var bg2 bw6633.G2Affine
		var abigint, bbigint big.Int
		// Get Points on EC
		ag1.ScalarMultiplication(&g1GenAff, &abigint)
		bg2.ScalarMultiplication(&g2GenAff, &bbigint)
		startProfile()
		for i := 0; i < *fCount; i++ {
			_, err := bw6633.Pair([]bw6633.G1Affine{ag1}, []bw6633.G2Affine{bg2})
			// Pair(api, []G1Affine{circuit.P}, []G2Affine{circuit.Q})
			// _, err := bw6633.Pair([]bw6633.G1Affine{*b)
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

func CurveOperationBW6761(operation string) time.Duration {
	switch operation {
	case "scalar-multiplication":
		// Scalar Multiplication in Jacobian coordinates as it's more efficient than affine coordinates in gnark
		if *fGroup == "g1" {
			var g1Jac bw6761.G1Jac
			var a bw6761fr.Element
			var b big.Int
			startProfile()
			for i := 0; i < *fCount; i++ {
				g1Jac.ScalarMultiplication(&g1Jac, a.BigInt(&b))
			}
			stopProfile()
			return took
		} else if *fGroup == "g2" {
			var g2Jac bw6761.G2Jac
			var test bw6761fr.Element
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
			var samplePoints [nbSamples]bw6761.G1Affine
			var sampleScalars [nbSamples]bw6761fr.Element
			for i := 0; i < len(sampleScalars); i++ {
				sampleScalars[i].SetRandom()
			}
			var g bw6761.G1Jac
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
			var samplePoints [nbSamples]bw6761.G2Affine
			var sampleScalars [nbSamples]bw6761fr.Element
			for i := 0; i < len(sampleScalars); i++ {
				sampleScalars[i].SetRandom()
			}
			var g bw6761.G2Jac
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
		_, _, g1GenAff, g2GenAff := bw6761.Generators()
		var ag1 bw6761.G1Affine
		var bg2 bw6761.G2Affine
		var abigint, bbigint big.Int
		// Get Points on EC
		ag1.ScalarMultiplication(&g1GenAff, &abigint)
		bg2.ScalarMultiplication(&g2GenAff, &bbigint)
		startProfile()
		for i := 0; i < *fCount; i++ {
			_, err := bw6761.Pair([]bw6761.G1Affine{ag1}, []bw6761.G2Affine{bg2})
			// Pair(api, []G1Affine{circuit.P}, []G2Affine{circuit.Q})
			// _, err := bw6761.Pair([]bw6761.G1Affine{*b)
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

	log := logger.Logger()
	log.Info().Msg("Benchmarking curve operations - gnark: " + *fCurve + " " + *fOperation + " " + *fInputPath)

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

		if *fOperation == "pairing" {
			operationString = string(*fOperation)
		} else {
			operationString = string(*fGroup) + "-" + string(*fOperation)
		}

		bDataArith := util.BenchDataCurve{
			Framework: "gnark",
			Category:  "ec",
			Curve:     curveID.String(),
			Operation: operationString,
			Input:     *fInputPath,
			MaxRAM:    (m.Sys),
            Count:     *fCount,
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

	switch curveID {
	case ecc.BN254:
		took := CurveOperation254(*fOperation)
		writeResults(took)
	case ecc.BLS12_381:
		took := CurveOperationBLS12377(*fOperation)
		writeResults(took)
	case ecc.BLS12_377:
		took := CurveOperationBLS12381(*fOperation)
		writeResults(took)
	case ecc.BLS24_315:
		took := CurveOperationBLS24315(*fOperation)
		writeResults(took)
	case ecc.BW6_633:
		took := CurveOperationBW6633(*fOperation)
		writeResults(took)
	case ecc.BW6_761:
		took := CurveOperationBW6761(*fOperation)
		writeResults(took)
	default:
		panic("field order not implemented")
	}
}

func init() {

	// Possible Operations: add, sub, mul, div, exp
	fGroup = curveCmd.Flags().String("group", "None", "group to benchmark")

	rootCmd.AddCommand(curveCmd)
}
