package sha256

import (
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/math/bits"
)

// uint8api performs binary operations on xuint8 variables. In the
// future possibly using lookup tables.
//
// TODO: we could possibly optimise using hints if working over many inputs. For
// example, if we OR many bits, then the result is 0 if the sum of the bits is
// larger than 1. And AND is 1 if the sum of bits is the number of inputs. BUt
// this probably helps only if we have a lot of similar operations in a row
// (more than 4). We could probably unroll the whole permutation and expand all
// the formulas to see. But long term tables are still better.
type uint8api struct {
	api frontend.API
}

func newUint8API(api frontend.API) *uint8api {
	return &uint8api{
		api: api,
	}
}

// varUint8 represents 8-bit unsigned integer. We use this type to ensure that
// we work over constrained bits. Do not initialize directly, use [wideBinaryOpsApi.asUint8].
type xuint8 [8]frontend.Variable

func constUint8(a uint8) xuint8 {
	var res xuint8
	for i := 0; i < 8; i++ {
		res[i] = (a >> i) & 1
	}
	return res
}

func (w *uint8api) asUint8(in frontend.Variable) xuint8 {
	bits := bits.ToBinary(w.api, in, bits.WithNbDigits(8))
	var res xuint8
	copy(res[:], bits)
	return res
}

func (w *uint8api) fromUint8(in xuint8) frontend.Variable {
	return bits.FromBinary(w.api, in[:], bits.WithUnconstrainedInputs())
}

func (w *uint8api) and(in ...xuint8) xuint8 {
	var res xuint8
	for i := range res {
		res[i] = 1
	}
	for i := range res {
		for _, v := range in {
			res[i] = w.api.And(res[i], v[i])
		}
	}
	return res
}

func (w *uint8api) xor(in ...xuint8) xuint8 {
	var res xuint8
	for i := range res {
		res[i] = 0
	}
	for i := range res {
		for _, v := range in {
			res[i] = w.api.Xor(res[i], v[i])
		}
	}
	return res
}

func (w *uint8api) lrot(in xuint8, shift int) xuint8 {
	var res xuint8
	for i := range res {
		res[i] = in[(i-shift+8)%8]
	}
	return res
}

func (w *uint8api) not(in xuint8) xuint8 {
	// TODO: it would be better to have separate method for it. If we have
	// native API support, then in R1CS would be free (1-X) and in PLONK 1
	// constraint (1-X). But if we do XOR, then we always have a constraint with
	// R1CS (not sure if 1-2 with PLONK). If we do 1-X ourselves, then compiler
	// marks as binary which is 1-2 (R1CS-PLONK).
	var res xuint8
	for i := range res {
		res[i] = w.api.Xor(in[i], 1)
	}
	return res
}

func (w *uint8api) assertEq(a, b xuint8) {
	for i := range a {
		w.api.AssertIsEqual(a[i], b[i])
	}
}

func (a xuint8) toUint64() xuint64 {
	var res xuint64
	for i := 0; i < 8; i++ {
		res[i] = a[i]
	}
	for i := 8; i < 64; i++ {
		res[i] = 0
	}
	return res
}
