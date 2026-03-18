package ports

import (
	"strings"
)

// InferTunnelScheme returns the likely scheme for a given port and optional service name.
func InferTunnelScheme(port uint16, serviceName string) string {
	if len(serviceName) > 0 {
		name := strings.ToLower(strings.TrimSpace(serviceName))
		switch name {
		case "rtsps", "https", "http":
			return name
		}
	}

	switch port {
	case 443, 8443:
		return "https"
	case 80, 8080:
		return "http"
	}

	return ""
}
