package ui

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/Ullaakut/cameradar/v6"
)

// Reporter defines the interface for cameradar UIs.
type Reporter interface {
	Start(step cameradar.Step, message string)
	Done(step cameradar.Step, message string)
	Progress(step cameradar.Step, message string)
	Debug(step cameradar.Step, message string)
	Error(step cameradar.Step, err error)
	Summary(streams []cameradar.Stream, err error)
	Close()
}

// NewReporter creates a Reporter based on the requested mode.
func NewReporter(mode cameradar.Mode, debug bool, out io.Writer, interactive bool, buildInfo BuildInfo, cancel context.CancelFunc) (Reporter, error) {
	if debug {
		return NewPlainReporter(out, debug), nil
	}

	switch mode {
	case cameradar.ModePlain:
		return NewPlainReporter(out, debug), nil
	case cameradar.ModeTUI:
		if !interactive {
			return nil, errors.New("tui mode requires an interactive terminal")
		}
		return NewTUIReporter(debug, out, buildInfo, cancel)
	case cameradar.ModeAuto:
		if interactive {
			return NewTUIReporter(debug, out, buildInfo, cancel)
		}
		return NewPlainReporter(out, debug), nil
	default:
		return nil, fmt.Errorf("unsupported ui mode %q", mode)
	}
}
