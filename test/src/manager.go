package main

import (
	"fmt"
	"sync"
)

type manager struct {
	Config

	Tests  []Result
	Result *TestCase
	DB     mysql_db
}

// Config needs refacto
type Config struct {
	Cameradar Service `json:"Cameradar"`

	Output string
}

func (m *manager) Init() bool {
	fmt.Println("- Parsing")
	if !m.parseConfig() {
		return false
	}

	fmt.Println("- Cleaning content")
	killService(&m.Config.Cameradar)

	return true
}

func (m *manager) Run() bool {
	var wg sync.WaitGroup

	fmt.Println("\n- Launching all tests")
	var newTest = new(TestCase)
	newTest.expected = m.Tests
	if m.generateConfig(m.Tests, &m.DB) {
		m.dropDB()
		wg.Add(1)
		go m.invokeTestCase(newTest, &wg)
		m.Result = newTest
	}
	wg.Wait()
	fmt.Printf("All tests completed\n")
	return true
}

func (m *manager) Stop() bool {
	killService(&m.Config.Cameradar)
	return true
}
