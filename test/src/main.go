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
	"os"
)

func main() {
	Tester := new(Tester)
	defer Tester.Stop()

	// Parse conf (streams should already be launched by Jenkins)
	fmt.Println("--- Initializing Cameradar Test Tool ... ---")
	if !Tester.Init() {
		fmt.Println("-> Cameradar Test Tool initialization FAILED")
		return
	}

	// Run tests
	if !Tester.Run() {
		fmt.Println("-> Cameradar Test Tool FAILED")
	}

	// Write results
	fmt.Println("--- Writing results... ---")
	if !Tester.WriteResults(*(Tester.Result), Tester.Output) {
		fmt.Println("-> Write results FAILED")
		os.Exit(1)
	}

	fmt.Println("--- Writing results done ---")
	os.Exit(0)
}
