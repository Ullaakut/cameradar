package nmap

import (
	"fmt"
	"log"
)

// A scanner can be instanciated with options to set the arguments
// that are given to nmap.
func ExampleScanner_simple() {
	s, err := NewScanner(
		WithTargets("google.com", "facebook.com", "youtube.com"),
		WithCustomDNSServers("8.8.8.8", "8.8.4.4"),
		WithTimingTemplate(TimingFastest),
		WithTCPScanFlags(FlagACK, FlagNULL, FlagRST),
	)
	if err != nil {
		log.Fatalf("unable to create nmap scanner: %v", err)
	}

	scanResult, err := s.Run()
	if err != nil {
		log.Fatalf("nmap encountered an error: %v", err)
	}

	fmt.Printf(
		"Scan successful: %d hosts up\n",
		scanResult.Stats.Hosts.Up,
	)
	// Output: Scan successful: 3 hosts up
}

// A scanner can be given custom idiomatic filters for both hosts
// and ports.
func ExampleScanner_filters() {
	s, err := NewScanner(
		WithTargets("google.com", "facebook.com"),
		WithPorts("843"),
		WithFilterHost(func(h Host) bool {
			// Filter out hosts with no open ports.
			for idx := range h.Ports {
				if h.Ports[idx].Status() == "open" {
					return true
				}
			}
			return false
		}),
	)
	if err != nil {
		log.Fatalf("unable to create nmap scanner: %v", err)
	}

	scanResult, err := s.Run()
	if err != nil {
		log.Fatalf("nmap encountered an error: %v", err)
	}

	fmt.Printf(
		"Filtered out hosts %d / Original number of hosts: %d\n",
		len(scanResult.Hosts),
		scanResult.Stats.Hosts.Total,
	)
	// Output: Filtered out hosts 1 / Original number of hosts: 2
}
