package cameradar_test

import (
	"net/netip"
	"strings"
	"testing"

	"github.com/Ullaakut/cameradar/v6"
	"github.com/stretchr/testify/require"
)

func TestStreamURL_SchemeNormalization(t *testing.T) {
	addr := netip.MustParseAddr("192.168.0.10")
	tests := []struct {
		name            string
		scheme          string
		wantParsedScheme string
	}{
		{name: "lowercase http maps to rtsp", scheme: "http", wantParsedScheme: "rtsp"},
		{name: "lowercase https maps to rtsps", scheme: "https", wantParsedScheme: "rtsps"},
		{name: "uppercase HTTP maps to rtsp", scheme: "HTTP", wantParsedScheme: "rtsp"},
		{name: "uppercase HTTPS maps to rtsps", scheme: "HTTPS", wantParsedScheme: "rtsps"},
		{name: "mixed-case Http maps to rtsp", scheme: "Http", wantParsedScheme: "rtsp"},
		{name: "mixed-case Https maps to rtsps", scheme: "Https", wantParsedScheme: "rtsps"},
		{name: "lowercase rtsp stays rtsp", scheme: "rtsp", wantParsedScheme: "rtsp"},
		{name: "uppercase RTSP stays rtsp", scheme: "RTSP", wantParsedScheme: "rtsp"},
		{name: "lowercase rtsps stays rtsps", scheme: "rtsps", wantParsedScheme: "rtsps"},
		{name: "uppercase RTSPS stays rtsps", scheme: "RTSPS", wantParsedScheme: "rtsps"},
		{name: "empty defaults to rtsp", scheme: "", wantParsedScheme: "rtsp"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s := cameradar.Stream{
				Address: addr,
				Port:    554,
				Routes:  []string{"stream"},
				Scheme:  test.scheme,
			}
			u, err := s.URL()
			require.NoError(t, err)
			require.Equal(t, test.wantParsedScheme, u.Scheme)
		})
	}
}

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
