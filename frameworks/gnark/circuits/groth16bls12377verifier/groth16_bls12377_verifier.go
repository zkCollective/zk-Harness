package groth16bls12377verifier

import (
	"github.com/consensys/gnark/frontend"
	groth16_bls12377 "github.com/consensys/gnark/std/groth16_bls12377"
)

type VerifierCircuit struct {
	InnerProof groth16_bls12377.Proof
	InnerVk    groth16_bls12377.VerifyingKey
	Witness    frontend.Variable
}

func (circuit *VerifierCircuit) Define(api frontend.API) error {
	// create the verifier cs
	groth16_bls12377.Verify(api, circuit.InnerVk, circuit.InnerProof, []frontend.Variable{circuit.Witness})
	return nil
}
