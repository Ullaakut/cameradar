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
