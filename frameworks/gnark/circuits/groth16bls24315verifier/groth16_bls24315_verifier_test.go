package groth16bls24315verifier

import (
	"testing"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
	"github.com/consensys/gnark/test"
	mimc "github.com/zkCollective/zk-Harness/gnark/circuits/prf/mimc"
	"github.com/zkCollective/zk-Harness/gnark/util"
)

const (
	preImage   = "4992816046196248432836492760315135318126925090839638585255611512962528270024"
	publicHash = "4875439939758844840941638351757981379945701574516438614845550995673793857363"
)

// FIXME - Currently only works with a single frontend.Variable
// E.g. SHA-256 uses [32]frontend.Variable
// Other circuits have >1 public input (exponentiate)
func TestRecursion(t *testing.T) {

	// create a mock cs: knowing the preimage of a hash using mimc
	var circuit mimc.MimcCircuit
	r1cs, err := frontend.Compile(ecc.BLS24_315.ScalarField(), r1cs.NewBuilder, &circuit)
	if err != nil {
		t.Fatal(err)
	}

	// build the witness
	var assignment mimc.MimcCircuit
	assignment.PreImage = preImage
	assignment.Hash = publicHash

	witness, err := frontend.NewWitness(&assignment, ecc.BLS24_315.ScalarField())
	if err != nil {
		t.Fatal(err)
	}

	innerPk, innerVk, err := groth16.Setup(r1cs)
	if err != nil {
		t.Fatal(err)
	}

	proof, err := groth16.Prove(r1cs, innerPk, witness)
	if err != nil {
		t.Fatal(err)
	}

	publicWitness, err := witness.Public()
	if err != nil {
		t.Fatal(err)
	}

	// Check that proof verifies before continuing
	if err := groth16.Verify(proof, innerVk, publicWitness); err != nil {
		t.Fatal(err)
	}

	var outerCircuit VerifierCircuit
	outerCircuit.InnerVk.Allocate(innerVk)

	var outerAssignment VerifierCircuit
	outerAssignment.InnerProof.Assign(proof)
	outerAssignment.InnerVk.Assign(innerVk)
	outerAssignment.Witness = util.PreCalcMIMC(ecc.BLS24_315, preImage)

	assert := test.NewAssert(t)

	assert.ProverSucceeded(&outerCircuit, &outerAssignment, test.WithCurves(ecc.BW6_633), test.WithBackends(backend.GROTH16))
}
