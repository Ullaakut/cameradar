package masscan

import (
	"context"
	"errors"
	"net/netip"
	"sync"
	"testing"

	"github.com/Ullaakut/cameradar/v6"
	masscanlib "github.com/Ullaakut/masscan"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRunScan(t *testing.T) {
	tests := []struct {
		name            string
		result          *masscanlib.Run
		err             error
		wantStreams     []cameradar.Stream
		wantDebug       []string
		wantProgress    []string
		wantErrContains string
	}{
		{
			name: "filters invalid addresses, closed and invalid ports",
			result: &masscanlib.Run{
				Hosts: []masscanlib.Host{
					{
						Address: "192.0.2.10",
						Ports: []masscanlib.Port{
							{Number: 554, Status: "open"},
							{Number: 8554, Status: "closed"},
							{Number: 0, Status: "open"},
						},
					},
					{Address: "not-an-ip", Ports: []masscanlib.Port{{Number: 8554, Status: "open"}}},
					{Address: "", Ports: []masscanlib.Port{{Number: 8554, Status: "open"}}},
				},
			},
			wantStreams: []cameradar.Stream{
				{Address: netip.MustParseAddr("192.0.2.10"), Port: 554},
			},
			wantProgress: []string{
				"Skipping invalid port 0 on 192.0.2.10",
				"Skipping invalid address \"not-an-ip\": ParseAddr(\"not-an-ip\"): unable to parse IP",
				"Skipping host with empty address",
				"Found 1 RTSP streams",
			},
		},
		{
			name: "collects streams from multiple hosts",
			result: &masscanlib.Run{
				Hosts: []masscanlib.Host{
					{Address: "192.0.2.10", Ports: []masscanlib.Port{{Number: 8554, Status: "open"}}},
					{Address: "198.51.100.9", Ports: []masscanlib.Port{{Number: 554, Status: "open"}}},
				},
			},
			wantStreams: []cameradar.Stream{
				{Address: netip.MustParseAddr("192.0.2.10"), Port: 8554},
				{Address: netip.MustParseAddr("198.51.100.9"), Port: 554},
			},
			wantProgress: []string{"Found 2 RTSP streams"},
		},
		{
			name:            "returns error when scan fails",
			err:             errors.New("scan failed"),
			wantErrContains: "scanning network",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			reporter := &recordingReporter{}

			streams, err := runScan(t.Context(), fakeRunner{result: test.result, err: test.err}, reporter)

			if test.wantErrContains != "" {
				require.Error(t, err)
				assert.ErrorContains(t, err, test.wantErrContains)
				assert.Empty(t, streams)
				assert.Empty(t, reporter.progress)
				assert.Equal(t, test.wantDebug, reporter.debug)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, test.wantStreams, streams)
			assert.Equal(t, test.wantDebug, reporter.debug)
			for _, progress := range test.wantProgress {
				assert.Contains(t, reporter.progress, progress)
			}
		})
	}
}

type fakeRunner struct {
	result *masscanlib.Run
	err    error
}

func (f fakeRunner) Run(context.Context) (*masscanlib.Run, error) {
	return f.result, f.err
}

type recordingReporter struct {
	mu       sync.Mutex
	debug    []string
	progress []string
}

func (r *recordingReporter) Start(cameradar.Step, string) {}

func (r *recordingReporter) Done(cameradar.Step, string) {}

func (r *recordingReporter) Progress(_ cameradar.Step, message string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.progress = append(r.progress, message)
}

func (r *recordingReporter) Debug(_ cameradar.Step, message string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.debug = append(r.debug, message)
}

func (r *recordingReporter) Error(cameradar.Step, error) {}

func (r *recordingReporter) Summary([]cameradar.Stream, error) {}

func (r *recordingReporter) Close() {}
