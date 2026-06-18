package cameradar

import (
	"context"
	"errors"
	"fmt"
)

// Reporter reports progress and results of the application.
type Reporter interface {
	Start(step Step, message string)
	Done(step Step, message string)
	Error(step Step, err error)
	Summary(streams []Stream, err error)
}

// App scans one or more targets and attacks all RTSP streams found to get their credentials.
type App struct {
	streamScanner StreamScanner
	attacker      StreamAttacker
	reporter      Reporter

	targets []string
	ports   []string
}

// StreamScanner discovers RTSP streams for the given inputs.
type StreamScanner interface {
	Scan(ctx context.Context) ([]Stream, error)
}

// StreamAttacker attacks streams to discover routes and credentials.
type StreamAttacker interface {
	Attack(ctx context.Context, streams []Stream) ([]Stream, error)
}

// New creates a new App with explicit dependencies.
func New(streamScanner StreamScanner, attacker StreamAttacker, targets, ports []string, reporter Reporter) (*App, error) {
	if streamScanner == nil {
		return nil, errors.New("stream scanner is required")
	}
	if attacker == nil {
		return nil, errors.New("stream attacker is required")
	}

	app := &App{
		streamScanner: streamScanner,
		attacker:      attacker,
		targets:       targets,
		ports:         ports,
		reporter:      reporter,
	}

	return app, nil
}

// Run runs the scan and prints the results.
func (a *App) Run(ctx context.Context) error {
	a.reporter.Start(StepScan, "Scanning targets for RTSP streams")
	streams, err := a.streamScanner.Scan(ctx)
	if err != nil {
		wrapped := fmt.Errorf("discovering devices: %w", err)
		a.reporter.Error(StepScan, wrapped)
		a.reporter.Summary(streams, wrapped)
		return wrapped
	}
	a.reporter.Done(StepScan, "Scan complete")

	streams, err = a.attacker.Attack(ctx, streams)
	if err != nil {
		wrapped := fmt.Errorf("attacking devices: %w", err)
		a.reporter.Summary(streams, wrapped)
		return wrapped
	}

	a.reporter.Summary(streams, nil)
	return nil
}
