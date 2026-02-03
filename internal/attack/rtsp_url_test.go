package attack

import (
	"net/netip"
	"testing"

	"github.com/Ullaakut/cameradar/v6"
	"github.com/stretchr/testify/require"
)

func TestBuildRTSPURL(t *testing.T) {
	stream := cameradar.Stream{
		Address: netip.MustParseAddr("192.168.0.10"),
		Port:    554,
	}

	tests := []struct {
		name     string
		route    string
		username string
		password string
		wantURL  string
	}{
		{
			name:    "empty route",
			wantURL: "rtsp://192.168.0.10:554/",
		},
		{
			name:    "root route",
			route:   "/",
			wantURL: "rtsp://192.168.0.10:554/",
		},
		{
			name:    "multiple leading slashes",
			route:   "////",
			wantURL: "rtsp://192.168.0.10:554/",
		},
		{
			name:    "route with no leading slash",
			route:   "stream",
			wantURL: "rtsp://192.168.0.10:554/stream",
		},
		{
			name:    "route with leading slash",
			route:   "/stream",
			wantURL: "rtsp://192.168.0.10:554/stream",
		},
		{
			name:    "route with trailing slash",
			route:   "stream/",
			wantURL: "rtsp://192.168.0.10:554/stream/",
		},
		{
			name:    "route with spaces",
			route:   "  /stream  ",
			wantURL: "rtsp://192.168.0.10:554/stream",
		},
		{
			name:     "username and password",
			route:    "stream",
			username: "admin",
			password: "admin123",
			wantURL:  "rtsp://admin:admin123@192.168.0.10:554/stream",
		},
		{
			name:     "empty username with password",
			route:    "stream",
			password: "pass",
			wantURL:  "rtsp://:pass@192.168.0.10:554/stream",
		},
		{
			name:     "username only",
			route:    "stream",
			username: "user",
			wantURL:  "rtsp://user:@192.168.0.10:554/stream",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, gotURL, err := buildRTSPURL(stream, test.route, test.username, test.password)
			require.NoError(t, err)
			require.Equal(t, test.wantURL, gotURL)
		})
	}
}
