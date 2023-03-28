package groth16verifier

import (
	"testing"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	"github.com/consensys/gnark-crypto/hash"
	"github.com/consensys/gnark/backend"
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

// Calculate the expected output of MIMC through plain invocation
func preComputeMimc(preImage frontend.Variable) interface{} {
	var expectedY fr.Element
	expectedY.SetInterface(preImage)
	// calc MiMC
	goMimc := hash.MIMC_BLS12_377.New()
	goMimc.Write(expectedY.Marshal())
	expectedh := goMimc.Sum(nil)
	return expectedh
}

func TestRecursion(t *testing.T) {

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
	outerCircuit.InnerVk.FillG1K(innerVk)

	var outerWitness VerifierCircuit
	outerWitness.InnerProof.Assign(proof)
	outerWitness.InnerVk.Assign(innerVk)
	outerWitness.Hash = preComputeMimc(preImage)

	assert := test.NewAssert(t)

	assert.ProverSucceeded(&outerCircuit, &outerWitness, test.WithCurves(ecc.BW6_761), test.WithBackends(backend.GROTH16))
}
