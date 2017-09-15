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
	"io/ioutil"
	"log"
	"os/exec"

	"github.com/pkg/errors"
	v "gopkg.in/go-playground/validator.v9"
)

func runNmap(targets, ports string, resultFilePath string) error {
	cmd := exec.Command("/usr/local/bin/nmap", "-T4", "-A", targets, "-p", ports, "-oX", resultFilePath)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return errors.Wrap(err, "Couldn't get stdout pipe")
	}

	if err := cmd.Start(); err != nil {
		return errors.Wrap(err, "Coudln't run nmap command")
	}

	in := bufio.NewScanner(stdout)

	log.Print("Reading nmap output...")
	for in.Scan() {
		log.Printf(in.Text())
	}
	log.Print("Done!")

	if err := in.Err(); err != nil {
		log.Printf("error: %s", err)
	}

	return nil
}

// ParseNmapResult returns a slice of streams from an NMap XML result file
// To generate one yourself, use the -X option when running NMap
func ParseNmapResult(nmapResultFilePath string) ([]Stream, error) {
	var streams []Stream
	validate := v.New()

	log.Println("Reading file...")
	content, err := ioutil.ReadFile(nmapResultFilePath)
	if err != nil {
		return streams, errors.Wrap(err, "Could not read nmap result file at "+nmapResultFilePath+":")
	}
	log.Println("OK!")

	log.Println("Parsing results...")
	result := &NmapResult{}
	err = xml.Unmarshal(content, &result)
	if err != nil {
		log.Println("Unmarshall error:", err)
	} else {
		log.Printf("Result:\n%v", result)
	}

	for idx, host := range result.Hosts {
		log.Println("Parsing host", idx, "...")
		log.Println("Found", host.Address.AddrType, "address", host.Address.Addr, "for host", idx)

		if host.Ports.Ports == nil {
			log.Println("No ports found for host", idx)
			continue
		}
		for _, port := range host.Ports.Ports {
			err = validate.Struct(port)
			if err != nil {
				log.Println("Invalid port found:", err)
				continue
			}
			log.Println(port.State.State, "port", port.PortID, "found running service", port.Service.Name)
			streams = append(streams, Stream{
				device:  port.Service.Product,
				address: host.Address.Addr,
				port:    port.PortID,
			})
		}
	}
	log.Println("OK!")

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
	nmapResultPath := "/tmp/cameradar_scan.xml"

	err := runNmap(targets, ports, nmapResultPath)
	if err != nil {
		return streams, err
	}

	streams, err = ParseNmapResult(nmapResultPath)
	if err != nil {
		return streams, err
	}

	return streams, nil
}
