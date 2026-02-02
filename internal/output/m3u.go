package output

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/Ullaakut/cameradar/v6"
	"github.com/Ullaakut/cameradar/v6/internal/ui"
)

type m3uReporter struct {
	delegate   ui.Reporter
	outputPath string
}

// NewM3UReporter wraps the provided reporter and writes an M3U playlist on summary.
func NewM3UReporter(delegate ui.Reporter, outputPath string) ui.Reporter {
	return &m3uReporter{
		delegate:   delegate,
		outputPath: strings.TrimSpace(outputPath),
	}
}

func (r *m3uReporter) Start(step cameradar.Step, message string) {
	r.delegate.Start(step, message)
}

func (r *m3uReporter) Done(step cameradar.Step, message string) {
	r.delegate.Done(step, message)
}

func (r *m3uReporter) Progress(step cameradar.Step, message string) {
	r.delegate.Progress(step, message)
}

func (r *m3uReporter) Debug(step cameradar.Step, message string) {
	r.delegate.Debug(step, message)
}

func (r *m3uReporter) Error(step cameradar.Step, err error) {
	r.delegate.Error(step, err)
}

func (r *m3uReporter) Summary(streams []cameradar.Stream, err error) {
	r.delegate.Summary(streams, err)
	if r.outputPath == "" {
		return
	}

	writeErr := writeM3UFile(r.outputPath, streams)
	if writeErr != nil {
		r.delegate.Error(cameradar.StepSummary, writeErr)
	}
}

func (r *m3uReporter) UpdateSummary(streams []cameradar.Stream) {
	updater, ok := r.delegate.(interface{ UpdateSummary([]cameradar.Stream) })
	if !ok {
		return
	}
	updater.UpdateSummary(streams)
}

func (r *m3uReporter) Close() {
	r.delegate.Close()
}

func writeM3UFile(path string, streams []cameradar.Stream) error {
	content := BuildM3U(streams)
	dir := filepath.Dir(path)
	if dir != "." {
		err := os.MkdirAll(dir, 0o750)
		if err != nil {
			return fmt.Errorf("creating output directory %q: %w", dir, err)
		}
	}

	err := os.WriteFile(path, []byte(content), 0o600)
	if err != nil {
		return fmt.Errorf("writing m3u output: %w", err)
	}
	return nil
}

// BuildM3U creates an M3U playlist with discovered streams.
func BuildM3U(streams []cameradar.Stream) string {
	var builder strings.Builder
	builder.WriteString("#EXTM3U\n")
	for _, stream := range streams {
		url := formatRTSPURL(stream)
		if url == "" {
			continue
		}
		builder.WriteString("#EXTINF:-1,")
		builder.WriteString(formatStreamLabel(stream))
		builder.WriteString("\n")
		builder.WriteString(url)
		builder.WriteString("\n")
	}
	return builder.String()
}

func formatStreamLabel(stream cameradar.Stream) string {
	label := stream.Address.String() + ":" + strconv.FormatUint(uint64(stream.Port), 10)
	if stream.Device == "" {
		return label
	}
	return label + " (" + stream.Device + ")"
}

func formatRTSPURL(stream cameradar.Stream) string {
	path := "/" + strings.TrimLeft(strings.TrimSpace(stream.Route()), "/")

	credentials := ""
	if stream.CredentialsFound && (stream.Username != "" || stream.Password != "") {
		credentials = stream.Username + ":" + stream.Password + "@"
	}

	return "rtsp://" + credentials + stream.Address.String() + ":" + strconv.FormatUint(uint64(stream.Port), 10) + path
}
