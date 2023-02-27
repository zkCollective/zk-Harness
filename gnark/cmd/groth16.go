/*
Copyright © 2021 ConsenSys Software Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/DmitriyVTitov/size"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/constraint"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
	"github.com/consensys/gnark/logger"
	"github.com/pkg/profile"
	"github.com/spf13/cobra"
	"github.com/tumberger/zk-compilers/gnark/util"
)

// groth16Cmd represents the groth16 command
var groth16Cmd = &cobra.Command{
	Use:   "groth16",
	Short: "runs benchmarks and profiles using Groth16 proof system",
	Run:   runGroth16,
}

func runGroth16(cmd *cobra.Command, args []string) {

	log := logger.Logger()
	log.Info().Msg("Benchmarking " + *fCircuit + " - gnark, groth16: " + *fAlgo + " " + *fCurve + " " + *fInputPath)

	var filename = "../benchmarks/gnark/gnark_" +
		"groth16" + "_" +
		*fCircuit + "." +
		*fFileType

	if err := parseFlags(); err != nil {
		fmt.Println("error: ", err.Error())
		cmd.Help()
		os.Exit(-1)
	}

	writeResults := func(took time.Duration, ccs constraint.ConstraintSystem, proof_size int) {

		// check memory usage, max ram requested from OS
		var m runtime.MemStats
		runtime.ReadMemStats(&m)

		_, secret, public := ccs.GetNbVariables()
		bData := util.BenchDataCircuit{
			Framework:         "gnark",
			Category:          "circuit",
			Backend:           "groth16",
			Curve:             curveID.String(),
			Circuit:           *fCircuit,
			Input:             *fInputPath,
			Operation:         *fAlgo,
			NbConstraints:     ccs.GetNbConstraints(),
			NbSecretVariables: secret,
			NbPublicVariables: public,
			ProofSize:         proof_size,
			MaxRAM:            (m.Sys),
			RunTime:           took.Milliseconds(),
		}

		if err := util.WriteData("csv", bData, filename); err != nil {
			panic(err)
		}
	}

	var (
		start time.Time
		took  time.Duration
		prof  interface{ Stop() }
	)

	startProfile := func() {
		start = time.Now()
		if p != nil {
			prof = profile.Start(p, profile.ProfilePath("."), profile.NoShutdownHook)
		}
	}

	stopProfile := func() {
		took = time.Since(start)
		if p != nil {
			prof.Stop()
		}
		took /= time.Duration(*fCount)
	}

	if *fAlgo == "compile" {
		var err error
		var ccs constraint.ConstraintSystem
		startProfile()
		for i := 0; i < *fCount; i++ {
			ccs, err = frontend.Compile(curveID.ScalarField(), r1cs.NewBuilder, c.Circuit(*fCircuitSize, *fCircuit), frontend.WithCapacity(*fCircuitSize))
		}
		stopProfile()
		assertNoError(err)
		writeResults(took, ccs, 0)
		return
	}

	ccs, err := frontend.Compile(curveID.ScalarField(), r1cs.NewBuilder, c.Circuit(*fCircuitSize, *fCircuit), frontend.WithCapacity(*fCircuitSize))
	assertNoError(err)

	if *fAlgo == "setup" {
		startProfile()
		var err error
		for i := 0; i < *fCount; i++ {
			_, _, err = groth16.Setup(ccs)
		}
		stopProfile()
		assertNoError(err)
		writeResults(took, ccs, 0)
		return
	}

	if *fAlgo == "witness" {
		startProfile()
		var err error
		for i := 0; i < *fCount; i++ {
			c.Witness(*fCircuitSize, curveID, *fCircuit, *fInputPath)
		}
		stopProfile()
		assertNoError(err)
		writeResults(took, ccs, 0)
		return
	}

	witness := c.Witness(*fCircuitSize, curveID, *fCircuit, *fInputPath)

	if *fAlgo == "prove" {
		pk, err := groth16.DummySetup(ccs)
		assertNoError(err)

		var proof interface{}
		startProfile()
		for i := 0; i < *fCount; i++ {
			proof, err = groth16.Prove(ccs, pk, witness)
		}
		stopProfile()
		assertNoError(err)
		proof_size := size.Of(proof)
		writeResults(took, ccs, proof_size)
		return
	}

	if *fAlgo != "verify" {
		panic("algo at this stage should be verify")
	}
	pk, vk, err := groth16.Setup(ccs)
	assertNoError(err)

	proof, err := groth16.Prove(ccs, pk, witness)
	assertNoError(err)

	// print(proof_size)
	// writeResults(took, ccs, proof_size)

	publicWitness, err := witness.Public()
	assertNoError(err)
	startProfile()
	for i := 0; i < *fCount; i++ {
		err = groth16.Verify(proof, vk, publicWitness)
	}
	stopProfile()
	assertNoError(err)
	writeResults(took, ccs, 0)

}

func assertNoError(err error) {
	if err != nil {
		panic(err)
	}
}

func init() {
	rootCmd.AddCommand(groth16Cmd)
}
