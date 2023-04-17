package gnark_sha256

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"runtime/debug"
	"time"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
	//cSha256 "github.com/tumberger/zk-compilers/gnarktest"
	cSha256 "gnark_sha256/sha256"
)

type Sha256Circuit struct {
	ExpectedResult [32]frontend.Variable `gnark:"data,public"`
	In             []frontend.Variable
}

func (circuit *Sha256Circuit) Define(api frontend.API) error {
	sha256 := cSha256.New(api)
	sha256.Write(circuit.In[:])
	result := sha256.Sum()
	for i := range result {
		api.AssertIsEqual(result[i], circuit.ExpectedResult[i])
	}
	return nil
}

func main() {

	isBenchSetup := flag.Bool("benchSetup", false, "set true only benchmark the setup proccess")
	inputLen := flag.Int("inputLen", 64, "")
	flag.Parse()
	fmt.Println("bench setup", *isBenchSetup)
	fmt.Printf("input length:%d\n", *inputLen)

	// 1byte
	var input = getInputData(*inputLen)
	shaOnGroth16(input, *isBenchSetup)

	// 2byte
	//input = getInputData(128)
	//shaGroth16(input)

}

// repeatN times "00"
func getInputData(N int) []byte {
	var base = "00"

	var inputStr = ""
	for i := 0; i < N; i++ {
		inputStr = inputStr + base
	}
	input, _ := hex.DecodeString(inputStr)
	return input
}

func shaOnGroth16(input []byte, isBenchSetup bool) {

	// witness values preparation
	assignment := Sha256Circuit{
		In:             make([]frontend.Variable, len(input)),
		ExpectedResult: [32]frontend.Variable{},
	}

	goSha256 := sha256.New()
	goSha256.Write(input[:])
	out := goSha256.Sum(nil)

	fmt.Printf("out>>%x\n", out)

	// assign values here because required to use make in assignment
	for i := 0; i < len(input); i++ {
		assignment.In[i] = input[i]
	}
	for i := 0; i < 32; i++ {
		assignment.ExpectedResult[i] = out[i]
	}

	witnessStart := time.Now()
	witness, err := frontend.NewWitness(&assignment, ecc.BN254.ScalarField())
	if err != nil {
		log.Fatal("witness creation failed")
	}
	publicWitness, _ := witness.Public()
	endWitness := time.Since(witnessStart)
	fmt.Printf("wintness generation time:%dms\n", endWitness/1000000)

	// var circuit SHA256
	circuit := Sha256Circuit{
		In: make([]frontend.Variable, len(input)),
	}

	// generate CompiledConstraintSystem
	ccs, err := frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, &circuit)
	if err != nil {
		log.Fatal("frontend.Compile")
	}

	// groth16 zkSNARK: Setup
	setupStart := time.Now()
	pk, vk, err := groth16.Setup(ccs)
	if err != nil {
		log.Fatal("groth16.Setup")
	}

	setUpDuration := time.Since(setupStart)
	fmt.Printf("setup duration:%ds\n", setUpDuration/1000000000)

	if isBenchSetup {
		fmt.Println("=============== end ===================")
		return
	}
	var pkBuf bytes.Buffer
	pkLen, _ := pk.WriteTo(&pkBuf)
	fmt.Printf("pk size:%d bytes\n", pkLen)

	// groth16: Prove & Verify
	proof, err := groth16.Prove(ccs, pk, witness)

	var buf bytes.Buffer

	n, _ := proof.WriteTo(&buf)

	fmt.Printf("proof size:%d bytes\n", n)

	if err != nil {
		debug.PrintStack()
		log.Fatal("prove computation failed...", err)
	}
	err = groth16.Verify(proof, vk, publicWitness)
	if err != nil {
		log.Fatal("groth16 verify failed...")
	}

	fmt.Println("=============== end ===================")

}
