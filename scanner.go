package cameradar

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/ullaakut/disgo"
	curl "github.com/ullaakut/go-curl"
)

// Scanner represents a cameradar scanner. It scans a network and
// attacks all streams found to get their RTSP credentials.
type Scanner struct {
	curl Curler
	term *disgo.Terminal

	targets                  []string
	ports                    []string
	debug                    bool
	speed                    int
	timeout                  time.Duration
	credentialDictionaryPath string
	routeDictionaryPath      string

	credentials Credentials
	routes      Routes
}

// New creates a new Cameradar Scanner and applies the given options.
func New(options ...func(*Scanner) error) (*Scanner, error) {
	err := curl.GlobalInit(curl.GLOBAL_ALL)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize curl library: %v", err)
	}

	handle := curl.EasyInit()
	if handle == nil {
		return nil, fmt.Errorf("unable to initialize curl handle: %v", err)
	}

	scanner := &Scanner{
		curl:                     &Curl{CURL: handle},
		credentialDictionaryPath: "<GOPATH>/src/github.com/ullaakut/cameradar/dictionaries/credentials.json",
		routeDictionaryPath:      "<GOPATH>/src/github.com/ullaakut/cameradar/dictionaries/route",
	}

	for _, option := range options {
		err := option(scanner)
		if err != nil {
			return nil, fmt.Errorf("unable to apply option to scanner: %v", err)
		}
	}

	gopath := os.Getenv("GOPATH")
	scanner.credentialDictionaryPath = strings.Replace(scanner.credentialDictionaryPath, "<GOPATH>", gopath, 1)
	scanner.routeDictionaryPath = strings.Replace(scanner.routeDictionaryPath, "<GOPATH>", gopath, 1)

	scanner.term = disgo.NewTerminal(
		disgo.WithDebug(scanner.debug),
		disgo.WithColors(!scanner.debug),
	)

	err = scanner.LoadTargets()
	if err != nil {
		return nil, fmt.Errorf("unable to parse target file: %v", err)
	}

	scanner.term.StartStepf("Loading credentials")
	err = scanner.LoadCredentials()
	if err != nil {
		return nil, scanner.term.FailStepf("unable to load credentials dictionary: %v", err)
	}

	scanner.term.StartStepf("Loading routes")
	err = scanner.LoadRoutes()
	if err != nil {
		return nil, scanner.term.FailStepf("unable to load credentials dictionary: %v", err)
	}

	disgo.EndStep()

	return scanner, nil
}
