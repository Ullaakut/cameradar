package main

import (
	"encoding/json"
	"fmt"
	"os"
)

// Launch it via goroutine
// Start read log of service
func getResult(test *[]Result, resultPath string) bool {
	// Load config
	resultFile, err := os.Open(resultPath)
	if err != nil {
		fmt.Printf("\nCan't open result file: %s\n", err)
		return false
	}
	dec := json.NewDecoder(resultFile)
	if err = dec.Decode(&test); err != nil {
		fmt.Printf("\nUnable to deserialize result file: %s\n", err)
		return false
	}
	return true
}
