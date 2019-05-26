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

func TestGetCameraRTSPURL(t *testing.T) {
	validStream := Stream{
		Address:  "1.2.3.4",
		Username: "ullaakut",
		Password: "ba69897483886f0d2b0afb6345b76c0c",
		Route:    "cameradar.sdp",
		Port:     1337,
	}

	testCases := []struct {
		stream Stream

		expectedRTSPURL string
	}{
		{
			stream: validStream,

			expectedRTSPURL: "rtsp://ullaakut:ba69897483886f0d2b0afb6345b76c0c@1.2.3.4:1337/cameradar.sdp",
		},
	}

	for _, test := range testCases {
		assert.Equal(t, test.expectedRTSPURL, GetCameraRTSPURL(test.stream))
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
