package ui_test

import (
	"errors"
	"net/netip"
	"strings"
	"testing"

	"github.com/Ullaakut/cameradar/v6"
	"github.com/Ullaakut/cameradar/v6/internal/ui"
	"github.com/stretchr/testify/assert"
)

func TestFormatSummary(t *testing.T) {
	tests := []struct {
		name            string
		streams         []cameradar.Stream
		err             error
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
			name: "mixed streams with error",
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
			err: errors.New("boom"),
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
				"RTSP URL: rtsp://user:pass@10.0.0.1:8554/stream1",
				"Admin panel: http://10.0.0.1/",
				"Admin panel: http://10.0.0.2/",
			},
			wantNotContains: []string{
				"RTSP URL: rtsp://10.0.0.2",
				"Error:",
			},
			orderedPairs: [][2]string{
				{"• 10.0.0.1:8554", "• 10.0.0.2:554"},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := ui.FormatSummary(test.streams, test.err)

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
