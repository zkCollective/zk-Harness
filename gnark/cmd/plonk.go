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
	"strings"
	"time"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/plonk"
	"github.com/consensys/gnark/constraint"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/scs"
	"github.com/consensys/gnark/test"
	"github.com/pkg/profile"
	"github.com/spf13/cobra"
	"github.com/tumberger/zk-compilers/gnark/util"
)

// plonkCmd represents the plonk command
var plonkCmd = &cobra.Command{
	Use:   "plonk",
	Short: "runs benchmarks and profiles using PlonK proof system",
	Run:   runPlonk,
}

func runPlonk(cmd *cobra.Command, args []string) {

	var filename = "../benchmarks/gnark/gnark_" +
		"plonk" + "_" +
		*fCircuit + "." +
		*fFileType

	if err := parseFlags(); err != nil {
		fmt.Println("error: ", err.Error())
		cmd.Help()
		os.Exit(-1)
	}

	writeResults := func(took time.Duration, ccs constraint.ConstraintSystem) {
		// check memory usage, max ram requested from OS
		var m runtime.MemStats
		runtime.ReadMemStats(&m)

		_, secret, public := ccs.GetNbVariables()
		bData := util.BenchDataCircuit{
			Framework:         "gnark",
			Category:          "circuit",
			Backend:           "plonk",
			Curve:             curveID.String(),
			Circuit:           *fCircuit,
			Input:             *fInputPath,
			Operation:         *fAlgo,
			NbConstraints:     ccs.GetNbConstraints(),
			NbSecretVariables: secret,
			NbPublicVariables: public,
			MaxRAM:            (m.Sys / 1024 / 1024),
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
		startProfile()
		var err error
		var ccs constraint.ConstraintSystem
		for i := 0; i < *fCount; i++ {
			ccs, err = frontend.Compile(curveID.ScalarField(), scs.NewBuilder, c.Circuit(*fCircuitSize, *fCircuit), frontend.WithCapacity(*fCircuitSize))
		}
		stopProfile()
		assertNoError(err)
		writeResults(took, ccs)
		return
	}

	ccs, err := frontend.Compile(curveID.ScalarField(), scs.NewBuilder, c.Circuit(*fCircuitSize, *fCircuit), frontend.WithCapacity(*fCircuitSize))
	assertNoError(err)

	// create srs
	srs, err := test.NewKZGSRS(ccs)
	assertNoError(err)

	if *fAlgo == "setup" {
		startProfile()
		var err error
		for i := 0; i < *fCount; i++ {
			_, _, err = plonk.Setup(ccs, srs)
		}
		stopProfile()
		assertNoError(err)
		writeResults(took, ccs)
		return
	}

	witness := c.Witness(*fCircuitSize, curveID, *fCircuit, *fInputPath)
	pk, vk, err := plonk.Setup(ccs, srs)
	assertNoError(err)

	if *fAlgo == "prove" {

		startProfile()
		for i := 0; i < *fCount; i++ {
			_, err = plonk.Prove(ccs, pk, witness)
		}
		stopProfile()
		assertNoError(err)
		writeResults(took, ccs)
		return
	}

	if *fAlgo != "verify" {
		panic("algo at this stage should be verify")
	}

	proof, err := plonk.Prove(ccs, pk, witness)
	assertNoError(err)

	publicWitness, err := witness.Public()
	assertNoError(err)

	startProfile()
	for i := 0; i < *fCount; i++ {
		err = plonk.Verify(proof, vk, publicWitness)
	}
	stopProfile()
	assertNoError(err)
	writeResults(took, ccs)

}

func init() {
	// Here the commands for the "circuit" category with Plonk are defined

	_curves := ecc.Implemented()
	curves := make([]string, len(_curves))
	for i := 0; i < len(_curves); i++ {
		curves[i] = strings.ToLower(_curves[i].String())
	}

	fCircuit = plonkCmd.Flags().String("circuit", "expo", "name of the circuit to use")
	fCircuitSize = plonkCmd.Flags().Int("size", 10000, "size of the circuit, parameter to circuit constructor")
	fCount = plonkCmd.Flags().Int("count", 2, "bench count (time is averaged on number of executions)")
	fAlgo = plonkCmd.Flags().String("algo", "prove", "name of the algorithm to benchmark. must be compile, setup, prove or verify")
	fProfile = plonkCmd.Flags().String("profile", "none", "type of profile. must be none, trace, cpu or mem")
	fCurve = plonkCmd.Flags().String("curve", "bn254", "curve name. must be "+fmt.Sprint(curves))
	fFileType = plonkCmd.Flags().String("filetype", "csv", "Type of file to output for benchmarks")

	rootCmd.AddCommand(plonkCmd)
}
