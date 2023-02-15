package util

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

func ReadFromInputPath(pathInput string) (map[string]interface{}, error) {

	// TODO read all input files from a single directory

	// Construct the absolute path to the file
	absPath := filepath.Join("../", pathInput, "/input_1.json")
	absPath, err := filepath.Abs(absPath)
	if err != nil {
		fmt.Println("Error constructing absolute path:", err)
		return nil, err
	}

	file, err := os.Open(absPath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var data map[string]interface{}
	err = json.NewDecoder(file).Decode(&data)
	if err != nil {
		panic(err)
	}

	return data, nil
}
