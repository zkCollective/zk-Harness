/*
Benchmarking Math Operations over a variety of curves in gnark
*/
package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"

	"github.com/consensys/gnark/logger"
	"github.com/spf13/cobra"
	"github.com/zkCollective/zk-Harness/gnark/util"
)

var mathCmd = &cobra.Command{
	Use:   "ec",
	Short: "runs benchmarks and profiles for the gnark arithmetic operations",
	Run:   benchCurveOperations,
}

// Result struct to hold values from txt files
type Result struct {
	Operation string
	Runtime   float64
	Count     int64
}

func benchCurveOperations(cmd *cobra.Command, args []string) {

	log := logger.Logger()
	log.Info().Msg("Benchmarking curve operations - gnark: " + *fCurve + " " + *fOperation + " " + *fInputPath)

	var filepath_zkalc = "../benchmarks/gnark/math/zkalc/gnark_" +
		"curve_" +
		*fCurve +
		"." + "txt"

	var filepath_zkHarness = "../benchmarks/gnark/math/zkHarness/gnark_" +
		"curve_" +
		*fCurve +
		"." + "txt"

	paths := []string{
		"../benchmarks/gnark/math",
		"../benchmarks/gnark/math/zkalc",
		"../benchmarks/gnark/math/zkHarness",
	}

	for _, path := range paths {
		err := os.MkdirAll(path, 0755)
		if err != nil {
			log.Printf("Failed to create directory '%s': %s\n", path, err)
		} else {
			log.Printf("Directory '%s' created successfully\n", path)
		}
	}

	if err := parseFlags(); err != nil {
		fmt.Println("error: ", err.Error())
		cmd.Help()
		os.Exit(-1)
	}

	// Clean up for new benchmarks
	cleanup()

	// Clone gnark-crypto repo
	cmdShell := exec.Command("git", "clone", "-b", "zkalc", "https://github.com/ConsenSys/gnark-crypto.git")
	cmdShell.Dir = ".."
	err := cmdShell.Run()
	if err != nil {
		log.Printf("git clone failed with %s\n", err)
	}

	// Benchmark for all curves in config
	bench_math(curveID.String(), filepath_zkalc)

	results := readBenchmarkFile(filepath_zkalc)
	for _, result := range results {
		writeResults(result, filepath_zkHarness)
	}

	cleanup()
}

func bench_math(curve string, filename string) {
	command := exec.Command("bash", "-c", "bash ./zkalc.sh "+curve+" | tee "+filename)
	command.Dir = "../gnark-crypto"
	err := command.Run()
	if err != nil {
		log.Printf("Command failed with error: %s\n", err)
	}
}

func writeResults(result Result, filename string) {
	// check memory usage, max ram requested from OS
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	operationString := result.Operation

	bDataArith := util.BenchDataCurve{
		Framework: "gnark",
		Category:  "ec",
		Curve:     curveID.String(),
		Operation: operationString,
		Input:     "", // This needs to be replaced with the appropriate value
		MaxRAM:    m.Sys,
		Count:     int(result.Count),
		RunTime:   int64(result.Runtime),
	}

	if err := util.WriteData("csv", bDataArith, filename); err != nil {
		panic(err)
	}
}

func readBenchmarkFile(filepath string) []Result {
	file, err := os.Open(filepath)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer file.Close()

	var results []Result
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "Benchmark") {
			splitLine := strings.Fields(line)
			if len(splitLine) < 3 {
				fmt.Printf("Invalid line format: %s\n", line)
				continue
			}
			operation := splitLine[0]
			count, _ := strconv.ParseInt(splitLine[1], 10, 64)
			runtimeStr := strings.TrimSuffix(splitLine[2], "ns/op")
			runtime, _ := strconv.ParseFloat(runtimeStr, 64)
			results = append(results, Result{
				Operation: operation,
				Runtime:   runtime,
				Count:     count,
			})
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}
	return results
}

func cleanup() {
	err := os.RemoveAll("../gnark-crypto")
	if err != nil {
		log.Printf("Failed to delete the cloned repository: %s\n", err)
	}
}

func init() {
	fGroup = mathCmd.Flags().String("group", "None", "group to benchmark")

	rootCmd.AddCommand(mathCmd)
}
