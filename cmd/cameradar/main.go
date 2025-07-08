package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"
	"time"

	"github.com/Ullaakut/disgo"
	"github.com/Ullaakut/disgo/style"
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
	flagVerbose           = "verbose"
	flagDebug             = "debug"
)

var version = "dev"

var flags = cmd.Flags{
	&cli.StringSliceFlag{
		Name:     flagTargets,
		Usage:    "The targets on which to scan for open RTSP streams in a network range format",
		Aliases:  []string{"t"},
		Sources:  cli.EnvVars(strcase.ToSNAKE(flagTargets)),
		Required: true,
	},
	&cli.UintSliceFlag{
		Name:    flagPorts,
		Usage:   "The ports on which to search for RTSP streams",
		Aliases: []string{"p"},
		Sources: cli.EnvVars(strcase.ToSNAKE(flagPorts)),
		Value:   []uint64{554, 5554, 8554},
	},
	&cli.StringFlag{
		Name:    flagCustomRoutes,
		Usage:   "The path on which to load a custom routes dictionary",
		Aliases: []string{"r"},
		Sources: cli.EnvVars(strcase.ToSNAKE(flagCustomRoutes)),
		Value:   "${GOPATH}/src/github.com/Ullaakut/cameradar/dictionaries/routes",
	},
	&cli.StringFlag{
		Name:    flagCustomCredentials,
		Usage:   "The path on which to load a custom credentials JSON dictionary",
		Aliases: []string{"c"},
		Sources: cli.EnvVars(strcase.ToSNAKE(flagCustomCredentials)),
		Value:   "${GOPATH}/src/github.com/Ullaakut/cameradar/dictionaries/credentials.json",
	},
	&cli.IntFlag{
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
		Name:    flagVerbose,
		Usage:   "Enable verbose logs",
		Aliases: []string{"v"},
		Sources: cli.EnvVars(strcase.ToSNAKE(flagVerbose)),
		Value:   false,
	},
	&cli.BoolFlag{
		Name:    flagDebug,
		Usage:   "Enable debug logs",
		Aliases: []string{"d"},
		Sources: cli.EnvVars(strcase.ToSNAKE(flagDebug)),
		Value:   false,
	},
}

func main() {
	os.Exit(realMain(os.Args))
}

func realMain(args []string) (code int) {
	defer func() {
		if v := recover(); v != nil {
			printErr(fmt.Errorf("Panic: %v\n%s\n", v, debug.Stack()))
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

	if err := app.Run(ctx, os.Args); err != nil {
		printErr(err)
		return 1
	}
	return 0
}

func printErr(err error) {
	disgo.Errorln(style.Failure(style.SymbolCross), err)
}
