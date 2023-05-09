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
	"github.com/consensys/gnark/backend/witness"
	"github.com/consensys/gnark/frontend"
	"github.com/tumberger/zk-compilers/gnark/circuits/prf/mimc"
	sha256 "github.com/tumberger/zk-compilers/gnark/circuits/prf/sha256"
	"github.com/tumberger/zk-compilers/gnark/circuits/toy/cubic"
	"github.com/tumberger/zk-compilers/gnark/circuits/toy/expo"
	"github.com/tumberger/zk-compilers/gnark/circuits/toy/exponentiate"
	"github.com/tumberger/zk-compilers/gnark/util"
)

var BenchCircuits map[string]BenchCircuit

type BenchCircuit interface {
	Circuit(size int, name string, path string) frontend.Circuit
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

func (d *defaultCircuit) Circuit(size int, name string, path string) frontend.Circuit {

	data, err := util.ReadFromInputPath(path)
	if err != nil {
		panic(err)
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
		return &sha256.Sha256Circuit{
			In: make([]frontend.Variable, (len(data["PreImage"].(string)) / 2)),
		}
	default:
		panic("not implemented")
	}
}

func (d *defaultCircuit) Witness(size int, curveID ecc.ID, name string, path string) witness.Witness {

	data, err := util.ReadFromInputPath(path)
	if err != nil {
		panic(err)
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

		// Needed for variable input!
		// circuit := sha256.Sha256Circuit{
		// 	PreImage: make([]frontend.Variable, inputByteLen),
		// }

		w, err := frontend.NewWitness(&witness, curveID.ScalarField())
		if err != nil {
			panic(err)
		}
		return w
	default:
		panic("not implemented")
	}
}
