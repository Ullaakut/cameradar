package attack

import (
	"context"
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/bluenviron/gortsplib/v5/pkg/base"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// serveFakeRTSP starts a loopback TCP server that, for each accepted connection,
// reads the DESCRIBE request then replies with the provided raw response bytes.
// It returns the listener address and registers cleanup via t.Cleanup.
func serveFakeRTSP(t *testing.T, response string) string {
	t.Helper()

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)

	go func() {
		conn, err := ln.Accept()
		if err != nil {
			return
		}
		defer conn.Close()

		_ = conn.SetDeadline(time.Now().Add(2 * time.Second))

		// Drain the incoming DESCRIBE request.
		buf := make([]byte, 4096)
		for {
			n, err := conn.Read(buf)
			if err != nil || n == 0 {
				break
			}
			if n >= 4 && string(buf[n-4:n]) == "\r\n\r\n" {
				break
			}
		}

		_, _ = fmt.Fprint(conn, response)
	}()

	t.Cleanup(func() { _ = ln.Close() })
	return ln.Addr().String()
}

// TestProbeDescribeHeaders_NonRTSPStatusLine verifies that
// probeDescribeHeaders rejects a response whose status-line protocol token
// is not "RTSP/1.0".  Before the fix, the parser used fields[1] as the
// status code and ignored fields[0] entirely, so an HTTP/1.1 401 response
// would be silently accepted and trigger false-positive auth detection.
func TestProbeDescribeHeaders_NonRTSPStatusLine(t *testing.T) {
	tests := []struct {
		name     string
		response string
	}{
		{
			name: "HTTP/1.1 401 response accepted as RTSP",
			response: "HTTP/1.1 401 Unauthorized\r\n" +
				"WWW-Authenticate: Basic realm=\"cam\"\r\n" +
				"Content-Length: 0\r\n\r\n",
		},
		{
			name: "HTTP/1.0 200 response accepted as RTSP",
			response: "HTTP/1.0 200 OK\r\n" +
				"Content-Length: 0\r\n\r\n",
		},
		{
			name: "SIP/2.0 response accepted as RTSP",
			response: "SIP/2.0 401 Unauthorized\r\n" +
				"WWW-Authenticate: Digest realm=\"cam\",nonce=\"abc\"\r\n" +
				"Content-Length: 0\r\n\r\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			addr := serveFakeRTSP(t, tt.response)

			a := Attacker{timeout: time.Second}
			u := &base.URL{Scheme: schemeRTSP, Host: addr, Path: "/"}

			statusCode, _, err := a.probeDescribeHeaders(context.Background(), u)

			// Pre-fix: err is nil and statusCode is 401/200 (false positive).
			// Post-fix: err must be non-nil and statusCode must be 0.
			require.Error(t, err, "expected error for non-RTSP status line, got statusCode=%d", statusCode)
			assert.Contains(t, err.Error(), "RTSP")
		})
	}
}

// TestProbeDescribeHeaders_AcceptsNonRFCRTSPVersion verifies that a status line
// whose protocol token is an RTSP version other than the exact "RTSP/1.0"
// spelling (e.g. RTSP/1.1) is still accepted. Many real cameras do not follow
// the RFC to the letter, so matching the RTSP/ family rather than a fixed
// version avoids false negatives while still rejecting non-RTSP protocols.
func TestProbeDescribeHeaders_AcceptsNonRFCRTSPVersion(t *testing.T) {
	tests := []struct {
		name     string
		response string
		wantCode base.StatusCode
	}{
		{
			name: "RTSP/1.1 401 accepted",
			response: "RTSP/1.1 401 Unauthorized\r\n" +
				"WWW-Authenticate: Digest realm=\"cam\",nonce=\"abc\"\r\n" +
				"Content-Length: 0\r\n\r\n",
			wantCode: 401,
		},
		{
			name: "RTSP/2.0 200 accepted",
			response: "RTSP/2.0 200 OK\r\n" +
				"Content-Length: 0\r\n\r\n",
			wantCode: 200,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			addr := serveFakeRTSP(t, tt.response)

			a := Attacker{timeout: time.Second}
			u := &base.URL{Scheme: schemeRTSP, Host: addr, Path: "/"}

			statusCode, _, err := a.probeDescribeHeaders(context.Background(), u)

			require.NoError(t, err)
			assert.Equal(t, tt.wantCode, statusCode)
		})
	}
}
