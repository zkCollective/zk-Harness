/*
Benchmarking Math Operations over a variety of curves in gnark
*/
package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"

	"github.com/consensys/gnark/logger"
	"github.com/spf13/cobra"
	"github.com/zkCollective/zk-Harness/gnark/util"
)

const (
	pathPrefix = "../benchmarks/gnark/math"
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

// benchCurveOperations runs benchmarks for a variety of curves in gnark
func benchCurveOperations(cmd *cobra.Command, args []string) {

	log := logger.Logger()
	log.Info().Msg("Benchmarking curve operations - gnark: " + *fCurve + " " + *fOperation + " " + *fInputPath)

	var filepath_zkalc = pathPrefix + "/zkalc/gnark_" +
		"curve_" +
		*fCurve +
		"." + "txt"

	var filepath_zkHarness = pathPrefix + "/zkHarness/gnark_" +
		"curve_" +
		*fCurve +
		"." + "txt"

	paths := []string{
		pathPrefix,
		pathPrefix + "/zkalc",
		pathPrefix + "/zkHarness",
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
	if err := cleanup(); err != nil {
		log.Printf("Cleanup failed: %s\n", err)
	}

	// Clone gnark-crypto repo
	cmdShell := exec.Command("git", "clone", "-b", "zkalc", "https://github.com/ConsenSys/gnark-crypto.git")
	cmdShell.Dir = ".."
	err := cmdShell.Run()
	if err != nil {
		log.Printf("git clone failed with %s\n", err)
	}

	// Benchmark for all curves in config
	if err := bench_math(curveID.String(), filepath_zkalc); err != nil {
		log.Printf("Benchmark failed: %s\n", err)
	}

	results, err := readBenchmarkFile(filepath_zkalc)
	if err != nil {
		log.Printf("Reading benchmark file failed: %s\n", err)
		return
	}
	for _, result := range results {
		if err := writeResults(result, filepath_zkHarness); err != nil {
			log.Printf("Writing results failed: %s\n", err)
			return
		}
	}

	if err := cleanup(); err != nil {
		log.Printf("Cleanup failed: %s\n", err)
	}
}

// bench_math runs a shell command to perform a math benchmark with a specified curve
func bench_math(curve string, filename string) error {
	curve = strings.Replace(curve, "_", "-", -1)
	command := exec.Command("bash", "-c", "bash ./zkalc.sh "+curve+" | grep -vE '^[^[:space:]]+/' | tee "+filename)
	command.Dir = "../gnark-crypto"
	err := command.Run()
	if err != nil {
		return fmt.Errorf("Command failed with error: %s", err)
	}
	return nil
}

// writeResults writes the results of a benchmark to a file
func writeResults(result Result, filename string) error {
	// check memory usage, max ram requested from OS
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	bDataArith := util.BenchDataCurve{
		Framework: "gnark",
		Category:  "ec",
		Curve:     curveID.String(),
		Operation: result.Operation,
		Input:     "", // This needs to be replaced with the appropriate value
		MaxRAM:    m.Sys,
		Count:     int(result.Count),
		RunTime:   int64(result.Runtime),
	}

	// Check if file exists
	if _, err := os.Stat(filename); !os.IsNotExist(err) {
		// If file exists, remove it
		if err := os.Remove(filename); err != nil {
			return fmt.Errorf("failed to remove existing file: %s", err)
		}
	}

	if err := util.WriteData("csv", bDataArith, filename); err != nil {
		return fmt.Errorf("failed to write data: %s", err)
	}
	return nil
}

// readBenchmarkFile reads a file of benchmarks and returns the results
func readBenchmarkFile(filepath string) ([]Result, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %s", err)
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
		return nil, fmt.Errorf("failed to scan file: %s", err)
	}
	return results, nil
}

// cleanup removes the gnark-crypto directory
func cleanup() error {
	err := os.RemoveAll("../gnark-crypto")
	if err != nil {
		return fmt.Errorf("Failed to delete the cloned repository: %s", err)
	}
	return nil
}

func init() {
	fGroup = mathCmd.Flags().String("group", "None", "group to benchmark")

	rootCmd.AddCommand(mathCmd)
}
