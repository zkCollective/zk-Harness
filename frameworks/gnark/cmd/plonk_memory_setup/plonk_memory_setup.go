package main

import (
	"fmt"
	"os"
	"strings"
	"log"
	"bytes"
	"reflect"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/plonk"
	"github.com/consensys/gnark/test"
	"github.com/spf13/cobra"
	"github.com/zkCollective/zk-Harness/frameworks/gnark/parser"
)

var plonkMemorySetupCmd = &cobra.Command{
	Use:   "plonkMemorySetup",
	Short: "runs benchmarks and profiles using Plonk proof system",
	Run:   runPlonkMemorySetup,
}

var cfg = parser.NewConfig()

func runPlonkMemorySetup(cmd *cobra.Command, args []string) {

	if err := parser.ParseFlagsMemory(cfg); err != nil {
		fmt.Println("error: ", err.Error())
		cmd.Help()
		os.Exit(-1)
	}

	_ccs, err := os.ReadFile("tmp/ccs.dat")
	if err != nil {
		log.Fatal(err)
	}
	_buf := *bytes.NewBuffer(_ccs)
	reconstructedCCS := plonk.NewCS(parser.CurveID)
	_, _ = reconstructedCCS.ReadFrom(&_buf)

	// create srs
	srs, err := test.NewKZGSRS(reconstructedCCS)
	fmt.Println(reflect.TypeOf(srs))
	if err != nil {
		panic("Failed to create srs: " + err.Error())
	}

	pk, vk, err := plonk.Setup(reconstructedCCS, srs)
	if err != nil {
		panic("Setup failed!")
	}

	// SERIALIZE - WRITE THE PK & VK
	var bufPK bytes.Buffer
	_, _ = pk.WriteTo(&bufPK)
	err = os.WriteFile("tmp/pk.dat", bufPK.Bytes(), 0644)
	if err != nil {
		log.Fatal(err)
	}
	var bufVK bytes.Buffer
	_, _ = vk.WriteTo(&bufVK)
	err = os.WriteFile("tmp/vk.dat", bufVK.Bytes(), 0644)
	if err != nil {
		log.Fatal(err)
	}
	return
}

func Execute() {
	if err := plonkMemorySetupCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {

	cobra.OnInitialize()

	_curves := ecc.Implemented()
	curves := make([]string, len(_curves))
	for i := 0; i < len(_curves); i++ {
		curves[i] = strings.ToLower(_curves[i].String())
	}

	cfg.InputPath = plonkMemorySetupCmd.PersistentFlags().String("input", "none", "input path to the dedicated input")
	plonkMemorySetupCmd.MarkPersistentFlagRequired("input")
	cfg.Circuit = plonkMemorySetupCmd.PersistentFlags().String("circuit", "expo", "name of the circuit to use")
	cfg.CircuitSize = plonkMemorySetupCmd.PersistentFlags().Int("size", 10000, "size of the circuit, parameter to circuit constructor")
	cfg.Count = plonkMemorySetupCmd.PersistentFlags().Int("count", 2, "bench count (time is averaged on number of executions)")
	cfg.Curve = plonkMemorySetupCmd.PersistentFlags().String("curve", "bn254", "curve name. must be "+fmt.Sprint(curves))
}
