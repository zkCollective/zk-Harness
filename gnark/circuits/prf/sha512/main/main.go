package main

import (
	"bytes"
	"crypto/sha512"
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"runtime/debug"
	"time"

	cSha512 "github.com/zkCollective/zk-Harness/gnark/circuits/prf/sha512"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
)

type sha512Circuit struct {
	ExpectedResult [64]frontend.Variable `gnark:"data,public"`
	In             []frontend.Variable
}

func (circuit *sha512Circuit) Define(api frontend.API) error {
	sha512 := cSha512.New(api)
	sha512.Write(circuit.In[:])
	result := sha512.Sum()
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
	shaGroth16(input, *isBenchSetup)

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

func shaGroth16(input []byte, isBenchSetup bool) {

	// witness values preparation
	assignment := sha512Circuit{
		In:             make([]frontend.Variable, len(input)),
		ExpectedResult: [64]frontend.Variable{},
	}

	goSha512 := sha512.New()
	goSha512.Write(input[:])
	out := goSha512.Sum(nil)

	fmt.Printf("out>>%x\n", out)

	// assign values here because required to use make in assignment
	for i := 0; i < len(input); i++ {
		assignment.In[i] = input[i]
	}
	for i := 0; i < 64; i++ {
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
	circuit := sha512Circuit{
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
