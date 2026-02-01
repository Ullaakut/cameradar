package ui

import "strings"

// BuildInfo represents build metadata injected at link time.
type BuildInfo struct {
	Version string
	Commit  string
	Date    string
}

// DisplayVersion returns the version prefixed with "v" when needed.
func (b BuildInfo) DisplayVersion() string {
	version := strings.TrimSpace(b.Version)
	if version == "" {
		version = "dev"
	}
	if strings.HasPrefix(version, "v") {
		return version
	}
	return "v" + version
}

// LogVersion returns the version without a leading "v".
func (b BuildInfo) LogVersion() string {
	version := strings.TrimSpace(b.Version)
	if version == "" {
		return "dev"
	}
	return strings.TrimPrefix(version, "v")
}

// ShortCommit returns a shortened commit hash suitable for display.
func (b BuildInfo) ShortCommit() string {
	commit := strings.TrimSpace(b.Commit)
	if commit == "" || commit == "none" || commit == "unknown" {
		return "unknown"
	}
	if len(commit) > 7 {
		return commit[:7]
	}
	return commit
}

// TUIHeader returns the header used by the TUI.
func (b BuildInfo) TUIHeader() string {
	return "Cameradar — " + b.DisplayVersion() + " (" + b.ShortCommit() + ")"
}
