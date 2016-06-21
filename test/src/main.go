package main

import (
	"fmt"
)

func main() {
	manager := new(manager)
	defer manager.Stop()

	// Parse conf (streams should already be launched by Jenkins)
	fmt.Println("--- Initializing Cameradar Test Tool ... ---")
	if !manager.Init() {
		fmt.Println("-> Cameradar Test Tool initialization FAILED")
		return
	}

	// Run tests
	if !manager.Run() {
		fmt.Println("-> Cameradar Test Tool FAILED")
	}

	// Write results
	fmt.Println("--- Writing results... ---")
	if !manager.WriteResults(*(manager.Result), manager.Config.Output) {
		fmt.Println("-> Write results FAILED")
		return
	}
	fmt.Println("--- Writing results done ---")
}
