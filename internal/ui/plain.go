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

// PrintStartup prints build metadata and configuration options.
func (r *PlainReporter) PrintStartup(buildInfo BuildInfo, options []string) {
	step := cameradar.Step("Startup")
	message := fmt.Sprintf("Running cameradar version %s, commit %s", buildInfo.LogVersion(), buildInfo.ShortCommit())
	r.print(step, "INFO", message)
	if len(options) == 0 {
		return
	}
	for _, option := range options {
		r.print(step, "INFO", option)
	}
}

// Start prints the beginning of a step.
func (r *PlainReporter) Start(step cameradar.Step, message string) {
	r.print(step, "STEP", message)
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
	r.print(step, "DBUG", message)
}

// Error prints an error message.
func (r *PlainReporter) Error(step cameradar.Step, err error) {
	if err == nil {
		return
	}
	r.print(step, "EROR", err.Error())
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

	level = normalizeLevel(level)
	_, _ = fmt.Fprintf(r.out, "%s [%s] %s: %s\n", time.Now().Format(time.RFC3339), level, cameradar.StepLabel(step), message)
}

func normalizeLevel(level string) string {
	switch level {
	case "DEBUG":
		return "DBUG"
	case "ERROR":
		return "EROR"
	case "START":
		return "STEP"
	case "STEP":
		return "STEP"
	}
	if len(level) >= 4 {
		return level[:4]
	}
	return fmt.Sprintf("%-4s", level)
}
