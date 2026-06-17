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

// TestProbeDescribeHeaders_EOFWithoutBlankLine tests that a peer closing the
// TCP connection right after the last header (without a trailing blank CRLF
// line, which is common for non-conforming RTSP cameras) does not cause
// probeDescribeHeaders to discard the already-parsed status code and headers.
func TestProbeDescribeHeaders_EOFWithoutBlankLine(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)

	done := make(chan struct{})
	t.Cleanup(func() {
		close(done)
		_ = ln.Close()
	})

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
			// Stop reading once we see the end of the RTSP request.
			data := string(buf[:n])
			if len(data) >= 4 && data[len(data)-4:] == "\r\n\r\n" {
				break
			}
		}

		// Send a 401 response with WWW-Authenticate, then close WITHOUT a
		// trailing blank CRLF line, simulating a non-conforming RTSP camera.
		resp := "RTSP/1.0 401 Unauthorized\r\nWWW-Authenticate: Digest realm=\"cam\", nonce=\"x\"\r\n"
		_, _ = fmt.Fprint(conn, resp)
		// Close immediately; no trailing \r\n.
	}()

	addr := ln.Addr().String()
	a := Attacker{timeout: time.Second}
	u := &base.URL{Scheme: schemeRTSP, Host: addr, Path: "/"}

	statusCode, headers, err := a.probeDescribeHeaders(context.Background(), u)

	require.NoError(t, err, "EOF without trailing blank line must not be returned as an error")
	assert.Equal(t, base.StatusCode(401), statusCode, "status code must be 401")
	require.NotNil(t, headers, "headers must not be nil")

	wwwAuth := headerValues(headers, "WWW-Authenticate")
	assert.NotEmpty(t, wwwAuth, "WWW-Authenticate header must be present")
}
