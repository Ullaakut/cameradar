package output

import (
	"net/netip"
	"net/url"
	"strings"
	"testing"

	"github.com/Ullaakut/cameradar/v6"
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

	playlist := BuildM3U([]cameradar.Stream{stream})

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
	require.Equal(t, "admin", u.User.Username())
	pass, _ := u.User.Password()
	require.Equal(t, "pass/word", pass)
	require.Equal(t, "192.0.2.10:554", u.Host)
	require.Equal(t, "/stream", u.Path)
}
