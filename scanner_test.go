package cameradar

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	curl "github.com/ullaakut/go-curl"
)

func TestNew(t *testing.T) {
	tests := []struct {
		description string

		targets           []string
		ports             []string
		debug             bool
		verbose           bool
		customCredentials string
		customRoutes      string
		speed             int
		timeout           time.Duration

		loadTargetsFail bool
		loadCredsFail   bool
		loadRoutesFail  bool

		curlGlobalFail bool
		curlEasyFail   bool

		expectedErr bool
	}{
		{
			description: "no error while loading dictionaries",

			targets: []string{"titi", "toto"},
			ports:   []string{"554"},
			debug:   true,
			verbose: false,
			speed:   3,
			timeout: time.Millisecond,
		},
		{
			description: "unable to load targets",

			loadTargetsFail: true,

			expectedErr: true,
		},
		{
			description: "unable to load credentials",

			loadCredsFail: true,

			expectedErr: true,
		},
		{
			description: "unable to load routes",

			loadRoutesFail: true,

			expectedErr: true,
		},
		{
			description: "curl fails to init",

			curlGlobalFail: true,

			expectedErr: true,
		},
		{
			description: "curl fails to create handle",

			curlEasyFail: true,

			expectedErr: true,
		},
		{
			description: "gopath not set and default dicts",

			customCredentials: defaultCredentialDictionaryPath,
			customRoutes:      defaultRouteDictionaryPath,

			expectedErr: true,
		},
	}

	// Temporarily empty the gopath for testing purposes.
	defer os.Setenv("GOPATH", os.Getenv("GOPATH"))

	for i, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			os.Setenv("GOPATH", "")

			if test.loadTargetsFail {
				test.targets = []string{generateTmpFileName(i, "targets")}
				ioutil.WriteFile(test.targets[0], []byte(`0.0.0.0`), 0000)
			}

			if !test.loadCredsFail && test.customCredentials == "" {
				test.customCredentials = generateTmpFileName(i, "creds")
				ioutil.WriteFile(test.customCredentials, []byte(`{"usernames":["admin"],"passwords":["admin"]}`), 0644)
			}

			if !test.loadRoutesFail && test.customRoutes == "" {
				test.customRoutes = generateTmpFileName(i, "routes")
				ioutil.WriteFile(test.customRoutes, []byte(`live.sdp`), 0644)
			}

			curl.TestGlobalFail = test.curlGlobalFail
			curl.TestEasyFail = test.curlEasyFail

			scanner, err := New(
				WithTargets(test.targets),
				WithPorts(test.ports),
				WithDebug(test.debug),
				WithVerbose(test.verbose),
				WithSpeed(test.speed),
				WithTimeout(test.timeout),
				WithCustomCredentials(test.customCredentials),
				WithCustomRoutes(test.customRoutes),
			)

			if test.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			if scanner != nil {
				assert.Equal(t, test.targets, scanner.targets)
				assert.Equal(t, test.ports, scanner.ports)
				assert.Equal(t, test.debug, scanner.debug)
				assert.Equal(t, test.verbose, scanner.verbose)
				assert.Equal(t, test.speed, scanner.speed)
				assert.Equal(t, test.timeout, scanner.timeout)
			}
		})
	}
}

func generateTmpFileName(iteration int, purpose string) string {
	return fmt.Sprintf("/tmp/cameradar_test_scanner_%s_%d_%d", purpose, time.Now().Unix(), iteration)
}
