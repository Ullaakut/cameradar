package attack

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/bluenviron/gortsplib/v5/pkg/base"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProbeDescribeHeadersTimeoutZeroDoesNotFailInstantly(t *testing.T) {
	addr := listenAcceptNoReply(t)

	a := Attacker{timeout: 0}
	u := &base.URL{Scheme: schemeRTSP, Host: addr, Path: "/"}

	start := time.Now()
	_, _, err := a.probeDescribeHeaders(context.Background(), u)
	elapsed := time.Since(start)

	require.Error(t, err, "expected an error from the unresponsive server")

	// With the fix, timeout==0 is clamped to the default (2s). Before the fix the
	// deadline was set to time.Now() and the probe failed in ~225µs.
	assert.Greater(t, elapsed, time.Second, "probe with timeout=0 failed too fast; expected to block ~2s on real deadline")
	assert.Less(t, elapsed, 5*time.Second, "probe with timeout=0 took too long")
}

func TestProbeDescribeHeadersTimeoutHonored(t *testing.T) {
	addr := listenAcceptNoReply(t)

	a := Attacker{timeout: 300 * time.Millisecond}
	u := &base.URL{Scheme: schemeRTSP, Host: addr, Path: "/"}

	start := time.Now()
	_, _, err := a.probeDescribeHeaders(context.Background(), u)
	elapsed := time.Since(start)

	require.Error(t, err, "expected an error from the unresponsive server")

	assert.Greater(t, elapsed, 150*time.Millisecond, "probe failed too fast; expected ~300ms")
	assert.Less(t, elapsed, time.Second, "probe took too long; expected ~300ms")
}

// listenAcceptNoReply starts a loopback TCP server that accepts connections but
// never sends any reply, so a DESCRIBE probe blocks until its deadline fires. It
// returns the server address and registers cleanup via t.Cleanup.
func listenAcceptNoReply(t *testing.T) string {
	t.Helper()

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err, "listen")

	done := make(chan struct{})
	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				return
			}
			// Hold the connection open without replying.
			go func(c net.Conn) {
				<-done
				_ = c.Close()
			}(conn)
		}
	}()

	t.Cleanup(func() {
		close(done)
		_ = ln.Close()
	})

	return ln.Addr().String()
}
