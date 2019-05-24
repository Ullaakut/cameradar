package main

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/ullaakut/cameradar"
	log "github.com/ullaakut/disgo"
	"github.com/ullaakut/disgo/style"
)

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
	term := log.NewTerminal()

	err := parseArguments()
	if err != nil {
		printErr(term, err)
	}

	// options.Credentials = viper.GetString("custom-credentials")
	// options.EnableLogs = viper.GetBool("log") || viper.GetBool("logging")
	// options.Ports = viper.GetStringSlice("ports")
	// options.Routes = viper.GetString("custom-routes")
	// options.Speed = viper.GetInt("speed")
	// options.Timeout = viper.GetInt("timeout")
	// options.Targets = viper.GetStringSlice("targets")

	cameradar := cameradar.New(
		cameradar.WithTargets(viper.GetStringSlice("targets")),
		cameradar.WithPorts(viper.GetStringSlice("ports")),
		cameradar.WithDebug(viper.GetBool("log") || viper.GetBool("logging")),
		cameradar.WithScanSummary(true),
	)

	cameradar.Scan()

	// TODO: Move this logic to cameradar library.
	// if len(options.Targets) == 1 {
	// 	options.Targets, err = cameradar.ParseTargetsFile(options.Targets[0])
	// 	if err != nil {
	// 		printErr(term, err)
	// 	}
	// }

	// err = curl.GlobalInit(curl.GLOBAL_ALL)
	// handle := curl.EasyInit()
	// if err != nil || handle == nil {
	// 	printErr(term, errors.New("libcurl initialization failed"))
	// }

	// c := &cameradar.Curl{CURL: handle}
	// defer curl.GlobalCleanup()

	// gopath := os.Getenv("GOPATH")
	// options.Credentials = strings.Replace(options.Credentials, "<GOPATH>", gopath, 1)
	// options.Routes = strings.Replace(options.Routes, "<GOPATH>", gopath, 1)

	// credentials, err := cameradar.LoadCredentials(options.Credentials)
	// if err != nil {
	// 	printErr(term, fmt.Errorf("Invalid credentials dictionary %q: %v", options.Credentials, err))
	// 	return
	// }

	// routes, err := cameradar.LoadRoutes(options.Routes)
	// if err != nil {
	// 	printErr(term, fmt.Errorf("Invalid routes dictionary %q: %v", options.Routes, err))
	// 	return
	// }

	// streams, err := cameradar.Discover(options.Targets, options.Ports, options.Speed)
	// if err != nil && len(streams) > 0 {
	// 	printErr(term, err)
	// }

	// Most cameras will be accessed successfully with these two attacks
	streams, err = cameradar.AttackRoute(c, streams, routes, time.Duration(options.Timeout)*time.Millisecond, options.EnableLogs)
	if err != nil && len(streams) > 0 {
		printErr(term, err)
	}

	streams, err = cameradar.DetectAuthMethods(c, streams, time.Duration(options.Timeout)*time.Millisecond, options.EnableLogs)

	streams, err = cameradar.AttackCredentials(c, streams, credentials, time.Duration(options.Timeout)*time.Millisecond, options.EnableLogs)
	if err != nil && len(streams) > 0 {
		printErr(term, err)
	}

	// But some cameras run GST RTSP Server which prioritizes 401 over 404 contrary to most cameras.
	// For these cameras, running another route attack will solve the problem.
	for _, stream := range streams {
		if !stream.RouteFound || !stream.CredentialsFound {
			streams, err = cameradar.AttackRoute(c, streams, routes, time.Duration(options.Timeout)*time.Millisecond, options.EnableLogs)
			if err != nil && len(streams) > 0 {
				printErr(term, err)
			}

			break
		}
	}

	streams, err = cameradar.ValidateStreams(c, streams, time.Duration(options.Timeout)*time.Millisecond, options.EnableLogs)
	if err != nil && len(streams) > 0 {
		printErr(term, err)
	}

	prettyPrint(term, streams)
}

func printErr(term *log.Terminal, err error) {
	term.Errorln(style.Failure(style.SymbolCross), err)
	os.Exit(1)
}
