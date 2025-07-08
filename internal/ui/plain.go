package ui

import (
	"fmt"
	"io"
	"time"

	"github.com/Ullaakut/cameradar/v6"
)

// PlainReporter renders a line-oriented UI for non-interactive terminals.
type PlainReporter struct {
	out   io.Writer
	debug bool
}

// NewPlainReporter creates a line-oriented reporter.
func NewPlainReporter(out io.Writer, debug bool) *PlainReporter {
	return &PlainReporter{
		out:   out,
		debug: debug,
	}
}

// Start prints the beginning of a step.
func (r *PlainReporter) Start(step cameradar.Step, message string) {
	r.print(step, "START", message)
}

// Done prints the completion of a step.
func (r *PlainReporter) Done(step cameradar.Step, message string) {
	r.print(step, "DONE", message)
}

// Progress prints a progress message.
func (r *PlainReporter) Progress(step cameradar.Step, message string) {
	if _, _, ok := cameradar.ParseProgressMessage(message); ok {
		return
	}
	r.print(step, "INFO", message)
}

// Debug prints a debug message when debug mode is enabled.
func (r *PlainReporter) Debug(step cameradar.Step, message string) {
	if !r.debug {
		return
	}
	r.print(step, "DEBUG", message)
}

// Error prints an error message.
func (r *PlainReporter) Error(step cameradar.Step, err error) {
	if err == nil {
		return
	}
	r.print(step, "ERROR", err.Error())
}

// Summary prints the final summary.
func (r *PlainReporter) Summary(streams []cameradar.Stream, err error) {
	_, _ = fmt.Fprintln(r.out, "Summary")
	_, _ = fmt.Fprintln(r.out, "-------")
	_, _ = fmt.Fprintln(r.out, FormatSummary(streams, err))
}

// Close is a no-op for the plain reporter.
func (r *PlainReporter) Close() {}

func (r *PlainReporter) print(step cameradar.Step, level, message string) {
	if message == "" {
		return
	}

	_, _ = fmt.Fprintf(r.out, "[%s] %s: %s (%s)\n", level, cameradar.StepLabel(step), message, time.Now().Format(time.RFC3339))
}
