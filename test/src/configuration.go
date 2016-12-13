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
	"encoding/json"
	"fmt"
	"os"
)

func (t *Tester) parseConfig() bool {
	confPath := "conf/cameratest.conf.json"
	av := len(os.Args)
	if av > 1 {
		confPath = os.Args[1]
	}

	// Load config
	fmt.Printf("Loading Tester configuration file: %s ... ", confPath)
	configFile, err := os.Open(confPath)
	if err != nil {
		fmt.Printf("\nCan't open Tester configuration file: %s\n", err)
		return false
	}
	dec := json.NewDecoder(configFile)
	if err = dec.Decode(&t); err != nil {
		fmt.Printf("\nUnable to deserialize Tester configuration file: %s\n", err)
		return false
	}
	fmt.Println("Tester configuration file successfully loaded")

	return true
}
