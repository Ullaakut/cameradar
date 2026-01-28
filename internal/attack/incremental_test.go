package attack

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDetectIncrementalRoute_ChannelID(t *testing.T) {
	route := "/StreamingSetting?version=1.0&action=getRTSPStream&ChannelID=01&ChannelName=Channel1"

	match, ok := detectIncrementalRoute(route)
	require.True(t, ok)
	assert.True(t, match.isChannel)
	assert.Equal(t, 1, match.number)
	assert.Equal(t, 2, match.width)
	assert.Equal(t, "/StreamingSetting?version=1.0&action=getRTSPStream&ChannelID=", match.prefix)
	assert.Equal(t, "&ChannelName=Channel1", match.suffix)

	next := buildIncrementedRoute(match, match.number+1)
	assert.Equal(t, "/StreamingSetting?version=1.0&action=getRTSPStream&ChannelID=02&ChannelName=Channel1", next)
}

func TestDetectIncrementalRoute_ChannelSuffix(t *testing.T) {
	route := "/path/to/channel2/stream"

	match, ok := detectIncrementalRoute(route)
	require.True(t, ok)
	assert.True(t, match.isChannel)
	assert.Equal(t, 2, match.number)
	assert.Equal(t, "/path/to/channel", match.prefix)
	assert.Equal(t, "/stream", match.suffix)
}

func TestDetectIncrementalRoute_LastNumber(t *testing.T) {
	route := "/foo/bar12/baz"

	match, ok := detectIncrementalRoute(route)
	require.True(t, ok)
	assert.False(t, match.isChannel)
	assert.Equal(t, 12, match.number)
	assert.Equal(t, 2, match.width)
	assert.Equal(t, "/foo/bar", match.prefix)
	assert.Equal(t, "/baz", match.suffix)

	next := buildIncrementedRoute(match, 13)
	assert.Equal(t, "/foo/bar13/baz", next)
}

func TestDetectIncrementalRoute_NoNumber(t *testing.T) {
	match, ok := detectIncrementalRoute("/no/number/here")
	assert.False(t, ok)
	assert.Equal(t, incrementalRoute{}, match)
}

func TestDetectIncrementalRoute_OverflowAtEndFallsBack(t *testing.T) {
	// The trailing token overflows strconv.Atoi, so we fall back to earlier numbers.
	route := "/foo1/bar999999999999999999999999999999"

	match, ok := detectIncrementalRoute(route)
	require.True(t, ok)
	assert.False(t, match.isChannel)
	assert.Equal(t, 1, match.number)
	assert.Equal(t, "/foo", match.prefix)
	assert.Equal(t, "/bar999999999999999999999999999999", match.suffix)
}

func TestDetectIncrementalRoute_ChannelKeywordShouldNotBindAcrossParams(t *testing.T) {
	// The channel keyword should not bind to digits in other query parameters.
	route := "/path?channelname=foo&version=12"

	match, ok := detectIncrementalRoute(route)
	require.True(t, ok)
	assert.False(t, match.isChannel)
	assert.Equal(t, 12, match.number)
	assert.Equal(t, "/path?channelname=foo&version=", match.prefix)
	assert.Equal(t, "", match.suffix)
}

func TestDetectIncrementalRoute_ChannelKeywordStopsAtDelimiter(t *testing.T) {
	// Digits after a delimiter should not be associated with a channel keyword.
	route := "/path/channel?channel=foo/7"

	match, ok := detectIncrementalRoute(route)
	require.True(t, ok)
	assert.False(t, match.isChannel)
	assert.Equal(t, 7, match.number)
	assert.Equal(t, "/path/channel?channel=foo/", match.prefix)
	assert.Equal(t, "", match.suffix)
}

func TestDetectIncrementalRoute_ChannelKeywordWithoutDigitsFallsBack(t *testing.T) {
	// channel keyword without digits should fall back to last numeric token.
	route := "/path/channel?channel=foo&stream=9"

	match, ok := detectIncrementalRoute(route)
	require.True(t, ok)
	assert.False(t, match.isChannel)
	assert.Equal(t, 9, match.number)
	assert.Equal(t, "/path/channel?channel=foo&stream=", match.prefix)
	assert.Equal(t, "", match.suffix)
}

func TestDetectIncrementalRoute_ChannelKeywordKeepsQueryDigits(t *testing.T) {
	// channel keyword with query param digits should be detected as channel.
	route := "/path?channel=03&other=1"

	match, ok := detectIncrementalRoute(route)
	require.True(t, ok)
	assert.True(t, match.isChannel)
	assert.Equal(t, 3, match.number)
	assert.Equal(t, 2, match.width)
	assert.Equal(t, "/path?channel=", match.prefix)
	assert.Equal(t, "&other=1", match.suffix)
}

func TestDetectIncrementalRoute_ChannelKeywordMultipleMatchesUsesKeywordPriority(t *testing.T) {
	// Keyword priority should win even if another keyword appears earlier in the route.
	route := "/path?channel=1&channelid=9"

	match, ok := detectIncrementalRoute(route)
	require.True(t, ok)
	assert.True(t, match.isChannel)
	assert.Equal(t, 9, match.number)
	assert.Equal(t, "/path?channel=1&channelid=", match.prefix)
	assert.Equal(t, "", match.suffix)
}

func TestDetectIncrementalRoute_ChannelKeywordSelectsLastMatchWithinKeyword(t *testing.T) {
	// The last match for a given keyword should be selected.
	route := "/path?channel=1&foo=bar&channel=4"

	match, ok := detectIncrementalRoute(route)
	require.True(t, ok)
	assert.True(t, match.isChannel)
	assert.Equal(t, 4, match.number)
	assert.Equal(t, "/path?channel=1&foo=bar&channel=", match.prefix)
	assert.Equal(t, "", match.suffix)
}

func TestBuildIncrementedRoute_ZeroPadding(t *testing.T) {
	match := incrementalRoute{
		prefix: "/channel",
		suffix: "/stream",
		number: 1,
		width:  3,
	}

	assert.Equal(t, "/channel002/stream", buildIncrementedRoute(match, 2))
}
