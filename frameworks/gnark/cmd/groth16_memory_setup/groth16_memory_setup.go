package main

import (
	"fmt"
	"os"
	"strings"
	"bytes"
	"log"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/spf13/cobra"
	"github.com/zkCollective/zk-Harness/frameworks/gnark/parser"
)

// groth16Cmd represents the groth16 command
var groth16MemorySetupCmd = &cobra.Command{
	Use:   "groth16MemorySetup",
	Short: "runs benchmarks and profiles using Groth16 proof system",
	Run:   runGroth16MemorySetup,
}

var cfg = parser.NewConfig()

func runGroth16MemorySetup(cmd *cobra.Command, args []string) {

	if err := parser.ParseFlagsMemory(cfg); err != nil {
		fmt.Println("error: ", err.Error())
		cmd.Help()
		os.Exit(-1)
	}

	// Open the file in read-only mode
	_ccs, err := os.ReadFile("tmp/ccs.dat")
	if err != nil {
		log.Fatal(err)
	}
	_buf := *bytes.NewBuffer(_ccs)
	reconstructedCCS := groth16.NewCS(parser.CurveID)
	_, _ = reconstructedCCS.ReadFrom(&_buf)

	pk, vk, err := groth16.Setup(reconstructedCCS)
	if err != nil {
		panic("Setup failed!")
	}

	// Open the file in write mode for pk
	var bufPK bytes.Buffer
	_, _ = pk.WriteTo(&bufPK)
	err = os.WriteFile("tmp/pk.dat", bufPK.Bytes(), 0644)
	if err != nil {
		log.Fatal(err)
	}
	

	// Open the file in write mode for vk
	var bufVK bytes.Buffer
	_, _ = vk.WriteTo(&bufVK)
	err = os.WriteFile("tmp/vk.dat", bufVK.Bytes(), 0644)
	if err != nil {
		log.Fatal(err)
	}
	return
}

func Execute() {
	if err := groth16MemorySetupCmd.Execute(); err != nil {
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

	cfg.InputPath = groth16MemorySetupCmd.PersistentFlags().String("input", "none", "input path to the dedicated input")
	groth16MemorySetupCmd.MarkPersistentFlagRequired("input")
	cfg.Circuit = groth16MemorySetupCmd.PersistentFlags().String("circuit", "expo", "name of the circuit to use")
	cfg.CircuitSize = groth16MemorySetupCmd.PersistentFlags().Int("size", 10000, "size of the circuit, parameter to circuit constructor")
	cfg.Count = groth16MemorySetupCmd.PersistentFlags().Int("count", 2, "bench count (time is averaged on number of executions)")
	cfg.Curve = groth16MemorySetupCmd.PersistentFlags().String("curve", "bn254", "curve name. must be "+fmt.Sprint(curves))
}
