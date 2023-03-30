package circuits

import (
	"encoding/hex"

	"github.com/consensys/gnark-crypto/ecc"
	bls12377fr "github.com/consensys/gnark-crypto/ecc/bls12-377/fr"
	"github.com/consensys/gnark-crypto/hash"

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
	Circuit(size int, name string, opts ...CircuitConfig) frontend.Circuit
	Witness(size int, curveID ecc.ID, name string, path string) witness.Witness
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

func preCalcMIMC(curveID ecc.ID, preImage frontend.Variable) interface{} {

	switch curveID {
	case ecc.BN254:
		// compute expected Y
		var expectedY bn254fr.Element
		expectedY.SetInterface(preImage)

		// running MiMC (Go)
		goMimc := hash.MIMC_BN254.New()
		goMimc.Write(expectedY.Marshal())
		expectedh := goMimc.Sum(nil)
		return expectedh

	case ecc.BLS12_377:
		// compute expected Y
		var expectedY bls12377fr.Element
		expectedY.SetInterface(preImage)

		// running MiMC (Go)
		goMimc := hash.MIMC_BLS12_377.New()
		goMimc.Write(expectedY.Marshal())
		expectedh := goMimc.Sum(nil)
		return expectedh

	case ecc.BLS24_315:
		// compute expected Y
		var expectedY bls24315fr.Element
		expectedY.SetInterface(preImage)

		// running MiMC (Go)
		goMimc := hash.MIMC_BLS24_315.New()
		goMimc.Write(expectedY.Marshal())
		expectedh := goMimc.Sum(nil)
		return expectedh

	case ecc.BW6_761:
		// compute expected Y
		var expectedY bw6761fr.Element
		expectedY.SetInterface(preImage)

		// running MiMC (Go)
		goMimc := hash.MIMC_BW6_761.New()
		goMimc.Write(expectedY.Marshal())
		expectedh := goMimc.Sum(nil)
		return expectedh

	case ecc.BW6_633:
		// compute expected Y
		var expectedY bw6633fr.Element
		expectedY.SetInterface(preImage)

		// running MiMC (Go)
		goMimc := hash.MIMC_BW6_633.New()
		goMimc.Write(expectedY.Marshal())
		expectedh := goMimc.Sum(nil)
		return expectedh
	default:
		panic("not implemented")
	}
}

type defaultCircuit struct {
}

func (d *defaultCircuit) Circuit(size int, name string, opts ...CircuitOption) frontend.Circuit {

	// Parse Options for input Path
	opt := CircuitConfig{}
	for _, o := range opts {
		if err := o(&opt); err != nil {
			panic(err)
		}
	}

	var data map[string]interface{}
	if opt.inputPath != "" {
		var err error
		data, err = util.ReadFromInputPath(opt.inputPath)
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
		return &groth16bls12377verifier.VerifierCircuit{}
	default:
		panic("not implemented")
	}
}

func (d *defaultCircuit) Witness(size int, curveID ecc.ID, name string, opts ...WitnessOption) witness.Witness {

	// Parse Options for input Path
	opt := WitnessConfig{}
	for _, o := range opts {
		if err := o(&opt); err != nil {
			panic(err)
		}
	}

	var data map[string]interface{}
	if opt.inputPath != "" {
		var err error
		data, err = util.ReadFromInputPath(opt.inputPath)
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
		// witness.PreImage = ("16130099170765464552823636852555369511329944820189892919423002775646948828469")
		witness.PreImage = (data["PreImage"].(string))
		witness.Hash = preCalcMIMC(curveID, witness.PreImage)

		w, err := frontend.NewWitness(&witness, curveID.ScalarField())
		if err != nil {
			panic(err)
		}
		return w
	case "sha256":
		input := (data["PreImage"].(string))
		output := (data["Hash"].(string))

		// 'hello-world-hello-world-hello-world-hello-world-hello-world-12345' as hex
		// input := "68656c6c6f2d776f726c642d68656c6c6f2d776f726c642d68656c6c6f2d776f726c642d68656c6c6f2d776f726c642d68656c6c6f2d776f726c642d3132333435"
		// output := "34caf9dcd6b137c56c59f81e071a4b77a11329f26c80d7023ac7dfc485dcd780"

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
		var outerCircuit groth16verifier.VerifierCircuit
		outerCircuit.InnerVk.FillG1K(opt.verifyingKey)

		var outerWitness groth16verifier.VerifierCircuit
		outerWitness.InnerProof.Assign(opt.proof)
		outerWitness.InnerVk.Assign(opt.verifyingKey)
		// TODO - Make variable for arbitrary circuits.
		outerWitness.Hash = preCalcMIMC(curveID, opt.witness)

		w, err := frontend.NewWitness(&outerWitness, ecc.BW6_633.ScalarField())
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
	inputPath string
}

// Optionally provide input path to Circuit
func WithInputCircuit(inputPath string) CircuitOption {
	return func(opt *CircuitConfig) error {
		opt.inputPath = inputPath
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

// Optionally provide input path to Witness def
func WithInputConfig(inputPath string) WitnessOption {
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
