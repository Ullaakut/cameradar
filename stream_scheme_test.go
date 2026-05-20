package cameradar_test

import (
	"testing"

	"github.com/Ullaakut/cameradar/v6"
	"github.com/stretchr/testify/require"
)

func TestStreamRTSPScheme(t *testing.T) {
	tests := []struct {
		name   string
		scheme string
		want   string
	}{
		{name: "empty defaults to rtsp", scheme: "", want: "rtsp"},
		{name: "rtsp stays rtsp", scheme: "rtsp", want: "rtsp"},
		{name: "http maps to rtsp", scheme: "http", want: "rtsp"},
		{name: "https maps to rtsps", scheme: "https", want: "rtsps"},
		{name: "rtsps stays rtsps", scheme: "rtsps", want: "rtsps"},
		{name: "unknown falls back to rtsp", scheme: "custom", want: "rtsp"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			stream := cameradar.Stream{Scheme: test.scheme}
			got := stream.RTSPScheme()
			require.Equal(t, test.want, got)
		})
	}
}
