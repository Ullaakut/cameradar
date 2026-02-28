package masscan

import (
	"context"
	"fmt"
	"net/netip"
	"strings"

	"github.com/Ullaakut/cameradar/v6"
	masscanlib "github.com/Ullaakut/masscan"
)

// Reporter reports scan progress and debug information.
type Reporter interface {
	Debug(step cameradar.Step, message string)
	Progress(step cameradar.Step, message string)
}

// Runner is something that can run a masscan scan.
type Runner interface {
	Run(ctx context.Context) (*masscanlib.Run, error)
}

// Scanner scans targets and ports for RTSP streams.
type Scanner struct {
	runner   Runner
	reporter Reporter
}

// New returns a Scanner configured with the provided targets and ports.
func New(targets, ports []string, reporter Reporter) (*Scanner, error) {
	runner, err := masscanlib.NewScanner(
		masscanlib.WithTargets(targets...),
		masscanlib.WithPorts(ports...),
		masscanlib.WithOpenOnly(),
	)
	if err != nil {
		return nil, fmt.Errorf("creating masscan scanner: %w", err)
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

func runScan(ctx context.Context, runner Runner, reporter Reporter) ([]cameradar.Stream, error) {
	results, err := runner.Run(ctx)
	if err != nil {
		return nil, fmt.Errorf("scanning network: %w", err)
	}

	for _, warning := range results.Warnings() {
		reporter.Debug(cameradar.StepScan, "masscan warning: "+warning)
	}

	var streams []cameradar.Stream
	for _, host := range results.Hosts {
		address := strings.TrimSpace(host.Address)
		if address == "" {
			reporter.Progress(cameradar.StepScan, "Skipping host with empty address")
			continue
		}

		addr, err := netip.ParseAddr(address)
		if err != nil {
			reporter.Progress(cameradar.StepScan, fmt.Sprintf("Skipping invalid address %q: %v", host.Address, err))
			continue
		}

		for _, port := range host.Ports {
			if port.Status != "open" {
				continue
			}

			if port.Number <= 0 || port.Number > 65535 {
				reporter.Progress(cameradar.StepScan, fmt.Sprintf("Skipping invalid port %d on %s", port.Number, host.Address))
				continue
			}

			streams = append(streams, cameradar.Stream{
				Address: addr,
				Port:    uint16(port.Number),
			})
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
