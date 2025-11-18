package cameradar

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReplace(t *testing.T) {
	validStream1 := Stream{
		Device:  "fakeDevice",
		Address: "fakeAddress",
		Port:    1,
	}

	validStream2 := Stream{
		Device:  "fakeDevice",
		Address: "differentFakeAddress",
		Port:    2,
	}

	invalidStream := Stream{
		Device:  "invalidDevice",
		Address: "anotherFakeAddress",
		Port:    3,
	}

	invalidStreamModified := Stream{
		Device:  "updatedDevice",
		Address: "anotherFakeAddress",
		Port:    3,
	}

	testCases := []struct {
		streams   []Stream
		newStream Stream

		expectedStreams []Stream
	}{
		{
			streams:   []Stream{validStream1, validStream2, invalidStream},
			newStream: invalidStreamModified,

			expectedStreams: []Stream{validStream1, validStream2, invalidStreamModified},
		},
	}

	for _, test := range testCases {
		streams := replace(test.streams, test.newStream)

		assert.Equal(t, len(test.expectedStreams), len(streams))

		for _, expectedStream := range test.expectedStreams {
			assert.Contains(t, streams, expectedStream)
		}
	}
}

func TestNormalizeRoute(t *testing.T) {
	testCases := []struct {
		route          string
		expectedRoute  string
		description    string
	}{
		{
			route:         "/live.sdp",
			expectedRoute: "live.sdp",
			description:   "Remove leading slash",
		},
		{
			route:         "live.sdp",
			expectedRoute: "live.sdp",
			description:   "No leading slash to remove",
		},
		{
			route:         "/",
			expectedRoute: "",
			description:   "Single slash becomes empty",
		},
		{
			route:         "//",
			expectedRoute: "/",
			description:   "Double slash becomes single slash",
		},
		{
			route:         "",
			expectedRoute: "",
			description:   "Empty route stays empty",
		},
		{
			route:         "/path/to/stream",
			expectedRoute: "path/to/stream",
			description:   "Remove leading slash from path",
		},
	}

	for _, test := range testCases {
		result := normalizeRoute(test.route)
		assert.Equal(t, test.expectedRoute, result, test.description)
	}
}

func TestGetCameraRTSPURL(t *testing.T) {
	testCases := []struct {
		stream Stream

		expectedRTSPURL string
		description     string
	}{
		{
			stream: Stream{
				Address:  "1.2.3.4",
				Username: "ullaakut",
				Password: "ba69897483886f0d2b0afb6345b76c0c",
				Routes:   []string{"cameradar.sdp"},
				Port:     1337,
			},
			expectedRTSPURL: "rtsp://ullaakut:ba69897483886f0d2b0afb6345b76c0c@1.2.3.4:1337/cameradar.sdp",
			description:     "Route without leading slash",
		},
		{
			stream: Stream{
				Address:  "1.2.3.4",
				Username: "ullaakut",
				Password: "ba69897483886f0d2b0afb6345b76c0c",
				Routes:   []string{"/cameradar.sdp"},
				Port:     1337,
			},
			expectedRTSPURL: "rtsp://ullaakut:ba69897483886f0d2b0afb6345b76c0c@1.2.3.4:1337/cameradar.sdp",
			description:     "Route with leading slash should be normalized",
		},
		{
			stream: Stream{
				Address:  "119.x.153.116",
				Username: "",
				Password: "",
				Routes:   []string{"/"},
				Port:     554,
			},
			expectedRTSPURL: "rtsp://:@119.x.153.116:554/",
			description:     "Route with single slash should not create double slash",
		},
		{
			stream: Stream{
				Address:  "119.x.153.116",
				Username: "",
				Password: "",
				Routes:   []string{"//"},
				Port:     554,
			},
			expectedRTSPURL: "rtsp://:@119.x.153.116:554//",
			description:     "Route with double slash should become single slash after base URL slash",
		},
	}

	for _, test := range testCases {
		result := GetCameraRTSPURL(test.stream)
		assert.Equal(t, test.expectedRTSPURL, result, test.description)
	}
}

func TestGetCameraAdminPanelURL(t *testing.T) {
	validStream := Stream{
		Address: "1.2.3.4",
	}

	testCases := []struct {
		stream Stream

		expectedRTSPURL string
	}{
		{
			stream: validStream,

			expectedRTSPURL: "http://1.2.3.4/",
		},
	}

	for _, test := range testCases {
		assert.Equal(t, test.expectedRTSPURL, GetCameraAdminPanelURL(test.stream))
	}
}
