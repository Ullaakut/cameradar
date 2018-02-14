package main

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/EtixLabs/cameradar"

	curl "github.com/andelf/go-curl"
	"github.com/fatih/color"
	"github.com/gernest/wow"
	"github.com/gernest/wow/spin"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type options struct {
	Target      string
	Ports       string
	OutputFile  string
	Routes      string
	Credentials string
	Speed       int
	Timeout     int
	EnableLogs  bool
}

func parseArguments() error {

	viper.BindEnv("target", "CAMERADAR_TARGET")
	viper.BindEnv("ports", "CAMERADAR_PORTS")
	viper.BindEnv("nmap-output", "CAMERADAR_NMAP_OUTPUT_FILE")
	viper.BindEnv("custom-routes", "CAMERADAR_CUSTOM_ROUTES")
	viper.BindEnv("custom-credentials", "CAMERADAR_CUSTOM_CREDENTIALS")
	viper.BindEnv("speed", "CAMERADAR_SPEED")
	viper.BindEnv("timeout", "CAMERADAR_TIMEOUT")
	viper.BindEnv("envlogs", "CAMERADAR_LOGS")

	pflag.StringP("target", "t", "", "The target on which to scan for open RTSP streams - required (ex: 172.16.100.0/24)")
	pflag.StringP("ports", "p", "554,8554", "The ports on which to search for RTSP streams")
	pflag.StringP("nmap-output", "o", "/tmp/cameradar_scan.xml", "The path where nmap will create its XML result file")
	pflag.StringP("custom-routes", "r", "<GOPATH>/src/github.com/EtixLabs/cameradar/dictionaries/routes", "The path on which to load a custom routes dictionary")
	pflag.StringP("custom-credentials", "c", "<GOPATH>/src/github.com/EtixLabs/cameradar/dictionaries/credentials.json", "The path on which to load a custom credentials JSON dictionary")
	pflag.IntP("speed", "s", 4, "The nmap speed preset to use")
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

	if viper.GetString("target") == "" {
		return errors.New("target (-t, --target) argument required\n    examples:\n      - 172.16.100.0/24\n      - localhost\n      - 8.8.8.8")
	}

	return nil
}

func main() {
	var options options

	err := parseArguments()
	if err != nil {
		printErr(err)
	}

	options.Credentials = viper.GetString("custom-credentials")
	options.EnableLogs = viper.GetBool("log") || viper.GetBool("envlogs")
	options.OutputFile = viper.GetString("nmap-output")
	options.Ports = viper.GetString("ports")
	options.Routes = viper.GetString("custom-routes")
	options.Speed = viper.GetInt("speed")
	options.Timeout = viper.GetInt("timeout")
	options.Target = viper.GetString("target")

	w := startSpinner(options.EnableLogs)

	err = curl.GlobalInit(curl.GLOBAL_ALL)
	c := curl.EasyInit()
	if err != nil || c == nil {
		printErr(errors.New("libcurl initialization failed"))
	}
	defer curl.GlobalCleanup()

	updateSpinner(w, "Loading dictionaries...", options.EnableLogs)
	gopath := os.Getenv("GOPATH")
	options.Credentials = strings.Replace(options.Credentials, "<GOPATH>", gopath, 1)
	options.Routes = strings.Replace(options.Routes, "<GOPATH>", gopath, 1)

	credentials, err := cmrdr.LoadCredentials(options.Credentials)
	if err != nil {
		color.Red("Invalid credentials dictionary: %s", err.Error())
		return
	}

	routes, err := cmrdr.LoadRoutes(options.Routes)
	if err != nil {
		color.Red("Invalid routes dictionary: %s", err.Error())
		return
	}

	updateSpinner(w, "Scanning the network...", options.EnableLogs)
	streams, err := cmrdr.Discover(options.Target, options.Ports, options.OutputFile, options.Speed, options.EnableLogs)
	if err != nil && len(streams) > 0 {
		printErr(err)
	}

	// Most cameras will be accessed successfully with these two attacks
	updateSpinner(w, "Found "+fmt.Sprint(len(streams))+" streams. Attacking their routes...", options.EnableLogs)
	streams, err = cmrdr.AttackRoute(c, streams, routes, time.Duration(options.Timeout)*time.Millisecond, options.EnableLogs)
	if err != nil && len(streams) > 0 {
		printErr(err)
	}

	updateSpinner(w, "Found "+fmt.Sprint(len(streams))+" streams. Attacking their credentials...", options.EnableLogs)
	streams, err = cmrdr.AttackCredentials(c, streams, credentials, time.Duration(options.Timeout)*time.Millisecond, options.EnableLogs)
	if err != nil && len(streams) > 0 {
		printErr(err)
	}

	// But some cameras run GST RTSP Server which prioritizes 401 over 404 contrary to most cameras.
	// For these cameras, running another route attack will solve the problem.
	for _, stream := range streams {
		if stream.RouteFound == false || stream.CredentialsFound == false {

			updateSpinner(w, "Found "+fmt.Sprint(len(streams))+" streams. Final attack...", options.EnableLogs)
			streams, err = cmrdr.AttackRoute(c, streams, routes, time.Duration(options.Timeout)*time.Millisecond, options.EnableLogs)
			if err != nil && len(streams) > 0 {
				printErr(err)
			}
			break
		}
	}

	clearOutput(w, options.EnableLogs)
	prettyPrint(streams)
}

func prettyPrint(streams []cmrdr.Stream) {
	yellow := color.New(color.FgYellow, color.Bold, color.Underline).SprintFunc()
	blue := color.New(color.FgBlue, color.Underline).SprintFunc()
	green := color.New(color.FgGreen, color.Bold).SprintFunc()
	red := color.New(color.FgRed, color.Bold).SprintFunc()
	white := color.New(color.Italic).SprintFunc()

	success := 0

	if len(streams) > 0 {
		for _, stream := range streams {
			if stream.CredentialsFound && stream.RouteFound {
				fmt.Printf("%s\tDevice RTSP URL:\t%s\n", green("\xE2\x96\xB6"), blue(cmrdr.GetCameraRTSPURL(stream)))
				success++
			} else {
				fmt.Printf("%s\tAdmin panel URL:\t%s %s\n", red("\xE2\x96\xB6"), yellow(cmrdr.GetCameraAdminPanelURL(stream)), white("You can use this URL to try attacking the camera's admin panel instead."))
			}

			fmt.Printf("\tDevice model:\t\t%s\n\n", stream.Device)
			fmt.Printf("\tIP address:\t\t%s\n", stream.Address)
			fmt.Printf("\tRTSP port:\t\t%d\n", stream.Port)
			if stream.CredentialsFound {
				fmt.Printf("\tUsername:\t\t%s\n", green(stream.Username))
				fmt.Printf("\tPassword:\t\t%s\n", green(stream.Password))
			} else {
				fmt.Printf("\tUsername:\t\t%s\n", red("not found"))
				fmt.Printf("\tPassword:\t\t%s\n", red("not found"))
			}
			if stream.RouteFound {
				fmt.Printf("\tRTSP route:\t\t%s\n\n\n", green("/"+stream.Route))
			} else {
				fmt.Printf("\tRTSP route:\t\t%s\n\n\n", red("not found"))
			}
		}
		if success > 1 {
			fmt.Printf("%s Successful attack: %s devices were accessed", green("\xE2\x9C\x94"), green(len(streams)))
		} else if success == 1 {
			fmt.Printf("%s Successful attack: %s device was accessed", green("\xE2\x9C\x94"), green(len(streams)))
		} else {
			fmt.Printf("%s Streams were found but none were accessed. They are most likely configured with secure credentials and routes. You can try adding entries to the dictionary or generating your own in order to attempt a bruteforce attack on the cameras.\n", red("\xE2\x9C\x96"))
		}
	} else {
		fmt.Printf("%s No streams were found. Please make sure that your target is on an accessible network.\n", red("\xE2\x9C\x96"))
	}
}

func printErr(err error) {
	red := color.New(color.FgRed, color.Bold).SprintFunc()
	fmt.Printf("%s %v\n", red("\xE2\x9C\x96"), err)
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
