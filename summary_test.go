package cameradar

import (
	"bytes"
	"testing"

	"github.com/Ullaakut/disgo"
	"github.com/stretchr/testify/assert"
)

var (
	unavailable = Stream{}

	available = Stream{
		Available: true,
	}

	deviceFound = Stream{
		Device: "devicename",
	}

	noAuth = Stream{
		AuthenticationType: 0,
	}

	basic = Stream{
		AuthenticationType: 1,
	}

	digest = Stream{
		AuthenticationType: 2,
	}

	credsFound = Stream{
		CredentialsFound: true,
		Username:         "us3r",
		Password:         "p4ss",
	}

	routeFound = Stream{
		RouteFound: true,
		Route:      "r0ute",
	}
)

func TestPrintStreams(t *testing.T) {
	tests := []struct {
		description string

		streams []Stream

		expectedLogs []string
	}{
		{
			description: "displays the proper message when no streams found",

			streams: nil,

			expectedLogs: []string{"No streams were found"},
		},
		{
			description: "displays the admin panel URL when a stream is not accessible",

			streams: []Stream{
				unavailable,
			},

			expectedLogs: []string{"Admin panel URL"},
		},
		{
			description: "displays the device name when it is found",

			streams: []Stream{
				deviceFound,
			},

			expectedLogs: []string{"Device model:"},
		},
		{
			description: "displays authentication type (no auth)",

			streams: []Stream{
				noAuth,
			},

			expectedLogs: []string{"This camera does not require authentication"},
		},
		{
			description: "displays authentication type (basic)",

			streams: []Stream{
				basic,
			},

			expectedLogs: []string{"basic"},
		},
		{
			description: "displays authentication type (digest)",

			streams: []Stream{
				digest,
			},

			expectedLogs: []string{"digest"},
		},
		{
			description: "displays credentials properly",

			streams: []Stream{
				credsFound,
			},

			expectedLogs: []string{
				"Username",
				"us3r",
				"Password",
				"p4ss",
			},
		},
		{
			description: "displays route properly",

			streams: []Stream{
				routeFound,
			},

			expectedLogs: []string{
				"RTSP route",
				"/r0ute",
			},
		},
		{
			description: "displays successes properly (no success)",

			streams: []Stream{
				unavailable,
			},

			expectedLogs: []string{
				"Streams were found but none were accessed",
			},
		},
		{
			description: "displays successes properly (1 success)",

			streams: []Stream{
				available,
			},

			expectedLogs: []string{
				"Successful attack",
				"device was accessed",
			},
		},
		{
			description: "displays successes properly (multiple successes)",

			streams: []Stream{
				available,
				available,
				available,
				available,
			},

			expectedLogs: []string{
				"Successful attack",
				"devices were accessed",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			writer := &bytes.Buffer{}
			scanner := &Scanner{
				term: disgo.NewTerminal(disgo.WithDefaultOutput(writer)),
			}

			scanner.PrintStreams(test.streams)

			for _, expectedLog := range test.expectedLogs {
				assert.Contains(t, writer.String(), expectedLog)
			}
		})
	}
}
