package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Ullaakut/cameradar/v6"
	"github.com/Ullaakut/cameradar/v6/internal/attack"
	"github.com/Ullaakut/cameradar/v6/internal/dict"
	"github.com/Ullaakut/cameradar/v6/internal/output"
	"github.com/Ullaakut/cameradar/v6/internal/scan"
	"github.com/Ullaakut/cameradar/v6/internal/ui"
	"github.com/urfave/cli/v3"
	"golang.org/x/term"
)

//nolint:cyclop // Splitting this function does not make it clearer.
func runCameradar(ctx context.Context, cmd *cli.Command) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	targetInputs := cmd.StringSlice(flagTargets)
	if len(targetInputs) == 0 {
		return errors.New("at least one target must be specified")
	}

	targets, err := loadTargets(targetInputs)
	if err != nil {
		return fmt.Errorf("loading targets: %w", err)
	}
	if len(targets) == 0 {
		return errors.New("no valid targets provided")
	}

	ports := cmd.StringSlice(flagPorts)
	if len(ports) == 0 {
		return errors.New("at least one port must be specified")
	}

	var credsPath, routesPath string
	if cmd.IsSet(flagCustomCredentials) {
		credsPath = os.ExpandEnv(cmd.String(flagCustomCredentials))
	}
	if cmd.IsSet(flagCustomRoutes) {
		routesPath = os.ExpandEnv(cmd.String(flagCustomRoutes))
	}

	dictionary, err := dict.New(credsPath, routesPath)
	if err != nil {
		return fmt.Errorf("loading dictionaries: %w", err)
	}

	mode, err := cameradar.ParseMode(cmd.String(flagUI))
	if err != nil {
		return err
	}

	var outputPath string
	if cmd.IsSet(flagOutput) {
		outputPath = os.ExpandEnv(cmd.String(flagOutput))
	}

	interactive := isInteractiveTerminal()
	buildInfo := ui.BuildInfo{Version: version, Commit: commit, Date: date}
	reporter, err := ui.NewReporter(mode, cmd.Bool(flagDebug), os.Stdout, interactive, buildInfo, cancel)
	if err != nil {
		return err
	}
	if plainReporter, ok := reporter.(*ui.PlainReporter); ok {
		resolvedMode := resolveMode(mode, interactive)
		plainReporter.PrintStartup(buildInfo, buildStartupOptions(
			targets,
			ports,
			routesPath,
			credsPath,
			outputPath,
			cmd.Int16(flagScanSpeed),
			cmd.Duration(flagAttackInterval),
			cmd.Duration(flagTimeout),
			cmd.Bool(flagSkipScan),
			cmd.Bool(flagDebug),
			resolvedMode,
		))
	}
	if outputPath != "" {
		reporter = output.NewM3UReporter(reporter, outputPath)
	}
	defer reporter.Close()

	config := scan.Config{
		SkipScan:  cmd.Bool(flagSkipScan),
		Targets:   targets,
		Ports:     ports,
		ScanSpeed: cmd.Int16(flagScanSpeed),
	}
	var scanner cameradar.StreamScanner
	scanner, err = scan.New(config, reporter)
	if err != nil {
		return fmt.Errorf("creating stream scanner: %w", err)
	}

	interval := cmd.Duration(flagAttackInterval)
	timeout := cmd.Duration(flagTimeout)
	attacker, err := attack.New(dictionary, interval, timeout, reporter)
	if err != nil {
		return fmt.Errorf("creating attacker: %w", err)
	}

	c, err := cameradar.New(
		scanner,
		attacker,
		targets,
		ports,
		reporter,
	)
	if err != nil {
		return fmt.Errorf("creating scanner: %w", err)
	}

	return c.Run(ctx)
}

func resolveMode(mode cameradar.Mode, interactive bool) cameradar.Mode {
	if mode != cameradar.ModeAuto {
		return mode
	}
	if interactive {
		return cameradar.ModeTUI
	}
	return cameradar.ModePlain
}

func buildStartupOptions(
	targets []string,
	ports []string,
	routesPath string,
	credsPath string,
	outputPath string,
	scanSpeed int16,
	attackInterval time.Duration,
	timeout time.Duration,
	skipScan bool,
	debug bool,
	mode cameradar.Mode,
) []string {
	options := []string{
		"targets: " + strings.Join(targets, ", "),
		"ports: " + strings.Join(ports, ", "),
		"custom-routes: " + fallbackValue(routesPath, "builtin"),
		"custom-credentials: " + fallbackValue(credsPath, "builtin"),
		"scan-speed: " + strconv.FormatInt(int64(scanSpeed), 10),
		"skip-scan: " + strconv.FormatBool(skipScan),
		"attack-interval: " + attackInterval.String(),
		"timeout: " + timeout.String(),
		"debug: " + strconv.FormatBool(debug),
		"ui: " + string(mode),
		"output: " + fallbackValue(outputPath, "disabled"),
	}
	return options
}

func fallbackValue(value, fallback string) string {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return fallback
	}
	return trimmed
}

func isInteractiveTerminal() bool {
	if !term.IsTerminal(int(os.Stdout.Fd())) {
		return false
	}
	if !term.IsTerminal(int(os.Stdin.Fd())) {
		return false
	}

	termEnv := strings.TrimSpace(os.Getenv("TERM"))
	if termEnv == "" || termEnv == "dumb" {
		return false
	}

	return true
}

// loadTargets merges targets from command line and file paths.
// Valid targets are:
//   - Single IP addresses (e.g., 192.168.1.10)
//   - CIDR notations      (e.g., 192.168.1.0/24)
//   - Hostnames           (e.g., localhost)
//   - IP Ranges           (e.g., 192.168.1.10-20)
func loadTargets(targets []string) ([]string, error) {
	if len(targets) == 0 {
		return nil, nil
	}

	var merged []string
	for _, target := range targets {
		trimmed := strings.TrimSpace(target)
		if trimmed == "" {
			continue
		}

		_, err := os.Stat(trimmed)
		if err != nil {
			merged = append(merged, trimmed)
			continue
		}

		bytes, err := os.ReadFile(trimmed)
		if err != nil {
			return nil, fmt.Errorf("reading targets file %q: %w", trimmed, err)
		}

		for line := range strings.SplitSeq(string(bytes), "\n") {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}
			merged = append(merged, line)
		}
	}

	return merged, nil
}
