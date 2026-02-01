package attack

import (
	"bufio"
	"fmt"
	"net"
	"net/netip"
	"strings"
	"testing"
	"time"

	"github.com/Ullaakut/cameradar/v6"
	"github.com/Ullaakut/cameradar/v6/internal/ui"
	"github.com/bluenviron/gortsplib/v5/pkg/base"
	"github.com/bluenviron/gortsplib/v5/pkg/headers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testDictionary struct {
	routes    []string
	usernames []string
	passwords []string
}

func (d testDictionary) Routes() []string {
	return d.routes
}

func (d testDictionary) Usernames() []string {
	return d.usernames
}

func (d testDictionary) Passwords() []string {
	return d.passwords
}

func TestAuthTypeFromHeaders(t *testing.T) {
	tests := []struct {
		name   string
		values base.HeaderValue
		want   cameradar.AuthType
	}{
		{
			name: "digest wins over basic",
			values: base.HeaderValue{
				headers.Authenticate{Method: headers.AuthMethodBasic, Realm: "cam"}.Marshal()[0],
				headers.Authenticate{Method: headers.AuthMethodDigest, Realm: "cam", Nonce: "nonce"}.Marshal()[0],
			},
			want: cameradar.AuthDigest,
		},
		{
			name:   "basic auth",
			values: headers.Authenticate{Method: headers.AuthMethodBasic, Realm: "cam"}.Marshal(),
			want:   cameradar.AuthBasic,
		},
		{
			name:   "digest auth",
			values: headers.Authenticate{Method: headers.AuthMethodDigest, Realm: "cam", Nonce: "nonce"}.Marshal(),
			want:   cameradar.AuthDigest,
		},
		{
			name:   "unknown with empty values",
			values: nil,
			want:   cameradar.AuthUnknown,
		},
		{
			name:   "unknown with unsupported header",
			values: base.HeaderValue{"Bearer abc"},
			want:   cameradar.AuthUnknown,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.want, authTypeFromHeaders(test.values))
		})
	}
}

func TestDetectAuthMethod(t *testing.T) {
	tests := []struct {
		name       string
		statusCode base.StatusCode
		headers    base.Header
		want       cameradar.AuthType
	}{
		{
			name:       "no auth when status ok",
			statusCode: base.StatusOK,
			headers: base.Header{
				"WWW-Authenticate": headers.Authenticate{Method: headers.AuthMethodBasic, Realm: "cam"}.Marshal(),
			},
			want: cameradar.AuthNone,
		},
		{
			name:       "basic auth on unauthorized",
			statusCode: base.StatusUnauthorized,
			headers: base.Header{
				"WWW-Authenticate": headers.Authenticate{Method: headers.AuthMethodBasic, Realm: "cam"}.Marshal(),
			},
			want: cameradar.AuthBasic,
		},
		{
			name:       "digest auth on unauthorized",
			statusCode: base.StatusUnauthorized,
			headers: base.Header{
				"WWW-Authenticate": headers.Authenticate{Method: headers.AuthMethodDigest, Realm: "cam", Nonce: "nonce"}.Marshal(),
			},
			want: cameradar.AuthDigest,
		},
		{
			name:       "unknown auth on unauthorized without www-authenticate",
			statusCode: base.StatusUnauthorized,
			headers:    nil,
			want:       cameradar.AuthUnknown,
		},
		{
			name:       "unknown auth on other status",
			statusCode: base.StatusNotFound,
			headers:    nil,
			want:       cameradar.AuthUnknown,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			addr, port := startRTSPProbeServer(t, test.statusCode, test.headers)

			attacker, err := New(testDictionary{}, 0, time.Second, ui.NopReporter{})
			require.NoError(t, err)

			stream := cameradar.Stream{
				Address: addr,
				Port:    port,
			}

			got, err := attacker.detectAuthMethod(t.Context(), stream)
			require.NoError(t, err)
			assert.Equal(t, test.want, got.AuthenticationType)
		})
	}
}

func startRTSPProbeServer(t *testing.T, statusCode base.StatusCode, headers base.Header) (netip.Addr, uint16) {
	t.Helper()

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)

	t.Cleanup(func() {
		_ = listener.Close()
	})

	go func() {
		conn, err := listener.Accept()
		if err != nil {
			return
		}
		defer conn.Close()

		_ = conn.SetDeadline(time.Now().Add(time.Second))

		reader := bufio.NewReader(conn)
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				return
			}
			if strings.TrimSpace(line) == "" {
				break
			}
		}

		statusText := statusTextFromCode(statusCode)

		var builder strings.Builder
		_, _ = fmt.Fprintf(&builder, "RTSP/1.0 %d %s\r\n", statusCode, statusText)
		builder.WriteString("CSeq: 1\r\n")
		for key, values := range headers {
			for _, value := range values {
				_, _ = fmt.Fprintf(&builder, "%s: %s\r\n", key, value)
			}
		}
		builder.WriteString("Content-Length: 0\r\n\r\n")

		_, _ = conn.Write([]byte(builder.String()))
	}()

	tcpAddr, ok := listener.Addr().(*net.TCPAddr)
	require.True(t, ok)

	return netip.MustParseAddr("127.0.0.1"), uint16(tcpAddr.Port)
}

func statusTextFromCode(code base.StatusCode) string {
	switch code {
	case base.StatusOK:
		return "OK"
	case base.StatusUnauthorized:
		return "Unauthorized"
	case base.StatusNotFound:
		return "Not Found"
	default:
		return "Unknown"
	}
}
