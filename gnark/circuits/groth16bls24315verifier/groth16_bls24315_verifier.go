package groth16bls24315verifier

import (
	"github.com/consensys/gnark/frontend"
	groth16_bls24315 "github.com/consensys/gnark/std/groth16_bls24315"
)

type VerifierCircuit struct {
	InnerProof groth16_bls24315.Proof
	InnerVk    groth16_bls24315.VerifyingKey
	Witness    frontend.Variable
}

func (circuit *VerifierCircuit) Define(api frontend.API) error {
	// create the verifier cs
	groth16_bls24315.Verify(api, circuit.InnerVk, circuit.InnerProof, []frontend.Variable{circuit.Witness})
	return nil
}
