package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"
	"time"

	"github.com/ettle/strcase"
	"github.com/hamba/cmd/v3"
	"github.com/urfave/cli/v3"
)

const (
	flagTargets           = "targets"
	flagPorts             = "ports"
	flagCustomRoutes      = "custom-routes"
	flagCustomCredentials = "custom-credentials"
	flagScanSpeed         = "scan-speed"
	flagAttackInterval    = "attack-interval"
	flagTimeout           = "timeout"
	flagSkipScan          = "skip-scan"
	flagDebug             = "debug"
	flagUI                = "ui"
	flagOutput            = "output"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

var flags = cmd.Flags{
	&cli.StringSliceFlag{
		Name:     flagTargets,
		Usage:    "The targets on which to scan for open RTSP streams in a network range format",
		Aliases:  []string{"t"},
		Sources:  cli.EnvVars(strcase.ToSNAKE(flagTargets)),
		Required: true,
	},
	&cli.StringSliceFlag{
		Name:    flagPorts,
		Usage:   "The ports on which to search for RTSP streams",
		Aliases: []string{"p"},
		Sources: cli.EnvVars(strcase.ToSNAKE(flagPorts)),
		Value:   []string{"554", "5554", "8554", "http"},
	},
	&cli.StringFlag{
		Name:    flagCustomRoutes,
		Usage:   "The path on which to load a custom routes dictionary",
		Aliases: []string{"r"},
		Sources: cli.EnvVars(strcase.ToSNAKE(flagCustomRoutes)),
	},
	&cli.StringFlag{
		Name:    flagCustomCredentials,
		Usage:   "The path on which to load a custom credentials JSON dictionary",
		Aliases: []string{"c"},
		Sources: cli.EnvVars(strcase.ToSNAKE(flagCustomCredentials)),
	},
	&cli.Int16Flag{
		Name:    flagScanSpeed,
		Usage:   "The nmap speed preset to use for scanning (lower is stealthier)",
		Aliases: []string{"s"},
		Sources: cli.EnvVars(strcase.ToSNAKE(flagScanSpeed)),
		Value:   4,
	},
	&cli.DurationFlag{
		Name:    flagAttackInterval,
		Usage:   "The interval between each attack (i.e: 2000ms, higher is stealthier)",
		Aliases: []string{"I"},
		Sources: cli.EnvVars(strcase.ToSNAKE(flagAttackInterval)),
		Value:   0,
	},
	&cli.DurationFlag{
		Name:    flagTimeout,
		Usage:   "The timeout to use for attack attempts (i.e: 2000ms)",
		Aliases: []string{"T"},
		Sources: cli.EnvVars(strcase.ToSNAKE(flagTimeout)),
		Value:   2000 * time.Millisecond,
	},
	&cli.BoolFlag{
		Name:    flagSkipScan,
		Usage:   "Skip discovery and treat every target and port as an RTSP stream",
		Sources: cli.EnvVars(strcase.ToSNAKE(flagSkipScan)),
		Value:   false,
	},
	&cli.BoolFlag{
		Name:    flagDebug,
		Usage:   "Enable debug logs",
		Aliases: []string{"d"},
		Sources: cli.EnvVars(strcase.ToSNAKE(flagDebug)),
		Value:   false,
	},
	&cli.StringFlag{
		Name:    flagUI,
		Usage:   "UI mode: auto, tui, or plain",
		Sources: cli.EnvVars(strcase.ToSNAKE(flagUI)),
		Value:   "auto",
	},
	&cli.StringFlag{
		Name:    flagOutput,
		Usage:   "Write discovered streams to an M3U file at the given path",
		Sources: cli.EnvVars(strcase.ToSNAKE(flagOutput)),
	},
}

func main() {
	os.Exit(realMain())
}

func realMain() (code int) {
	defer func() {
		if v := recover(); v != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Panic: %v\n%s\n", v, debug.Stack())
			code = 1
		}
	}()

	app := &cli.Command{
		Name:    "Cameradar",
		Version: version,
		Flags:   flags,
		Action:  runCameradar,
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	err := app.Run(ctx, os.Args)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
		return 1
	}
	return 0
}
