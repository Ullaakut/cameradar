package ports

import (
	"strings"
)

// InferTunnelScheme returns the likely scheme for a given port and optional service name.
func InferTunnelScheme(port uint16, serviceName string) string {
	if len(serviceName) > 0 {
		name := strings.ToLower(strings.TrimSpace(serviceName))
		switch name {
		case "rtsps":
			return "rtsps"
		case "https":
			return "https"
		case "http":
			return "http"
		}
	}

	if port != 80 && port != 443 && port != 8080 && port != 8443 {
		return ""
	}
	switch port {
	case 443, 8443:
		return "https"
	default:
		return "http"
	}
}
