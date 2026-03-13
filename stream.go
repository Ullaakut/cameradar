package cameradar

import (
	"net"
	"net/netip"
	"net/url"
	"strconv"
	"strings"

	"github.com/bluenviron/gortsplib/v5/pkg/base"
)

// AuthType represents the RTSP authentication method.
type AuthType int

// Supported authentication methods.
const (
	AuthUnknown AuthType = iota
	AuthNone
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

	CredentialsFound bool   `json:"credentials_found"`
	RouteFound       bool   `json:"route_found"`
	Available        bool   `json:"available"`
	Scheme           string `json:"scheme"`

	AuthenticationType AuthType `json:"authentication_type"`
}

func (s Stream) resolvedScheme() string {
	scheme := s.Scheme
	if scheme == "" {
		return "rtsp"
	}
	return scheme
}

func parseScheme(scheme string) string {
	switch scheme {
	case "http":
		return "rtsp"
	case "https":
		return "rtsps"
	default:
		return scheme
	}
}

// Route returns this stream's route if there is one.
func (s Stream) Route() string {
	if len(s.Routes) > 0 {
		return s.Routes[0]
	}
	return ""
}

// String builds the RTSP URL for this stream.
func (s Stream) String() string {
	scheme := s.resolvedScheme()

	host := net.JoinHostPort(s.Address.String(), strconv.Itoa(int(s.Port)))
	path := "/" + strings.TrimLeft(strings.TrimSpace(s.Route()), "/")

	u := &url.URL{
		Scheme: scheme,
		Host:   host,
		Path:   path,
	}
	if s.Username != "" || s.Password != "" {
		u.User = url.UserPassword(s.Username, s.Password)
	}

	return u.String()
}

// URL parses the stream's RTSP URL into a *base.URL.
func (s Stream) URL() (*base.URL, error) {
	s.Scheme = parseScheme(s.resolvedScheme())
	return base.ParseURL(s.String())
}
