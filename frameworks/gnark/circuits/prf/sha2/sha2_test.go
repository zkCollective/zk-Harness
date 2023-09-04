package sha2

import (
	"crypto/sha256"
	"testing"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/test"
	"github.com/consensys/gnark/std/math/uints"
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