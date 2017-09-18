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
	"log"
	"os/exec"

	"github.com/pkg/errors"
	v "gopkg.in/go-playground/validator.v9"
)

// These constants detail the different level of nmap aggressivity
// that determine the timeout values and wether or not nmap makes use of parallelism
const (
	// PARANOID 	NO PARALLELISM | 5min  timeout | 100ms to 10s    round-trip time timeout	 |  5mn   scan delay
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

// RunNmap runs nmap on the specified targets's specified ports, using the given nmapAggressivity
func RunNmap(targets, ports string, resultFilePath string, nmapAggressivity uint8) error {
	// Prepare nmap command
	cmd := exec.Command(
		"nmap",
		fmt.Sprintf("-T%d", nmapAggressivity),
		"-A",
		targets,
		"-p",
		ports,
		"-oX",
		resultFilePath,
	)

	// Pipe stdout to be able to write the logs in realtime
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return errors.Wrap(err, "Couldn't get stdout pipe")
	}

	// Execute the nmap command
	if err := cmd.Start(); err != nil {
		return errors.Wrap(err, "Coudln't run nmap command")
	}

	// Scan the pipe until an end of file or an error occurs
	in := bufio.NewScanner(stdout)
	for in.Scan() {
		log.Printf(in.Text())
	}
	if err := in.Err(); err != nil {
		log.Printf("error: %s", err)
	}

	return nil
}

// ParseNmapResult returns a slice of streams from an NMap XML result file
// To generate one yourself, use the -X option when running NMap
func ParseNmapResult(nmapResultFilePath string) ([]Stream, error) {
	var streams []Stream

	// Open & Read XML file
	content, err := ioutil.ReadFile(nmapResultFilePath)
	if err != nil {
		return streams, errors.Wrap(err, "Could not read nmap result file at "+nmapResultFilePath+":")
	}

	// Unmarshal content of XML file into data structure
	result := &NmapResult{}
	err = xml.Unmarshal(content, &result)
	if err != nil {
		log.Println("Unmarshall error:", err)
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
				device:  port.Service.Product,
				address: host.Address.Addr,
				port:    port.PortID,
			})
		}
	}

	return streams, nil
}

// Discover scans the target networks and tries to find RTSP streams within them
// targets - string: The addresses
//    - a subnet (e.g.: 172.16.100.0/24)
//    - an IP (e.g.: 172.16.100.10)
//    - a hostname (e.g.: localhost)
//    - a range of IPs (e.g.: 172.16.100.10-172.16.100.20)
//    - a mix of all those separated by commas (e.g.: localhost,172.17.100.0/24,172.16.100.10-172.16.100.20,0.0.0.0).
// ports - string :
//    - one or multiple ports and port ranges separated by commas (e.g.: 554,8554-8560,18554-28554)
func Discover(targets string, ports string) ([]Stream, error) {
	var streams []Stream

	// TODO: Provide configuration for this
	nmapResultPath := "/tmp/cameradar_scan.xml"

	// Run nmap command to discover open ports on the specified targets & ports
	err := RunNmap(targets, ports, nmapResultPath, 4)
	if err != nil {
		return streams, err
	}

	// Get found streams from nmap results
	streams, err = ParseNmapResult(nmapResultPath)
	if err != nil {
		return streams, err
	}

	return streams, nil
}
