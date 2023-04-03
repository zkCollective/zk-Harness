package main

import (
	goEd25519 "crypto/ed25519"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"math/big"

	ed25519 "github.com/zkCollective/zk-Harness/gnark/circuits/signature/ed25519"
	ed25519test "github.com/zkCollective/zk-Harness/gnark/circuits/signature/ed25519/test"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
	"github.com/rs/zerolog/log"
)

type Ed25519Circuit struct {
	Message    [5]frontend.Variable `gnark:",public"`
	PublicKey  [2]frontend.Variable `gnark:",public"`
	PbPoint    ed25519.PublicKey
	Signatures ed25519.Signature
}

func (c *Ed25519Circuit) Define(api frontend.API) error {
	ed25519, _ := ed25519.NewEd25519(api)
	err := ed25519.Verify(c.PublicKey[:], c.Signatures, c.Message, c.PbPoint)
	return err
}

func main() {
	pub, priv, _ := goEd25519.GenerateKey(nil)
	fmt.Printf("pub:%x\n", pub)

	A, _ := (&ed25519test.Point{}).SetBytes(pub)
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

	var assigment = &Ed25519Circuit{}
	eSig := &ed25519.Signature{
		R: *ed25519.NewEmulatedPoint(R),
		S: new(big.Int).SetBytes(ed25519.PutBigEndian(S.Bytes())),
	}
	ePublicKey := &ed25519.PublicKey{
		A: *ed25519.NewEmulatedPoint(A),
	}
	assigment.PbPoint = *ePublicKey
	assigment.Signatures = *eSig

	var pubBytes [32]byte
	for i := 0; i < 32; i++ {
		pubBytes[i] = pub[i]
	}

	// compress 32byte public key to 2 frontend variable for bn254 field
	assigment.PublicKey[0] = pubBytes[:16]
	assigment.PublicKey[1] = pubBytes[16:]

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

	assigment.Message = msgFv

	ccs, _ := frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, &Ed25519Circuit{})

	pk, vk, _ := groth16.Setup(ccs)

	witness, _ := frontend.NewWitness(assigment, ecc.BN254.ScalarField())

	proof, _ := groth16.Prove(ccs, pk, witness)

	publicWitness, _ := witness.Public()
	err := groth16.Verify(proof, vk, publicWitness)

	log.Err(err)

}

func getInputData() []byte {

	var inputStr = "79080211c7b933020000000022480a2055552f335dcaead25b5b7abaab85dcaf82e15ff729b7e97681da9997ba7784da12240a2057a8f6f0a1022ec94eda4eb3f0b8f2a975ad1667ecbc6709b3a661b7b82cf3a710012a0c0880fd99a00610f7d5d5dc01321442696e616e63652d436861696e2d47616e676573"
	input, _ := hex.DecodeString(inputStr)

	return input
}
