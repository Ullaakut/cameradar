package main

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gernest/wow"
	"github.com/gernest/wow/spin"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	cmrdr "github.com/ullaakut/cameradar"
	"github.com/ullaakut/disgo/logger"
	log "github.com/ullaakut/disgo/logger"
	"github.com/ullaakut/disgo/symbol"
	curl "github.com/ullaakut/go-curl"
)

type options struct {
	Targets     []string
	Ports       []string
	Routes      string
	Credentials string
	Speed       int
	Timeout     int
	EnableLogs  bool
}

func parseArguments() error {
	viper.SetEnvPrefix("cameradar")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))

	pflag.StringSliceP("targets", "t", []string{}, "The targets on which to scan for open RTSP streams - required (ex: 172.16.100.0/24)")
	pflag.StringSliceP("ports", "p", []string{"554", "5554", "8554"}, "The ports on which to search for RTSP streams")
	pflag.StringP("custom-routes", "r", "<GOPATH>/src/github.com/ullaakut/cameradar/dictionaries/routes", "The path on which to load a custom routes dictionary")
	pflag.StringP("custom-credentials", "c", "<GOPATH>/src/github.com/ullaakut/cameradar/dictionaries/credentials.json", "The path on which to load a custom credentials JSON dictionary")
	pflag.IntP("speed", "s", 4, "The nmap speed preset to use for discovery")
	pflag.IntP("timeout", "T", 2000, "The timeout in miliseconds to use for attack attempts")
	pflag.BoolP("log", "l", false, "Enable the logs for nmap's output to stdout")
	pflag.BoolP("help", "h", false, "displays this help message")

	viper.AutomaticEnv()

	pflag.Parse()

	err := viper.BindPFlags(pflag.CommandLine)
	if err != nil {
		return err
	}

	if viper.GetBool("help") {
		pflag.Usage()
		fmt.Println("\nExamples of usage:")
		fmt.Println("\tScanning your home network for RTSP streams:\tcameradar -t 192.168.0.0/24")
		fmt.Println("\tScanning a remote camera on a specific port:\tcameradar -t 172.178.10.14 -p 18554 -s 2")
		fmt.Println("\tScanning an unstable remote network: \t\tcameradar -t 172.178.10.14/24 -s 1 --timeout 10000 -l")
		os.Exit(0)
	}

	if viper.GetStringSlice("targets") == nil {
		return errors.New("targets (-t, --targets) argument required\n    examples:\n      - 172.16.100.0/24\n      - localhost\n      - 8.8.8.8")
	}

	return nil
}

func main() {
	var options options
	logger, err := log.New(os.Stdout, log.WithErrorOutput(os.Stderr))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create logger: %v", err)
	}

	err = parseArguments()
	if err != nil {
		printErr(logger, err)
	}

	options.Credentials = viper.GetString("custom-credentials")
	options.EnableLogs = viper.GetBool("log") || viper.GetBool("logging")
	options.Ports = viper.GetStringSlice("ports")
	options.Routes = viper.GetString("custom-routes")
	options.Speed = viper.GetInt("speed")
	options.Timeout = viper.GetInt("timeout")
	options.Targets = viper.GetStringSlice("targets")

	w := startSpinner(options.EnableLogs)

	if len(options.Targets) == 1 {
		options.Targets, err = cmrdr.ParseTargetsFile(options.Targets[0])
		if err != nil {
			printErr(logger, err)
		}
	}

	err = curl.GlobalInit(curl.GLOBAL_ALL)
	handle := curl.EasyInit()
	if err != nil || handle == nil {
		printErr(logger, errors.New("libcurl initialization failed"))
	}

	c := &cmrdr.Curl{CURL: handle}
	defer curl.GlobalCleanup()

	updateSpinner(w, "Loading dictionaries...", options.EnableLogs)
	gopath := os.Getenv("GOPATH")
	options.Credentials = strings.Replace(options.Credentials, "<GOPATH>", gopath, 1)
	options.Routes = strings.Replace(options.Routes, "<GOPATH>", gopath, 1)

	credentials, err := cmrdr.LoadCredentials(options.Credentials)
	if err != nil {
		printErr(logger, fmt.Errorf("Invalid credentials dictionary %q: %v", options.Credentials, err))
		return
	}

	routes, err := cmrdr.LoadRoutes(options.Routes)
	if err != nil {
		printErr(logger, fmt.Errorf("Invalid routes dictionary %q: %v", options.Routes, err))
		return
	}

	updateSpinner(w, "Scanning the network...", options.EnableLogs)
	streams, err := cmrdr.Discover(options.Targets, options.Ports, options.Speed)
	if err != nil && len(streams) > 0 {
		printErr(logger, err)
	}

	// Most cameras will be accessed successfully with these two attacks
	updateSpinner(w, "Found "+fmt.Sprint(len(streams))+" streams. Attacking their routes...", options.EnableLogs)
	streams, err = cmrdr.AttackRoute(c, streams, routes, time.Duration(options.Timeout)*time.Millisecond, options.EnableLogs)
	if err != nil && len(streams) > 0 {
		printErr(logger, err)
	}

	updateSpinner(w, "Found "+fmt.Sprint(len(streams))+" streams. Attacking their credentials...", options.EnableLogs)
	streams, err = cmrdr.AttackCredentials(c, streams, credentials, time.Duration(options.Timeout)*time.Millisecond, options.EnableLogs)
	if err != nil && len(streams) > 0 {
		printErr(logger, err)
	}

	// But some cameras run GST RTSP Server which prioritizes 401 over 404 contrary to most cameras.
	// For these cameras, running another route attack will solve the problem.
	for _, stream := range streams {
		if !stream.RouteFound || !stream.CredentialsFound {
			updateSpinner(w, "Found "+fmt.Sprint(len(streams))+" streams. Final attack...", options.EnableLogs)
			streams, err = cmrdr.AttackRoute(c, streams, routes, time.Duration(options.Timeout)*time.Millisecond, options.EnableLogs)
			if err != nil && len(streams) > 0 {
				printErr(logger, err)
			}

			break
		}
	}

	updateSpinner(w, "Found "+fmt.Sprint(len(streams))+" streams. Validating their availability...", options.EnableLogs)
	streams, err = cmrdr.ValidateStreams(c, streams, time.Duration(options.Timeout)*time.Millisecond, options.EnableLogs)
	if err != nil && len(streams) > 0 {
		printErr(logger, err)
	}

	clearOutput(w, options.EnableLogs)

	prettyPrint(logger, streams)
}

func prettyPrint(logger *logger.Logger, streams []cmrdr.Stream) {
	success := 0
	if len(streams) > 0 {
		for _, stream := range streams {
			if stream.CredentialsFound && stream.RouteFound && stream.Available {
				logger.Infof("%s\tDevice RTSP URL:\t%s\n", log.Success(symbol.RightTriangle), log.Link(cmrdr.GetCameraRTSPURL(stream)))
				success++
			} else {
				logger.Infof("%s\tAdmin panel URL:\t%s You can use this URL to try attacking the camera's admin panel instead.\n", log.Failure(symbol.Cross), log.Link(cmrdr.GetCameraAdminPanelURL(stream)))
			}

			logger.Infof("\tDevice model:\t\t%s\n\n", stream.Device)

			if stream.Available {
				logger.Infof("\tAvailable:\t\t%s\n", log.Success(symbol.Check))
			} else {
				logger.Infof("\tAvailable:\t\t%s\n", log.Failure(symbol.Cross))
			}

			logger.Infof("\tIP address:\t\t%s\n", stream.Address)
			logger.Infof("\tRTSP port:\t\t%d\n", stream.Port)
			if stream.CredentialsFound {
				logger.Infof("\tUsername:\t\t%s\n", log.Success(stream.Username))
				logger.Infof("\tPassword:\t\t%s\n", log.Success(stream.Password))
			} else {
				logger.Infof("\tUsername:\t\t%s\n", log.Failure("not found"))
				logger.Infof("\tPassword:\t\t%s\n", log.Failure("not found"))
			}
			if stream.RouteFound {
				logger.Infof("\tRTSP route:\t\t%s\n\n\n", log.Success("/"+stream.Route))
			} else {
				logger.Infof("\tRTSP route:\t\t%s\n\n\n", log.Failure("not found"))
			}
		}
		if success > 1 {
			logger.Infof("%s Successful attack: %s devices were accessed", log.Success(symbol.Check), log.Success(len(streams)))
		} else if success == 1 {
			logger.Infof("%s Successful attack: %s device was accessed", log.Success(symbol.Check), log.Success(len(streams)))
		} else {
			logger.Infof("%s Streams were found but none were accessed. They are most likely configured with secure credentials and routes. You can try adding entries to the dictionary or generating your own in order to attempt a bruteforce attack on the cameras.\n", log.Failure("\xE2\x9C\x96"))
		}
	} else {
		logger.Infof("%s No streams were found. Please make sure that your target is on an accessible network.\n", log.Failure(symbol.Cross))
	}
}

func printErr(logger *logger.Logger, err error) {
	logger.Errorln(log.Failure(symbol.Cross), err)
	os.Exit(1)
}

func updateSpinner(w *wow.Wow, text string, disabled bool) {
	if !disabled {
		w.Text(" " + text)
	}
}

func startSpinner(disabled bool) *wow.Wow {
	if !disabled {
		w := wow.New(os.Stdout, spin.Get(spin.Dots), " Loading dictionaries...")
		w.Start()
		return w
	}
	return nil
}

// HACK: Waiting for a fix to issue
// https://github.com/gernest/wow/issues/5
func clearOutput(w *wow.Wow, disabled bool) {
	if !disabled {
		w.Text("\b")
		time.Sleep(80 * time.Millisecond)
		w.Stop()
	}
}
