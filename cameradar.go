// Package cameradar provides methods to discover and remotely access RTSP streams.
package cameradar

import (
	"context"
	"fmt"
	"net/netip"
	"os"
	"time"

	"github.com/Ullaakut/disgo"
	"github.com/Ullaakut/disgo/style"
	"github.com/Ullaakut/go-curl"
)

// Scanner represents a cameradar scanner. It scans a network and
// attacks all streams found to get their RTSP credentials.
type Scanner struct {
	curl *curl.CURL
	term *disgo.Terminal

	targets                  []netip.Prefix
	ports                    []uint16
	debug                    bool
	verbose                  int
	scanSpeed                int
	attackInterval           time.Duration
	timeout                  time.Duration
	credentialDictionaryPath string
	routeDictionaryPath      string

	credentials Credentials
	routes      Routes
}

type Option func(s *Scanner)

// WithDebug specifies whether to enable debug logs.
func WithDebug(debug bool) func(s *Scanner) {
	return func(s *Scanner) {
		s.debug = debug
	}
}

// WithVerbose specifies whether to enable verbose logs.
func WithVerbose(verbose bool) func(s *Scanner) {
	return func(s *Scanner) {
		if verbose {
			s.verbose = 1
			return
		}

		s.verbose = 0
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

// New creates a new Cameradar Scanner and applies the given options.
func New(targets []netip.Prefix, ports []uint16, opts ...Option) (*Scanner, error) {
	err := curl.GlobalInit(curl.GLOBAL_ALL)
	if err != nil {
		return nil, fmt.Errorf("initializing curl library: %v", err)
	}

	handle := curl.EasyInit()
	if handle == nil {
		return nil, fmt.Errorf("initializing curl handle: %v", err)
	}

	scanner := Scanner{
		targets:                  targets,
		ports:                    ports,
		curl:                     handle,
		credentialDictionaryPath: "https://raw.githubusercontent.com/Ullaakut/cameradar/refs/heads/master/dictionaries/credentials.json",
		routeDictionaryPath:      "https://raw.githubusercontent.com/Ullaakut/cameradar/refs/heads/master/dictionaries/routes",
	}

	for _, opt := range opts {
		opt(&scanner)
	}

	// FIXME: Load default dictionaries without relying on GOPATH. Parse contents/file/url contents instead.

	scanner.credentialDictionaryPath = os.ExpandEnv(scanner.credentialDictionaryPath)
	scanner.routeDictionaryPath = os.ExpandEnv(scanner.routeDictionaryPath)

	scanner.term = disgo.NewTerminal(
		disgo.WithDebug(scanner.debug),
	)

	err = scanner.LoadTargets()
	if err != nil {
		return nil, fmt.Errorf("fetching targets: %v", err)
	}

	scanner.term.StartStepf("Loading credentials")
	err = scanner.LoadCredentials()
	if err != nil {
		return nil, scanner.term.FailStepf("loading credentials: %v", err)
	}

	scanner.term.StartStepf("Loading routes")
	err = scanner.LoadRoutes()
	if err != nil {
		return nil, scanner.term.FailStepf("loading routes: %v", err)
	}

	disgo.EndStep()
	return &scanner, nil
}

func (s *Scanner) Run(ctx context.Context) error {
	discovered, err := s.Scan(ctx)
	if err != nil {
		return fmt.Errorf("discovering devices: %w", err)
	}

	streams, err := s.Attack(ctx, discovered)
	if err != nil {
		return fmt.Errorf("attacking devices: %w", err)
	}

	s.PrintStreams(streams)
	return nil
}
