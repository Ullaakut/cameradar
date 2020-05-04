package cameradar

import "fmt"

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

// GetCameraRTSPURL generates a stream's RTSP URL.
func GetCameraRTSPURL(stream Stream) string {
	return "rtsp://" + stream.Username + ":" + stream.Password + "@" + stream.Address + ":" + fmt.Sprint(stream.Port) + "/" + stream.Route()
}

// GetCameraAdminPanelURL returns the URL to the camera's admin panel.
func GetCameraAdminPanelURL(stream Stream) string {
	return "http://" + stream.Address + "/"
}
