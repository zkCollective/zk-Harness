package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/spf13/cobra"
	"github.com/zkCollective/zk-Harness/gnark/parser"
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
	f, err := os.Open("tmp/ccs.dat")
	if err != nil {
		panic("Failed to open file: " + err.Error())
	}

	ccs := groth16.NewCS(parser.CurveID)
	_, err = ccs.ReadFrom(f)
	if err != nil {
		panic("Failed to read from file: " + err.Error())
	}

	// Close the file after reading
	f.Close()

	pk, vk, err := groth16.Setup(ccs)
	if err != nil {
		panic("Setup failed!")
	}

	// Open the file in write mode for pk
	fPK, err := os.OpenFile("tmp/pk.dat", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		panic("Failed to open file for writing: " + err.Error())
	}
	defer fPK.Close()

	_, err = pk.WriteRawTo(fPK)
	if err != nil {
		panic("Failed to write pk to file: " + err.Error())
	}

	// Open the file in write mode for vk
	fVK, err := os.OpenFile("tmp/vk.dat", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		panic("Failed to open file for writing: " + err.Error())
	}
	defer fVK.Close()

	_, err = vk.WriteRawTo(fVK)
	if err != nil {
		panic("Failed to write vk to file: " + err.Error())
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
	cfg.FileType = groth16MemorySetupCmd.PersistentFlags().String("filetype", "csv", "Type of file to output for benchmarks")
}