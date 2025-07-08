package cameradar_test

import (
	"context"
	"errors"
	"sync"
	"testing"

	"github.com/Ullaakut/cameradar/v6"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name     string
		scanner  cameradar.StreamScanner
		attacker cameradar.StreamAttacker
		wantErr  require.ErrorAssertionFunc
		wantMsg  string
	}{
		{
			name:     "missing scanner",
			scanner:  nil,
			attacker: &fakeAttacker{},
			wantErr:  require.Error,
			wantMsg:  "stream scanner is required",
		},
		{
			name:     "missing attacker",
			scanner:  &fakeScanner{},
			attacker: nil,
			wantErr:  require.Error,
			wantMsg:  "stream attacker is required",
		},
		{
			name:     "valid",
			scanner:  &fakeScanner{},
			attacker: &fakeAttacker{},
			wantErr:  require.NoError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			app, err := cameradar.New(test.scanner, test.attacker, []string{"target"}, []string{"554"}, &recordingReporter{})
			test.wantErr(t, err)
			if test.wantMsg != "" {
				assert.ErrorContains(t, err, test.wantMsg)
			}
			if err == nil {
				require.NotNil(t, app)
			}
		})
	}
}

func TestApp_Run(t *testing.T) {
	ctx := t.Context()
	streams := []cameradar.Stream{{Port: 554}}
	attacked := []cameradar.Stream{{Port: 8554}}

	tests := []struct {
		name            string
		scanner         *fakeScanner
		attacker        *fakeAttacker
		wantErrContains string
		wantErrorCalls  int
		wantDoneCalls   int
		wantSummaryErr  string
		wantSummary     []cameradar.Stream
	}{
		{
			name: "success",
			scanner: &fakeScanner{
				streams: streams,
			},
			attacker: &fakeAttacker{
				streams: attacked,
			},
			wantDoneCalls:  1,
			wantSummary:    attacked,
			wantSummaryErr: "",
		},
		{
			name: "scan error",
			scanner: &fakeScanner{
				streams: streams,
				err:     errors.New("scan failed"),
			},
			attacker:        &fakeAttacker{},
			wantErrContains: "discovering devices",
			wantErrorCalls:  1,
			wantSummary:     streams,
			wantSummaryErr:  "discovering devices",
		},
		{
			name: "attack error",
			scanner: &fakeScanner{
				streams: streams,
			},
			attacker: &fakeAttacker{
				err: errors.New("attack failed"),
			},
			wantErrContains: "attacking devices",
			wantDoneCalls:   1,
			wantSummary:     streams,
			wantSummaryErr:  "attacking devices",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			reporter := &recordingReporter{}
			scanner := test.scanner
			attacker := test.attacker

			app, err := cameradar.New(scanner, attacker, []string{"target"}, []string{"554"}, reporter)
			require.NoError(t, err)

			err = app.Run(ctx)
			if test.wantErrContains != "" {
				require.Error(t, err)
				assert.ErrorContains(t, err, test.wantErrContains)
			} else {
				require.NoError(t, err)
			}

			assert.Equal(t, 1, scanner.calls)
			assert.Same(t, ctx, scanner.gotCtx)

			if test.wantErrContains == "discovering devices" {
				assert.Equal(t, 0, attacker.calls)
			} else {
				assert.Equal(t, 1, attacker.calls)
				assert.Equal(t, streams, attacker.gotStreams)
			}

			assert.Equal(t, 1, reporter.startCalls)
			assert.Equal(t, test.wantDoneCalls, reporter.doneCalls)
			assert.Equal(t, test.wantErrorCalls, reporter.errorCalls)
			require.Equal(t, 1, reporter.summaryCalls)
			assert.Equal(t, test.wantSummary, reporter.summaryStreams)
			if test.wantSummaryErr == "" {
				assert.NoError(t, reporter.summaryErr)
			} else {
				require.Error(t, reporter.summaryErr)
				assert.ErrorContains(t, reporter.summaryErr, test.wantSummaryErr)
			}
		})
	}
}

type fakeScanner struct {
	streams []cameradar.Stream
	err     error

	calls      int
	gotCtx     context.Context
	gotTargets []string
	gotPorts   []string
}

func (f *fakeScanner) Scan(ctx context.Context) ([]cameradar.Stream, error) {
	f.calls++
	f.gotCtx = ctx
	return f.streams, f.err
}

type fakeAttacker struct {
	streams []cameradar.Stream
	err     error

	calls      int
	gotStreams []cameradar.Stream
}

func (f *fakeAttacker) Attack(_ context.Context, streams []cameradar.Stream) ([]cameradar.Stream, error) {
	f.calls++
	f.gotStreams = append([]cameradar.Stream(nil), streams...)
	if f.err != nil {
		return streams, f.err
	}
	if f.streams != nil {
		return f.streams, nil
	}
	return streams, nil
}

type recordingReporter struct {
	mu             sync.Mutex
	startCalls     int
	doneCalls      int
	errorCalls     int
	summaryCalls   int
	summaryStreams []cameradar.Stream
	summaryErr     error
}

func (r *recordingReporter) Start(cameradar.Step, string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.startCalls++
}

func (r *recordingReporter) Done(cameradar.Step, string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.doneCalls++
}

func (r *recordingReporter) Progress(cameradar.Step, string) {}

func (r *recordingReporter) Debug(cameradar.Step, string) {}

func (r *recordingReporter) Error(cameradar.Step, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.errorCalls++
}

func (r *recordingReporter) Summary(streams []cameradar.Stream, err error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.summaryCalls++
	r.summaryStreams = append([]cameradar.Stream(nil), streams...)
	r.summaryErr = err
}

func (r *recordingReporter) Close() {}
