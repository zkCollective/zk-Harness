package circuits

import (
	"encoding/hex"
	"fmt"
	"strconv"

	"github.com/consensys/gnark-crypto/ecc"

	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/backend/witness"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/math/emulated"
	"github.com/zkCollective/zk-Harness/gnark/circuits/groth16bls12377verifier"
	"github.com/zkCollective/zk-Harness/gnark/circuits/groth16bls24315verifier"
	"github.com/zkCollective/zk-Harness/gnark/circuits/prf/mimc"
	"github.com/zkCollective/zk-Harness/gnark/circuits/prf/sha256"
	"github.com/zkCollective/zk-Harness/gnark/circuits/toy/cubic"
	emulate "github.com/zkCollective/zk-Harness/gnark/circuits/toy/emulate"
	"github.com/zkCollective/zk-Harness/gnark/circuits/toy/exponentiate"
	"github.com/zkCollective/zk-Harness/gnark/circuits/toy/exponentiate_opt"
	"github.com/zkCollective/zk-Harness/gnark/util"
)

var err error

var BenchCircuits map[string]BenchCircuit

type BenchCircuit interface {
	Circuit(size int, name string, opts ...CircuitOption) frontend.Circuit
	Witness(size int, curveID ecc.ID, name string, opts ...WitnessOption) witness.Witness
}

func init() {
	BenchCircuits = make(map[string]BenchCircuit)

	// Exponentiate Circuit
	BenchCircuits["exponentiate"] = &defaultCircuit{}

	// Toy Circuits
	BenchCircuits["cubic"] = &defaultCircuit{}
	BenchCircuits["exponentiate_opt"] = &defaultCircuit{}
	BenchCircuits["emulate"] = &defaultCircuit{}

	// Hashes
	BenchCircuits["mimc"] = &defaultCircuit{}
	BenchCircuits["sha256"] = &defaultCircuit{}

	// Recursion
	BenchCircuits["groth16_bls12377"] = &defaultCircuit{}
	BenchCircuits["groth16_bls24315"] = &defaultCircuit{}
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
	if optCircuit.inputPath != "none" && optCircuit.inputPath != "" {
		var err error
		data, err = util.ReadFromInputPath(optCircuit.inputPath)
		if err != nil {
			panic(err)
		}
	}

	fmt.Printf(optCircuit.inputPath)

	switch name {
	case "exponentiate":
		return &exponentiate.ExponentiateCircuit{E: size}
	case "cubic":
		return &cubic.CubicCircuit{}
	case "exponentiate_opt":
		return &exponentiate_opt.ExponentiateOptCircuit{}
	case "emulate":
		return &emulate.Circuit{}
	case "mimc":
		return &mimc.MimcCircuit{}
	case "sha256":
		if data == nil || data["PreImage"] == nil {
			panic("Input for PreImage is not defined")
		}
		return &sha256.Sha256Circuit{
			In: make([]frontend.Variable, (len(data["PreImage"].(string)) / 2)),
		}
	case "groth16_bls12377":
		outerCircuit := groth16bls12377verifier.VerifierCircuit{}
		outerCircuit.InnerVk.Allocate(optCircuit.verifyingKey)
		return &outerCircuit
	case "groth16_bls24315":
		outerCircuit := groth16bls24315verifier.VerifierCircuit{}
		outerCircuit.InnerVk.Allocate(optCircuit.verifyingKey)
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
	if optWitness.inputPath != "none" && optWitness.inputPath != "" {
		data, err = util.ReadFromInputPath(optWitness.inputPath)
		if err != nil {
			panic(err)
		}
	}

	switch name {
	case "exponentiate":
		witness := exponentiate.ExponentiateCircuit{}
		witness.X = (data["X"].(string))
		witness.Y = (data["Y"].(string))
		strVal, ok := data["E"].(string)
		if !ok {
			panic("E is not a string")
		}
		witness.E, err = strconv.Atoi(strVal)
		if err != nil {
			panic("Error converting exponent to int")
		}
		w, err := frontend.NewWitness(&witness, curveID.ScalarField())
		if err != nil {
			panic(err)
		}
		return w
	case "cubic":
		witness := cubic.CubicCircuit{}
		witness.X = (data["X"].(string))
		witness.Y = (data["Y"].(string))
		w, err := frontend.NewWitness(&witness, curveID.ScalarField())
		if err != nil {
			panic(err)
		}
		return w
	case "exponentiate_opt":
		witness := exponentiate_opt.ExponentiateOptCircuit{}
		witness.X = (2)
		witness.E = (12)
		witness.Y = (4096)

		w, err := frontend.NewWitness(&witness, curveID.ScalarField())
		if err != nil {
			panic(err)
		}
		return w
	case "emulate":
		witness := emulate.Circuit{}
		witness.X = emulated.ValueOf[emulated.Secp256k1Fp](data["X"].(string))
		witness.Y = emulated.ValueOf[emulated.Secp256k1Fp](data["Y"].(string))
		witness.Res = emulated.ValueOf[emulated.Secp256k1Fp](data["Res"].(string))

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
		preImageAssign := sha256.StrToByteSlice(input, true)
		outputAssign := sha256.StrToByteSlice(output, true)

		// witness values preparation
		witness := sha256.Sha256Circuit{
			In:             make([]frontend.Variable, inputByteLen),
			ExpectedResult: [32]frontend.Variable{},
		}

		// assign values here because required to use make in assignment
		for i := 0; i < inputByteLen; i++ {
			witness.In[i] = preImageAssign[i]
		}
		for i := 0; i < outputByteLen; i++ {
			witness.ExpectedResult[i] = outputAssign[i]
		}

		w, err := frontend.NewWitness(&witness, curveID.ScalarField())
		if err != nil {
			panic(err)
		}
		return w
	case "groth16_bls12377":
		var outerAssignment groth16bls12377verifier.VerifierCircuit
		outerAssignment.InnerProof.Assign(optWitness.proof)
		outerAssignment.InnerVk.Assign(optWitness.verifyingKey)
		outerAssignment.Witness = optWitness.witness
		w, err := frontend.NewWitness(&outerAssignment, ecc.BW6_761.ScalarField())
		if err != nil {
			panic(err)
		}
		return w
	case "groth16_bls24315":
		var outerAssignment groth16bls24315verifier.VerifierCircuit
		outerAssignment.InnerProof.Assign(optWitness.proof)
		outerAssignment.InnerVk.Assign(optWitness.verifyingKey)
		outerAssignment.Witness = optWitness.witness
		w, err := frontend.NewWitness(&outerAssignment, ecc.BW6_633.ScalarField())
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
