package main

import (
	"fmt"
	"os/exec"
	"strings"
	"sync"
)

// Service needs refacto
type Service struct {
	Path       string `json:"Path"`
	Args       string `json:"Args"`
	Ports      string `json:"Ports"`
	IdsPath    string `json:"IdsPath"`
	RoutesPath string `json:"RoutesPath"`
	DbHost     string `json:"dbHost"`
	DbPort     int    `json:"dbPort"`
	DbUser     string `json:"dbUser"`
	DbPassword string `json:"dbPassword"`
	DbName     string `json:"dbName"`
	ThumbPath  string `json:"ThumbPath"`
	Console    bool   `json:"Console"`

	Logs   []string
	Active bool // Based on io.ReadCloser status

	Mutex sync.Mutex
	cmd   *exec.Cmd // Go handler of the service
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

	// Sending SIGABORT, more reliable for VLC
	sigAbort := []string{service.Path, "-s", "SIGABRT"}
	fmt.Printf("Executing: killall %s -s SIGABRT\n", service.Path)
	cmd = exec.Command("killall", sigAbort...)
	err = cmd.Run()
	if err != nil {
		fmt.Println(err)
	}
}
