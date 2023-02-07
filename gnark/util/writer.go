package util

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
)

// WriteData writes the data to a file in either CSV or JSON format, based on the file format specified.
func WriteData(fileFormat string, data BenchData, filename ...string) error {

	var writer *csv.Writer
	var jsonEncoder *json.Encoder
	var file *os.File
	var err error
	var exists bool = false

	// if len(filename) == 0 {
	// 	file = os.Stdout
	// } else {
	// 	file, err = os.Create(filename[0])
	// 	if err != nil {
	// 		return err
	// 	}
	// 	defer file.Close()
	// }

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

func writeDataToCSV(exists bool, data BenchData, writer *csv.Writer) error {
	if exists == false {
		writer.Write(data.Headers())
		writer.Write(data.Values())
	} else {
		writer.Write(data.Values())
	}
	writer.Flush()
	return nil
}

func writeDataToJSON(data BenchData, encoder *json.Encoder) error {
	// TODO
	return nil
}
