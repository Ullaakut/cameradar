package cameradar_test

import (
	"testing"

	"github.com/Ullaakut/cameradar/v6"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseMode(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		want           cameradar.Mode
		wantErr        require.ErrorAssertionFunc
		wantErrMessage string
	}{
		{
			name:    "auto",
			input:   "auto",
			want:    cameradar.ModeAuto,
			wantErr: require.NoError,
		},
		{
			name:    "tui",
			input:   "TUI",
			want:    cameradar.ModeTUI,
			wantErr: require.NoError,
		},
		{
			name:    "plain",
			input:   "plain",
			want:    cameradar.ModePlain,
			wantErr: require.NoError,
		},
		{
			name:    "empty",
			input:   "  ",
			want:    cameradar.ModeAuto,
			wantErr: require.NoError,
		},
		{
			name:           "invalid",
			input:          "nope",
			want:           cameradar.ModeAuto,
			wantErr:        require.Error,
			wantErrMessage: "invalid ui mode",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := cameradar.ParseMode(test.input)
			test.wantErr(t, err)
			if test.wantErrMessage != "" {
				assert.ErrorContains(t, err, test.wantErrMessage)
			}
			assert.Equal(t, test.want, got)
		})
	}
}

func TestStepLabel(t *testing.T) {
	tests := []struct {
		step cameradar.Step
		want string
	}{
		{step: cameradar.StepScan, want: "Scan targets"},
		{step: cameradar.StepAttackRoutes, want: "Attack routes"},
		{step: cameradar.StepDetectAuth, want: "Detect authentication"},
		{step: cameradar.StepAttackCredentials, want: "Attack credentials"},
		{step: cameradar.StepValidateStreams, want: "Validate streams"},
		{step: cameradar.StepSummary, want: "Summary"},
		{step: cameradar.Step("custom"), want: "custom"},
	}

	for _, test := range tests {
		t.Run(test.want, func(t *testing.T) {
			assert.Equal(t, test.want, cameradar.StepLabel(test.step))
		})
	}
}

func TestSteps(t *testing.T) {
	assert.Equal(t, []cameradar.Step{
		cameradar.StepScan,
		cameradar.StepAttackRoutes,
		cameradar.StepDetectAuth,
		cameradar.StepAttackCredentials,
		cameradar.StepValidateStreams,
		cameradar.StepSummary,
	}, cameradar.Steps())
}
