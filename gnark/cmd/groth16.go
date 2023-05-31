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
	"github.com/zkCollective/zk-Harness/gnark/parser"
	"github.com/zkCollective/zk-Harness/gnark/util"
)

// groth16Cmd represents the groth16 command
var groth16Cmd = &cobra.Command{
	Use:   "groth16",
	Short: "runs benchmarks and profiles using Groth16 proof system",
	Run:   runGroth16,
}

func runGroth16(cmd *cobra.Command, args []string) {

	log := logger.Logger()
	log.Info().Msg("Benchmarking " + *cfg.Circuit + " - gnark, groth16: " + *cfg.Algo + " " + *cfg.Curve + " " + *cfg.InputPath)

	var filename = "../benchmarks/gnark/gnark_" +
		"groth16" + "_" +
		*cfg.Circuit + "." +
		*cfg.FileType

	if err := parser.ParseFlags(cfg); err != nil {
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
			Curve:             parser.CurveID.String(),
			Circuit:           *cfg.Circuit,
			Input:             *cfg.InputPath,
			Operation:         *cfg.Algo,
			NbConstraints:     ccs.GetNbConstraints(),
			NbSecretVariables: secret,
			NbPublicVariables: public,
			ProofSize:         proof_size,
			MaxRAM:            (m.Sys / 1024 / 1024),
			Count:             *cfg.Count,
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
		if parser.P != nil {
			prof = profile.Start(parser.P, profile.ProfilePath("."), profile.NoShutdownHook)
		}
	}

	stopProfile := func() {
		took = time.Since(start)
		if parser.P != nil {
			prof.Stop()
		}
		took /= time.Duration(*cfg.Count)
	}

	if *cfg.Algo == "compile" {
		var err error
		var ccs constraint.ConstraintSystem
		startProfile()
		for i := 0; i < *cfg.Count; i++ {
			ccs, err = frontend.Compile(parser.CurveID.ScalarField(), r1cs.NewBuilder, parser.C.Circuit(*cfg.CircuitSize, *cfg.Circuit, *cfg.InputPath), frontend.WithCapacity(*cfg.CircuitSize))
		}
		stopProfile()
		assertNoError(err)
		// Set compile time to 1 ms, otherwise 0 in frontend
		if took < (1024 * 1024) {
			took = (1024 * 1024)
		}
		writeResults(took, ccs, 0)
		return
	}

	ccs, err := frontend.Compile(parser.CurveID.ScalarField(), r1cs.NewBuilder, parser.C.Circuit(*cfg.CircuitSize, *cfg.Circuit, *cfg.InputPath), frontend.WithCapacity(*cfg.CircuitSize))
	assertNoError(err)

	if *cfg.Algo == "setup" {
		startProfile()
		var err error
		for i := 0; i < *cfg.Count; i++ {
			_, _, err = groth16.Setup(ccs)
		}
		stopProfile()
		assertNoError(err)
		writeResults(took, ccs, 0)
		return
	}

	if *cfg.Algo == "witness" {
		startProfile()
		var err error
		for i := 0; i < *cfg.Count; i++ {
			parser.C.Witness(*cfg.CircuitSize, parser.CurveID, *cfg.Circuit, *cfg.InputPath)
		}
		stopProfile()
		assertNoError(err)
		// Set compile time to 1 ms, otherwise 0 in frontend
		if took < (1024 * 1024) {
			took = (1024 * 1024)
		}
		writeResults(took, ccs, 0)
		return
	}

	witness := parser.C.Witness(*cfg.CircuitSize, parser.CurveID, *cfg.Circuit, *cfg.InputPath)

	if *cfg.Algo == "prove" {
		pk, err := groth16.DummySetup(ccs)
		assertNoError(err)

		var proof interface{}
		startProfile()
		for i := 0; i < *cfg.Count; i++ {
			proof, err = groth16.Prove(ccs, pk, witness)
		}
		stopProfile()
		assertNoError(err)
		proof_size := size.Of(proof)
		writeResults(took, ccs, proof_size)
		return
	}

	if *cfg.Algo != "verify" {
		panic("algo at this stage should be verify")
	}
	pk, vk, err := groth16.Setup(ccs)
	assertNoError(err)

	proof, err := groth16.Prove(ccs, pk, witness)
	assertNoError(err)

	publicWitness, err := witness.Public()
	assertNoError(err)
	startProfile()
	for i := 0; i < *cfg.Count; i++ {
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
