/*
Benchmarking Recursion for 1-Step, currently over 2-chains of elliptic curves
*/
package cmd

import (
	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
	gnark_r1cs "github.com/consensys/gnark/frontend/cs/r1cs"
	"github.com/consensys/gnark/std/groth16_bls12377"
	"github.com/spf13/cobra"
	"github.com/tumberger/zk-compilers/gnark/circuits"
)

// groth16Cmd represents the groth16 command
var oneStepCommand = &cobra.Command{
	Use:   "oneStep",
	Short: "runs benchmarks and profiles using with one step of recursion",
	Run:   runOneStep,
}

func runOneStep(cmd *cobra.Command, args []string) {

	// Set filename

	// Parse flags

	// Write Results

	// Timing Function

	// Get inner proof circuit -> From ./circuits
	// create a mock cs: knowing the preimage of a hash using mimc
	// var circuit mimcCircuit

	// Compile inner proof circuit -> PASS INNER CURVE
	r1cs, err := frontend.Compile(ecc.BLS12_377.ScalarField(), gnark_r1cs.NewBuilder, c.Circuit(*fCircuitSize, *fCircuit, *fInputPath))
	if err != nil {
		panic("Circuit doesn't compile")
	}

	// Do Setup --> Either Groth16 / Plonk / PlonkFRI -> PASS INNER CURVE
	// Gets the initial verifier key
	pk, vk, err := groth16.Setup(r1cs)
	assertNoError(err)

	// Generate inner proof witness --> CASE Groth16 / Plonk / PlonkFRI
	witness := c.Witness(*fCircuitSize, curveID, *fCircuit, *fInputPath)

	// Generate inner proof --> CASE Groth16 / Plonk / PlonkFRI
	proof, err := groth16.Prove(r1cs, pk, witness)
	assertNoError(err)

	// Get public Witness of inner proof
	publicWitness, err := witness.Public()
	// Check whether the computed proof verifies that the proof passes on bls12377
	if err := groth16.Verify(proof, vk, publicWitness); err != nil {
		panic("Computed Proof doesn't verify!‚")
	}

	// Create dummy recursion circuit without passing argument
	c_rec := circuits.BenchCircuits["recursion"]

	// Assign verifier key

	// Compile outer proof circuit -> PASS OUTER CURVE
	r1cs_outer, err := frontend.Compile(ecc.BW6_633.ScalarField(), gnark_r1cs.NewBuilder, c_rec.Circuit(*fCircuitSize, *fCircuit, *fInputPath))
	if err != nil {
		panic("Circuit doesn't compile")
	}

	// Do Setup --> Either Groth16 / Plonk / PlonkFRI -> PASS INNER CURVE
	// Gets the initial verifier key
	// opk, ovk, err := groth16.Setup(r1cs_outer)
	// assertNoError(err)

	// Assign the verifier Key
	ovk := groth16_bls12377.VerifyingKey{}
	ovk.Assign(vk)

	// Time verification of inner proof on outer curve (Run verifier algorithm of one proof inside the other)
	// groth16_bls12377.Verify(nil, circuit.InnerVk, circuit.InnerProof, []frontend.Variable{circuit.Hash})
	groth16_bls12377.Verify(proof, vk, publicWitness.Vector().(fr.Vector))
}
