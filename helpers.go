package cameradar

import (
	"fmt"
	"strings"
)

func replace(streams []Stream, new Stream) []Stream {
	var updatedSlice []Stream

	for _, old := range streams {
		if old.Address == new.Address && old.Port == new.Port {
			updatedSlice = append(updatedSlice, new)
		} else {
			updatedSlice = append(updatedSlice, old)
		}
	}

	return updatedSlice
}

// normalizeRoute ensures route has proper formatting to prevent double slashes
func normalizeRoute(route string) string {
	// Remove leading slash if present to avoid double slashes when concatenating
	return strings.TrimPrefix(route, "/")
}

// GetCameraRTSPURL generates a stream's RTSP URL.
func GetCameraRTSPURL(stream Stream) string {
	route := normalizeRoute(stream.Route())
	return "rtsp://" + stream.Username + ":" + stream.Password + "@" + stream.Address + ":" + fmt.Sprint(stream.Port) + "/" + route
}

// GetCameraAdminPanelURL returns the URL to the camera's admin panel.
func GetCameraAdminPanelURL(stream Stream) string {
	return "http://" + stream.Address + "/"
}
