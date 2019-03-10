package cmrdr

import (
	"errors"
	"os"
	"testing"

	"github.com/ullaakut/nmap"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type nmapMock struct {
	mock.Mock
}

func (m *nmapMock) Run() (*nmap.Run, error) {
	args := m.Called()

	if args.Get(0) != nil {
		return args.Get(0).(*nmap.Run), args.Error(1)
	}
	return nil, args.Error(1)
}

func TestDiscover(t *testing.T) {
	tests := []struct {
		description string

		targets    []string
		ports      []string
		speed      int
		removePath bool

		expectedErr    error
		expectedResult []Stream
	}{
		{
			description: "create new scanner and call scan, no error",

			targets: []string{"localhost"},
			ports:   []string{"80"},
			speed:   5,
		},
		{
			description: "create new scanner with missing nmap installation",

			removePath: true,
			ports:      []string{"80"},

			expectedErr: errors.New("'nmap' binary was not found"),
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			if test.removePath {
				os.Setenv("PATH", "")
			}

			result, err := Discover(test.targets, test.ports, test.speed)

			assert.Equal(t, test.expectedErr, err)
			assert.Equal(t, test.expectedResult, result)
		})
	}
}

func TestScan(t *testing.T) {
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

	invalidStreamNoAddress := Stream{
		Device:  "invalidDevice",
		Address: "",
		Port:    1337,
	}

	testCases := []struct {
		description string
		nmapResult  *nmap.Run
		nmapError   error

		expectedStreams []Stream
		expectedErr     error
	}{
		{
			description: "valid streams",

			nmapResult: &nmap.Run{
				Hosts: []nmap.Host{
					{
						Addresses: []nmap.Address{
							{
								Addr: validStream1.Address,
							},
						},
						Ports: []nmap.Port{
							{
								State: nmap.State{
									State: "open",
								},
								ID: validStream1.Port,
								Service: nmap.Service{
									Name:    "rtsp",
									Product: validStream1.Device,
								},
							},
						},
					},
					{
						Addresses: []nmap.Address{
							{
								Addr: validStream2.Address,
							},
						},
						Ports: []nmap.Port{
							{
								State: nmap.State{
									State: "open",
								},
								ID: validStream2.Port,
								Service: nmap.Service{
									Name:    "rtsp-alt",
									Product: validStream2.Device,
								},
							},
						},
					},
				},
			},

			expectedStreams: []Stream{validStream1, validStream2},
		},
		{
			description: "two invalid targets, no error",

			nmapResult: &nmap.Run{
				Hosts: []nmap.Host{
					{
						Addresses: []nmap.Address{
							{
								Addr: invalidStreamNoPort.Address,
							},
						},
					},
					{
						Addresses: []nmap.Address{},
						Ports: []nmap.Port{
							{
								State: nmap.State{
									State: "open",
								},
								ID: validStream2.Port,
								Service: nmap.Service{
									Name:    "rtsp-alt",
									Product: invalidStreamNoAddress.Device,
								},
							},
						},
					},
				},
			},

			expectedStreams: nil,
		},
		{
			description: "different port states, no error",

			nmapResult: &nmap.Run{
				Hosts: []nmap.Host{
					{
						Addresses: []nmap.Address{
							{
								Addr: invalidStreamNoPort.Address,
							}},
						Ports: []nmap.Port{
							{
								State: nmap.State{
									State: "closed",
								},
								ID: validStream2.Port,
								Service: nmap.Service{
									Name:    "rtsp-alt",
									Product: invalidStreamNoAddress.Device,
								},
							},
						},
					},
					{
						Addresses: []nmap.Address{
							{
								Addr: invalidStreamNoPort.Address,
							}},
						Ports: []nmap.Port{
							{
								State: nmap.State{
									State: "unfiltered",
								},
								ID: validStream2.Port,
								Service: nmap.Service{
									Name:    "rtsp-alt",
									Product: invalidStreamNoAddress.Device,
								},
							},
						},
					},
					{
						Addresses: []nmap.Address{
							{
								Addr: invalidStreamNoPort.Address,
							}},
						Ports: []nmap.Port{
							{
								State: nmap.State{
									State: "filtered",
								},
								ID: validStream2.Port,
								Service: nmap.Service{
									Name:    "rtsp-alt",
									Product: invalidStreamNoAddress.Device,
								},
							},
						},
					},
				},
			},

			expectedStreams: nil,
		},
		{
			description: "not rtsp, no error",

			nmapResult: &nmap.Run{
				Hosts: []nmap.Host{
					{
						Addresses: []nmap.Address{
							{
								Addr: invalidStreamNoPort.Address,
							}},
						Ports: []nmap.Port{
							{
								State: nmap.State{
									State: "open",
								},
								ID: validStream2.Port,
								Service: nmap.Service{
									Name:    "tcp",
									Product: invalidStreamNoAddress.Device,
								},
							},
						},
					},
				},
			},

			expectedStreams: nil,
		},
		{
			description: "no hosts found",

			nmapResult:      &nmap.Run{},
			expectedStreams: nil,
		},
		{
			description: "scan failed",

			nmapError:   errors.New("scan failed"),
			expectedErr: errors.New("scan failed"),
		},
	}

	for _, test := range testCases {
		t.Run(test.description, func(t *testing.T) {
			nmapMock := &nmapMock{}

			nmapMock.On("Run").Return(test.nmapResult, test.nmapError)

			results, err := scan(nmapMock)

			assert.Equal(t, test.expectedErr, err)
			assert.Equal(t, test.expectedStreams, results, "wrong streams parsed")
			assert.Equal(t, len(test.expectedStreams), len(results), "wrong streams parsed")

			nmapMock.AssertExpectations(t)
		})
	}
}
