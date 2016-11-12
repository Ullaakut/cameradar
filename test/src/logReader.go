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
	"bufio"
	"fmt"
	"io"
)

// Launch it via goroutine
// Start read log of service
func readLog(service *Service, reader io.ReadCloser) {
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		str := scanner.Text()
		if service.Console {
			fmt.Printf("[%s] %s\n", service.Path, str)
		}
		fmt.Printf("%s\n", str)
		service.Mutex.Lock()
		service.Logs = append(service.Logs, str)
		service.Mutex.Unlock()
	}
	if err := scanner.Err(); err != nil {
		fmt.Printf("[%s] Service failed: %s\n", service.Path, err)
	}
	fmt.Printf("Logger of service: [%s] stopped\n", service.Path)
	service.Active = false
}
