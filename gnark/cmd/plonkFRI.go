package cmd

import (
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/DmitriyVTitov/size"
	"github.com/consensys/gnark/backend/plonkfri"
	"github.com/consensys/gnark/constraint"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/scs"
	"github.com/consensys/gnark/logger"
	"github.com/pkg/profile"
	"github.com/spf13/cobra"
	"github.com/zkCollective/zk-Harness/gnark/circuits"
	"github.com/zkCollective/zk-Harness/gnark/util"
)

// plonkCmd represents the plonk command
var plonkFRIcmd = &cobra.Command{
	Use:   "plonkFRI",
	Short: "runs benchmarks and profiles using PlonK proof system",
	Run:   runPlonkFRI,
}

func runPlonkFRI(plonkCmd *cobra.Command, args []string) {
	log := logger.Logger()
	log.Info().Msg("Benchmarking " + *fCircuit + " - gnark, plonk FRI: " + *fAlgo + " " + *fCurve + " " + *fInputPath)

	var filename = "../benchmarks/gnark/gnark_" +
		"plonkFRI" + "_" +
		*fCircuit + "." +
		*fFileType

	if err := parseFlags(); err != nil {
		fmt.Println("error: ", err.Error())
		plonkCmd.Help()
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
			Backend:           "plonk",
			Curve:             curveID.String(),
			Circuit:           *fCircuit,
			Input:             *fInputPath,
			Operation:         *fAlgo,
			NbConstraints:     ccs.GetNbConstraints(),
			NbSecretVariables: secret,
			NbPublicVariables: public,
			ProofSize:         proof_size,
			MaxRAM:            (m.Sys / 1024 / 1024),
			Count:             *fCount,
			RunTime:           took.Milliseconds(),
		}

		if err := util.WriteData("csv", bData, filename); err != nil {
			panic(err)
		}
	}

	// Run Benchmarks for Groth16 on given specification
	benchPlonkFRI(writeResults, *fAlgo, *fCount, *fCircuitSize, *fCircuit, util.WithInput(*fInputPath))
}

func benchPlonkFRI(fnWrite writeFunction, falgo string, fcount int, fcircuitSize int, fcircuit string, opts ...util.BenchOption) {
	fmt.Println("BENCHMARKING PLONK WITH FRI")

	// Parse Options, if no option is provided it runs plain G16 benches
	opt := util.BenchConfig{}
	for _, o := range opts {
		if err := o(&opt); err != nil {
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

	circuit := c.Circuit(fcircuitSize,
		fcircuit,
		circuits.WithInputCircuit(opt.InputPath),
		circuits.WithVKCircuit(opt.VerifyingKey))

	if *fAlgo == "compile" {
		startProfile()
		var err error
		var ccs constraint.ConstraintSystem
		for i := 0; i < *fCount; i++ {
			ccs, err = frontend.Compile(curveID.ScalarField(), scs.NewBuilder, circuit, frontend.WithCapacity(*fCircuitSize))
		}
		stopProfile()
		assertNoError(err)
		// Set compile time to 1 ms, otherwise 0 in frontend
		if took < (1024 * 1024) {
			took = (1024 * 1024)
		}
		fnWrite(took, ccs, 0)
		return
	}

	ccs, err := frontend.Compile(curveID.ScalarField(), scs.NewBuilder, circuit, frontend.WithCapacity(*fCircuitSize))
	assertNoError(err)

	if *fAlgo == "setup" {
		startProfile()
		var err error
		for i := 0; i < *fCount; i++ {
			_, _, err = plonkfri.Setup(ccs)
		}
		stopProfile()
		assertNoError(err)
		fnWrite(took, ccs, 0)
		return
	}

	if *fAlgo == "witness" {
		startProfile()
		var err error
		for i := 0; i < *fCount; i++ {
			c.Witness(
				fcircuitSize,
				curveID,
				fcircuit,
				circuits.WithInputWitness(opt.InputPath),
				circuits.WithVK(opt.VerifyingKey),
				circuits.WithProof(opt.Proof),
				circuits.WithWitness(opt.Witness))
		}
		stopProfile()
		assertNoError(err)
		// Set compile time to 1 ms, otherwise 0 in frontend
		if took < (1024 * 1024) {
			took = (1024 * 1024)
		}
		fnWrite(took, ccs, 0)
		return
	}

	validWitness := c.Witness(
		fcircuitSize,
		curveID,
		fcircuit,
		circuits.WithInputWitness(opt.InputPath),
		circuits.WithVK(opt.VerifyingKey),
		circuits.WithProof(opt.Proof),
		circuits.WithWitness(opt.Witness))

	pk, vk, err := plonkfri.Setup(ccs)
	assertNoError(err)

	if *fAlgo == "prove" {
		fmt.Println("BENCHMARK PROOF GENERATION PLONK FRI")
		var proof interface{}
		startProfile()
		for i := 0; i < *fCount; i++ {
			proof, err = plonkfri.Prove(ccs, pk, validWitness)
		}
		stopProfile()
		assertNoError(err)
		proof_size := size.Of(proof)
		fnWrite(took, ccs, proof_size)
		return
	}

	if *fAlgo != "verify" {
		panic("algo at this stage should be verify")
	}

	correctProof, err := plonkfri.Prove(ccs, pk, validWitness)
	assertNoError(err)

	validPublicWitness, err := validWitness.Public()
	assertNoError(err)

	fmt.Println("BENCHMARK PROOF VERIFICATION")
	startProfile()
	for i := 0; i < *fCount; i++ {
		err = plonkfri.Verify(correctProof, vk, validPublicWitness)
	}
	stopProfile()
	assertNoError(err)
	fnWrite(took, ccs, 0)
}

func init() {
	rootCmd.AddCommand(plonkFRIcmd)
}
