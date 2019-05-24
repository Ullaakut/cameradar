package cameradar

import (
	"time"

	"github.com/ullaakut/disgo"
)

// Scanner represents a cameradar scanner. It scans a network and
// attacks all streams found to get their RTSP credentials.
type Scanner struct {
	curl    Curler
	term    *disgo.Terminal
	debug   bool
	timeout time.Duration
}
