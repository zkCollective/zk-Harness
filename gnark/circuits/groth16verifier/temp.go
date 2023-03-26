package groth16verifier

import (
	"testing"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
	"github.com/consensys/gnark/test"
	mimc "github.com/zkCollective/zk-Harness/gnark/circuits/prf/mimc"
)

const (
	preImage   = "4992816046196248432836492760315135318126925090839638585255611512962528270024"
	publicHash = "7831393781387060555412927989411398077996792073838215843928284475008119358174"
)

func TestRecursion(t *testing.T) {
	assert := test.NewAssert(t)

	proof := groth16.NewProof(ecc.BLS12_377)
	vk := groth16.NewVerifyingKey(ecc.BLS12_377)
	pk := groth16.NewProvingKey(ecc.BLS12_377)

	// create a mock cs: knowing the preimage of a hash using mimc
	var circuit mimc.MimcCircuit
	r1cs, err := frontend.Compile(ecc.BLS12_377.ScalarField(), r1cs.NewBuilder, &circuit)
	if err != nil {
		t.Fatal(err)
	}

	// build the witness
	var assignment mimc.MimcCircuit
	assignment.PreImage = preImage
	assignment.Hash = publicHash

	witness, err := frontend.NewWitness(&assignment, ecc.BLS12_377.ScalarField())
	if err != nil {
		t.Fatal(err)
	}

	publicWitness, err := witness.Public()
	if err != nil {
		t.Fatal(err)
	}

	// Do Setup --> Either Groth16 / Plonk / PlonkFRI -> PASS INNER CURVE
	// Gets the initial verifier key
	pk, vk, err = groth16.Setup(r1cs)
	if err != nil {
		t.Fatal(err)
	}

	// Generate inner proof --> CASE Groth16 / Plonk / PlonkFRI
	proof, err = groth16.Prove(r1cs, pk, witness)
	if err != nil {
		t.Fatal(err)
	}

	// Check whether the computed proof verifies that the proof passes on bls12377
	if err := groth16.Verify(proof, vk, publicWitness); err != nil {
		panic("Computed Proof doesn't verify!‚")
	}

	// get the data
	// var innerVk groth16_bls12377.VerifyingKey
	// var innerProof groth16_bls12377.Proof

	var groth16VerifierCircuit VerifierCircuit

	assert.ProverFailed(&groth16VerifierCircuit, &VerifierCircuit{
		InnerProof: proof,
		InnerVk:    vk,
		Hash:       "8674594860895598770446879254410848023850744751986836044725552747672873438975",
	})

	// assert.ProverSucceeded(&groth16VerifierCircuit, &VerifierCircuit{
	// 	PreImage: "16130099170765464552823636852555369511329944820189892919423002775646948828469",
	// 	Hash:     "8674594860895598770446879254410848023850744751986836044725552747672873438975",
	// }, test.WithCurves(ecc.BN254))

}
