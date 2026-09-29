package skip

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStreamCapacityEnforcesLimit(t *testing.T) {
	capacity, err := streamCapacity(maxExpandedStreams/4, 4)
	require.NoError(t, err)
	require.Equal(t, maxExpandedStreams, capacity)

	_, err = streamCapacity(maxExpandedStreams/4+1, 4)
	require.ErrorContains(t, err, "stream candidates")
}

func TestParseIPv4RangeRejectsExcessiveExpansion(t *testing.T) {
	addrs, ok, err := parseIPv4Range("0-255.0-255.0-255.0-255")
	require.True(t, ok)
	require.ErrorContains(t, err, "too large to expand")
	require.Nil(t, addrs)
}
