package ui

import (
	"errors"
	"strings"
	"testing"

	"github.com/Ullaakut/cameradar/v6"
	"github.com/stretchr/testify/assert"
)

func TestFinalViewRendersSummaryError(t *testing.T) {
	t.Run("appends error after summary", func(t *testing.T) {
		m := &modelState{
			steps:      cameradar.Steps(),
			status:     summaryStatusAllDone(),
			buildInfo:  BuildInfo{Version: "dev", Commit: "none"},
			summaryErr: errors.New("boom"),
		}

		view := m.FinalView()

		summaryIdx := strings.Index(view, "Summary - Streams")
		errIdx := strings.Index(view, "Error: boom")
		assert.GreaterOrEqual(t, summaryIdx, 0)
		assert.GreaterOrEqual(t, errIdx, 0)
		assert.Greater(t, errIdx, summaryIdx)
	})

	t.Run("omits error line when none set", func(t *testing.T) {
		m := &modelState{
			steps:     cameradar.Steps(),
			status:    summaryStatusAllDone(),
			buildInfo: BuildInfo{Version: "dev", Commit: "none"},
		}

		view := m.FinalView()

		assert.NotContains(t, view, "Error:")
	})
}
