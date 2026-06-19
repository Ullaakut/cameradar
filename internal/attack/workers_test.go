package attack

import (
	"context"
	"net/netip"
	"sync/atomic"
	"testing"

	"github.com/Ullaakut/cameradar/v6"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestRunParallel_CancelledContextSurfacesError verifies that when the caller's
// context is cancelled before all targets have been queued, runParallel returns
// a non-nil error rather than silently returning incomplete (stale pre-attack)
// results.  Pre-fix the function returned nil error even though some targets
// were never processed.
func TestRunParallel_CancelledContextSurfacesError(t *testing.T) {
	const total = 10

	ctx, cancel := context.WithCancel(t.Context())

	var processed atomic.Int32
	fn := func(fnCtx context.Context, s cameradar.Stream) (cameradar.Stream, error) {
		// Cancel after the first job starts so that subsequent targets are
		// never queued.  The fn itself still succeeds for this one job.
		if processed.Add(1) == 1 {
			cancel()
		}
		s.RouteFound = true
		return s, nil
	}

	targets := make([]cameradar.Stream, total)
	for i := range targets {
		targets[i] = cameradar.Stream{
			Address: netip.MustParseAddr("127.0.0.1"),
			Port:    uint16(8554 + i),
		}
	}

	_, err := runParallel(ctx, targets, fn)

	// Pre-fix: err == nil even though only 1 of 10 targets was processed.
	// Post-fix: err wraps context.Canceled and tells the caller results are partial.
	require.Error(t, err, "cancelled context must surface an error so callers know results are incomplete")
	assert.ErrorIs(t, err, context.Canceled)
}
