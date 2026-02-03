package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/urfave/cli/v3"
)

func printVersion(ctx context.Context, _ *cli.Command) error {
	nmapVersion := getNmapVersion(ctx)
	_, err := fmt.Fprintf(os.Stdout, "Version:\tv%s\nCommit:\t\t%s\nBuild date:\t%s\nNmap:\t\t%s\n", version, commit, date, nmapVersion)
	return err
}

const unknownVersion = "unknown"

func getNmapVersion(ctx context.Context) string {
	output, err := exec.CommandContext(ctx, "nmap", "--version").Output()
	if err != nil {
		return unknownVersion
	}

	lines := strings.SplitN(string(output), "\n", 2)
	if len(lines) == 0 {
		return unknownVersion
	}

	firstLine := strings.TrimSpace(lines[0])
	const prefix = "Nmap version "
	if !strings.HasPrefix(firstLine, prefix) {
		return unknownVersion
	}

	versionPart := strings.TrimSpace(strings.TrimPrefix(firstLine, prefix))
	fields := strings.Fields(versionPart)
	if len(fields) == 0 {
		return unknownVersion
	}
	return fields[0]
}
