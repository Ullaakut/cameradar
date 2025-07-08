package scan

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExpandTargetsForScan_ExpandsFullIPv4Range(t *testing.T) {
	targets := []string{
		"192.0.2.10-192.0.2.12",
		"192.168.1.140-255",
		"192.0.2.0/30",
		"localhost",
		"",
	}

	got, err := expandTargetsForScan(targets)
	require.NoError(t, err)

	assert.ElementsMatch(t, []string{
		"192.0.2.10/31",
		"192.0.2.12",
		"192.168.1.140-255",
		"192.0.2.0/30",
		"localhost",
	}, got)
}

func TestExpandTargetsForScan_ReturnsErrorOnInvalidRange(t *testing.T) {
	t.Run("inverted range", func(t *testing.T) {
		_, err := expandTargetsForScan([]string{"192.0.2.12-192.0.2.10"})
		require.Error(t, err)
		assert.ErrorContains(t, err, "invalid range")
	})

	t.Run("invalid range", func(t *testing.T) {
		_, err := expandTargetsForScan([]string{"192.0.2.12-foo"})
		require.Error(t, err)
		assert.ErrorContains(t, err, "invalid range")
	})

	t.Run("hostname with dash", func(t *testing.T) {
		tgts, err := expandTargetsForScan([]string{"my-host.com"})
		require.NoError(t, err)
		assert.Equal(t, []string{"my-host.com"}, tgts)
	})

	t.Run("ends with dash", func(t *testing.T) {
		tgts, err := expandTargetsForScan([]string{"a-"})
		require.NoError(t, err)
		assert.Equal(t, []string{"a-"}, tgts)
	})

	t.Run("starts with dash", func(t *testing.T) {
		tgts, err := expandTargetsForScan([]string{"-a"})
		require.NoError(t, err)
		assert.Equal(t, []string{"-a"}, tgts)
	})

	t.Run("only a dash", func(t *testing.T) {
		tgts, err := expandTargetsForScan([]string{"-"})
		require.NoError(t, err)
		assert.Equal(t, []string{"-"}, tgts)
	})

	t.Run("nmap format", func(t *testing.T) {
		tgts, err := expandTargetsForScan([]string{"192.168.1.10-255"})
		require.NoError(t, err)
		assert.Equal(t, []string{"192.168.1.10-255"}, tgts)
	})
}
