package attack

import (
	"fmt"
	"strconv"
	"strings"
)

type incrementalRoute struct {
	prefix    string
	suffix    string
	number    int
	width     int
	isChannel bool
}

// detectIncrementalRoute identifies routes that can be incremented.
// It prioritizes channel-like patterns to enable sequential scanning when possible.
//
// Examples of supported patterns:
// - /StreamingSetting?ChannelID=01&other=params -> /StreamingSetting?ChannelID=02&other=params
// - /path/to/channel2/stream -> /path/to/channel3/stream
// - /foo/bar12/baz -> /foo/bar13/baz
//
// It returns false if no incrementable pattern is found.
func detectIncrementalRoute(route string) (incrementalRoute, bool) {
	if strings.TrimSpace(route) == "" {
		return incrementalRoute{}, false
	}

	if match, ok := findChannelIncrement(route); ok {
		match.isChannel = true
		return match, true
	}

	match, ok := findLastNumber(route)
	if !ok {
		return incrementalRoute{}, false
	}
	return match, true
}

// findChannelIncrement locates a numeric segment tied to channel-like keywords.
// It returns the last matching segment so we increment the most specific part.
//
// Supported keywords include: channel_id, channelid, channelno, channel, channelname.
func findChannelIncrement(route string) (incrementalRoute, bool) {
	patterns := []string{"channel_id", "channelid", "channelno", "channel", "channelname"}
	lower := strings.ToLower(route)

	for _, pattern := range patterns {
		var lastMatch incrementalRoute
		found := false
		index := 0

		for {
			pos := strings.Index(lower[index:], pattern)
			if pos == -1 {
				break
			}
			pos += index

			start, end, ok := firstNumberAfter(route, pos+len(pattern))
			if ok {
				num, width, parseOK := parseNumber(route, start, end)
				if parseOK {
					lastMatch = incrementalRoute{
						prefix: route[:start],
						suffix: route[end:],
						number: num,
						width:  width,
					}
					found = true
				}
			}
			index = pos + len(pattern)
		}
		if found {
			return lastMatch, true
		}
	}

	return incrementalRoute{}, false
}

// findLastNumber finds the last numeric token in the route so it can be incremented.
// This supports routes where the channel number is not the final component.
func findLastNumber(route string) (incrementalRoute, bool) {
	for i := len(route) - 1; i >= 0; {
		if !isDigit(route[i]) {
			i--
			continue
		}

		end := i + 1
		start := i
		for start >= 0 && isDigit(route[start]) {
			start--
		}
		start++

		num, width, ok := parseNumber(route, start, end)
		if !ok {
			i = start - 1
			continue
		}

		return incrementalRoute{
			prefix: route[:start],
			suffix: route[end:],
			number: num,
			width:  width,
		}, true
	}

	return incrementalRoute{}, false
}

// parseNumber reads the numeric token and returns its integer value and width.
func parseNumber(route string, start, end int) (int, int, bool) {
	if start < 0 || end > len(route) || start >= end {
		return 0, 0, false
	}

	value := route[start:end]
	num, err := strconv.Atoi(value)
	if err != nil {
		return 0, 0, false
	}

	return num, len(value), true
}

// firstNumberAfter returns the first numeric token after a given index.
func firstNumberAfter(route string, after int) (start, end int, ok bool) {
	if after < 0 {
		after = 0
	}

	for i := after; i < len(route); i++ {
		if !isDigit(route[i]) {
			continue
		}

		end := i + 1
		for end < len(route) && isDigit(route[end]) {
			end++
		}
		return i, end, true
	}
	return 0, 0, false
}

// buildIncrementedRoute formats the route with the new numeric value.
// It preserves zero padding when the original token had a fixed width.
func buildIncrementedRoute(match incrementalRoute, number int) string {
	if match.width <= 0 {
		return match.prefix + strconv.Itoa(number) + match.suffix
	}
	return match.prefix + fmt.Sprintf("%0*d", match.width, number) + match.suffix
}

func isDigit(b byte) bool {
	return b >= '0' && b <= '9'
}
