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

import "encoding/xml"

// NmapResult is the structure that holds all the information from an NMap scan
type nmapResult struct {
	XMLName xml.Name `xml:"nmaprun"`
	Hosts   []host   `xml:"host" validate:"required"`
}

// Host represents a host discovered during a scan
type host struct {
	XMLName xml.Name `xml:"host"`
	Address address  `xml:"address"`
	Ports   ports    `xml:"ports"`
}

// Address is a host's address discovered during a scan
type address struct {
	XMLName  xml.Name `xml:"address"`
	Addr     string   `xml:"addr,attr"`
	AddrType string   `xml:"addrType,attr"`
}

// Ports is the list of openned ports on a host
type ports struct {
	XMLName xml.Name `xml:"ports"`
	Ports   []port   `xml:"port"`
}

// Port is a port found on a host during a scan
type port struct {
	XMLName xml.Name `xml:"port"`
	PortID  uint     `xml:"portid,attr"`
	State   state    `xml:"state"`
	Service service  `xml:"service"`
}

// State is the state of a port
type state struct {
	XMLName xml.Name `xml:"state"`
	State   string   `xml:"state,attr" validate:"required,eq=open"`
}

// Service represents the service that a port provides
type service struct {
	XMLName xml.Name `xml:"service"`
	Name    string   `xml:"name,attr" validate:"required,eq=rtsp"`
	Product string   `xml:"product,attr"`
}
