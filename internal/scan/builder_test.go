package scan_test

import (
	"net/netip"
	"testing"

	"github.com/Ullaakut/cameradar/v6"
	"github.com/Ullaakut/cameradar/v6/internal/scan"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew_UsesSkipScanner(t *testing.T) {
	config := scan.Config{
		SkipScan: true,
		Targets: []string{
			"192.0.2.0/30",
			"192.0.2.10-11",
		},
		Ports:     []string{"554", "8554-8555"},
		ScanSpeed: 4,
	}

	scanner, err := scan.New(config, nil)
	require.NoError(t, err)

	streams, err := scanner.Scan(t.Context())
	require.NoError(t, err)

	addrs := []netip.Addr{
		netip.MustParseAddr("192.0.2.0"),
		netip.MustParseAddr("192.0.2.1"),
		netip.MustParseAddr("192.0.2.2"),
		netip.MustParseAddr("192.0.2.3"),
		netip.MustParseAddr("192.0.2.10"),
		netip.MustParseAddr("192.0.2.11"),
	}
	portsExpected := []uint16{554, 8554, 8555}

	var expected []cameradar.Stream
	for _, addr := range addrs {
		for _, port := range portsExpected {
			expected = append(expected, cameradar.Stream{
				Address: addr,
				Port:    port,
			})
		}
	}

	assert.Equal(t, expected, streams)
}

func TestNew_SkipScanPropagatesErrors(t *testing.T) {
	config := scan.Config{
		SkipScan: true,
		Targets:  []string{"192.0.2.1"},
		Ports:    []string{"8555-8554"},
	}

	scanner, err := scan.New(config, nil)
	require.NoError(t, err)

	_, err = scanner.Scan(t.Context())
	require.Error(t, err)
	assert.ErrorContains(t, err, "invalid port range")
}
