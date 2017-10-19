// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/EtixLabs/cameradar"
	"github.com/gernest/wow"
	"github.com/gernest/wow/spin"

	"github.com/fatih/color"
	"github.com/jessevdk/go-flags"
)

type options struct {
	Target      string `short:"t" long:"target" description:"The target on which to scan for open RTSP streams - required (ex: 172.16.100.0/24)" required:"true"`
	Ports       string `short:"p" long:"ports" description:"The ports on which to search for RTSP streams" default:"554,8554"`
	OutputFile  string `short:"o" long:"nmap-output" description:"The path where nmap will create its XML result file" default:"/tmp/cameradar_scan.xml"`
	Routes      string `short:"r" long:"custom-routes" description:"The path on which to load a custom routes dictionary" default:"<GOPATH>/src/github.com/EtixLabs/cameradar/dictionaries/routes"`
	Credentials string `short:"c" long:"custom-credentials" description:"The path on which to load a custom credentials JSON dictionary" default:"<GOPATH>/src/github.com/EtixLabs/cameradar/dictionaries/credentials.json"`
	Speed       int    `short:"s" long:"speed" description:"The nmap speed preset to use" default:"4"`
	Timeout     int    `short:"T" long:"timeout" description:"The timeout in miliseconds to use for attack attempts" default:"2000"`
	EnableLogs  bool   `short:"l" long:"log" description:"Enable the logs for nmap's output to stdout"`
}

func main() {
	var options options
	_, err := flags.ParseArgs(&options, os.Args[1:])
	if err != nil {
		os.Exit(0)
	}

	w := startSpinner(options.EnableLogs)

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
	streams, _ := cmrdr.Discover(options.Target, options.Ports, options.OutputFile, options.Speed, options.EnableLogs)

	// Most cameras will be accessed successfully with these two attacks

	updateSpinner(w, "Found "+fmt.Sprint(len(streams))+" streams. Attacking their routes...", options.EnableLogs)
	streams, _ = cmrdr.AttackRoute(streams, routes, time.Duration(options.Timeout)*time.Millisecond, options.EnableLogs)

	updateSpinner(w, "Found "+fmt.Sprint(len(streams))+" streams. Attacking their credentials...", options.EnableLogs)
	streams, _ = cmrdr.AttackCredentials(streams, credentials, time.Duration(options.Timeout)*time.Millisecond, options.EnableLogs)

	// But some cameras run GST RTSP Server which prioritizes 401 over 404 contrary to most cameras.
	// For these cameras, running another route attack will solve the problem.
	for _, stream := range streams {
		if stream.RouteFound == false || stream.CredentialsFound == false {

			updateSpinner(w, "Found "+fmt.Sprint(len(streams))+" streams. Final attack...", options.EnableLogs)
			streams, _ = cmrdr.AttackRoute(streams, routes, time.Duration(options.Timeout)*time.Millisecond, options.EnableLogs)
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
