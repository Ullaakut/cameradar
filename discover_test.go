package cmrdr

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

// HACK: See https://golang.org/src/os/exec/exec_test.go
func fakeExecCommand(command string, args ...string) *exec.Cmd {
	cs := []string{"-test.run=TestExecCommandHelper", "--", command}
	cs = append(cs, args...)
	cmd := exec.Command(os.Args[0], cs...)
	cmd.Env = []string{"GO_WANT_HELPER_PROCESS=1",
		"STDOUT= ",
		"EXIT_STATUS=0"}
	return cmd
}

func TestExecCommandHelper(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}

	fmt.Fprintf(os.Stdout, os.Getenv("STDOUT"))
	i, _ := strconv.Atoi(os.Getenv("EXIT_STATUS"))
	os.Exit(i)
}

func TestNmapRun(t *testing.T) {
	execCommand = fakeExecCommand
	defer func() { execCommand = exec.Command }()

	testCases := []struct {
		targets        string
		ports          string
		resultFilePath string
		nmapSpeed      int
		enableLogs     bool

		expectedErrMsg string
	}{
		// Valid baseline with logs enabled
		{
			targets:        "localhost",
			ports:          "554",
			resultFilePath: "/tmp/results.xml",
			nmapSpeed:      PARANOIAC,
			enableLogs:     true,
		},
		// Invalid speed
		{
			targets:        "localhost",
			ports:          "554",
			resultFilePath: "/tmp/results.xml",
			nmapSpeed:      INSANE + 1,
			enableLogs:     false,

			expectedErrMsg: "invalid nmap speed value",
		},
	}
	for _, test := range testCases {
		err := NmapRun(test.targets, test.ports, test.resultFilePath, test.nmapSpeed, test.enableLogs)
		if len(test.expectedErrMsg) > 0 {
			if err == nil {
				fmt.Printf("unexpected success. expected error: %s\n", test.expectedErrMsg)
				os.Exit(1)
			}
			assert.Contains(t, err.Error(), test.expectedErrMsg, "wrong error message")
		} else {
			if err != nil {
				fmt.Printf("unexpected error: %v\n", err)
				os.Exit(1)
			}
		}
	}
}

func TestNmapParseResults(t *testing.T) {
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
		fileExists bool
		streamsXML *nmapResult

		expectedStreams []Stream
		expectedErrMsg  string
	}{
		// File exists
		// Two valid streams, no error
		{
			expectedStreams: []Stream{validStream1, validStream2},
			streamsXML: &nmapResult{
				Hosts: []host{
					{
						Address: address{
							Addr:     validStream1.Address,
							AddrType: "ipv4",
						},
						Ports: ports{
							Ports: []port{
								{
									PortID: validStream1.Port,
									State: state{
										State: "open",
									},
									Service: service{
										Name:    "rtsp",
										Product: validStream1.Device,
									},
								},
							},
						},
					},
					{
						Address: address{
							Addr:     validStream2.Address,
							AddrType: "ipv4",
						},
						Ports: ports{
							Ports: []port{
								{
									PortID: validStream2.Port,
									State: state{
										State: "open",
									},
									Service: service{
										Name:    "rtsp",
										Product: validStream2.Device,
									},
								},
							},
						},
					},
				},
			},
			fileExists: true,
		},
		// File exists
		// Two invalid targets, no error
		{
			fileExists:      true,
			expectedStreams: []Stream{invalidStreamNoPort, invalidStreamNoAddress},
			streamsXML: &nmapResult{
				Hosts: []host{
					{
						Address: address{
							Addr:     invalidStreamNoAddress.Address,
							AddrType: "ipv4",
						},
						Ports: ports{
							Ports: []port{
								{
									PortID: invalidStreamNoAddress.Port,
									State: state{
										State: "open",
									},
									Service: service{
										Name:    "rtsp",
										Product: invalidStreamNoAddress.Device,
									},
								},
							},
						},
					},
					{
						Address: address{
							Addr:     invalidStreamNoPort.Address,
							AddrType: "ipv4",
						},
						Ports: ports{
							Ports: []port{
								{
									PortID: invalidStreamNoPort.Port,
									State: state{
										State: "open",
									},
									Service: service{
										Name:    "rtsp",
										Product: invalidStreamNoPort.Device,
									},
								},
							},
						},
					},
				},
			},
		},
		// File does not exist, error
		{
			fileExists:     false,
			expectedErrMsg: "could not read nmap result file",
		},
		// No valid streams found
		{
			fileExists:      true,
			expectedStreams: []Stream{},
			streamsXML: &nmapResult{
				Hosts: []host{
					{
						Address: address{
							Addr:     "Camera with closed ports",
							AddrType: "ipv4",
						},
						Ports: ports{
							Ports: []port{
								{
									PortID: 0,
									State: state{
										State: "closed",
									},
									Service: service{
										Name:    "rtsp",
										Product: "Camera without closed ports",
									},
								},
							},
						},
					},
					{
						Address: address{
							Addr:     "Camera with closed ports",
							AddrType: "ipv4",
						},
					},
				},
			},
		},
		// XML Unmarshal error
		{
			fileExists:      true,
			expectedStreams: []Stream{},
			expectedErrMsg:  "expected element type <nmaprun> but have <failure>",
		},
	}
	for i, test := range testCases {
		filePath := "/tmp/cameradar_test_parse_results_" + fmt.Sprint(i) + ".xml"

		// create file
		if test.fileExists {
			_, err := os.Create(filePath)
			if err != nil {
				fmt.Printf("could not create xml file for NmapParseResults: %v. iteration: %d. file path: %s\n", err, i, filePath)
				os.Exit(1)
			}

			// marshal and write
			if test.streamsXML != nil {
				streams, err := xml.Marshal(test.streamsXML)
				if err != nil {
					fmt.Printf("invalid targets for NmapParseResults: %v. iteration: %d. streams: %v\n", err, i, test.streamsXML)
					os.Exit(1)
				}

				err = ioutil.WriteFile(filePath, streams, 0644)
				if err != nil {
					fmt.Printf("could not write xml file for NmapParseResults: %v. iteration: %d. file path: %s\n", err, i, filePath)
					os.Exit(1)
				}
			} else {
				err := ioutil.WriteFile(filePath, []byte("<?xml version=\"1.0\" encoding=\"UTF-8\"?><failure>"), 0644)
				if err != nil {
					fmt.Printf("could not write xml file for NmapParseResults: %v. iteration: %d. file path: %s\n", err, i, filePath)
					os.Exit(1)
				}
			}
		}

		results, err := NmapParseResults(filePath)
		if len(test.expectedErrMsg) > 0 {
			if err == nil {
				fmt.Printf("unexpected success. expected error: %s\n", test.expectedErrMsg)
				os.Exit(1)
			}
			assert.Contains(t, err.Error(), test.expectedErrMsg, "wrong error message")
		} else {
			if err != nil {
				fmt.Printf("unexpected error: %v\n", err)
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
	}
}

func TestDiscover(t *testing.T) {
	execCommand = fakeExecCommand
	defer func() { execCommand = exec.Command }()

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
		targets        string
		ports          string
		resultFilePath string
		nmapSpeed      int
		enableLogs     bool
		fileExists     bool
		streamsXML     *nmapResult

		expectedStreams []Stream
		expectedErrMsg  string
	}{
		// Valid baseline
		{
			expectedStreams: []Stream{validStream1, validStream2},
			streamsXML: &nmapResult{
				Hosts: []host{
					{
						Address: address{
							Addr:     validStream1.Address,
							AddrType: "ipv4",
						},
						Ports: ports{
							Ports: []port{
								{
									PortID: validStream1.Port,
									State: state{
										State: "open",
									},
									Service: service{
										Name:    "rtsp",
										Product: validStream1.Device,
									},
								},
							},
						},
					},
					{
						Address: address{
							Addr:     validStream2.Address,
							AddrType: "ipv4",
						},
						Ports: ports{
							Ports: []port{
								{
									PortID: validStream2.Port,
									State: state{
										State: "open",
									},
									Service: service{
										Name:    "rtsp",
										Product: validStream2.Device,
									},
								},
							},
						},
					},
				},
			},
			fileExists:     true,
			targets:        "localhost",
			ports:          "554",
			resultFilePath: "/tmp/results.xml",
			nmapSpeed:      PARANOIAC,
			enableLogs:     false,
		},
		// Invalid speed
		{
			expectedStreams: []Stream{},
			streamsXML: &nmapResult{
				Hosts: []host{
					{
						Address: address{
							Addr:     validStream1.Address,
							AddrType: "ipv4",
						},
						Ports: ports{
							Ports: []port{
								{
									PortID: validStream1.Port,
									State: state{
										State: "open",
									},
									Service: service{
										Name:    "rtsp",
										Product: validStream1.Device,
									},
								},
							},
						},
					},
					{
						Address: address{
							Addr:     validStream2.Address,
							AddrType: "ipv4",
						},
						Ports: ports{
							Ports: []port{
								{
									PortID: validStream2.Port,
									State: state{
										State: "open",
									},
									Service: service{
										Name:    "rtsp",
										Product: validStream2.Device,
									},
								},
							},
						},
					},
				},
			},
			fileExists:     true,
			targets:        "localhost",
			ports:          "554",
			resultFilePath: "/tmp/results.xml",
			nmapSpeed:      INSANE + 1,
			enableLogs:     false,

			expectedErrMsg: "invalid nmap speed value",
		},
		// File exists
		// Two valid streams, no error
		{
			expectedStreams: []Stream{validStream1, validStream2},
			streamsXML: &nmapResult{
				Hosts: []host{
					{
						Address: address{
							Addr:     validStream1.Address,
							AddrType: "ipv4",
						},
						Ports: ports{
							Ports: []port{
								{
									PortID: validStream1.Port,
									State: state{
										State: "open",
									},
									Service: service{
										Name:    "rtsp",
										Product: validStream1.Device,
									},
								},
							},
						},
					},
					{
						Address: address{
							Addr:     validStream2.Address,
							AddrType: "ipv4",
						},
						Ports: ports{
							Ports: []port{
								{
									PortID: validStream2.Port,
									State: state{
										State: "open",
									},
									Service: service{
										Name:    "rtsp",
										Product: validStream2.Device,
									},
								},
							},
						},
					},
				},
			},
			fileExists:     true,
			targets:        "localhost",
			ports:          "554",
			resultFilePath: "/tmp/results.xml",
			nmapSpeed:      PARANOIAC,
			enableLogs:     false,
		},
		// File exists
		// Two invalid targets, no error
		{
			fileExists:      true,
			expectedStreams: []Stream{invalidStreamNoPort, invalidStreamNoAddress},
			streamsXML: &nmapResult{
				Hosts: []host{
					{
						Address: address{
							Addr:     invalidStreamNoAddress.Address,
							AddrType: "ipv4",
						},
						Ports: ports{
							Ports: []port{
								{
									PortID: invalidStreamNoAddress.Port,
									State: state{
										State: "open",
									},
									Service: service{
										Name:    "rtsp",
										Product: invalidStreamNoAddress.Device,
									},
								},
							},
						},
					},
					{
						Address: address{
							Addr:     invalidStreamNoPort.Address,
							AddrType: "ipv4",
						},
						Ports: ports{
							Ports: []port{
								{
									PortID: invalidStreamNoPort.Port,
									State: state{
										State: "open",
									},
									Service: service{
										Name:    "rtsp",
										Product: invalidStreamNoPort.Device,
									},
								},
							},
						},
					},
				},
			},
			targets:        "localhost",
			ports:          "554",
			resultFilePath: "/tmp/results.xml",
			nmapSpeed:      PARANOIAC,
			enableLogs:     false,
		},
		// File does not exist, error
		{
			fileExists:     false,
			expectedErrMsg: "could not read nmap result file",
			targets:        "localhost",
			ports:          "554",
			resultFilePath: "/tmp/results.xml",
			nmapSpeed:      PARANOIAC,
			enableLogs:     false,
		},
		// No valid streams found
		{
			fileExists:      true,
			expectedStreams: []Stream{},
			streamsXML: &nmapResult{
				Hosts: []host{
					{
						Address: address{
							Addr:     "Camera with closed ports",
							AddrType: "ipv4",
						},
						Ports: ports{
							Ports: []port{
								{
									PortID: 0,
									State: state{
										State: "closed",
									},
									Service: service{
										Name:    "rtsp",
										Product: "Camera without closed ports",
									},
								},
							},
						},
					},
					{
						Address: address{
							Addr:     "Camera with closed ports",
							AddrType: "ipv4",
						},
					},
				},
			},
			targets:        "localhost",
			ports:          "554",
			resultFilePath: "/tmp/results.xml",
			nmapSpeed:      PARANOIAC,
			enableLogs:     false,
		},
		// XML Unmarshal error
		{
			fileExists:      true,
			expectedStreams: []Stream{},
			expectedErrMsg:  "expected element type <nmaprun> but have <failure>",
			targets:         "localhost",
			ports:           "554",
			resultFilePath:  "/tmp/results.xml",
			nmapSpeed:       PARANOIAC,
			enableLogs:      false,
		},
	}
	for i, test := range testCases {
		filePath := "/tmp/cameradar_test_discover_" + fmt.Sprint(i) + ".xml"

		// create file
		if test.fileExists {
			_, err := os.Create(filePath)
			if err != nil {
				fmt.Printf("could not create xml file for Discover: %v. iteration: %d. file path: %s\n", err, i, filePath)
				os.Exit(1)
			}

			// marshal and write
			if test.streamsXML != nil {
				streams, err := xml.Marshal(test.streamsXML)
				if err != nil {
					fmt.Printf("invalid targets for Discover: %v. iteration: %d. streams: %v\n", err, i, test.streamsXML)
					os.Exit(1)
				}

				err = ioutil.WriteFile(filePath, streams, 0644)
				if err != nil {
					fmt.Printf("could not write xml file for Discover: %v. iteration: %d. file path: %s\n", err, i, filePath)
					os.Exit(1)
				}
			} else {
				err := ioutil.WriteFile(filePath, []byte("<?xml version=\"1.0\" encoding=\"UTF-8\"?><failure>"), 0644)
				if err != nil {
					fmt.Printf("could not write xml file for Discover: %v. iteration: %d. file path: %s\n", err, i, filePath)
					os.Exit(1)
				}
			}
		}

		results, err := Discover(test.targets, test.ports, filePath, test.nmapSpeed, test.enableLogs)

		if len(test.expectedErrMsg) > 0 {
			if err == nil {
				fmt.Printf("unexpected success in Discover test, iteration %d. expected error: %s\n", i, test.expectedErrMsg)
				os.Exit(1)
			}
			assert.Contains(t, err.Error(), test.expectedErrMsg, "wrong error message")
		} else {
			if err != nil {
				fmt.Printf("unexpected error in Discover test, iteration %d: %v\n", i, err)
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
	}
}
