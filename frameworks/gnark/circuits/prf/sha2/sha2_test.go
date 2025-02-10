package sha2

import (
	"crypto/sha256"
	"fmt"
	"testing"
	"encoding/hex"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/test"
	"github.com/consensys/gnark/std/math/uints"
	"github.com/consensys/gnark/backend"
)

func TestSHA2(t *testing.T) {
	bts := make([]byte, 310)
	dgst := sha256.Sum256(bts)
	witness := Sha2Circuit{
		In: uints.NewU8Array(bts),
	}
	copy(witness.Expected[:], uints.NewU8Array(dgst[:]))
	err := test.IsSolved(&Sha2Circuit{In: make([]uints.U8, len(bts))}, &witness, ecc.BN254.ScalarField())
	if err != nil {
		t.Fatal(err)
	}
}

func TestPreimage(t *testing.T) {
	assert := test.NewAssert(t)

	input := "68656c6c6f20776f726c64"
	// output := "b94d27b9934d3e08a52e52d7da7dabfac484efe37a5380ee9088f7ace2efcde9"

	bts, err := hex.DecodeString(input)
	// bts := make([]byte, 1)
	if err != nil {
		fmt.Println("Error decoding hex:", err)
		return
	}
	dgst := sha256.Sum256(bts)
	randomDigest := make([]byte, 32)

	witness := Sha2Circuit{
		In: uints.NewU8Array(bts),
	}
	copy(witness.Expected[:], uints.NewU8Array(dgst[:]))

	assert.ProverSucceeded(&Sha2Circuit{In: make([]uints.U8, len(bts))}, &witness, test.WithCurves(ecc.BN254), test.WithBackends(backend.GROTH16))

	wrongWitness := Sha2Circuit{
		In: uints.NewU8Array(bts),
	}
	copy(wrongWitness.Expected[:], uints.NewU8Array(randomDigest[:]))

	assert.ProverFailed(&Sha2Circuit{In: make([]uints.U8, len(bts))}, &wrongWitness, test.WithCurves(ecc.BN254), test.WithBackends(backend.GROTH16))
}