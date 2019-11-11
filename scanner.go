package cameradar

import (
	"fmt"
	"os"
	"time"

	"github.com/Ullaakut/disgo"
	"github.com/Ullaakut/disgo/style"
	curl "github.com/Ullaakut/go-curl"
)

const (
	defaultCredentialDictionaryPath = "${GOPATH}/src/github.com/Ullaakut/cameradar/dictionaries/credentials.json"
	defaultRouteDictionaryPath      = "${GOPATH}/src/github.com/Ullaakut/cameradar/dictionaries/routes"
)

// Scanner represents a cameradar scanner. It scans a network and
// attacks all streams found to get their RTSP credentials.
type Scanner struct {
	curl Curler
	term *disgo.Terminal

	targets                  []string
	ports                    []string
	debug                    bool
	verbose                  bool
	scanSpeed                int
	attackInterval           time.Duration
	timeout                  time.Duration
	credentialDictionaryPath string
	routeDictionaryPath      string

	credentials Credentials
	routes      Routes
}

// New creates a new Cameradar Scanner and applies the given options.
func New(options ...func(*Scanner)) (*Scanner, error) {
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
		credentialDictionaryPath: defaultCredentialDictionaryPath,
		routeDictionaryPath:      defaultRouteDictionaryPath,
	}

	for _, option := range options {
		option(scanner)
	}

	gopath := os.Getenv("GOPATH")
	if gopath == "" && scanner.credentialDictionaryPath == defaultCredentialDictionaryPath && scanner.routeDictionaryPath == defaultRouteDictionaryPath {
		disgo.Errorln(style.Failure("No $GOPATH was found.\nDictionaries may not be loaded properly, please set your $GOPATH to use the default dictionaries."))
	}

	scanner.credentialDictionaryPath = os.ExpandEnv(scanner.credentialDictionaryPath)
	scanner.routeDictionaryPath = os.ExpandEnv(scanner.routeDictionaryPath)

	scanner.term = disgo.NewTerminal(
		disgo.WithDebug(scanner.debug),
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

// WithTargets specifies the targets to scan and attack.
func WithTargets(targets []string) func(s *Scanner) {
	return func(s *Scanner) {
		s.targets = targets
	}
}

// WithPorts specifies the ports to scan and attack.
func WithPorts(ports []string) func(s *Scanner) {
	return func(s *Scanner) {
		s.ports = ports
	}
}

// WithDebug specifies whether or not to enable debug logs.
func WithDebug(debug bool) func(s *Scanner) {
	return func(s *Scanner) {
		s.debug = debug
	}
}

// WithVerbose specifies whether or not to enable verbose logs.
func WithVerbose(verbose bool) func(s *Scanner) {
	return func(s *Scanner) {
		s.verbose = verbose
	}
}

// WithCustomCredentials specifies a custom credential dictionary
// to use for the attacks.
func WithCustomCredentials(dictionaryPath string) func(s *Scanner) {
	return func(s *Scanner) {
		s.credentialDictionaryPath = dictionaryPath
	}
}

// WithCustomRoutes specifies a custom route dictionary
// to use for the attacks.
func WithCustomRoutes(dictionaryPath string) func(s *Scanner) {
	return func(s *Scanner) {
		s.routeDictionaryPath = dictionaryPath
	}
}

// WithScanSpeed specifies the speed at which the scan should be executed. Faster
// means easier to detect, slower has bigger timeout values and is more silent.
func WithScanSpeed(speed int) func(s *Scanner) {
	return func(s *Scanner) {
		s.scanSpeed = speed
	}
}

// WithAttackInterval specifies the interval of time during which Cameradar
// should wait between each attack attempt during bruteforcing.
// Setting a high value for this obviously makes attacks much slower.
func WithAttackInterval(interval time.Duration) func(s *Scanner) {
	return func(s *Scanner) {
		s.attackInterval = interval
	}
}

// WithTimeout specifies the amount of time after which attack requests should
// timeout. This should be high if the network you are attacking has a poor
// connectivity or that you are located far away from it.
func WithTimeout(timeout time.Duration) func(s *Scanner) {
	return func(s *Scanner) {
		s.timeout = timeout
	}
}
