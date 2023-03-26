package groth16verifier

import (
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
	groth16_bls12377 "github.com/consensys/gnark/std/groth16_bls12377"
)

type VerifierCircuit struct {
	InnerProof groth16_bls12377.Proof
	InnerVk    groth16_bls12377.VerifyingKey
	Hash       frontend.Variable
}

func (circuit *VerifierCircuit) Define(api frontend.API) error {
	// create the verifier cs
	groth16_bls12377.Verify(api, circuit.InnerVk, circuit.InnerProof, []frontend.Variable{circuit.Hash})
	return nil
}

// Verify implements the verification function of Groth16.
// Notation follows Figure 4. in DIZK paper https://eprint.iacr.org/2018/691.pdf
// publicInputs do NOT contain the ONE_WIRE
func (vk *VerifierCircuit) AssignVK(_ovk groth16.VerifyingKey) {
	vk.InnerVk.Assign(_ovk)
}
