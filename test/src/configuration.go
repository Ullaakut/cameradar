package main

import (
	"encoding/json"
	"fmt"
	"os"
)

func (m *manager) parseConfig() bool {
	// Get config file path
	confPath := "conf/cameratest.conf.json"
	av := len(os.Args)
	if av == 2 {
		confPath = os.Args[1]
	}

	// Load config
	fmt.Printf("Loading config file: %s ... ", confPath)
	configFile, err := os.Open(confPath)
	if err != nil {
		fmt.Printf("\nCan't open config file: %s\n", err)
		return false
	}
	dec := json.NewDecoder(configFile)
	if err = dec.Decode(&m); err != nil {
		fmt.Printf("\nUnable to deserialize config file: %s\n", err)
		return false
	}
	fmt.Println("Configuration file successfully loaded")

	return true
}
