package ui_test

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/Ullaakut/cameradar/v6"
	"github.com/Ullaakut/cameradar/v6/internal/ui"
	"github.com/stretchr/testify/assert"
)

func TestPlainReporter_Outputs(t *testing.T) {
	t.Run("prints events", func(t *testing.T) {
		out := &bytes.Buffer{}
		reporter := ui.NewPlainReporter(out, true)

		reporter.Start(cameradar.StepScan, "starting")
		reporter.Progress(cameradar.StepScan, "working")
		reporter.Debug(cameradar.StepScan, "details")
		reporter.Done(cameradar.StepScan, "finished")
		reporter.Error(cameradar.StepScan, errors.New("boom"))
		reporter.Summary([]cameradar.Stream{}, nil)

		content := out.String()
		assert.Contains(t, content, " [STEP] Scan targets: starting")
		assert.Contains(t, content, " [INFO] Scan targets: working")
		assert.Contains(t, content, " [DBUG] Scan targets: details")
		assert.Contains(t, content, " [DONE] Scan targets: finished")
		assert.Contains(t, content, " [EROR] Scan targets: boom")
		assert.Contains(t, content, "Summary\n-------\nAccessible streams: 0")
	})

	t.Run("respects debug flag and empty input", func(t *testing.T) {
		out := &bytes.Buffer{}
		reporter := ui.NewPlainReporter(out, false)

		reporter.Debug(cameradar.StepScan, "hidden")
		reporter.Progress(cameradar.StepScan, "")
		reporter.Error(cameradar.StepScan, nil)

		content := out.String()
		assert.NotContains(t, content, "DBUG")
		assert.Equal(t, "", strings.TrimSpace(content))
	})
}
