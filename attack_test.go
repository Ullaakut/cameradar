package cmrdr

import (
	"errors"
	"fmt"
	"os"
	"testing"
	"time"

	curl "github.com/andelf/go-curl"
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

func TestAttackCredentials(t *testing.T) {
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

	invalidStream := Stream{
		Device: "InvalidDevice",
	}

	fakeTargets := []Stream{validStream1, validStream2}
	invalidTargets := []Stream{invalidStream}
	fakeCredentials := Credentials{
		Usernames: []string{"admin", "root"},
		Passwords: []string{"12345", "root"},
	}

	testCases := []struct {
		targets     []Stream
		credentials Credentials
		timeout     time.Duration
		log         bool

		status int

		performErr     error
		getInfoErr     error
		invalidTargets bool

		expectedStreams []Stream
		expectedErrMsg  string
	}{
		// Credentials found
		{
			targets:     fakeTargets,
			credentials: fakeCredentials,
			timeout:     1 * time.Millisecond,

			status: 404,

			expectedStreams: fakeTargets,
		},
		// Camera accessed
		{
			targets:     fakeTargets,
			credentials: fakeCredentials,
			timeout:     1 * time.Millisecond,

			status: 200,

			expectedStreams: fakeTargets,
		},
		// Invalid targets
		{
			targets:     invalidTargets,
			credentials: fakeCredentials,
			timeout:     1 * time.Millisecond,

			invalidTargets: true,

			expectedErrMsg:  "invalid targets",
			expectedStreams: invalidTargets,
		},
		// curl perform fails
		{
			targets:     fakeTargets,
			credentials: fakeCredentials,
			timeout:     1 * time.Millisecond,

			performErr: errors.New("dummy error"),

			expectedStreams: fakeTargets,
		},
		// curl getinfo fails
		{
			targets:     fakeTargets,
			credentials: fakeCredentials,
			timeout:     1 * time.Millisecond,

			getInfoErr: errors.New("dummy error"),

			expectedStreams: fakeTargets,
		},
		// Logging disabled
		{
			targets:     fakeTargets,
			credentials: fakeCredentials,
			timeout:     1 * time.Millisecond,
			log:         false,

			status: 403,

			expectedStreams: fakeTargets,
		},
		// Logging enabled
		{
			targets:     fakeTargets,
			credentials: fakeCredentials,
			timeout:     1 * time.Millisecond,
			log:         true,

			status: 403,

			expectedStreams: fakeTargets,
		},
	}
	for i, test := range testCases {
		curlerMock := &CurlerMock{}

		if !test.invalidTargets {
			curlerMock.On("Setopt", mock.Anything, mock.Anything).Return(nil)
			curlerMock.On("Perform").Return(test.performErr)
			if test.performErr == nil {
				curlerMock.On("Getinfo", mock.Anything).Return(test.status, test.getInfoErr)
			}
		}

		results, err := AttackCredentials(curlerMock, test.targets, test.credentials, test.timeout, test.log)

		if len(test.expectedErrMsg) > 0 {
			if err == nil {
				fmt.Printf("unexpected success in AttackCredentials test, iteration %d. expected error: %s\n", i, test.expectedErrMsg)
				os.Exit(1)
			}
			assert.Contains(t, err.Error(), test.expectedErrMsg, "wrong error message")
		} else {
			if err != nil {
				fmt.Printf("unexpected error in AttackCredentials test, iteration %d: %v\n", i, err)
				os.Exit(1)
			}
			for _, stream := range test.expectedStreams {
				foundStream := false
				for _, result := range results {
					if result.Address == stream.Address && result.Device == stream.Device && result.Port == stream.Port {
						foundStream = true
					}
				}
				assert.Equal(t, true, foundStream, "wrong streams parsed")
			}
		}
		assert.Equal(t, len(test.expectedStreams), len(results), "wrong streams parsed")
		curlerMock.AssertExpectations(t)
	}
}

func TestAttackRoute(t *testing.T) {
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

	invalidStream := Stream{
		Device: "InvalidDevice",
	}

	fakeTargets := []Stream{validStream1, validStream2}
	fakeRoutes := Routes{"live.sdp", "media.amp"}
	invalidTargets := []Stream{invalidStream}

	testCases := []struct {
		targets []Stream
		routes  Routes
		timeout time.Duration
		log     bool

		status int

		performErr     error
		getInfoErr     error
		invalidTargets bool

		expectedStreams []Stream
		expectedErrMsg  string
	}{
		// Route found
		{
			targets: fakeTargets,
			routes:  fakeRoutes,
			timeout: 1 * time.Millisecond,

			status: 403,

			expectedStreams: fakeTargets,
		},
		// Route found
		{
			targets: fakeTargets,
			routes:  fakeRoutes,
			timeout: 1 * time.Millisecond,

			status: 401,

			expectedStreams: fakeTargets,
		},
		// Camera accessed
		{
			targets: fakeTargets,
			routes:  fakeRoutes,
			timeout: 1 * time.Millisecond,

			status: 200,

			expectedStreams: fakeTargets,
		},
		// Invalid targets
		{
			targets:        invalidTargets,
			routes:         fakeRoutes,
			timeout:        1 * time.Millisecond,
			invalidTargets: true,

			expectedErrMsg:  "invalid targets",
			expectedStreams: invalidTargets,
		},
		// curl perform fails
		{
			targets: fakeTargets,
			routes:  fakeRoutes,
			timeout: 1 * time.Millisecond,

			performErr: errors.New("dummy error"),

			expectedStreams: fakeTargets,
		},
		// curl getinfo fails
		{
			targets: fakeTargets,
			routes:  fakeRoutes,
			timeout: 1 * time.Millisecond,

			getInfoErr: errors.New("dummy error"),

			expectedStreams: fakeTargets,
		},
		// Logs disabled
		{
			targets: fakeTargets,
			routes:  fakeRoutes,
			timeout: 1 * time.Millisecond,
			log:     false,

			expectedStreams: fakeTargets,
		},
		// Logs enabled
		{
			targets: fakeTargets,
			routes:  fakeRoutes,
			timeout: 1 * time.Millisecond,
			log:     true,

			expectedStreams: fakeTargets,
		},
	}

	for i, test := range testCases {
		curlerMock := &CurlerMock{}

		if !test.invalidTargets {
			curlerMock.On("Setopt", mock.Anything, mock.Anything).Return(nil)
			curlerMock.On("Perform").Return(test.performErr)
			if test.performErr == nil {
				curlerMock.On("Getinfo", mock.Anything).Return(test.status, test.getInfoErr)
			}
		}

		results, err := AttackRoute(curlerMock, test.targets, test.routes, test.timeout, test.log)

		if len(test.expectedErrMsg) > 0 {
			if err == nil {
				fmt.Printf("unexpected success in AttackRoute test, iteration %d. expected error: %s\n", i, test.expectedErrMsg)
				os.Exit(1)
			}

			assert.Contains(t, err.Error(), test.expectedErrMsg, "wrong error message")
		} else {
			if err != nil {
				fmt.Printf("unexpected error in AttackRoute test, iteration %d: %v\n", i, err)
				os.Exit(1)
			}

			for _, stream := range test.expectedStreams {
				foundStream := false
				for _, result := range results {
					if result.Address == stream.Address && result.Device == stream.Device && result.Port == stream.Port {
						foundStream = true
					}
				}

				assert.Equal(t, true, foundStream, "wrong streams parsed")
			}
		}

		assert.Equal(t, len(test.expectedStreams), len(results), "wrong streams parsed")

		curlerMock.AssertExpectations(t)
	}
}

func TestValidateStreams(t *testing.T) {
	validStream1 := Stream{
		Device:    "fakeDevice",
		Address:   "fakeAddress",
		Port:      1337,
		Available: true,
	}

	validStream2 := Stream{
		Device:    "fakeDevice",
		Address:   "differentFakeAddress",
		Port:      1337,
		Available: true,
	}

	unavailableStream := Stream{
		Device:    "fakeDevice",
		Available: false,
	}

	fakeTargets := []Stream{validStream1, validStream2}
	unavailableTargets := []Stream{unavailableStream}

	testCases := []struct {
		desc string

		targets []Stream
		timeout time.Duration
		log     bool

		status int

		performErr error
		getInfoErr error

		expectedStreams []Stream
		expectedErrMsg  string
	}{
		// Route found
		{
			desc: "route found",

			targets: fakeTargets,
			timeout: 1 * time.Millisecond,

			status: 403,

			expectedStreams: fakeTargets,
		},
		// Route found
		{
			desc: "route found",

			targets: fakeTargets,
			timeout: 1 * time.Millisecond,

			status: 401,

			expectedStreams: fakeTargets,
		},
		// Camera accessed
		{
			desc: "camera accessed",

			targets: fakeTargets,
			timeout: 1 * time.Millisecond,

			status: 200,

			expectedStreams: fakeTargets,
		},
		// Unavailable stream
		{
			desc: "unavailable stream",

			targets: unavailableTargets,
			timeout: 1 * time.Millisecond,

			status: 400,

			expectedStreams: unavailableTargets,
		},
		// curl perform fails
		{
			desc: "curl perform fails",

			targets: fakeTargets,
			timeout: 1 * time.Millisecond,

			performErr: errors.New("dummy error"),

			expectedStreams: fakeTargets,
		},
		// curl getinfo fails
		{
			desc: "curl getinfo fails",

			targets: fakeTargets,
			timeout: 1 * time.Millisecond,

			getInfoErr: errors.New("dummy error"),

			expectedStreams: fakeTargets,
		},
		// Logs disabled
		{
			desc: "logs disabled",

			targets: fakeTargets,
			timeout: 1 * time.Millisecond,
			log:     false,

			expectedStreams: fakeTargets,
		},
		// Logs enabled
		{
			desc: "logs enabled",

			targets: fakeTargets,
			timeout: 1 * time.Millisecond,
			log:     true,

			expectedStreams: fakeTargets,
		},
	}
	for i, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			curlerMock := &CurlerMock{}

			curlerMock.On("Setopt", mock.Anything, mock.Anything).Return(nil)
			curlerMock.On("Perform").Return(tC.performErr)
			if tC.performErr == nil {
				curlerMock.On("Getinfo", mock.Anything).Return(tC.status, tC.getInfoErr)
			}

			results, err := ValidateStreams(curlerMock, tC.targets, tC.timeout, tC.log)

			if len(tC.expectedErrMsg) > 0 {
				if err == nil {
					fmt.Printf("unexpected success in ValidateStream test, iteration %d. expected error: %s\n", i, tC.expectedErrMsg)
					os.Exit(1)
				}

				assert.Contains(t, err.Error(), tC.expectedErrMsg, "wrong error message")
			} else {
				if err != nil {
					fmt.Printf("unexpected error in ValidateStream test, iteration %d: %v\n", i, err)
					os.Exit(1)
				}

				for _, stream := range tC.expectedStreams {
					foundStream := false
					for _, result := range results {
						if result.Address == stream.Address && result.Device == stream.Device && result.Port == stream.Port {
							foundStream = true
						}
					}

					assert.Equal(t, true, foundStream, "wrong streams parsed")
				}
			}

			assert.Equal(t, len(tC.expectedStreams), len(results), "wrong streams parsed")

			curlerMock.AssertExpectations(t)
		})
	}
}

func TestDoNotWrite(t *testing.T) {
	assert.Equal(t, true, doNotWrite(nil, nil))
}
