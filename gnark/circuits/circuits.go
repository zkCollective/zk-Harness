package circuits

import (
	"encoding/hex"

	"github.com/consensys/gnark-crypto/ecc"
	bls12377fr "github.com/consensys/gnark-crypto/ecc/bls12-377/fr"

	bls12381fr "github.com/consensys/gnark-crypto/ecc/bls12-381/fr"
	bls24315fr "github.com/consensys/gnark-crypto/ecc/bls24-315/fr"
	bn254fr "github.com/consensys/gnark-crypto/ecc/bn254/fr"
	bw6633fr "github.com/consensys/gnark-crypto/ecc/bw6-633/fr"
	bw6761fr "github.com/consensys/gnark-crypto/ecc/bw6-761/fr"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/backend/witness"
	"github.com/consensys/gnark/frontend"
	"github.com/zkCollective/zk-Harness/gnark/circuits/groth16bls12377verifier"
	groth16verifier "github.com/zkCollective/zk-Harness/gnark/circuits/groth16bls12377verifier"
	"github.com/zkCollective/zk-Harness/gnark/circuits/prf/mimc"
	"github.com/zkCollective/zk-Harness/gnark/circuits/prf/sha256"
	"github.com/zkCollective/zk-Harness/gnark/circuits/toy/cubic"
	"github.com/zkCollective/zk-Harness/gnark/circuits/toy/expo"
	"github.com/zkCollective/zk-Harness/gnark/circuits/toy/exponentiate"
	"github.com/zkCollective/zk-Harness/gnark/util"
)

var BenchCircuits map[string]BenchCircuit

type BenchCircuit interface {
	Circuit(size int, name string, opts ...CircuitOption) frontend.Circuit
	Witness(size int, curveID ecc.ID, name string, opts ...WitnessOption) witness.Witness
}

func init() {
	BenchCircuits = make(map[string]BenchCircuit)

	// Toy Circuits
	BenchCircuits["cubic"] = &defaultCircuit{}
	BenchCircuits["expo"] = &defaultCircuit{}
	BenchCircuits["exponentiate"] = &defaultCircuit{}

	// Hashes
	BenchCircuits["mimc"] = &defaultCircuit{}
	BenchCircuits["sha256"] = &defaultCircuit{}

	// Recursion
	BenchCircuits["groth16_bls12377"] = &defaultCircuit{}
}

func preCalc(size int, curveID ecc.ID) interface{} {
	switch curveID {
	case ecc.BN254:
		// compute expected Y
		var expectedY bn254fr.Element
		expectedY.SetInterface(2)
		for i := 0; i < size; i++ {
			expectedY.Mul(&expectedY, &expectedY)
		}
		return expectedY
	case ecc.BLS12_381:
		// compute expected Y
		var expectedY bls12381fr.Element
		expectedY.SetInterface(2)
		for i := 0; i < size; i++ {
			expectedY.Mul(&expectedY, &expectedY)
		}

		return expectedY
	case ecc.BLS12_377:
		// compute expected Y
		var expectedY bls12377fr.Element
		expectedY.SetInterface(2)
		for i := 0; i < size; i++ {
			expectedY.Mul(&expectedY, &expectedY)
		}

		return expectedY
	case ecc.BLS24_315:
		// compute expected Y
		var expectedY bls24315fr.Element
		expectedY.SetInterface(2)
		for i := 0; i < size; i++ {
			expectedY.Mul(&expectedY, &expectedY)
		}

		return expectedY
	case ecc.BW6_761:
		// compute expected Y
		var expectedY bw6761fr.Element
		expectedY.SetInterface(2)
		for i := 0; i < size; i++ {
			expectedY.Mul(&expectedY, &expectedY)
		}

		return expectedY
	case ecc.BW6_633:
		// compute expected Y
		var expectedY bw6633fr.Element
		expectedY.SetInterface(2)
		for i := 0; i < size; i++ {
			expectedY.Mul(&expectedY, &expectedY)
		}

		return expectedY
	default:
		panic("not implemented")
	}
}

type defaultCircuit struct {
}

func (d *defaultCircuit) Circuit(size int, name string, opts ...CircuitOption) frontend.Circuit {

	// Parse Options for input Path
	optCircuit := CircuitConfig{}
	for _, o := range opts {
		if err := o(&optCircuit); err != nil {
			panic(err)
		}
	}

	var data map[string]interface{}
	if optCircuit.inputPath != "" {
		var err error
		data, err = util.ReadFromInputPath(optCircuit.inputPath)
		if err != nil {
			panic(err)
		}
	}

	switch name {
	case "cubic":
		return &cubic.CubicCircuit{}
	case "expo":
		return &expo.BenchCircuit{N: size}
	case "exponentiate":
		return &exponentiate.ExponentiateCircuit{}
	case "mimc":
		return &mimc.MimcCircuit{}
	case "sha256":
		if data == nil || data["PreImage"] == nil {
			panic("Input for PreImage is not defined")
		}
		return &sha256.Sha256Circuit{
			PreImage: make([]frontend.Variable, (len(data["PreImage"].(string)) / 2)),
		}
	case "groth16_bls12377":
		outerCircuit := groth16bls12377verifier.VerifierCircuit{}
		outerCircuit.InnerVk.FillG1K(optCircuit.verifyingKey)
		return &outerCircuit
	default:
		panic("not implemented")
	}
}

func (d *defaultCircuit) Witness(size int, curveID ecc.ID, name string, opts ...WitnessOption) witness.Witness {

	// Parse Options for input Path
	optWitness := WitnessConfig{}
	for _, o := range opts {
		if err := o(&optWitness); err != nil {
			panic(err)
		}
	}

	var data map[string]interface{}
	if optWitness.inputPath != "" {
		var err error
		data, err = util.ReadFromInputPath(optWitness.inputPath)
		if err != nil {
			panic(err)
		}
	}

	switch name {
	case "cubic":
		witness := cubic.CubicCircuit{}
		witness.X = (data["X"].(string))
		witness.Y = (data["Y"].(string))

		w, err := frontend.NewWitness(&witness, curveID.ScalarField())
		if err != nil {
			panic(err)
		}
		return w
	case "expo":
		witness := expo.BenchCircuit{N: size}
		witness.X = (2)
		witness.Y = preCalc(size, curveID)

		w, err := frontend.NewWitness(&witness, curveID.ScalarField())
		if err != nil {
			panic(err)
		}
		return w
	case "exponentiate":
		witness := exponentiate.ExponentiateCircuit{}
		witness.X = (2)
		witness.E = (12)
		witness.Y = (4096)

		w, err := frontend.NewWitness(&witness, curveID.ScalarField())
		if err != nil {
			panic(err)
		}
		return w
	case "mimc":
		witness := mimc.MimcCircuit{}
		witness.PreImage = (data["PreImage"].(string))
		witness.Hash = util.PreCalcMIMC(curveID, witness.PreImage)

		w, err := frontend.NewWitness(&witness, curveID.ScalarField())
		if err != nil {
			panic(err)
		}
		return w
	case "sha256":
		input := (data["PreImage"].(string))
		output := (data["Hash"].(string))

		byteSlice, _ := hex.DecodeString(input)
		inputByteLen := len(byteSlice)

		byteSlice, _ = hex.DecodeString(output)
		outputByteLen := len(byteSlice)

		// witness definition
		preImageAssign := sha256.StrToIntSlice(input, true)
		outputAssign := sha256.StrToIntSlice(output, true)

		// witness values preparation
		witness := sha256.Sha256Circuit{
			PreImage:       make([]frontend.Variable, inputByteLen),
			ExpectedResult: [32]frontend.Variable{},
		}

		// assign values here because required to use make in assignment
		for i := 0; i < inputByteLen; i++ {
			witness.PreImage[i] = preImageAssign[i]
		}
		for i := 0; i < outputByteLen; i++ {
			witness.ExpectedResult[i] = outputAssign[i]
		}

		// Needed for variable input!
		// circuit := sha256.Sha256Circuit{
		// 	PreImage: make([]frontend.Variable, inputByteLen),
		// }

		w, err := frontend.NewWitness(&witness, curveID.ScalarField())
		if err != nil {
			panic(err)
		}
		return w
	case "groth16_bls12377":
		// Witness is already provided in this case (pre-computed proof)
		var outerAssignment groth16verifier.VerifierCircuit
		outerAssignment.InnerProof.Assign(optWitness.proof)
		outerAssignment.InnerVk.Assign(optWitness.verifyingKey)
		outerAssignment.Witness = optWitness.witness

		w, err := frontend.NewWitness(&outerAssignment, ecc.BW6_761.ScalarField())
		if err != nil {
			panic(err)
		}
		return w
	default:
		panic("not implemented")
	}
}

// Optional Parameters Circuit
type CircuitOption func(opt *CircuitConfig) error

type CircuitConfig struct {
	inputPath    string
	verifyingKey groth16.VerifyingKey
}

func WithInputCircuit(inputPath string) CircuitOption {
	return func(opt *CircuitConfig) error {
		opt.inputPath = inputPath
		return nil
	}
}

func WithVKCircuit(verifyingKey groth16.VerifyingKey) CircuitOption {
	return func(opt *CircuitConfig) error {
		opt.verifyingKey = verifyingKey
		return nil
	}
}

// Optional Parameters Witness
type WitnessOption func(opt *WitnessConfig) error

type WitnessConfig struct {
	inputPath    string
	proof        groth16.Proof
	verifyingKey groth16.VerifyingKey
	witness      frontend.Variable
}

func WithInputWitness(inputPath string) WitnessOption {
	return func(opt *WitnessConfig) error {
		opt.inputPath = inputPath
		return nil
	}
}

func WithProof(proof groth16.Proof) WitnessOption {
	return func(opt *WitnessConfig) error {
		opt.proof = proof
		return nil
	}
}

func WithVK(verifyingKey groth16.VerifyingKey) WitnessOption {
	return func(opt *WitnessConfig) error {
		opt.verifyingKey = verifyingKey
		return nil
	}
}

func WithWitness(witness frontend.Variable) WitnessOption {
	return func(opt *WitnessConfig) error {
		opt.witness = witness
		return nil
	}
}
