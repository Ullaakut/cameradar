package cameradar

import (
	"net/netip"
)

// AuthType represents the RTSP authentication method.
type AuthType int

// Supported authentication methods.
const (
	AuthNone AuthType = iota
	AuthBasic
	AuthDigest
)

// Stream represents a camera's RTSP stream.
type Stream struct {
	Device   string     `json:"device"`
	Username string     `json:"username"`
	Password string     `json:"password"`
	Routes   []string   `json:"route"`
	Address  netip.Addr `json:"address"  validate:"required"`
	Port     uint16     `json:"port"     validate:"required"`

	CredentialsFound bool `json:"credentials_found"`
	RouteFound       bool `json:"route_found"`
	Available        bool `json:"available"`

	AuthenticationType AuthType `json:"authentication_type"`
}

// Route returns this stream's route if there is one.
func (s Stream) Route() string {
	if len(s.Routes) > 0 {
		return s.Routes[0]
	}
	return ""
}
