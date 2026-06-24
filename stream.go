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

const (
	schemeRTSP  = "rtsp"
	schemeRTSPS = "rtsps"
	schemeHTTP  = "http"
	schemeHTTPS = "https"
)

// Stream represents a camera's stream, typically accessed over RTSP/RTSPS.
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
		return schemeRTSP
	}
	return scheme
}

func parseScheme(scheme string) string {
	switch scheme {
	case schemeHTTP:
		return schemeRTSP
	case schemeHTTPS:
		return schemeRTSPS
	default:
		return scheme
	}
}

// RTSPScheme returns the normalized scheme to use when rendering RTSP URLs.
// It returns "rtsps" for secure schemes ("rtsps" and "https"), and "rtsp" otherwise.
func (s Stream) RTSPScheme() string {
	scheme := parseScheme(strings.ToLower(strings.TrimSpace(s.resolvedScheme())))
	if scheme == schemeRTSPS {
		return schemeRTSPS
	}
	return schemeRTSP
}

// Route returns this stream's route if there is one.
func (s Stream) Route() string {
	if len(s.Routes) > 0 {
		return s.Routes[0]
	}
	return ""
}

// String builds the stream URL using the configured scheme, defaulting to rtsp.
func (s Stream) String() string {
	scheme := s.resolvedScheme()

	host := net.JoinHostPort(s.Address.String(), strconv.Itoa(int(s.Port)))
	route := strings.TrimLeft(strings.TrimSpace(s.Route()), "/")
	pathPart := "/" + route
	rawQuery := ""
	if i := strings.IndexByte(route, '?'); i >= 0 {
		pathPart = "/" + route[:i]
		rawQuery = route[i+1:]
	}

	u := &url.URL{
		Scheme:   scheme,
		Host:     host,
		Path:     pathPart,
		RawQuery: rawQuery,
	}
	if s.Username != "" || s.Password != "" {
		u.User = url.UserPassword(s.Username, s.Password)
	}

	return u.String()
}

// URL parses the stream URL into a *base.URL, normalizing http/https to rtsp/rtsps.
// The scheme is case-folded before normalization so that "HTTP" and "HTTPS" are
// treated identically to their lowercase equivalents, matching the behavior of
// RTSPScheme.
func (s Stream) URL() (*base.URL, error) {
	s.Scheme = parseScheme(strings.ToLower(strings.TrimSpace(s.resolvedScheme())))
	return base.ParseURL(s.String())
}
