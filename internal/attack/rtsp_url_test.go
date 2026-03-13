package attack

import (
	"net/netip"
	"net/url"
	"testing"

	"github.com/Ullaakut/cameradar/v6"
	"github.com/stretchr/testify/require"
)

func TestStreamURL(t *testing.T) {
	tests := []struct {
		name    string
		stream  cameradar.Stream
		wantURL string
	}{
		{
			name: "empty route",
			stream: cameradar.Stream{
				Address: netip.MustParseAddr("192.168.0.10"),
				Port:    554,
			},
			wantURL: "rtsp://192.168.0.10:554/",
		},
		{
			name: "root route",
			stream: cameradar.Stream{
				Address: netip.MustParseAddr("192.168.0.10"),
				Port:    554,
				Routes:  []string{"/"},
			},
			wantURL: "rtsp://192.168.0.10:554/",
		},
		{
			name: "multiple leading slashes",
			stream: cameradar.Stream{
				Address: netip.MustParseAddr("192.168.0.10"),
				Port:    554,
				Routes:  []string{"////"},
			},
			wantURL: "rtsp://192.168.0.10:554/",
		},
		{
			name: "route with no leading slash",
			stream: cameradar.Stream{
				Address: netip.MustParseAddr("192.168.0.10"),
				Port:    554,
				Routes:  []string{"stream"},
			},
			wantURL: "rtsp://192.168.0.10:554/stream",
		},
		{
			name: "route with leading slash",
			stream: cameradar.Stream{
				Address: netip.MustParseAddr("192.168.0.10"),
				Port:    554,
				Routes:  []string{"/stream"},
			},
			wantURL: "rtsp://192.168.0.10:554/stream",
		},
		{
			name: "route with trailing slash",
			stream: cameradar.Stream{
				Address: netip.MustParseAddr("192.168.0.10"),
				Port:    554,
				Routes:  []string{"stream/"},
			},
			wantURL: "rtsp://192.168.0.10:554/stream/",
		},
		{
			name: "route with spaces",
			stream: cameradar.Stream{
				Address: netip.MustParseAddr("192.168.0.10"),
				Port:    554,
				Routes:  []string{"  /stream  "},
			},
			wantURL: "rtsp://192.168.0.10:554/stream",
		},
		{
			name: "username and password",
			stream: cameradar.Stream{
				Address:  netip.MustParseAddr("192.168.0.10"),
				Port:     554,
				Routes:   []string{"stream"},
				Username: "admin",
				Password: "admin123",
			},
			wantURL: "rtsp://admin:admin123@192.168.0.10:554/stream",
		},
		{
			name: "empty username with password",
			stream: cameradar.Stream{
				Address:  netip.MustParseAddr("192.168.0.10"),
				Port:     554,
				Routes:   []string{"stream"},
				Password: "pass",
			},
			wantURL: "rtsp://:pass@192.168.0.10:554/stream",
		},
		{
			name: "username only",
			stream: cameradar.Stream{
				Address:  netip.MustParseAddr("192.168.0.10"),
				Port:     554,
				Routes:   []string{"stream"},
				Username: "user",
			},
			wantURL: "rtsp://user:@192.168.0.10:554/stream",
		},
		{
			name: "http scheme",
			stream: cameradar.Stream{
				Address: netip.MustParseAddr("192.168.0.10"),
				Port:    554,
				Routes:  []string{"stream"},
				Scheme:  "http",
			},
			wantURL: "http://192.168.0.10:554/stream",
		},
		{
			name: "https scheme",
			stream: cameradar.Stream{
				Address: netip.MustParseAddr("192.168.0.10"),
				Port:    554,
				Routes:  []string{"stream"},
				Scheme:  "https",
			},
			wantURL: "https://192.168.0.10:554/stream",
		},
		{
			name: "rtsps scheme",
			stream: cameradar.Stream{
				Address: netip.MustParseAddr("192.168.0.10"),
				Port:    554,
				Routes:  []string{"stream"},
				Scheme:  "rtsps",
			},
			wantURL: "rtsps://192.168.0.10:554/stream",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			gotURL := test.stream.String()
			require.Equal(t, test.wantURL, gotURL)

			parsedURL, err := test.stream.URL()
			require.NoError(t, err)

			expectedURL, err := url.Parse(test.wantURL)
			require.NoError(t, err)
			require.Equal(t, expectedURL.Scheme, parsedURL.Scheme)
		})
	}
}
