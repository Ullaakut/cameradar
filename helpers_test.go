package cmrdr

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReplace(t *testing.T) {
	validStream1 := Stream{
		Device:  "fakeDevice",
		Address: "fakeAddress",
		Port:    1337,
	}

	validStream2 := Stream{
		Device:  "fakeDevice",
		Address: "differentFakeAddress",
		Port:    1337,
	}

	invalidStreamNoPort := Stream{
		Device:  "invalidDevice",
		Address: "fakeAddress",
		Port:    0,
	}

	invalidStreamNoPortModified := Stream{
		Device:  "updatedDevice",
		Address: "fakeAddress",
		Port:    1337,
	}

	testCases := []struct {
		streams   []Stream
		newStream Stream

		expectedStreams []Stream
	}{
		// Valid baseline
		{
			streams:   []Stream{validStream1, validStream2, invalidStreamNoPort},
			newStream: invalidStreamNoPortModified,

			expectedStreams: []Stream{validStream1, validStream2, invalidStreamNoPortModified},
		},
	}
	for _, test := range testCases {
		streams := replace(test.streams, test.newStream)

		for _, stream := range test.streams {
			foundStream := false
			for _, result := range streams {
				if result.Address == stream.Address && result.Device == stream.Device && result.Port == stream.Port {
					foundStream = true
				}
			}
			assert.Equal(t, true, foundStream, "wrong streams parsed")
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
		// Valid baseline
		{
			stream: validStream,

			expectedRTSPURL: "rtsp://ullaakut:ba69897483886f0d2b0afb6345b76c0c@1.2.3.4:1337/cameradar.sdp",
		},
	}
	for _, test := range testCases {
		output := GetCameraRTSPURL(test.stream)
		assert.Equal(t, test.expectedRTSPURL, output, "wrong RTSP URL generated")
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
		// Valid baseline
		{
			stream: validStream,

			expectedRTSPURL: "http://1.2.3.4/",
		},
	}
	for _, test := range testCases {
		output := GetCameraAdminPanelURL(test.stream)
		assert.Equal(t, test.expectedRTSPURL, output, "wrong Admin Panel URL generated")
	}
}
