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
	"time"

	"github.com/EtixLabs/cameradar/cameradar"
	"github.com/fatih/color"
	"github.com/jessevdk/go-flags"
)

type options struct {
	Target     string `short:"t" long:"target" description:"The target on which to scan for open RTSP streams - required" required:"true"`
	Ports      string `short:"p" long:"ports" description:"The ports on which to search for RTSP streams" default:"554,8554"`
	OutputFile string `short:"o" long:"nmap-output" description:"The path where nmap will create its XML result file" default:"/tmp/cameradar_scan.xml"`
	Speed      int    `short:"s" long:"speed" description:"The nmap speed preset to use" default:"4"`
	Timeout    int    `short:"T" long:"timeout" description:"The timeout in miliseconds to use for attack attempts" default:"1000"`
	EnableLogs bool   `short:"l" long:"log" description:"Enable the logs for nmap's output to stdout"`
}

func main() {
	var options options
	_, err := flags.ParseArgs(&options, os.Args[1:])
	if err != nil {
		os.Exit(0)
	}

	streams, err := cmrdr.Discover(options.Target, options.Ports, options.OutputFile, options.Speed, options.EnableLogs)
	if err != nil {
		color.Red("Cloud not discover")
	}

	credentials, err := cmrdr.LoadCredentials("./dictionaries/credentials.json")
	if err != nil {
		color.Red("Invalid credentials dictionary: %s", err.Error())
		return
	}

	routes, err := cmrdr.LoadRoutes("./dictionaries/routes")
	if err != nil {
		color.Red("Invalid routes dictionary: %s", err.Error())
		return
	}

	streams, err = cmrdr.AttackRoute(streams, routes, time.Duration(options.Timeout)*time.Millisecond, options.EnableLogs)
	if err != nil {
		color.Red("Could not attack routes")
	}

	streams, err = cmrdr.AttackCredentials(streams, credentials, time.Duration(options.Timeout)*time.Millisecond, options.EnableLogs)
	if err != nil {
		color.Red("Cloud not attack credentials")
	}

	prettyPrint(streams)
}

func prettyPrint(streams []cmrdr.Stream) {
	blue := color.New(color.FgBlue, color.Underline).SprintFunc()
	green := color.New(color.FgGreen, color.Bold).SprintFunc()
	red := color.New(color.FgRed, color.Bold).SprintFunc()
	if len(streams) > 0 {
		for _, stream := range streams {
			fmt.Printf("Device RTSP URL:\t%s\n", blue(cmrdr.RTSPURL(stream)))
			fmt.Printf("Device model:\t\t%s\n\n", stream.Device)
			fmt.Printf("IP address:\t\t%s\n", stream.Address)
			fmt.Printf("RTSP port:\t\t%d\n", stream.Port)
			fmt.Printf("Username:\t\t%s\n", green(stream.Username))
			fmt.Printf("Password:\t\t%s\n", green(stream.Password))
			fmt.Printf("RTSP route:\t\t%s\n\n\n", green("/"+stream.Route))
		}
		if len(streams) > 1 {
			fmt.Printf("%s Successful attack: %s devices were accessed", green("\xE2\x9C\x94"), green(len(streams)))
		} else {
			fmt.Printf("%s Successful attack: %s device was accessed", green("\xE2\x9C\x94"), green(len(streams)))
		}
	} else {
		fmt.Printf("%s No streams were found. Please make sure that your target is on an accessible network.", red("\xE2\x9C\x96"))
	}
}
