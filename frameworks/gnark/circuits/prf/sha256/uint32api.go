package sha256

import (
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/math/bits"
)

// uint32api performs binary operations on xuint32 variables. In the
// future possibly using lookup tables.
//
// TODO: we could possibly optimise using hints if working over many inputs. For
// example, if we OR many bits, then the result is 0 if the sum of the bits is
// larger than 1. And AND is 1 if the sum of bits is the number of inputs. BUt
// this probably helps only if we have a lot of similar operations in a row
// (more than 4). We could probably unroll the whole permutation and expand all
// the formulas to see. But long term tables are still better.
type uint32api struct {
	api frontend.API
}

func newUint32API(api frontend.API) *uint32api {
	return &uint32api{
		api: api,
	}
}

// varUint8 represents 32-bit unsigned integer. We use this type to ensure that
// we work over constrained bits. Do not initialize directly, use [wideBinaryOpsApi.asUint8].
type xuint32 [32]frontend.Variable

func constUint32(a uint32) xuint32 {
	var res xuint32
	for i := 0; i < 32; i++ {
		res[i] = (a >> i) & 1
	}
	return res
}

func (w *uint32api) asUint32(in frontend.Variable) xuint32 {
	bits := bits.ToBinary(w.api, in, bits.WithNbDigits(32))
	var res xuint32
	copy(res[:], bits)
	return res
}

func (w *uint32api) fromUint32(in xuint32) frontend.Variable {
	return bits.FromBinary(w.api, in[:], bits.WithUnconstrainedInputs())
}

func (w *uint32api) and(in ...xuint32) xuint32 {
	var res xuint32
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

func (w *uint32api) or(in ...xuint32) xuint32 {
	var res xuint32
	for i := range res {
		res[i] = 0
	}
	for i := range res {
		for _, v := range in {
			res[i] = w.api.Or(res[i], v[i])
		}
	}
	return res
}

func (w *uint32api) xor(in ...xuint32) xuint32 {
	var res xuint32
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

func (w *uint32api) lrot(in xuint32, shift int) xuint32 {
	var res xuint32
	for i := range res {
		res[i] = in[(i-shift+32)%32]
	}
	return res
}

func (w *uint32api) not(in xuint32) xuint32 {
	// TODO: it would be better to have separate method for it. If we have
	// native API support, then in R1CS would be free (1-X) and in PLONK 1
	// constraint (1-X). But if we do XOR, then we always have a constraint with
	// R1CS (not sure if 1-2 with PLONK). If we do 1-X ourselves, then compiler
	// marks as binary which is 1-2 (R1CS-PLONK).
	var res xuint32
	for i := range res {
		res[i] = w.api.Xor(in[i], 1)
	}
	return res
}

func (w *uint32api) rshift(in xuint32, shift int) xuint32 {
	var res xuint32
	for i := 0; i < 32-shift; i++ {
		res[i] = in[i+shift]
	}
	for i := 32 - shift; i < 32; i++ {
		res[i] = 0
	}
	return res
}

func (w *uint32api) lshift(in xuint32, shift int) xuint32 {
	var res xuint32
	for i := 0; i < shift; i++ {
		res[i] = 0
	}
	for i := shift; i < 32; i++ {
		res[i] = in[i-shift]
	}
	return res
}

func (w *uint32api) add(i1, i2 xuint32, in ...xuint32) xuint32 {
	var v []frontend.Variable
	for _, i := range in {
		v = append(v, w.fromUint32(i))
	}
	sum := w.api.Add(w.fromUint32(i1), w.fromUint32(i2), v...)

	b := bits.ToBinary(w.api, sum, bits.WithNbDigits(33+len(in)))
	var res xuint32
	copy(res[:], b)

	return res
}

func (w *uint32api) assertEq(a, b xuint32) {
	for i := range a {
		w.api.AssertIsEqual(a[i], b[i])
	}
}

func (in xuint32) toUnit8() xuint8 {
	var res xuint8
	for i := 0; i < 8; i++ {
		res[i] = in[i]
	}
	return res
}
