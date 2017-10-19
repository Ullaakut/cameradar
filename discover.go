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

package cmrdr

import (
	"bufio"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os/exec"

	"github.com/pkg/errors"
	v "gopkg.in/go-playground/validator.v9"
)

// These constants detail the different level of nmap speed presets
// that determine the timeout values and wether or not nmap makes use of parallelism
const (
	// PARANOIAC 	NO PARALLELISM | 5min  timeout | 100ms to 10s    round-trip time timeout	 |  5mn   scan delay
	PARANOIAC = 0
	// SNEAKY 		NO PARALLELISM | 15sec timeout | 100ms to 10s    round-trip time timeout	 |  15s   scan delay
	SNEAKY = 1
	// POLITE 		NO PARALLELISM | 1sec  timeout | 100ms to 10s    round-trip time timeout	 |  400ms scan delay
	POLITE = 2
	// NORMAL 		PARALLELISM	   | 1sec  timeout | 100ms to 10s    round-trip time timeout	 |  0s    scan delay
	NORMAL = 3
	// AGGRESSIVE 	PARALLELISM	   | 500ms timeout | 100ms to 1250ms round-trip time timeout	 |  0s    scan delay
	AGGRESSIVE = 4
	// INSANE 		PARALLELISM	   | 250ms timeout |  50ms to 300ms  round-trip time timeout	 |  0s    scan delay
	INSANE = 5
)

// Allows unit tests to override the exec function to avoid launching a real command
// during the tests. The NmapRun method will soon be refactored with an adaptor in order
// to make it possible to mock all external calls.
var execCommand = exec.Command

// NmapRun runs nmap on the specified targets's specified ports, using the given nmap speed.
func NmapRun(targets, ports, resultFilePath string, nmapSpeed int, enableLogs bool) error {
	if nmapSpeed < PARANOIAC || nmapSpeed > INSANE {
		return fmt.Errorf("invalid nmap speed value '%d'. Should be between '%d' and '%d'", nmapSpeed, PARANOIAC, INSANE)
	}

	// Prepare nmap command
	cmd := execCommand(
		"nmap",
		fmt.Sprintf("-T%d", nmapSpeed),
		"-A",
		"-p",
		ports,
		"-oX",
		resultFilePath,
		targets,
	)

	// Pipe stdout to be able to write the logs in realtime
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return errors.Wrap(err, "couldn't get stdout pipe")
	}

	// Execute the nmap command
	if err := cmd.Start(); err != nil {
		return errors.Wrap(err, "coudln't run nmap command")
	}

	// Scan the pipe until an end of file or an error occurs
	in := bufio.NewScanner(stdout)
	for in.Scan() {
		if enableLogs {
			fmt.Println(in.Text())
		}
	}
	if err := in.Err(); err != nil {
		if enableLogs {
			fmt.Printf("error: %s\n", err)
		}
	}

	return nil
}

// NmapParseResults returns a slice of streams from an NMap XML result file.
// To generate one yourself, use the -X option when running NMap.
func NmapParseResults(nmapResultFilePath string) ([]Stream, error) {
	var streams []Stream

	// Open & Read XML file
	content, err := ioutil.ReadFile(nmapResultFilePath)
	if err != nil {
		return streams, errors.Wrap(err, "could not read nmap result file at "+nmapResultFilePath+":")
	}

	// Unmarshal content of XML file into data structure
	result := &nmapResult{}
	err = xml.Unmarshal(content, &result)
	if err != nil {
		return streams, err
	}

	// Iterate on hosts to try to find hosts with ports that
	//     - serve RTSP
	//     - are open
	validate := v.New()
	for _, host := range result.Hosts {
		if host.Ports.Ports == nil {
			continue
		}
		for _, port := range host.Ports.Ports {
			err = validate.Struct(port)
			if err != nil {
				continue
			}
			streams = append(streams, Stream{
				Device:  port.Service.Product,
				Address: host.Address.Addr,
				Port:    port.PortID,
			})
		}
	}

	return streams, nil
}

// Discover scans the target networks and tries to find RTSP streams within them.
//
// targets can be:
//
//    - a subnet (e.g.: 172.16.100.0/24)
//    - an IP (e.g.: 172.16.100.10)
//    - a hostname (e.g.: localhost)
//    - a range of IPs (e.g.: 172.16.100.10-20)
//
// ports can be:
//
//    - one or multiple ports and port ranges separated by commas (e.g.: 554,8554-8560,18554-28554)
func Discover(targets, ports, nmapResultPath string, speed int, log bool) ([]Stream, error) {
	var streams []Stream

	// Run nmap command to discover open ports on the specified targets & ports
	err := NmapRun(targets, ports, nmapResultPath, speed, log)
	if err != nil {
		return streams, err
	}

	// Get found streams from nmap results
	streams, err = NmapParseResults(nmapResultPath)
	if err != nil {
		return streams, err
	}

	return streams, nil
}
