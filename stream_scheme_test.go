package cameradar_test

import (
	"net/netip"
	"strings"
	"testing"

	"github.com/Ullaakut/cameradar/v6"
	"github.com/stretchr/testify/require"
)

func TestStreamString_QueryRoute(t *testing.T) {
	addr := netip.MustParseAddr("192.168.1.1")
	tests := []struct {
		name      string
		route     string
		wantURL   string
		wantQuery string
	}{
		{
			name:      "route with query string",
			route:     "cam/realmonitor?channel=0&subtype=0",
			wantURL:   "rtsp://192.168.1.1:554/cam/realmonitor?channel=0&subtype=0",
			wantQuery: "channel=0&subtype=0",
		},
		{
			name:    "plain route unchanged",
			route:   "live/ch0",
			wantURL: "rtsp://192.168.1.1:554/live/ch0",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s := cameradar.Stream{
				Address: addr,
				Port:    554,
				Routes:  []string{test.route},
			}
			got := s.String()
			require.Equal(t, test.wantURL, got)
			require.False(t, strings.Contains(got, "%3F"), "URL must not percent-encode '?': %s", got)
		})
	}
}

func TestStreamRTSPScheme(t *testing.T) {
	tests := []struct {
		name   string
		scheme string
		want   string
	}{
		{name: "empty defaults to rtsp", scheme: "", want: "rtsp"},
		{name: "rtsp stays rtsp", scheme: "rtsp", want: "rtsp"},
		{name: "http maps to rtsp", scheme: "http", want: "rtsp"},
		{name: "https maps to rtsps", scheme: "https", want: "rtsps"},
		{name: "rtsps stays rtsps", scheme: "rtsps", want: "rtsps"},
		{name: "unknown falls back to rtsp", scheme: "custom", want: "rtsp"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			stream := cameradar.Stream{Scheme: test.scheme}
			got := stream.RTSPScheme()
			require.Equal(t, test.want, got)
		})
	}
}
