package ui_test

import (
	"net/netip"
	"net/url"
	"strings"
	"testing"

	"github.com/Ullaakut/cameradar/v6"
	"github.com/Ullaakut/cameradar/v6/internal/ui"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFormatSummaryIPv6BracketsHost(t *testing.T) {
	streams := []cameradar.Stream{
		{
			Address:            netip.MustParseAddr("fe80::1"),
			Port:               554,
			Available:          true,
			RouteFound:         true,
			Routes:             []string{"stream"},
			CredentialsFound:   true,
			Username:           "user",
			Password:           "pass",
			AuthenticationType: cameradar.AuthBasic,
		},
	}

	got := ui.FormatSummary(streams)

	// IPv6 hosts must be bracketed to form valid URLs.
	assert.Contains(t, got, "RTSP URL: rtsp://user:pass@[fe80::1]:554/stream")
	assert.Contains(t, got, "Admin panel: http://[fe80::1]/")
}

func TestFormatSummary(t *testing.T) {
	tests := []struct {
		name            string
		streams         []cameradar.Stream
		wantContains    []string
		wantNotContains []string
		orderedPairs    [][2]string
	}{
		{
			name:    "empty",
			streams: nil,
			wantContains: []string{
				"Accessible streams: 0",
				"• None",
			},
			wantNotContains: []string{
				"Other discovered streams",
				"Error:",
			},
		},
		{
			name: "mixed streams",
			streams: []cameradar.Stream{
				{
					Device:             "Model B",
					Address:            netip.MustParseAddr("10.0.0.2"),
					Port:               554,
					Available:          true,
					AuthenticationType: cameradar.AuthNone,
				},
				{
					Device:             "Model A",
					Address:            netip.MustParseAddr("10.0.0.1"),
					Port:               8554,
					Available:          true,
					Routes:             []string{"stream1", "stream2"},
					RouteFound:         true,
					CredentialsFound:   true,
					Username:           "user",
					Password:           "pass",
					AuthenticationType: cameradar.AuthBasic,
				},
				{
					Address:            netip.MustParseAddr("10.0.0.3"),
					Port:               554,
					Available:          false,
					AuthenticationType: cameradar.AuthDigest,
				},
			},
			wantContains: []string{
				"Accessible streams: 2",
				"Other discovered streams: 1",
				"• 10.0.0.1:8554 (Model A)",
				"• 10.0.0.2:554 (Model B)",
				"• 10.0.0.3:554",
				"Authentication: basic",
				"Authentication: none",
				"Authentication: digest",
				"Routes: stream1, stream2",
				"Credentials: user:pass",
				"Credentials: none",
				"RTSP URL: rtsp://user:pass@10.0.0.1:8554/stream1",
				"Admin panel: http://10.0.0.1/",
				"Admin panel: http://10.0.0.2/",
			},
			wantNotContains: []string{
				"Error:",
			},
			orderedPairs: [][2]string{
				{"• 10.0.0.1:8554", "• 10.0.0.2:554"},
			},
		},
		{
			name: "empty discovered credentials render as none",
			streams: []cameradar.Stream{
				{
					Address:            netip.MustParseAddr("10.0.0.4"),
					Port:               554,
					Available:          true,
					RouteFound:         true,
					Routes:             []string{"stream"},
					CredentialsFound:   true,
					AuthenticationType: cameradar.AuthNone,
				},
			},
			wantContains: []string{
				"Accessible streams: 1",
				"Credentials: none",
				"RTSP URL: rtsp://10.0.0.4:554/stream",
			},
			wantNotContains: []string{
				"Credentials: :",
				"rtsp://:@10.0.0.4:554/stream",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := ui.FormatSummary(test.streams)

			for _, expected := range test.wantContains {
				assert.Contains(t, got, expected)
			}
			for _, unexpected := range test.wantNotContains {
				assert.NotContains(t, got, unexpected)
			}
			for _, pair := range test.orderedPairs {
				first := strings.Index(got, pair[0])
				second := strings.Index(got, pair[1])
				assert.True(t, first >= 0 && second >= 0 && first < second)
			}
		})
	}
}

func TestFormatSummaryEncodesCredentials(t *testing.T) {
	// The default credentials dictionary ships passwords with characters that
	// are not URL-safe (for example "admin pass"). The summary RTSP URL must
	// percent-encode them so it stays parseable.
	streams := []cameradar.Stream{
		{
			Address:            netip.MustParseAddr("10.0.0.1"),
			Port:               554,
			Available:          true,
			RouteFound:         true,
			Routes:             []string{"stream"},
			CredentialsFound:   true,
			Username:           "ad@min",
			Password:           "p@ss word",
			AuthenticationType: cameradar.AuthBasic,
		},
	}

	got := ui.FormatSummary(streams)

	const prefix = "RTSP URL: "
	idx := strings.Index(got, prefix)
	require.GreaterOrEqual(t, idx, 0, "summary missing RTSP URL line:\n%s", got)
	rest := got[idx+len(prefix):]
	rawURL := strings.TrimSpace(rest[:strings.IndexByte(rest, '\n')])

	u, err := url.Parse(rawURL)
	require.NoError(t, err, "RTSP URL is not parseable: %q", rawURL)
	require.NotNil(t, u.User, "RTSP URL has no userinfo: %q", rawURL)
	pw, _ := u.User.Password()
	assert.Equal(t, "ad@min", u.User.Username())
	assert.Equal(t, "p@ss word", pw)
}
