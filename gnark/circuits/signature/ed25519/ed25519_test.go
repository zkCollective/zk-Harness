package ed25519

import (
	"bytes"
	goEd25519 "crypto/ed25519"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"math/big"
	"testing"

	ed25519test "github.com/zkCollective/zk-Harness/gnark/circuits/signature/ed25519/test"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
	"github.com/consensys/gnark/test"
)

type Ed25519Circuit struct {
	Message    [5]frontend.Variable `gnark:",public"`
	PublicKey  [2]frontend.Variable `gnark:",public"`
	PbPoint    PublicKey
	Signatures Signature
}

func (c *Ed25519Circuit) Define(api frontend.API) error {
	ed25519, _ := NewEd25519(api)
	err := ed25519.Verify(c.PublicKey[:], c.Signatures, c.Message, c.PbPoint)
	return err
}

func getInputData() []byte {

	var inputStr = "79080211c7b933020000000022480a2055552f335dcaead25b5b7abaab85dcaf82e15ff729b7e97681da9997ba7784da12240a2057a8f6f0a1022ec94eda4eb3f0b8f2a975ad1667ecbc6709b3a661b7b82cf3a710012a0c0880fd99a00610f7d5d5dc01321442696e616e63652d436861696e2d47616e676573"
	input, _ := hex.DecodeString(inputStr)

	return input
}

func TestVerify(t *testing.T) {
	assert := test.NewAssert(t)

	pub, priv, err := goEd25519.GenerateKey(nil)
	fmt.Printf("pub:%x\n", pub)
	assert.NoError(err)

	A, err := (&ed25519test.Point{}).SetBytes(pub)
	assert.NoError(err)
	_A := (&ed25519test.Point{}).Negate(A)

	msg := getInputData()
	goSig := goEd25519.Sign(priv, msg)

	goSha512 := sha512.New()
	goSha512.Write(goSig[:32])
	goSha512.Write(pub)
	goSha512.Write(msg)
	hramDigest := goSha512.Sum(nil)

	//split 64 byte signature into two 32byte halves, first halve as point R, second half as S(integer)
	k, _ := ed25519test.NewScalar().SetUniformBytes(hramDigest)
	S, _ := ed25519test.NewScalar().SetCanonicalBytes(goSig[32:])
	R := (&ed25519test.Point{}).VarTimeDoubleScalarBaseMult(k, _A, S)

	var witness = &Ed25519Circuit{}
	eSig := &Signature{
		R: *NewEmulatedPoint(R),
		S: new(big.Int).SetBytes(PutBigEndian(S.Bytes())),
	}
	ePublicKey := &PublicKey{
		A: *NewEmulatedPoint(A),
	}
	witness.PbPoint = *ePublicKey
	witness.Signatures = *eSig

	var pubBytes [32]byte
	for i := 0; i < 32; i++ {
		pubBytes[i] = pub[i]
	}

	// compress 32byte public key to 2 frontend variable for bn254 field
	witness.PublicKey[0] = pubBytes[:16]
	witness.PublicKey[1] = pubBytes[16:]

	//122 byte -> 5
	var msgFv [5]frontend.Variable

	// compress message to 5 frontend variable for bn254 field
	for i := 0; i < 5; i++ {
		if i == 4 {
			msgFv[i] = msg[i*25:]
		} else {
			msgFv[i] = msg[i*25 : (i+1)*25]
		}
	}

	witness.Message = msgFv

	var result = bytes.Equal(goSig[:32], R.Bytes())
	assert.True(result)

	err = test.IsSolved(&Ed25519Circuit{}, witness, ecc.BN254.ScalarField())
	assert.NoError(err)

	ccs, err := frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, &Ed25519Circuit{})
	assert.NoError(err)
	fmt.Println(ccs.GetNbConstraints())
}
