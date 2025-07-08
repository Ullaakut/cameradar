package nmap

import (
	"context"
	"errors"
	"net/netip"
	"sync"
	"testing"

	"github.com/Ullaakut/cameradar/v6"
	nmaplib "github.com/Ullaakut/nmap/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestScanner_Scan(t *testing.T) {
	ctx := context.WithValue(t.Context(), contextKey("trace"), "scan")

	tests := []struct {
		name            string
		result          *nmaplib.Run
		err             error
		wantStreams     []cameradar.Stream
		wantDebug       []string
		wantProgress    string
		wantErrContains string
	}{
		{
			name: "filters non-rtsp and closed ports",
			result: buildRun(nmaplib.Host{
				Addresses: []nmaplib.Address{
					{Addr: "127.0.0.1"},
					{Addr: "not-an-ip"},
				},
				Ports: []nmaplib.Port{
					openPort(8554, "rtsp", "ACME"),
					closedPort(554, "rtsp", "ACME"),
					openPort(80, "http", "ACME"),
				},
			}),
			wantStreams: []cameradar.Stream{
				{
					Device:  "ACME",
					Address: netip.MustParseAddr("127.0.0.1"),
					Port:    8554,
				},
			},
			wantProgress: "Found 1 RTSP streams",
		},
		{
			name: "collects multiple hosts",
			result: buildRun(
				nmaplib.Host{
					Addresses: []nmaplib.Address{{Addr: "192.0.2.10"}, {Addr: "192.0.2.11"}},
					Ports: []nmaplib.Port{
						openPort(8554, "rtsp-alt", "Model A"),
					},
				},
				nmaplib.Host{
					Addresses: []nmaplib.Address{{Addr: "198.51.100.9"}},
					Ports: []nmaplib.Port{
						openPort(554, "rtsp", "Model B"),
					},
				},
			),
			wantStreams: []cameradar.Stream{
				{
					Device:  "Model A",
					Address: netip.MustParseAddr("192.0.2.10"),
					Port:    8554,
				},
				{
					Device:  "Model A",
					Address: netip.MustParseAddr("192.0.2.11"),
					Port:    8554,
				},
				{
					Device:  "Model B",
					Address: netip.MustParseAddr("198.51.100.9"),
					Port:    554,
				},
			},
			wantProgress: "Found 3 RTSP streams",
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

			scanner, err := New(4, []string{"192.0.2.1"}, []string{"554", "8554"}, reporter)
			require.NoError(t, err)

			scanner.runner = fakeRunner{result: test.result, err: test.err}

			streams, err := scanner.Scan(ctx)

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
			assert.Contains(t, reporter.progress, test.wantProgress)
		})
	}
}

type contextKey string

type fakeRunner struct {
	result *nmaplib.Run
	err    error
}

func (f fakeRunner) Run(context.Context) (*nmaplib.Run, error) {
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

func buildRun(hosts ...nmaplib.Host) *nmaplib.Run {
	return &nmaplib.Run{Hosts: hosts}
}

func openPort(id uint16, serviceName, product string) nmaplib.Port {
	return nmaplib.Port{
		ID: id,
		State: nmaplib.State{
			State: string(nmaplib.Open),
		},
		Service: nmaplib.Service{
			Name:    serviceName,
			Product: product,
		},
	}
}

func closedPort(id uint16, serviceName, product string) nmaplib.Port {
	return nmaplib.Port{
		ID: id,
		State: nmaplib.State{
			State: string(nmaplib.Closed),
		},
		Service: nmaplib.Service{
			Name:    serviceName,
			Product: product,
		},
	}
}
