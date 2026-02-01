package ui_test

import (
	"testing"

	"github.com/Ullaakut/cameradar/v6/internal/ui"
	"github.com/stretchr/testify/assert"
)

func TestBuildInfo_DisplayVersion(t *testing.T) {
	tests := []struct {
		name    string
		version string
		want    string
	}{
		{
			name:    "empty defaults to dev with prefix",
			version: "",
			want:    "vdev",
		},
		{
			name:    "dev without prefix",
			version: "dev",
			want:    "vdev",
		},
		{
			name:    "already prefixed",
			version: "v1.2.3",
			want:    "v1.2.3",
		},
		{
			name:    "adds prefix",
			version: "1.2.3",
			want:    "v1.2.3",
		},
		{
			name:    "trims spaces with prefix",
			version: " v2.0 ",
			want:    "v2.0",
		},
		{
			name:    "trims spaces without prefix",
			version: " 2.0 ",
			want:    "v2.0",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			info := ui.BuildInfo{Version: test.version}
			assert.Equal(t, test.want, info.DisplayVersion())
		})
	}
}

func TestBuildInfo_LogVersion(t *testing.T) {
	tests := []struct {
		name    string
		version string
		want    string
	}{
		{
			name:    "empty defaults to dev",
			version: "",
			want:    "dev",
		},
		{
			name:    "removes leading v",
			version: "v1.2.3",
			want:    "1.2.3",
		},
		{
			name:    "keeps version without prefix",
			version: "1.2.3",
			want:    "1.2.3",
		},
		{
			name:    "trims spaces and removes prefix",
			version: " v2.0 ",
			want:    "2.0",
		},
		{
			name:    "removes only first prefix",
			version: "vv1",
			want:    "v1",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			info := ui.BuildInfo{Version: test.version}
			assert.Equal(t, test.want, info.LogVersion())
		})
	}
}

func TestBuildInfo_ShortCommit(t *testing.T) {
	tests := []struct {
		name   string
		commit string
		want   string
	}{
		{
			name:   "empty defaults to unknown",
			commit: "",
			want:   "unknown",
		},
		{
			name:   "none defaults to unknown",
			commit: "none",
			want:   "unknown",
		},
		{
			name:   "unknown defaults to unknown",
			commit: "unknown",
			want:   "unknown",
		},
		{
			name:   "short commit preserved",
			commit: "abcdef",
			want:   "abcdef",
		},
		{
			name:   "seven chars preserved",
			commit: "abcdefg",
			want:   "abcdefg",
		},
		{
			name:   "long commit shortened",
			commit: "abcdefghi",
			want:   "abcdefg",
		},
		{
			name:   "trims spaces before shortening",
			commit: " 1234567890 ",
			want:   "1234567",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			info := ui.BuildInfo{Commit: test.commit}
			assert.Equal(t, test.want, info.ShortCommit())
		})
	}
}

func TestBuildInfo_TUIHeader(t *testing.T) {
	tests := []struct {
		name    string
		version string
		commit  string
		want    string
	}{
		{
			name:    "uses display version and short commit",
			version: "1.2.3",
			commit:  "abcdefghi",
			want:    "Cameradar — v1.2.3 (abcdefg)",
		},
		{
			name:    "uses defaults for empty values",
			version: "",
			commit:  "",
			want:    "Cameradar — vdev (unknown)",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			info := ui.BuildInfo{Version: test.version, Commit: test.commit}
			assert.Equal(t, test.want, info.TUIHeader())
		})
	}
}
