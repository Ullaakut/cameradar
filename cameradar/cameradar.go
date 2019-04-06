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
	log "github.com/ullaakut/disgo"
	"github.com/ullaakut/disgo/style"
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
	term := log.NewTerminal()

	err := parseArguments()
	if err != nil {
		printErr(term, err)
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
			printErr(term, err)
		}
	}

	err = curl.GlobalInit(curl.GLOBAL_ALL)
	handle := curl.EasyInit()
	if err != nil || handle == nil {
		printErr(term, errors.New("libcurl initialization failed"))
	}

	c := &cmrdr.Curl{CURL: handle}
	defer curl.GlobalCleanup()

	updateSpinner(w, "Loading dictionaries...", options.EnableLogs)
	gopath := os.Getenv("GOPATH")
	options.Credentials = strings.Replace(options.Credentials, "<GOPATH>", gopath, 1)
	options.Routes = strings.Replace(options.Routes, "<GOPATH>", gopath, 1)

	credentials, err := cmrdr.LoadCredentials(options.Credentials)
	if err != nil {
		printErr(term, fmt.Errorf("Invalid credentials dictionary %q: %v", options.Credentials, err))
		return
	}

	routes, err := cmrdr.LoadRoutes(options.Routes)
	if err != nil {
		printErr(term, fmt.Errorf("Invalid routes dictionary %q: %v", options.Routes, err))
		return
	}

	updateSpinner(w, "Scanning the network...", options.EnableLogs)
	streams, err := cmrdr.Discover(options.Targets, options.Ports, options.Speed)
	if err != nil && len(streams) > 0 {
		printErr(term, err)
	}

	// Most cameras will be accessed successfully with these two attacks
	updateSpinner(w, "Found "+fmt.Sprint(len(streams))+" streams. Attacking their routes...", options.EnableLogs)
	streams, err = cmrdr.AttackRoute(c, streams, routes, time.Duration(options.Timeout)*time.Millisecond, options.EnableLogs)
	if err != nil && len(streams) > 0 {
		printErr(term, err)
	}

	updateSpinner(w, "Found "+fmt.Sprint(len(streams))+" streams. Attacking their credentials...", options.EnableLogs)
	streams, err = cmrdr.AttackCredentials(c, streams, credentials, time.Duration(options.Timeout)*time.Millisecond, options.EnableLogs)
	if err != nil && len(streams) > 0 {
		printErr(term, err)
	}

	// But some cameras run GST RTSP Server which prioritizes 401 over 404 contrary to most cameras.
	// For these cameras, running another route attack will solve the problem.
	for _, stream := range streams {
		if !stream.RouteFound || !stream.CredentialsFound {
			updateSpinner(w, "Found "+fmt.Sprint(len(streams))+" streams. Final attack...", options.EnableLogs)
			streams, err = cmrdr.AttackRoute(c, streams, routes, time.Duration(options.Timeout)*time.Millisecond, options.EnableLogs)
			if err != nil && len(streams) > 0 {
				printErr(term, err)
			}

			break
		}
	}

	updateSpinner(w, "Found "+fmt.Sprint(len(streams))+" streams. Validating their availability...", options.EnableLogs)
	streams, err = cmrdr.ValidateStreams(c, streams, time.Duration(options.Timeout)*time.Millisecond, options.EnableLogs)
	if err != nil && len(streams) > 0 {
		printErr(term, err)
	}

	clearOutput(w, options.EnableLogs)

	prettyPrint(term, streams)
}

func prettyPrint(term *log.Terminal, streams []cmrdr.Stream) {
	success := 0
	if len(streams) > 0 {
		for _, stream := range streams {
			if stream.CredentialsFound && stream.RouteFound && stream.Available {
				term.Infof("%s\tDevice RTSP URL:\t%s\n", style.Success(style.SymbolRightTriangle), style.Link(cmrdr.GetCameraRTSPURL(stream)))
				success++
			} else {
				term.Infof("%s\tAdmin panel URL:\t%s You can use this URL to try attacking the camera's admin panel instead.\n", style.Failure(style.SymbolCross), style.Link(cmrdr.GetCameraAdminPanelURL(stream)))
			}

			term.Infof("\tDevice model:\t\t%s\n\n", stream.Device)

			if stream.Available {
				term.Infof("\tAvailable:\t\t%s\n", style.Success(style.SymbolCheck))
			} else {
				term.Infof("\tAvailable:\t\t%s\n", style.Failure(style.SymbolCross))
			}

			term.Infof("\tIP address:\t\t%s\n", stream.Address)
			term.Infof("\tRTSP port:\t\t%d\n", stream.Port)
			if stream.CredentialsFound {
				term.Infof("\tUsername:\t\t%s\n", style.Success(stream.Username))
				term.Infof("\tPassword:\t\t%s\n", style.Success(stream.Password))
			} else {
				term.Infof("\tUsername:\t\t%s\n", style.Failure("not found"))
				term.Infof("\tPassword:\t\t%s\n", style.Failure("not found"))
			}
			if stream.RouteFound {
				term.Infof("\tRTSP route:\t\t%s\n\n\n", style.Success("/"+stream.Route))
			} else {
				term.Infof("\tRTSP route:\t\t%s\n\n\n", style.Failure("not found"))
			}
		}
		if success > 1 {
			term.Infof("%s Successful attack: %s devices were accessed", style.Success(style.SymbolCheck), style.Success(len(streams)))
		} else if success == 1 {
			term.Infof("%s Successful attack: %s device was accessed", style.Success(style.SymbolCheck), style.Success(len(streams)))
		} else {
			term.Infof("%s Streams were found but none were accessed. They are most likely configured with secure credentials and routes. You can try adding entries to the dictionary or generating your own in order to attempt a bruteforce attack on the cameras.\n", style.Failure("\xE2\x9C\x96"))
		}
	} else {
		term.Infof("%s No streams were found. Please make sure that your target is on an accessible network.\n", style.Failure(style.SymbolCross))
	}
}

func printErr(term *log.Terminal, err error) {
	term.Errorln(style.Failure(style.SymbolCross), err)
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
