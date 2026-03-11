package ports

// IsCommonHTTPPort returns true if the given port is a well-known HTTP port.
func IsCommonHTTPPort(port uint16) bool {
	switch port {
	case 80, 443, 8080, 8443:
		return true
	default:
		return false
	}
}
