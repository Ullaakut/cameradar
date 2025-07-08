package cameradar

import (
	"context"
	"strconv"
	"strings"

	"github.com/Ullaakut/nmap"
)

// Scan scans the target networks and tries to find RTSP streams within them.
//
// targets can be:
//
//   - a subnet (e.g.: 172.16.100.0/24)
//   - an IP (e.g.: 172.16.100.10)
//   - a hostname (e.g.: localhost)
//   - a range of IPs (e.g.: 172.16.100.10-20)
//
// ports can be:
//
//   - one or multiple ports and port ranges separated by commas (e.g.: 554,8554-8560,18554-28554)
func (s *Scanner) Scan(ctx context.Context) ([]Stream, error) {
	s.term.StartStep("Scanning the network")

	targets := make([]string, len(s.targets))
	for i, target := range s.targets {
		targets[i] = target.String()
	}

	ports := make([]string, len(s.ports))
	for i, port := range s.ports {
		ports[i] = strconv.Itoa(int(port))
	}

	// Run nmap command to discover open ports on the specified targets & ports.
	nmapScanner, err := nmap.NewScanner(
		nmap.WithContext(ctx),
		nmap.WithTargets(targets...),
		nmap.WithPorts(ports...),
		nmap.WithServiceInfo(),
		nmap.WithTimingTemplate(nmap.Timing(s.scanSpeed)),
	)
	if err != nil {
		return nil, s.term.FailStepf("preparing network scanner: %v", err)
	}

	return s.scan(nmapScanner)
}

func (s *Scanner) scan(nmapScanner nmap.ScanRunner) ([]Stream, error) {
	results, warnings, err := nmapScanner.Run()
	for _, warning := range warnings {
		s.term.Infoln("[Nmap Warning]", warning)
	}
	if err != nil {
		return nil, s.term.FailStepf("scanning network: %v", err)
	}

	// Get streams from nmap results.
	var streams []Stream
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

	s.term.Debugf("Found %d RTSP streams\n", len(streams))
	s.term.EndStep()

	return streams, nil
}
