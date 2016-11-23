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
	"errors"
	"fmt"
	"os"
)

// Result contains the data of a Cameradar result, plus an error field in order to add error messages to the JUnit report
type Result struct {
	Address     string `json:"address"`
	IDsFound    bool   `json:"ids_found"`
	PathFound   bool   `json:"path_found"`
	Password    string `json:"password"`
	Port        int    `json:"port"`
	Route       string `json:"route"`
	ServiceName string `json:"service_name"`
	Protocol    string `json:"protocol"`
	State       string `json:"state"`
	Username    string `json:"username"`
	Valid       bool   `json:"valid"`
	Thumb       string `json:"thumbnail_path"`
	err         error  // in case of a fail, add a message
}

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

func isValid(e *Result, r Result) bool {
	if e.Username != r.Username {
		e.err = errors.New(e.Address + " had a different username than " + r.Username)
		return false
	}
	if e.Password != r.Password {
		e.err = errors.New(e.Address + " had a different password than " + r.Password)
		return false
	}
	if e.Port != r.Port {
		e.err = errors.New(e.Address + " had a different port than expected")
		return false
	}
	if e.Valid != r.Valid {
		e.err = errors.New(e.Address + " had a different validity than expected")
		return false
	}
	return true
}

// Extend needs refacto
func Extend(slice []Result, element Result) []Result {
	n := len(slice)
	if n == cap(slice) {
		// Slice is full; must grow.
		// We double its size and add 1, so if the size is zero we still grow.
		newSlice := make([]Result, len(slice), 2*len(slice)+1)
		copy(newSlice, slice)
		slice = newSlice
	}
	slice = slice[0 : n+1]
	slice[n] = element
	return slice
}
