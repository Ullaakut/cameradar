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
