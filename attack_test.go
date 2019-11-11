package cameradar

import (
	"errors"
	"io/ioutil"
	"testing"
	"time"

	"github.com/Ullaakut/disgo"
	curl "github.com/Ullaakut/go-curl"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type CurlerMock struct {
	mock.Mock
}

func (m *CurlerMock) Setopt(opt int, param interface{}) error {
	args := m.Called(opt, param)
	return args.Error(0)
}

func (m *CurlerMock) Perform() error {
	args := m.Called()
	return args.Error(0)
}

func (m *CurlerMock) Getinfo(info curl.CurlInfo) (interface{}, error) {
	args := m.Called(info)
	return args.Int(0), args.Error(1)
}

func (m *CurlerMock) Duphandle() Curler {
	return m
}

func TestAttack(t *testing.T) {
	var (
		stream1 = Stream{
			Device:  "fakeDevice",
			Address: "fakeAddress",
			Port:    1337,
		}

		stream2 = Stream{
			Device:  "fakeDevice",
			Address: "differentFakeAddress",
			Port:    1337,
		}

		fakeTargets     = []Stream{stream1, stream2}
		fakeRoutes      = Routes{"live.sdp", "media.amp"}
		fakeCredentials = Credentials{
			Usernames: []string{"admin", "root"},
			Passwords: []string{"12345", "root"},
		}
	)

	tests := []struct {
		description string

		targets []Stream

		performErr error

		expectedStreams []Stream
		expectedErr     error
	}{
		{
			description: "inverted RTSP RFC",

			targets: fakeTargets,

			performErr: errors.New("dummy error"),

			expectedStreams: fakeTargets,
		},
		{
			description: "attack works",

			targets: fakeTargets,

			expectedStreams: fakeTargets,
		},
		{
			description: "no targets",

			targets: nil,

			expectedStreams: nil,
			expectedErr:     errors.New("unable to attack empty list of targets"),
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			curlerMock := &CurlerMock{}

			if len(test.targets) != 0 {
				curlerMock.On("Setopt", mock.Anything, mock.Anything).Return(nil)
				curlerMock.On("Perform").Return(test.performErr)
				if test.performErr == nil {
					curlerMock.On("Getinfo", mock.Anything).Return(200, nil)
				}
			}

			scanner := &Scanner{
				term:        disgo.NewTerminal(disgo.WithDefaultOutput(ioutil.Discard)),
				curl:        curlerMock,
				timeout:     time.Millisecond,
				verbose:     false,
				credentials: fakeCredentials,
				routes:      fakeRoutes,
			}

			results, err := scanner.Attack(test.targets)

			assert.Equal(t, test.expectedErr, err)

			assert.Len(t, results, len(test.expectedStreams))

			curlerMock.AssertExpectations(t)
		})
	}
}

func TestAttackCredentials(t *testing.T) {
	var (
		stream1 = Stream{
			Device:    "fakeDevice",
			Address:   "fakeAddress",
			Port:      1337,
			Available: true,
		}

		stream2 = Stream{
			Device:    "fakeDevice",
			Address:   "differentFakeAddress",
			Port:      1337,
			Available: true,
		}

		fakeTargets     = []Stream{stream1, stream2}
		fakeCredentials = Credentials{
			Usernames: []string{"admin", "root"},
			Passwords: []string{"12345", "root"},
		}
	)

	tests := []struct {
		description string

		targets     []Stream
		credentials Credentials
		timeout     time.Duration
		verbose     bool

		status int

		performErr     error
		getInfoErr     error
		invalidTargets bool

		expectedStreams []Stream
	}{
		{
			description: "Credentials found",

			targets:     fakeTargets,
			credentials: fakeCredentials,
			timeout:     1 * time.Millisecond,

			status: 404,

			expectedStreams: fakeTargets,
		},
		{
			description: "Camera accessed",

			targets:     fakeTargets,
			credentials: fakeCredentials,
			timeout:     1 * time.Millisecond,

			status: 200,

			expectedStreams: fakeTargets,
		},
		{
			description: "curl perform fails",

			targets:     fakeTargets,
			credentials: fakeCredentials,
			timeout:     1 * time.Millisecond,

			performErr: errors.New("dummy error"),

			expectedStreams: fakeTargets,
		},
		{
			description: "curl getinfo fails",

			targets:     fakeTargets,
			credentials: fakeCredentials,
			timeout:     1 * time.Millisecond,

			getInfoErr: errors.New("dummy error"),

			expectedStreams: fakeTargets,
		},
		{
			description: "Verbose mode disabled",

			targets:     fakeTargets,
			credentials: fakeCredentials,
			timeout:     1 * time.Millisecond,
			verbose:     false,

			status: 403,

			expectedStreams: fakeTargets,
		},
		{
			description: "Verbose mode enabled",

			targets:     fakeTargets,
			credentials: fakeCredentials,
			timeout:     1 * time.Millisecond,
			verbose:     true,

			status: 403,

			expectedStreams: fakeTargets,
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			curlerMock := &CurlerMock{}

			if !test.invalidTargets {
				curlerMock.On("Setopt", mock.Anything, mock.Anything).Return(nil)
				curlerMock.On("Perform").Return(test.performErr)
				if test.performErr == nil {
					curlerMock.On("Getinfo", mock.Anything).Return(test.status, test.getInfoErr)
				}
			}

			scanner := &Scanner{
				term:        disgo.NewTerminal(disgo.WithDefaultOutput(ioutil.Discard)),
				curl:        curlerMock,
				timeout:     test.timeout,
				verbose:     test.verbose,
				credentials: test.credentials,
			}

			results := scanner.AttackCredentials(test.targets)

			assert.Len(t, results, len(test.expectedStreams))

			curlerMock.AssertExpectations(t)
		})
	}
}

func TestAttackRoute(t *testing.T) {
	var (
		stream1 = Stream{
			Device:    "fakeDevice",
			Address:   "fakeAddress",
			Port:      1337,
			Available: true,
		}

		stream2 = Stream{
			Device:    "fakeDevice",
			Address:   "differentFakeAddress",
			Port:      1337,
			Available: true,
		}

		fakeTargets = []Stream{stream1, stream2}
		fakeRoutes  = Routes{"live.sdp", "media.amp"}
	)

	tests := []struct {
		description string

		targets []Stream
		routes  Routes
		timeout time.Duration
		verbose bool

		status int

		performErr     error
		getInfoErr     error
		invalidTargets bool

		expectedStreams []Stream
		expectedErr     error
	}{
		{
			description: "Route found",

			targets: fakeTargets,
			routes:  fakeRoutes,
			timeout: 1 * time.Millisecond,

			status: 403,

			expectedStreams: fakeTargets,
		},
		{
			description: "Route found",

			targets: fakeTargets,
			routes:  fakeRoutes,
			timeout: 1 * time.Millisecond,

			status: 401,

			expectedStreams: fakeTargets,
		},
		{
			description: "Camera accessed",

			targets: fakeTargets,
			routes:  fakeRoutes,
			timeout: 1 * time.Millisecond,

			status: 200,

			expectedStreams: fakeTargets,
		},
		{
			description: "curl perform fails",

			targets: fakeTargets,
			routes:  fakeRoutes,
			timeout: 1 * time.Millisecond,

			performErr: errors.New("dummy error"),

			expectedStreams: fakeTargets,
		},
		{
			description: "curl getinfo fails",

			targets: fakeTargets,
			routes:  fakeRoutes,
			timeout: 1 * time.Millisecond,

			getInfoErr: errors.New("dummy error"),

			expectedStreams: fakeTargets,
		},
		{
			description: "verbose mode disabled",

			targets: fakeTargets,
			routes:  fakeRoutes,
			timeout: 1 * time.Millisecond,
			verbose: false,

			expectedStreams: fakeTargets,
		},
		{
			description: "verbose mode enabled",

			targets: fakeTargets,
			routes:  fakeRoutes,
			timeout: 1 * time.Millisecond,
			verbose: true,

			expectedStreams: fakeTargets,
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			curlerMock := &CurlerMock{}

			if !test.invalidTargets {
				curlerMock.On("Setopt", mock.Anything, mock.Anything).Return(nil)
				curlerMock.On("Perform").Return(test.performErr)
				if test.performErr == nil {
					curlerMock.On("Getinfo", mock.Anything).Return(test.status, test.getInfoErr)
				}
			}

			scanner := &Scanner{
				term:    disgo.NewTerminal(disgo.WithDefaultOutput(ioutil.Discard)),
				curl:    curlerMock,
				timeout: test.timeout,
				verbose: test.verbose,
				routes:  test.routes,
			}

			results := scanner.AttackRoute(test.targets)

			assert.Len(t, results, len(test.expectedStreams))

			curlerMock.AssertExpectations(t)
		})
	}
}

func TestValidateStreams(t *testing.T) {
	var (
		stream1 = Stream{
			Device:    "fakeDevice",
			Address:   "fakeAddress",
			Port:      1337,
			Available: true,
		}

		stream2 = Stream{
			Device:    "fakeDevice",
			Address:   "differentFakeAddress",
			Port:      1337,
			Available: true,
		}

		fakeTargets = []Stream{stream1, stream2}
	)

	tests := []struct {
		description string

		targets []Stream
		timeout time.Duration
		verbose bool

		status int

		performErr error
		getInfoErr error

		expectedStreams []Stream
	}{
		{
			description: "route found",

			targets: fakeTargets,
			timeout: 1 * time.Millisecond,

			status: 403,

			expectedStreams: fakeTargets,
		},
		{
			description: "route found",

			targets: fakeTargets,
			timeout: 1 * time.Millisecond,

			status: 401,

			expectedStreams: fakeTargets,
		},
		{
			description: "camera accessed",

			targets: fakeTargets,
			timeout: 1 * time.Millisecond,

			status: 200,

			expectedStreams: fakeTargets,
		},
		{
			description: "unavailable stream",

			targets: fakeTargets,
			timeout: 1 * time.Millisecond,

			status: 400,

			expectedStreams: fakeTargets,
		},
		{
			description: "curl perform fails",

			targets: fakeTargets,
			timeout: 1 * time.Millisecond,

			performErr: errors.New("dummy error"),

			expectedStreams: fakeTargets,
		},
		{
			description: "curl getinfo fails",

			targets: fakeTargets,
			timeout: 1 * time.Millisecond,

			getInfoErr: errors.New("dummy error"),

			expectedStreams: fakeTargets,
		},
		{
			description: "verbose disabled",

			targets: fakeTargets,
			timeout: 1 * time.Millisecond,
			verbose: false,

			expectedStreams: fakeTargets,
		},
		{
			description: "verbose enabled",

			targets: fakeTargets,
			timeout: 1 * time.Millisecond,
			verbose: true,

			expectedStreams: fakeTargets,
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			curlerMock := &CurlerMock{}

			curlerMock.On("Setopt", mock.Anything, mock.Anything).Return(nil)
			curlerMock.On("Perform").Return(test.performErr)
			if test.performErr == nil {
				curlerMock.On("Getinfo", mock.Anything).Return(test.status, test.getInfoErr)
			}

			scanner := &Scanner{
				term:    disgo.NewTerminal(disgo.WithDefaultOutput(ioutil.Discard)),
				curl:    curlerMock,
				timeout: test.timeout,
				verbose: test.verbose,
			}

			results := scanner.ValidateStreams(test.targets)

			assert.Equal(t, len(test.expectedStreams), len(results))

			for _, expectedStream := range test.expectedStreams {
				assert.Contains(t, results, expectedStream)
			}

			curlerMock.AssertExpectations(t)
		})
	}
}

func TestDetectAuthenticationType(t *testing.T) {
	var (
		stream1 = Stream{
			Device:    "fakeDevice",
			Address:   "fakeAddress",
			Port:      1337,
			Available: true,
		}

		stream2 = Stream{
			Device:    "fakeDevice",
			Address:   "differentFakeAddress",
			Port:      1337,
			Available: true,
		}

		fakeTargets = []Stream{stream1, stream2}
	)

	tests := []struct {
		description string

		targets []Stream
		timeout time.Duration
		verbose bool

		status int

		performErr error
		getInfoErr error

		expectedStreams []Stream
	}{
		{
			description: "no auth enabled",

			targets: fakeTargets,
			timeout: 1 * time.Millisecond,

			status: 0,

			expectedStreams: fakeTargets,
		},
		{
			description: "basic auth enabled",

			targets: fakeTargets,
			timeout: 1 * time.Millisecond,

			status: 1,

			expectedStreams: fakeTargets,
		},
		{
			description: "digest auth enabled",

			targets: fakeTargets,
			timeout: 1 * time.Millisecond,

			status: 2,

			expectedStreams: fakeTargets,
		},
		{
			description: "curl getinfo fails",

			targets: fakeTargets,
			timeout: 1 * time.Millisecond,

			getInfoErr: errors.New("dummy error"),

			expectedStreams: fakeTargets,
		},
		{
			description: "curl perform fails",

			targets: fakeTargets,
			timeout: 1 * time.Millisecond,

			performErr: errors.New("dummy error"),

			expectedStreams: fakeTargets,
		},
		{
			description: "verbose disabled",

			targets: fakeTargets,
			timeout: 1 * time.Millisecond,
			verbose: false,

			expectedStreams: fakeTargets,
		},
		{
			description: "verbose enabled",

			targets: fakeTargets,
			timeout: 1 * time.Millisecond,
			verbose: true,

			expectedStreams: fakeTargets,
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			curlerMock := &CurlerMock{}

			curlerMock.On("Setopt", mock.Anything, mock.Anything).Return(nil)
			curlerMock.On("Perform").Return(test.performErr)
			if test.performErr == nil {
				curlerMock.On("Getinfo", mock.Anything).Return(test.status, test.getInfoErr)
			}

			scanner := &Scanner{
				term:    disgo.NewTerminal(disgo.WithDefaultOutput(ioutil.Discard)),
				curl:    curlerMock,
				timeout: test.timeout,
				verbose: test.verbose,
			}

			results := scanner.DetectAuthMethods(test.targets)

			assert.Equal(t, len(test.expectedStreams), len(results))

			for _, expectedStream := range test.expectedStreams {
				assert.Contains(t, results, expectedStream)
			}

			curlerMock.AssertExpectations(t)
		})
	}
}

func TestDoNotWrite(t *testing.T) {
	assert.Equal(t, true, doNotWrite(nil, nil))
}
