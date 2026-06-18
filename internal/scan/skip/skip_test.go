package skip_test

import (
	"net/netip"
	"strconv"
	"testing"

	"github.com/Ullaakut/cameradar/v6/internal/scan/skip"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew_ExpandsTargetsAndPorts(t *testing.T) {
	targets := []string{
		"192.0.2.0/30",
		"192.0.2.15",
		"192.0.2.10-11",
	}
	ports := []string{"554", "322", "8554-8555"}

	scanner := skip.New(targets, ports)

	streams, err := scanner.Scan(t.Context())
	require.NoError(t, err)

	addrs := []netip.Addr{
		netip.MustParseAddr("192.0.2.0"),
		netip.MustParseAddr("192.0.2.1"),
		netip.MustParseAddr("192.0.2.2"),
		netip.MustParseAddr("192.0.2.3"),
		netip.MustParseAddr("192.0.2.10"),
		netip.MustParseAddr("192.0.2.11"),
		netip.MustParseAddr("192.0.2.15"),
	}
	portsExpected := []uint16{554, 322, 8554, 8555}

	var want []string
	for _, addr := range addrs {
		for _, port := range portsExpected {
			want = append(want, addr.String()+":"+strconv.Itoa(int(port)))
		}
	}

	var got []string
	for _, stream := range streams {
		got = append(got, stream.Address.String()+":"+strconv.Itoa(int(stream.Port)))

		if stream.Port == 322 || stream.Port == 8322 {
			assert.Equal(t, "rtsps", stream.Scheme)
		} else {
			assert.Empty(t, stream.Scheme)
		}

	}

	assert.ElementsMatch(t, want, got)
}

func TestNew_ExpandsPortRangeEndingAt65535(t *testing.T) {
	scanner := skip.New([]string{"192.0.2.1"}, []string{"65534-65535"})

	streams, err := scanner.Scan(t.Context())
	require.NoError(t, err)
	require.Len(t, streams, 2)

	var got []uint16
	for _, stream := range streams {
		got = append(got, stream.Port)
	}
	assert.ElementsMatch(t, []uint16{65534, 65535}, got)
}

func TestNew_ReturnsErrorOnInvalidPortRange(t *testing.T) {
	scanner := skip.New([]string{"192.0.2.1"}, []string{"8555-8554"})

	_, err := scanner.Scan(t.Context())
	require.Error(t, err)
	assert.ErrorContains(t, err, "invalid port range")
}

func TestNew_ReturnsErrorOnEmptyTargets(t *testing.T) {
	scanner := skip.New([]string{}, []string{"554"})

	_, err := scanner.Scan(t.Context())
	require.Error(t, err)
	assert.ErrorContains(t, err, "no valid target addresses resolved")
}

func TestNew_ResolvesServicePorts(t *testing.T) {
	scanner := skip.New([]string{"127.0.0.1"}, []string{"http"})

	streams, err := scanner.Scan(t.Context())
	require.NoError(t, err)
	require.Len(t, streams, 1)

	assert.Equal(t, netip.MustParseAddr("127.0.0.1"), streams[0].Address)
	assert.Equal(t, uint16(80), streams[0].Port)
}

func TestNew_ResolvesHyphenatedServicePortNames(t *testing.T) {
	// Service names that contain a hyphen (e.g. "http-alt") must be resolved via
	// LookupPort rather than rejected as malformed numeric ranges.
	scanner := skip.New([]string{"127.0.0.1"}, []string{"http-alt"})

	streams, err := scanner.Scan(t.Context())
	require.NoError(t, err)
	require.Len(t, streams, 1)
	assert.Equal(t, uint16(8080), streams[0].Port)
}

func TestNew_ReturnsErrorOnUnknownServicePort(t *testing.T) {
	scanner := skip.New([]string{"127.0.0.1"}, []string{"not-a-service"})

	_, err := scanner.Scan(t.Context())
	require.Error(t, err)
	assert.ErrorContains(t, err, "invalid port")
}

func TestNew_ResolvesHostnames(t *testing.T) {
	scanner := skip.New([]string{"localhost"}, []string{"554"})

	streams, err := scanner.Scan(t.Context())
	require.NoError(t, err)
	require.NotEmpty(t, streams)
	addr := streams[0].Address
	assert.True(t,
		addr == netip.MustParseAddr("127.0.0.1") || addr == netip.MustParseAddr("::1"),
		"expected localhost to resolve to 127.0.0.1 or ::1, got %s",
		addr.String(),
	)
}

func TestNew_ReturnsErrorOnHostnameLookupFailure(t *testing.T) {
	scanner := skip.New([]string{"does-not-exist.invalid"}, []string{"554"})

	_, err := scanner.Scan(t.Context())
	require.Error(t, err)
	assert.ErrorContains(t, err, "resolving hostname")
}

func TestNew_ExpandsPrefixesWithinBounds(t *testing.T) {
	tests := []struct {
		name      string
		target    string
		wantAddrs int
		wantErr   string
	}{
		{
			name:      "bounded ipv6 /120 expands to 256 addresses",
			target:    "2001:db8::/120",
			wantAddrs: 256,
		},
		{
			name:    "oversized ipv6 /64 returns an error",
			target:  "2001:db8::/64",
			wantErr: "too large to expand",
		},
		{
			name:      "ipv4 /24 expands to 256 addresses",
			target:    "192.0.2.0/24",
			wantAddrs: 256,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			scanner := skip.New([]string{test.target}, []string{"554"})

			streams, err := scanner.Scan(t.Context())
			if test.wantErr != "" {
				require.Error(t, err)
				assert.ErrorContains(t, err, test.wantErr)
				return
			}

			require.NoError(t, err)
			assert.Len(t, streams, test.wantAddrs)
		})
	}
}
