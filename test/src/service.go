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
	"os/exec"
	"strings"
	"sync"
)

// Service needs refacto
type Service struct {
	Path       string `json:"path"`
	Args       string `json:"args"`
	Ports      string `json:"ports"`
	IdsPath    string `json:"ids_path"`
	RoutesPath string `json:"routes_path"`
	ThumbPath  string `json:"thumb_path"`
	DbHost     string `json:"db_host"`
	DbPort     int    `json:"db_port"`
	DbUser     string `json:"db_user"`
	DbPassword string `json:"db_password"`
	DbName     string `json:"db_name"`
	Console    bool   `json:"console"`

	Logs   []string
	Active bool // Based on io.ReadCloser status
	Mutex  sync.Mutex
	cmd    *exec.Cmd // Go handler of the service
}

func startService(service *Service) bool {
	// Launch service
	args := strings.Fields(service.Args)
	service.cmd = exec.Command(service.Path, args...)

	handler, err := service.cmd.StdoutPipe()
	if err != nil {
		fmt.Println(err)
		return false
	}
	errHandler, err := service.cmd.StderrPipe()
	if err != nil {
		fmt.Println(err)
		return false
	}
	// Launch
	err = service.cmd.Start()
	if err != nil {
		fmt.Println(err)
		return false
	}

	fmt.Printf("Service: [%s] started\n", service.Path)
	service.Active = true

	// Read service logs and update service status
	// Set pipes
	go readLog(service, handler)
	go readLog(service, errHandler)

	return true
}

// Stop only specified service instance
func stopService(service *Service) {
	service.cmd.Process.Kill()
}

// Kill all instances of specified service
func killService(service *Service) {
	// Sending SIGTERM
	fmt.Printf("Executing: killall %s\n", service.Path)
	cmd := exec.Command("killall", service.Path)
	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
	}

	sigAbort := []string{service.Path, "-s", "SIGABRT"}
	fmt.Printf("Executing: killall %s -s SIGABRT\n", service.Path)
	cmd = exec.Command("killall", sigAbort...)
	err = cmd.Run()
	if err != nil {
		fmt.Println(err)
	}
}
