package cmrdr

import (
	"strings"

	"github.com/Ullaakut/nmap"
)

// Discover scans the target networks and tries to find RTSP streams within them.
//
// targets can be:
//
//    - a subnet (e.g.: 172.16.100.0/24)
//    - an IP (e.g.: 172.16.100.10)
//    - a hostname (e.g.: localhost)
//    - a range of IPs (e.g.: 172.16.100.10-20)
//
// ports can be:
//
//    - one or multiple ports and port ranges separated by commas (e.g.: 554,8554-8560,18554-28554)
func Discover(targets, ports []string, speed int) ([]Stream, error) {
	// Run nmap command to discover open ports on the specified targets & ports
	scanner, err := nmap.NewScanner(
		nmap.WithTargets(targets...),
		nmap.WithPorts(ports...),
		nmap.WithTimingTemplate(nmap.Timing(speed)),
	)
	if err != nil {
		return nil, err
	}

	return scan(scanner)
}

func scan(scanner nmap.ScanRunner) ([]Stream, error) {
	results, err := scanner.Run()
	if err != nil {
		return nil, err
	}

	var streams []Stream
	// Get streams from nmap results
	for _, host := range results.Hosts {
		for _, port := range host.Ports {
			if port.Status() != "open" {
				continue
			}

			if !strings.Contains(port.Service.Name, "rtsp") {
				continue
			}

			for _, address := range host.Addresses {
				streams = append(streams, Stream{
					Device:  port.Service.Product,
					Address: address.Addr,
					Port:    port.ID,
				})
			}
		}
	}

	return streams, nil
}
