package ui

import (
	"github.com/Ullaakut/cameradar/v6"
)

// NopReporter discards all UI events.
type NopReporter struct{}

// Start implements Reporter.
func (NopReporter) Start(cameradar.Step, string) {}

// Done implements Reporter.
func (NopReporter) Done(cameradar.Step, string) {}

// Progress implements Reporter.
func (NopReporter) Progress(cameradar.Step, string) {}

// Debug implements Reporter.
func (NopReporter) Debug(cameradar.Step, string) {}

// Error implements Reporter.
func (NopReporter) Error(cameradar.Step, error) {}

// Summary implements Reporter.
func (NopReporter) Summary([]cameradar.Stream, error) {}

// Close implements Reporter.
func (NopReporter) Close() {}
