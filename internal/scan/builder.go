package scan

import (
	"github.com/Ullaakut/cameradar/v6"
	"github.com/Ullaakut/cameradar/v6/internal/scan/nmap"
	"github.com/Ullaakut/cameradar/v6/internal/scan/skip"
)

// Config configures how Cameradar discovers RTSP streams.
type Config struct {
	SkipScan  bool
	Targets   []string
	Ports     []string
	ScanSpeed int16
}

// Reporter reports scan progress and debug information.
type Reporter interface {
	Debug(step cameradar.Step, message string)
	Progress(step cameradar.Step, message string)
}

// New builds a stream scanner based on the provided configuration.
func New(config Config, reporter Reporter) (cameradar.StreamScanner, error) {
	expandedTargets, err := expandTargetsForScan(config.Targets)
	if err != nil {
		return nil, err
	}

	if config.SkipScan {
		return skip.New(expandedTargets, config.Ports), nil
	}

	return nmap.New(config.ScanSpeed, expandedTargets, config.Ports, reporter)
}
