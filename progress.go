package cameradar

import (
	"strconv"
	"strings"
)

const progressMessagePrefix = "\x00progress:"

// ProgressTotalMessage returns a progress control message that sets the total units for a step.
func ProgressTotalMessage(total int) string {
	return progressMessagePrefix + "total=" + strconv.Itoa(total)
}

// ProgressTickMessage returns a progress control message that increments a step's progress by one unit.
func ProgressTickMessage() string {
	return progressMessagePrefix + "tick"
}

// ParseProgressMessage parses a progress control message.
// It returns a kind of "total" or "tick" and an optional value.
func ParseProgressMessage(message string) (string, int, bool) {
	if !strings.HasPrefix(message, progressMessagePrefix) {
		return "", 0, false
	}

	payload := strings.TrimPrefix(message, progressMessagePrefix)
	if payload == "tick" {
		return "tick", 1, true
	}
	if valuePart, ok := strings.CutPrefix(payload, "total="); ok {
		value, err := strconv.Atoi(valuePart)
		if err != nil {
			return "", 0, false
		}
		return "total", value, true
	}

	return "", 0, false
}
