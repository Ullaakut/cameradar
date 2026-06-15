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
