package output_test

import (
	"net/netip"
	"net/url"
	"strings"
	"testing"

	"github.com/Ullaakut/cameradar/v6"
	"github.com/Ullaakut/cameradar/v6/internal/output"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuildM3U_EncodesCredentials(t *testing.T) {
	stream := cameradar.Stream{
		Address:          netip.MustParseAddr("192.0.2.10"),
		Port:             554,
		Routes:           []string{"stream"},
		Username:         "admin",
		Password:         "pass/word",
		CredentialsFound: true,
	}

	playlist := output.BuildM3U([]cameradar.Stream{stream})

	var rtspLine string
	for _, line := range strings.Split(playlist, "\n") {
		if strings.HasPrefix(line, "rtsp") {
			rtspLine = line
			break
		}
	}
	require.NotEmpty(t, rtspLine)

	u, err := url.Parse(rtspLine)
	require.NoError(t, err)
	assert.Equal(t, "admin", u.User.Username())
	pass, _ := u.User.Password()
	assert.Equal(t, "pass/word", pass)
	assert.Equal(t, "192.0.2.10:554", u.Host)
	assert.Equal(t, "/stream", u.Path)
}

func TestBuildM3U_SanitizesDeviceLabelNewlines(t *testing.T) {
	stream := cameradar.Stream{
		Address: netip.MustParseAddr("192.0.2.20"),
		Port:    554,
		Routes:  []string{"stream"},
		Device:  "Cam\r\n#EXTINF:-1,Injected\nrtsp://attacker.example/evil\r",
	}

	playlist := output.BuildM3U([]cameradar.Stream{stream})

	extinfCount := 0
	for _, line := range strings.Split(playlist, "\n") {
		if strings.HasPrefix(line, "#EXTINF") {
			extinfCount++
		}
	}
	assert.Equal(t, 1, extinfCount, "device newlines must not inject extra #EXTINF entries")

	for _, line := range strings.Split(playlist, "\n") {
		assert.NotEqual(t, "rtsp://attacker.example/evil", strings.TrimSpace(line),
			"device newlines must not inject a standalone rtsp line")
	}
}

func TestBuildM3U_RendersNormalDeviceLabel(t *testing.T) {
	stream := cameradar.Stream{
		Address: netip.MustParseAddr("192.0.2.30"),
		Port:    554,
		Routes:  []string{"stream"},
		Device:  "Hikvision DS-2CD",
	}

	playlist := output.BuildM3U([]cameradar.Stream{stream})

	assert.Contains(t, playlist, "192.0.2.30:554 (Hikvision DS-2CD)")
}
