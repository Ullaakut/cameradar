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

// Stream represents a camera's RTSP stream
type Stream struct {
	device   string
	username string
	password string
	route    string
	address  string `validate:"required"`
	port     uint   `validate:"required"`
}

// Attack represents the state of the attack on a stream
type Attack struct {
	stream           Stream
	credentialsFound bool
	routeFound       bool
}

// Credentials is a map of credentials
// usernames are keys and passwords are values
// creds['admin'] -> 'secure_password'
type Credentials map[string]string

// Routes is a slice of Routes
// ['/live.sdp', '/media.amp', ...]
type Routes []string

// NmapResult is the structure that holds all the information from an NMap scan
type NmapResult struct {
	XMLName xml.Name `xml:"nmaprun"`
	Hosts   []Host   `xml:"host" validate:"required"`
}

// Host represents a host discovered during a scan
type Host struct {
	XMLName xml.Name `xml:"host"`
	Address Address  `xml:"address"`
	Ports   Ports    `xml:"ports"`
}

// Address is a host's address discovered during a scan
type Address struct {
	XMLName  xml.Name `xml:"address"`
	Addr     string   `xml:"addr,attr"`
	AddrType string   `xml:"addrType,attr"`
}

// Ports is the list of openned ports on a host
type Ports struct {
	XMLName xml.Name `xml:"ports"`
	Ports   []Port   `xml:"port"`
}

// Port is a port found on a host during a scan
type Port struct {
	XMLName xml.Name `xml:"port"`
	PortID  uint     `xml:"portid,attr"`
	State   State    `xml:"state"`
	Service Service  `xml:"service"`
}

// State is the state of a port
type State struct {
	XMLName xml.Name `xml:"state"`
	State   string   `xml:"state,attr" validate:"required,eq=open"`
}

// Service represents the service that a port provides
type Service struct {
	XMLName xml.Name `xml:"service"`
	Name    string   `xml:"name,attr" validate:"required,eq=rtsp"`
	Product string   `xml:"product,attr"`
}
