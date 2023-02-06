package sha256

import (
	"encoding/hex"
	"log"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
)

type SHA256 struct {
	PreImage []frontend.Variable
	Output   [32]frontend.Variable `gnark:",public"`
}

// Define declares the circuit's constraints
func (circuit *SHA256) Define(api frontend.API) error {

	// K values
	var K = [64]frontend.Variable{0x428a2f98, 0x71374491, 0xb5c0fbcf, 0xe9b5dba5, 0x3956c25b, 0x59f111f1, 0x923f82a4, 0xab1c5ed5, 0xd807aa98, 0x12835b01, 0x243185be, 0x550c7dc3, 0x72be5d74, 0x80deb1fe, 0x9bdc06a7, 0xc19bf174, 0xe49b69c1, 0xefbe4786, 0x0fc19dc6, 0x240ca1cc, 0x2de92c6f, 0x4a7484aa, 0x5cb0a9dc, 0x76f988da, 0x983e5152, 0xa831c66d, 0xb00327c8, 0xbf597fc7, 0xc6e00bf3, 0xd5a79147, 0x06ca6351, 0x14292967, 0x27b70a85, 0x2e1b2138, 0x4d2c6dfc, 0x53380d13, 0x650a7354, 0x766a0abb, 0x81c2c92e, 0x92722c85, 0xa2bfe8a1, 0xa81a664b, 0xc24b8b70, 0xc76c51a3, 0xd192e819, 0xd6990624, 0xf40e3585, 0x106aa070, 0x19a4c116, 0x1e376c08, 0x2748774c, 0x34b0bcb5, 0x391c0cb3, 0x4ed8aa4a, 0x5b9cca4f, 0x682e6ff3, 0x748f82ee, 0x78a5636f, 0x84c87814, 0x8cc70208, 0x90befffa, 0xa4506ceb, 0xbef9a3f7, 0xc67178f2}

	// H values
	var H [8]frontend.Variable
	H[0] = frontend.Variable(0x6A09E667)
	H[1] = frontend.Variable(0xBB67AE85)
	H[2] = frontend.Variable(0x3C6EF372)
	H[3] = frontend.Variable(0xA54FF53A)
	H[4] = frontend.Variable(0x510E527F)
	H[5] = frontend.Variable(0x9B05688C)
	H[6] = frontend.Variable(0x1F83D9AB)
	H[7] = frontend.Variable(0x5BE0CD19)

	// padding
	paddedInput := padding(api, circuit.PreImage)

	// chunk processing of padded input
	numberChunks := int(len(paddedInput) / 64)
	for epoch := 0; epoch < numberChunks; epoch++ {

		eIndex := epoch * 64

		// w values init
		var w [64]frontend.Variable

		// first 16 w values is set based on input data
		for i := 0; i < 16; i++ {

			j := i * 4

			// same as in go except that | is replaced with api.Add for multi-bit operation
			leftShift24 := shiftLeft(api, paddedInput[eIndex+j], 24)
			leftShift16 := shiftLeft(api, paddedInput[eIndex+j+1], 16)
			leftShift8 := shiftLeft(api, paddedInput[eIndex+j+2], 8)
			leftShiftNone := api.FromBinary(api.ToBinary(paddedInput[eIndex+j+3], 32)...)
			w[i] = trimBits(api, api.Add(api.Add(api.Add(leftShift24, leftShift16), leftShift8), leftShiftNone), 34)
		}

		// remaining w values computation
		for i := 16; i < 64; i++ {

			// t1 := (bits.RotateLeft32(v1, -17)) ^ (bits.RotateLeft32(v1, -19)) ^ (v1 >> 10)
			v1 := w[i-2]

			rotateRight17 := rotateRight(api, v1, 17)
			rotateRight19 := rotateRight(api, v1, 19)
			rightShift10 := shiftRight(api, v1, 10)
			t1Slice := api.ToBinary(0, 32)
			for l := 0; l < 32; l++ {
				t1Slice[l] = api.Xor(api.Xor(rotateRight17[l], rotateRight19[l]), rightShift10[l])
			}
			t1 := api.FromBinary(t1Slice...)

			// t2 := (bits.RotateLeft32(v2, -7)) ^ (bits.RotateLeft32(v2, -18)) ^ (v2 >> 3)
			v2 := w[i-15]
			rotateRight7 := rotateRight(api, v2, 7)
			rotateRight18 := rotateRight(api, v2, 18)
			rightShift3 := shiftRight(api, v2, 3) // api.Div(v1, 3)
			t2Slice := api.ToBinary(0, 32)
			for l := 0; l < 32; l++ {
				t2Slice[l] = api.Xor(api.Xor(rotateRight7[l], rotateRight18[l]), rightShift3[l])
			}
			t2 := api.FromBinary(t2Slice...)

			w7 := w[i-7]
			w16 := w[i-16]
			w[i] = trimBits(api, api.Add(api.Add(api.Add(t1, w7), t2), w16), 34) // addition mod 2^32 ==> cut number to 32 bit
		}

		// a to h values
		var a frontend.Variable
		var b frontend.Variable
		var c frontend.Variable
		var d frontend.Variable
		var e frontend.Variable
		var f frontend.Variable
		var g frontend.Variable
		var h frontend.Variable

		a = H[0]
		b = H[1]
		c = H[2]
		d = H[3]
		e = H[4]
		f = H[5]
		g = H[6]
		h = H[7]

		// computation of alphabet values
		for i := 0; i < 64; i++ {

			// ^e is a not
			// t1 := h + ((bits.RotateLeft32(e, -6)) ^ (bits.RotateLeft32(e, -11)) ^ (bits.RotateLeft32(e, -25))) + ((e & f) ^ (^e & g)) + _K[i] + w[i]
			rotateRight6 := rotateRight(api, e, 6)
			rotateRight11 := rotateRight(api, e, 11)
			rotateRight25 := rotateRight(api, e, 25)
			tmp1Slice := api.ToBinary(0, 32)
			for k := 0; k < 32; k++ {
				tmp1Slice[k] = api.Xor(api.Xor(rotateRight6[k], rotateRight11[k]), rotateRight25[k])
			}
			tmp1 := api.FromBinary(tmp1Slice...)

			tmp2Slice := api.ToBinary(0, 32)
			eBits := api.ToBinary(e, 32)
			fBits := api.ToBinary(f, 32)
			gBits := api.ToBinary(g, 32)
			for k := 0; k < 32; k++ {
				tmp2Slice[k] = api.Xor(api.And(eBits[k], fBits[k]), api.And(api.Xor(eBits[k], frontend.Variable(1)), gBits[k]))
			}
			tmp2 := api.FromBinary(tmp2Slice...)

			// t1 computation
			t1 := api.Add(api.Add(api.Add(api.Add(h, tmp1), tmp2), K[i]), w[i])

			// t2 := ((bits.RotateLeft32(a, -2)) ^ (bits.RotateLeft32(a, -13)) ^ (bits.RotateLeft32(a, -22))) + ((a & b) ^ (a & c) ^ (b & c))
			rotateRight2 := rotateRight(api, a, 2)
			rotateRight13 := rotateRight(api, a, 13)
			rotateRight22 := rotateRight(api, a, 22)
			tmp3Slice := api.ToBinary(0, 32)
			for l := 0; l < 32; l++ {
				tmp3Slice[l] = api.Xor(api.Xor(rotateRight2[l], rotateRight13[l]), rotateRight22[l])
			}
			tmp3 := api.FromBinary(tmp3Slice...)

			// TODO: modulo from here: https://github.com/akosba/jsnark/blob/master/JsnarkCircuitBuilder/src/examples/gadgets/hash/SHA256Gadget.java
			// since after each iteration, SHA256 does c = b; and b = a;, we can make use of that to save multiplications in maj computation.
			// To do this, we make use of the caching feature, by just changing the order of wires sent to maj(). Caching will take care of the rest.
			minusTwo := [32]frontend.Variable{0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1} // -2 in little endian, of size 32
			tmp4Bits := api.ToBinary(0, 32)
			var x, y, z []frontend.Variable
			// var x, y, z frontend.Variable
			if i%2 == 1 {
				// x = c
				// y = b
				// z = a
				x = api.ToBinary(c, 32)
				y = api.ToBinary(b, 32)
				z = api.ToBinary(a, 32)
			} else {
				// x = a
				// y = b
				// z = c
				x = api.ToBinary(a, 32)
				y = api.ToBinary(b, 32)
				z = api.ToBinary(c, 32)
			}

			// least complexity (saves 20k constraints) but gives wrong values
			// t4t1 := api.Mul(x, y)
			// api.Println("t4t1:", t4t1)
			// t4t2 := api.Add(api.Add(x, y), api.Mul(t4t1, api.Neg(2)))
			// tmp4 := trimBits(api, api.Add(t4t1, api.Mul(z, t4t2)), 35)

			// working with less complexity compared to uncommented tmp4 calculation below works
			for j := 0; j < 32; j++ {
				t4t1 := api.And(x[j], y[j])
				t4t2 := api.Or(api.Or(x[j], y[j]), api.And(t4t1, minusTwo[j]))
				tmp4Bits[j] = api.Or(t4t1, api.And(z[j], t4t2))
			}
			tmp4 := api.FromBinary(tmp4Bits...)

			// this calculation of tmp4 works as well:
			// aBits := api.ToBinary(a, 32)
			// bBits := api.ToBinary(b, 32)
			// cBits := api.ToBinary(c, 32)
			// tmp4Slice := api.ToBinary(0, 32)
			// for l := 0; l < 32; l++ {
			// 	tmp4Slice[l] = api.Xor(api.Xor(api.And(aBits[l], bBits[l]), api.And(aBits[l], cBits[l])), api.And(bBits[l], cBits[l]))
			// }
			// tmp4 = api.FromBinary(tmp4Slice...)

			// t2 computation
			t2 := api.Add(tmp3, tmp4)

			h = g
			g = f
			f = e
			e = trimBits(api, api.Add(t1, d), 35)
			d = c
			c = b
			b = a
			a = trimBits(api, api.Add(t1, t2), 35)
		}

		// updating H values
		H[0] = trimBits(api, api.Add(H[0], a), 33)
		H[1] = trimBits(api, api.Add(H[1], b), 33)
		H[2] = trimBits(api, api.Add(H[2], c), 33)
		H[3] = trimBits(api, api.Add(H[3], d), 33)
		H[4] = trimBits(api, api.Add(H[4], e), 33)
		H[5] = trimBits(api, api.Add(H[5], f), 33)
		H[6] = trimBits(api, api.Add(H[6], g), 33)
		H[7] = trimBits(api, api.Add(H[7], h), 33)

	}

	// reorder bits
	var out [32]frontend.Variable
	ctr := 0
	for i := 0; i < 8; i++ {
		bits := api.ToBinary(H[i], 32)
		for j := 3; j >= 0; j-- {
			start := 8 * j
			// little endian order chunk parsing from back to front
			out[ctr] = api.FromBinary(bits[start : start+8]...)
			ctr += 1
		}
	}

	// constraints check
	for i := 0; i < len(circuit.Output); i++ {
		api.AssertIsEqual(circuit.Output[i], out[i])
	}

	return nil
}

func trimBits(api frontend.API, a frontend.Variable, size int) frontend.Variable {

	requiredSize := 32
	aBits := api.ToBinary(a, size)
	x := make([]frontend.Variable, requiredSize)

	for i := requiredSize; i < size; i++ {
		aBits[i] = 0
		// api.AssertIsEqual(aBits[i], 0)
	}
	for i := 0; i < requiredSize; i++ {
		x[i] = aBits[i]
	}

	return api.FromBinary(x...)
}

func shiftRight(api frontend.API, a frontend.Variable, shift int) []frontend.Variable {

	bits := api.ToBinary(a, 32)
	x := api.ToBinary(0, 32)

	for i := 0; i < 32; i++ {
		if i >= 32-shift {
			x[i] = 0
		} else {
			x[i] = bits[i+shift]
		}
	}
	return x
}

func shiftLeft(api frontend.API, a frontend.Variable, shift int) frontend.Variable {

	bits := api.ToBinary(a, 32)
	x := api.ToBinary(0, 32)

	for i := 0; i < 32; i++ {
		if i >= shift {
			x[i] = bits[i-shift]
		} else {
			x[i] = 0
		}
	}

	return api.FromBinary(x...)
}

func rotateRight(api frontend.API, a frontend.Variable, rotation int) []frontend.Variable {

	bits := api.ToBinary(a, 32)
	x := api.ToBinary(0, 32)
	split := 32 - rotation

	for i := 0; i < 32; i++ {
		if i >= split {
			x[i] = bits[i-split]
		} else {
			x[i] = bits[i+rotation]
		}
	}

	return x
}

func padding(api frontend.API, a []frontend.Variable) []frontend.Variable {

	// helpers
	inputLen := len(a)
	paddingLen := inputLen % 64

	// t is start index of intputBitLen encoding
	var t int
	if inputLen%64 < 56 {
		t = 56 - inputLen%64
	} else {
		t = 64 + 56 - inputLen%64
	}

	// total length of padded input
	totalLen := inputLen + t + 8

	// encode every byte in frontend.Variable
	out := make([]frontend.Variable, totalLen)

	// return if no padding required
	if paddingLen == 0 {

		// overwrite into fixed size slice
		for i := 0; i < inputLen; i++ {
			out[i] = a[i]
		}
		return out
	}

	// existing bytes into out
	for i := 0; i < inputLen; i++ {
		out[i] = a[i]
	}

	// padding, first byte is always a 128=2^7=10000000
	out[inputLen] = frontend.Variable(128) // api.FromBinary(0, 0, 0, 0, 0, 0, 0, 1) // input as little endian

	// zero padding
	for i := 0; i < t; i++ {
		out[inputLen+1+i] = frontend.Variable(0)
	}

	// bit size of number of input bytes
	inputBitLen := inputLen << 3
	// value, _ := api.Compiler().ConstantValue(inputBitLen) // inputBitLen = inputLen * 8#

	// fill last 8 byte in reverse because of little endian
	bits := api.ToBinary(inputBitLen, 64) // 64 bit = 8 byte
	ctr := inputLen + t
	for i := 7; i >= 0; i-- {
		start := i * 8
		out[ctr] = api.FromBinary(bits[start : start+8]...)
		ctr += 1
	}

	return out
}

// main function of program
func main() {

	// 'hello world' as hex
	input := "68656c6c6f20776f726c64"
	output := "b94d27b9934d3e08a52e52d7da7dabfac484efe37a5380ee9088f7ace2efcde9"

	// 'hello-world-hello-world-hello-world-hello-world-hello-world-12345' as hex
	// input := "68656c6c6f2d776f726c642d68656c6c6f2d776f726c642d68656c6c6f2d776f726c642d68656c6c6f2d776f726c642d68656c6c6f2d776f726c642d3132333435"
	// output := "34caf9dcd6b137c56c59f81e071a4b77a11329f26c80d7023ac7dfc485dcd780"

	byteSlice, _ := hex.DecodeString(input)
	inputByteLen := len(byteSlice)

	byteSlice, _ = hex.DecodeString(output)
	outputByteLen := len(byteSlice)

	// witness definition
	preImageAssign := strToIntSlice(input, true)
	outputAssign := strToIntSlice(output, true)

	// witness values preparation
	assignment := SHA256{
		PreImage: make([]frontend.Variable, inputByteLen),
		Output:   [32]frontend.Variable{},
	}

	// assign values here because required to use make in assignment
	for i := 0; i < inputByteLen; i++ {
		assignment.PreImage[i] = preImageAssign[i]
	}
	for i := 0; i < outputByteLen; i++ {
		assignment.Output[i] = outputAssign[i]
	}

	witness, err := frontend.NewWitness(&assignment, ecc.BN254)
	if err != nil {
		log.Fatal("witness creation failed")
	}
	publicWitness, _ := witness.Public()

	// var circuit SHA256
	circuit := SHA256{
		PreImage: make([]frontend.Variable, inputByteLen),
	}

	// generate CompiledConstraintSystem
	ccs, err := frontend.Compile(ecc.BN254, r1cs.NewBuilder, &circuit)
	if err != nil {
		log.Fatal("frontend.Compile")
	}

	// groth16 zkSNARK: Setup
	pk, vk, err := groth16.Setup(ccs)
	if err != nil {
		log.Fatal("groth16.Setup")
	}

	// groth16: Prove & Verify
	proof, err := groth16.Prove(ccs, pk, witness)
	if err != nil {
		log.Fatal("prove computation failed...")
	}
	err = groth16.Verify(proof, vk, publicWitness)
	if err != nil {
		log.Fatal("groth16 verify failed...")
	}

}

func strToIntSlice(inputData string, hexRepresentation bool) []int {

	// check if inputData in hex representation
	var byteSlice []byte
	if hexRepresentation {
		hexBytes, err := hex.DecodeString(inputData)
		if err != nil {
			log.Fatal("hex.DecodeString error.")
		}
		byteSlice = hexBytes
	} else {
		byteSlice = []byte(inputData)
	}

	// convert byte slice to int numbers which can be passed to gnark frontend.Variable
	var data []int
	for i := 0; i < len(byteSlice); i++ {
		data = append(data, int(byteSlice[i]))
	}

	return data
}

The script is working when you copy it into an empty folder and run go mod init, then go mod tidy, and when you change the go.mod file gnark and gnark-crypto versions as follows (go.mod):

module sha256

go 1.19

require (
	github.com/consensys/gnark v0.7.1
	github.com/consensys/gnark-crypto v0.7.0
)

require (
	github.com/fxamacker/cbor/v2 v2.2.0 // indirect
	github.com/mmcloughlin/addchain v0.4.0 // indirect
	github.com/rs/zerolog v1.26.1 // indirect
	github.com/x448/float16 v0.8.4 // indirect
	golang.org/x/crypto v0.0.0-20220321153916-2c7772ba3064 // indirect
	golang.org/x/sys v0.0.0-20220727055044-e65921a090b8 // indirect
)