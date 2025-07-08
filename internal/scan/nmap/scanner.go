package nmap

import (
	"context"
	"fmt"
	"net/netip"
	"strings"

	"github.com/Ullaakut/cameradar/v6"
	nmaplib "github.com/Ullaakut/nmap/v4"
)

// Reporter reports scan progress and debug information.
type Reporter interface {
	Debug(step cameradar.Step, message string)
	Progress(step cameradar.Step, message string)
}

// Runner is something that can run an nmap scan.
type Runner interface {
	Run(ctx context.Context) (*nmaplib.Run, error)
}

// Scanner scans targets and ports for RTSP streams.
type Scanner struct {
	runner   Runner
	reporter Reporter
}

// New returns a Scanner configured with the provided terminal and scan speed.
func New(scanSpeed int16, targets, ports []string, reporter Reporter) (*Scanner, error) {
	runner, err := nmaplib.NewScanner(
		nmaplib.WithTargets(targets...),
		nmaplib.WithPorts(ports...),
		nmaplib.WithServiceInfo(),
		nmaplib.WithTimingTemplate(nmaplib.Timing(scanSpeed)),
	)
	if err != nil {
		return nil, fmt.Errorf("creating nmap scanner: %w", err)
	}

	return &Scanner{
		runner:   runner,
		reporter: reporter,
	}, nil
}

// Scan discovers RTSP streams on the configured targets and ports.
func (s *Scanner) Scan(ctx context.Context) ([]cameradar.Stream, error) {
	return runScan(ctx, s.runner, s.reporter)
}

func runScan(ctx context.Context, nmap Runner, reporter Reporter) ([]cameradar.Stream, error) {
	results, err := nmap.Run(ctx)
	if err != nil {
		return nil, fmt.Errorf("scanning network: %w", err)
	}

	for _, warning := range results.Warnings() {
		reporter.Debug(cameradar.StepScan, "nmap warning: "+warning)
	}

	var streams []cameradar.Stream
	for _, host := range results.Hosts {
		for _, port := range host.Ports {
			if port.Status() != "open" {
				continue
			}

			if !strings.Contains(port.Service.Name, "rtsp") {
				continue
			}

			for _, address := range host.Addresses {
				addr, err := netip.ParseAddr(address.Addr)
				if err != nil {
					reporter.Progress(cameradar.StepScan, fmt.Sprintf("Skipping invalid address %q: %v", address.Addr, err))
					continue
				}

				streams = append(streams, cameradar.Stream{
					Device:  port.Service.Product,
					Address: addr,
					Port:    port.ID,
				})
			}
		}
	}

	reporter.Progress(cameradar.StepScan, fmt.Sprintf("Found %d RTSP streams", len(streams)))
	updateSummary(reporter, streams)

	return streams, nil
}

type summaryUpdater interface {
	UpdateSummary(streams []cameradar.Stream)
}

func updateSummary(reporter Reporter, streams []cameradar.Stream) {
	updater, ok := reporter.(summaryUpdater)
	if !ok {
		return
	}
	updater.UpdateSummary(streams)
}
