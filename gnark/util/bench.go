package util

import (
	"github.com/consensys/gnark/constraint"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
)

// These Options are used for recursive Groth16 verifier
type BenchOption func(opt *BenchConfig) error

type BenchConfig struct {
	InputPath    string
	Proof        groth16.Proof
	VerifyingKey groth16.VerifyingKey
	Witness      frontend.Variable
	CCS          constraint.ConstraintSystem
	InnerCurve   ecc.ID
	OuterCurve   ecc.ID
}

// Optionally provide input path to Witness def
func WithInput(inputPath string) BenchOption {
	return func(opt *BenchConfig) error {
		opt.InputPath = inputPath
		return nil
	}
}

func WithProof(proof groth16.Proof) BenchOption {
	return func(opt *BenchConfig) error {
		opt.Proof = proof
		return nil
	}
}

func WithVK(verifyingKey groth16.VerifyingKey) BenchOption {
	return func(opt *BenchConfig) error {
		opt.VerifyingKey = verifyingKey
		return nil
	}
}

func WithWitness(witness frontend.Variable) BenchOption {
	return func(opt *BenchConfig) error {
		opt.Witness = witness
		return nil
	}
}

func WithInnerCCS(ccs constraint.ConstraintSystem) BenchOption {
	return func(opt *BenchConfig) error {
		opt.CCS = ccs
		return nil
	}
}

func WithInnerCurve(innerCurve ecc.ID) BenchOption {
	return func(opt *BenchConfig) error {
		opt.InnerCurve = innerCurve
		return nil
	}
}

func WithOuterCurve(outerCurve ecc.ID) BenchOption {
	return func(opt *BenchConfig) error {
		opt.OuterCurve = outerCurve
		return nil
	}
}
