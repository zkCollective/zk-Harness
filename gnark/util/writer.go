package util

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
)

// WriteData writes the data to a file in either CSV or JSON format, based on the file format specified.
func WriteData(fileFormat string, data interface{}, filename ...string) error {

	var writer *csv.Writer
	var jsonEncoder *json.Encoder
	var file *os.File
	var err error
	var exists bool = false

	if len(filename) == 0 {
		file = os.Stdout
	} else {
		if _, err := os.Stat(filename[0]); os.IsNotExist(err) {
			file, err = os.Create(filename[0])
			if err != nil {
				return err
			}
		} else {
			exists = true
			// file, err = os.Open(filename[0])
			file, err = os.OpenFile(filename[0], os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
			if err != nil {
				return err
			}
		}
		defer file.Close()
	}

	// Write the data to the file or stdout in either CSV or JSON format
	switch fileFormat {
	case "csv":
		writer = csv.NewWriter(file)
		err = writeDataToCSV(exists, data, writer)
	case "json":
		jsonEncoder = json.NewEncoder(file)
		err = writeDataToJSON(data, jsonEncoder)
	default:
		return fmt.Errorf("unsupported file format %s", fileFormat)
	}

	return err
}

func writeDataToCSV(exists bool, data interface{}, writer *csv.Writer) error {
	if exists == false {
		if headers, ok := data.(HeadersProvider); ok {
			writer.Write(headers.Headers())
		}
		writer.Write(data.(ValuesProvider).Values())
	} else {
		writer.Write(data.(ValuesProvider).Values())
	}
	writer.Flush()
	return nil
}

func writeDataToJSON(data interface{}, encoder *json.Encoder) error {
	// TODO
	return nil
}

// func writeResultCircuit(took time.Duration, ccs constraint.ConstraintSystem, data util.BenchDataCircuit, proof_size int) {
// 	// check memory usage, max ram requested from OS
// 	var m runtime.MemStats
// 	runtime.ReadMemStats(&m)

// 	_, secret, public := ccs.GetNbVariables()
// 	bData := util.BenchDataCircuit{
// 		Framework:         data.Framework,
// 		Category:          "circuit",
// 		Backend:           "groth16",
// 		Curve:             curveID.String(),
// 		Circuit:           *fCircuit,
// 		Input:             *fInputPath,
// 		Operation:         *fAlgo,
// 		NbConstraints:     ccs.GetNbConstraints(),
// 		NbSecretVariables: secret,
// 		NbPublicVariables: public,
// 		ProofSize:         proof_size,
// 		MaxRAM:            (m.Sys / 1024 / 1024),
// 		Count:             *fCount,
// 		RunTime:           took.Milliseconds(),
// 	}

// 	if err := util.WriteData("csv", bData, filename); err != nil {
// 		panic(err)
// 	}
// }
