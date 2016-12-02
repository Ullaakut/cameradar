// Copyright 2016 Etix Labs
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"sync"
)

// Tester is the structure that will manage the whole testing
type Tester struct {
	ServiceConf ServiceConfig `json:"cameradar"`
	Output      string        `json:"output"`
	Tests       []Result      `json:"tests"`

	Cameradar Service // Runs the command and manages the logs
	Result    *Test   // Results of the testing
	DB        MysqlDB // Access to the database to make sure it's empty
}

// Init gets the testing configuration and makes sure that no other Cameradar service is running at the moment
func (t *Tester) Init() bool {
	fmt.Println("- Parsing")
	if !t.parseConfig() {
		return false
	}

	fmt.Println("- Cleaning content")
	killService(&t.Cameradar)

	return true
}

// Run launches the tests that have been set up by the init method
func (t *Tester) Run() bool {
	var wg sync.WaitGroup

	fmt.Println("\n- Launching all tests")
	var newTest = new(Test)
	newTest.expected = t.Tests

	if t.configureDatabase(&t.DB) {
		t.dropDB()
		wg.Add(1)
		go t.invokeTestCase(newTest, &wg)
		t.Result = newTest
	}

	wg.Wait()
	fmt.Println("All tests completed")
	return true
}

// Stop kills the service launched by the tester
func (t *Tester) Stop() bool {
	killService(&t.Cameradar)
	return true
}
