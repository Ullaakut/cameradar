package ui_test

import (
	"bytes"
	"testing"

	"github.com/Ullaakut/cameradar/v6"
	"github.com/Ullaakut/cameradar/v6/internal/ui"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewReporter(t *testing.T) {
	tests := []struct {
		name            string
		mode            cameradar.Mode
		interactive     bool
		wantType        string
		wantErrContains string
	}{
		{
			name:        "plain",
			mode:        cameradar.ModePlain,
			interactive: false,
			wantType:    "plain",
		},
		{
			name:        "auto non-interactive",
			mode:        cameradar.ModeAuto,
			interactive: false,
			wantType:    "plain",
		},
		{
			name:            "tui non-interactive",
			mode:            cameradar.ModeTUI,
			interactive:     false,
			wantErrContains: "interactive terminal",
		},
		{
			name:            "unsupported",
			mode:            cameradar.Mode("unknown"),
			interactive:     false,
			wantErrContains: "unsupported ui mode",
		},
		{
			name:        "auto interactive",
			mode:        cameradar.ModeAuto,
			interactive: true,
			wantType:    "tui",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			out := &bytes.Buffer{}

				reporter, err := ui.NewReporter(test.mode, false, out, test.interactive, ui.BuildInfo{Version: "dev", Commit: "none"})

			if test.wantErrContains != "" {
				require.Error(t, err)
				assert.ErrorContains(t, err, test.wantErrContains)
				assert.Nil(t, reporter)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, reporter)

			switch test.wantType {
			case "plain":
				_, ok := reporter.(*ui.PlainReporter)
				assert.True(t, ok)
			case "tui":
				_, ok := reporter.(*ui.TUIReporter)
				assert.True(t, ok)
			}

			reporter.Close()
		})
	}
}

func TestNopReporter_DoesNotPanic(t *testing.T) {
	reporter := ui.NopReporter{}
	assert.NotPanics(t, func() {
		reporter.Start(cameradar.StepScan, "start")
		reporter.Done(cameradar.StepScan, "done")
		reporter.Progress(cameradar.StepScan, "progress")
		reporter.Debug(cameradar.StepScan, "debug")
		reporter.Error(cameradar.StepScan, assert.AnError)
		reporter.Summary(nil, nil)
		reporter.Close()
	})
}
