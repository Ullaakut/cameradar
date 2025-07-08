package ui

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/Ullaakut/cameradar/v6"
)

// FormatSummary builds a human-readable summary of discovered streams.
func FormatSummary(streams []cameradar.Stream, _ error) string {
	accessible, others := partitionStreams(streams)

	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("Accessible streams: %d\n", len(accessible)))
	if len(accessible) == 0 {
		builder.WriteString("• None\n")
	} else {
		for _, stream := range accessible {
			builder.WriteString(formatStream(stream))
		}
	}

	if len(others) > 0 {
		builder.WriteString("\n")
		builder.WriteString(fmt.Sprintf("Other discovered streams: %d\n", len(others)))
		for _, stream := range others {
			builder.WriteString(formatStream(stream))
		}
	}

	return builder.String()
}

func partitionStreams(streams []cameradar.Stream) ([]cameradar.Stream, []cameradar.Stream) {
	var accessible []cameradar.Stream
	var others []cameradar.Stream
	for _, stream := range streams {
		if stream.Available {
			accessible = append(accessible, stream)
		} else {
			others = append(others, stream)
		}
	}

	// Sort streams by address and port.
	sort.Slice(accessible, func(i, j int) bool {
		if accessible[i].Address.String() == accessible[j].Address.String() {
			return accessible[i].Port < accessible[j].Port
		}
		return accessible[i].Address.String() < accessible[j].Address.String()
	})
	sort.Slice(others, func(i, j int) bool {
		if others[i].Address.String() == others[j].Address.String() {
			return others[i].Port < others[j].Port
		}
		return others[i].Address.String() < others[j].Address.String()
	})

	return accessible, others
}

func formatStream(stream cameradar.Stream) string {
	var builder strings.Builder
	builder.WriteString("• ")
	builder.WriteString(stream.Address.String())
	builder.WriteString(":")
	builder.WriteString(strconv.FormatUint(uint64(stream.Port), 10))

	if stream.Device != "" {
		builder.WriteString(" (")
		builder.WriteString(stream.Device)
		builder.WriteString(")")
	}
	builder.WriteString("\n")

	builder.WriteString("  Authentication: ")
	builder.WriteString(authTypeLabel(stream.AuthenticationType))
	builder.WriteString("\n")

	if len(stream.Routes) > 0 {
		builder.WriteString("  Routes: ")
		builder.WriteString(strings.Join(stream.Routes, ", "))
		builder.WriteString("\n")
	} else {
		builder.WriteString("  Routes: not found\n")
	}

	if stream.CredentialsFound {
		builder.WriteString("  Credentials: ")
		builder.WriteString(stream.Username)
		builder.WriteString(":")
		builder.WriteString(stream.Password)
		builder.WriteString("\n")
	} else {
		builder.WriteString("  Credentials: not found\n")
	}

	builder.WriteString("  Availability: ")
	if stream.Available {
		builder.WriteString("yes\n")
	} else {
		builder.WriteString("no\n")
	}

	if stream.RouteFound && stream.CredentialsFound {
		builder.WriteString("  RTSP URL: ")
		builder.WriteString(formatRTSPURL(stream))
		builder.WriteString("\n")
	}

	builder.WriteString("  Admin panel: ")
	builder.WriteString(formatAdminPanelURL(stream))
	builder.WriteString("\n")

	return builder.String()
}

func formatRTSPURL(stream cameradar.Stream) string {
	path := stream.Route()
	if path != "" && !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	credentials := ""
	if stream.Username != "" || stream.Password != "" {
		credentials = stream.Username + ":" + stream.Password + "@"
	}

	return fmt.Sprintf("rtsp://%s%s:%d%s", credentials, stream.Address.String(), stream.Port, path)
}

func formatAdminPanelURL(stream cameradar.Stream) string {
	return fmt.Sprintf("http://%s/", stream.Address.String())
}

func authTypeLabel(auth cameradar.AuthType) string {
	switch auth {
	case cameradar.AuthNone:
		return "none"
	case cameradar.AuthBasic:
		return "basic"
	case cameradar.AuthDigest:
		return "digest"
	default:
		return fmt.Sprintf("unknown(%d)", auth)
	}
}
