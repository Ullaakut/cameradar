package attack

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/bluenviron/gortsplib/v5/pkg/base"
)

// listenAcceptNoReply starts a loopback TCP server that accepts connections but
// never sends any reply, so a DESCRIBE probe blocks until its deadline fires.
func listenAcceptNoReply(t *testing.T) (string, func()) {
	t.Helper()

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}

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

	return ln.Addr().String(), func() {
		close(done)
		_ = ln.Close()
	}
}

func TestProbeDescribeHeadersTimeoutZeroDoesNotFailInstantly(t *testing.T) {
	addr, cleanup := listenAcceptNoReply(t)
	defer cleanup()

	a := Attacker{timeout: 0}
	u := &base.URL{Scheme: schemeRTSP, Host: addr, Path: "/"}

	start := time.Now()
	_, _, err := a.probeDescribeHeaders(context.Background(), u)
	elapsed := time.Since(start)

	if err == nil {
		t.Fatal("expected an error from the unresponsive server, got nil")
	}

	// With the fix, timeout==0 is clamped to the default (2s). Before the fix the
	// deadline was set to time.Now() and the probe failed in ~225µs.
	if elapsed < time.Second {
		t.Fatalf("probe with timeout=0 failed too fast (%v); expected to block ~2s on real deadline", elapsed)
	}
	if elapsed > 5*time.Second {
		t.Fatalf("probe with timeout=0 took too long (%v)", elapsed)
	}
}

func TestProbeDescribeHeadersTimeoutHonored(t *testing.T) {
	addr, cleanup := listenAcceptNoReply(t)
	defer cleanup()

	a := Attacker{timeout: 300 * time.Millisecond}
	u := &base.URL{Scheme: schemeRTSP, Host: addr, Path: "/"}

	start := time.Now()
	_, _, err := a.probeDescribeHeaders(context.Background(), u)
	elapsed := time.Since(start)

	if err == nil {
		t.Fatal("expected an error from the unresponsive server, got nil")
	}

	if elapsed < 150*time.Millisecond {
		t.Fatalf("probe failed too fast (%v); expected ~300ms", elapsed)
	}
	if elapsed > time.Second {
		t.Fatalf("probe took too long (%v); expected ~300ms", elapsed)
	}
}
