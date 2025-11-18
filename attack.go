package cameradar

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Ullaakut/go-curl"
)

// HTTP responses.
const (
	httpOK                = 200
	httpUnauthorized      = 401
	httpForbidden         = 403
	httpNotFound          = 404
	httpServiceUnavailable = 503
)

// CURL RTSP request types.
const (
	rtspDescribe = 2
	rtspSetup    = 4
)

// Authentication types.
const (
	authNone   = 0
	authBasic  = 1
	authDigest = 2
)

// Route that should never be a constructor default.
const dummyRoute = "/0x8b6c42"

// Attack attacks the given targets and returns the accessed streams.
func (s *Scanner) Attack(targets []Stream) ([]Stream, error) {
	if len(targets) == 0 {
		return nil, fmt.Errorf("no stream found")
	}

	// Most cameras will be accessed successfully with these two attacks.
	s.term.StartStepf("Attacking routes of %d streams", len(targets))
	streams := s.AttackRoute(targets)

	s.term.StartStepf("Attempting to detect authentication methods of %d streams", len(targets))
	streams = s.DetectAuthMethods(streams)

	s.term.StartStepf("Attacking credentials of %d streams", len(targets))
	streams = s.AttackCredentials(streams)

	s.term.StartStep("Validating that streams are accessible")
	streams = s.ValidateStreams(streams)

	// But some cameras run GST RTSP Server which prioritizes 401 over 404 contrary to most cameras.
	// For these cameras, running another route attack will solve the problem.
	for _, stream := range streams {
		if !stream.RouteFound || !stream.CredentialsFound || !stream.Available {
			s.term.StartStepf("Second round of attacks")
			streams = s.AttackRoute(streams)

			s.term.StartStep("Validating that streams are accessible")
			streams = s.ValidateStreams(streams)

			break
		}
	}

	s.term.EndStep()

	return streams, nil
}

// ValidateStreams tries to setup the stream to validate whether or not it is available.
func (s *Scanner) ValidateStreams(targets []Stream) []Stream {
	for i := range targets {
		err := s.validateStream(targets[i])
		targets[i].Available = (err == nil)
		time.Sleep(s.attackInterval)
	}

	return targets
}

// AttackCredentials attempts to guess the provided targets' credentials using the given
// dictionary or the default dictionary if none was provided by the user.
func (s *Scanner) AttackCredentials(targets []Stream) []Stream {
	resChan := make(chan Stream)
	defer close(resChan)

	for i := range targets {
		go s.attackCameraCredentials(targets[i], resChan)
	}

	for range targets {
		attackResult := <-resChan
		if attackResult.CredentialsFound {
			targets = replace(targets, attackResult)
		}
	}

	return targets
}

// AttackRoute attempts to guess the provided targets' streaming routes using the given
// dictionary or the default dictionary if none was provided by the user.
func (s *Scanner) AttackRoute(targets []Stream) []Stream {
	resChan := make(chan Stream)
	defer close(resChan)

	for i := range targets {
		go s.attackCameraRoute(targets[i], resChan)
	}

	for range targets {
		attackResult := <-resChan
		if attackResult.RouteFound {
			targets = replace(targets, attackResult)
		}
	}

	return targets
}

// DetectAuthMethods attempts to guess the provided targets' authentication types, between
// digest, basic auth or none at all.
func (s *Scanner) DetectAuthMethods(targets []Stream) []Stream {
	for i := range targets {
		targets[i].AuthenticationType = s.detectAuthMethod(targets[i])
		time.Sleep(s.attackInterval)

		var authMethod string
		switch targets[i].AuthenticationType {
		case authNone:
			authMethod = "no"
		case authBasic:
			authMethod = "basic"
		case authDigest:
			authMethod = "digest"
		default:
			authMethod = "unknown:" + string(targets[i].AuthenticationType)
		}

		s.term.Debugf("Stream %s uses %s authentication method\n", GetCameraRTSPURL(targets[i]), authMethod)
	}

	return targets
}

func (s *Scanner) attackCameraCredentials(target Stream, resChan chan<- Stream) {
	consecutiveErrors := 0
	maxConsecutiveErrors := 5 // Stop after 5 consecutive connection failures
	
	for _, username := range s.credentials.Usernames {
		for _, password := range s.credentials.Passwords {
			err := s.credAttack(target, username, password)
			if err == nil {
				target.CredentialsFound = true
				target.Username = username
				target.Password = password
				resChan <- target
				return
			}
			
			// Track consecutive connection errors
			if isConnectionError(err) {
				consecutiveErrors++
				if consecutiveErrors >= maxConsecutiveErrors {
					s.term.Errorf("Stream %s: Too many consecutive connection failures (%d), server may be blocking requests", GetCameraRTSPURL(target), consecutiveErrors)
					break
				}
			} else {
				consecutiveErrors = 0 // Reset on non-connection errors
			}
			
			time.Sleep(s.attackInterval)
		}
		if consecutiveErrors >= maxConsecutiveErrors {
			break // Exit outer loop as well
		}
	}

	target.CredentialsFound = false
	resChan <- target
}

func (s *Scanner) attackCameraRoute(target Stream, resChan chan<- Stream) {
	// If the stream responds positively to the dummy route, it means
	// it doesn't require (or respect the RFC) a route and the attack
	// can be skipped.
	err := s.routeAttack(target, dummyRoute)
	if err == nil {
		target.RouteFound = true
		target.Routes = append(target.Routes, "/")
		resChan <- target
		return
	}

	consecutiveErrors := 0
	maxConsecutiveErrors := 5 // Stop after 5 consecutive connection failures
	
	// Otherwise, bruteforce the routes.
	for _, route := range s.routes {
		err := s.routeAttack(target, route)
		if err == nil {
			target.RouteFound = true
			target.Routes = append(target.Routes, route)
			consecutiveErrors = 0 // Reset on success
		} else if isConnectionError(err) {
			// Track consecutive connection errors
			consecutiveErrors++
			if consecutiveErrors >= maxConsecutiveErrors {
				s.term.Errorf("Stream %s: Too many consecutive errors (%d), server may be blocking requests", GetCameraRTSPURL(target), consecutiveErrors)
				break
			}
		} else {
			consecutiveErrors = 0 // Reset on non-connection errors
		}
		
		time.Sleep(s.attackInterval)
	}

	resChan <- target
}

func (s *Scanner) detectAuthMethod(stream Stream) int {
	c := s.curl.Duphandle()

	route := normalizeRoute(stream.Route())
	attackURL := fmt.Sprintf(
		"rtsp://%s:%d/%s",
		stream.Address,
		stream.Port,
		route,
	)

	s.setCurlOptions(c)

	// Send a request to the URL of the stream we want to attack.
	_ = c.Setopt(curl.OPT_URL, attackURL)
	// Set the RTSP STREAM URI as the stream URL.
	_ = c.Setopt(curl.OPT_RTSP_STREAM_URI, attackURL)
	// Add Accept header for proper SDP negotiation with strict RTSP servers
	_ = c.Setopt(curl.OPT_HTTPHEADER, []string{"Accept: application/sdp"})
	_ = c.Setopt(curl.OPT_RTSP_REQUEST, rtspDescribe)

	// Perform the request.
	err := c.Perform()
	if err != nil {
		s.term.Errorf("Perform failed for %q (auth %d): %v", attackURL, stream.AuthenticationType, err)
		return -1
	}

	authType, err := c.Getinfo(curl.INFO_HTTPAUTH_AVAIL)
	if err != nil {
		s.term.Errorf("Getinfo failed: %v", err)
		return -1
	}

	if s.debug {
		s.term.Debugln("DESCRIBE", attackURL, "RTSP/1.0 >", authType)
	}

	return authType.(int)
}

func (s *Scanner) routeAttack(stream Stream, route string) error {
	c := s.curl.Duphandle()

	normalizedRoute := normalizeRoute(route)
	attackURL := fmt.Sprintf(
		"rtsp://%s:%d/%s",
		stream.Address,
		stream.Port,
		normalizedRoute,
	)

	s.setCurlOptions(c)

	// Set proper authentication type.
	_ = c.Setopt(curl.OPT_HTTPAUTH, stream.AuthenticationType)
	_ = c.Setopt(curl.OPT_USERNAME, stream.Username)
	_ = c.Setopt(curl.OPT_PASSWORD, stream.Password)

	// Send a request to the URL of the stream we want to attack.
	_ = c.Setopt(curl.OPT_URL, attackURL)
	// Set the RTSP STREAM URI as the stream URL.
	_ = c.Setopt(curl.OPT_RTSP_STREAM_URI, attackURL)
	// Add Accept header for proper SDP negotiation with strict RTSP servers
	_ = c.Setopt(curl.OPT_HTTPHEADER, []string{"Accept: application/sdp"})
	_ = c.Setopt(curl.OPT_RTSP_REQUEST, rtspDescribe)

	// Perform the request.
	err := c.Perform()
	if err != nil {
		// Check if it's a connection error (reset, timeout, etc.)
		if isConnectionError(err) {
			return fmt.Errorf("%w: %w", ErrConnection, err)
		}
		if !s.verbose {
			s.term.Errorf("Perform failed for %q (auth %d): %v", attackURL, stream.AuthenticationType, err)
		}
		return err
	}

	// Get return code for the request.
	rc, err := c.Getinfo(curl.INFO_RESPONSE_CODE)
	if err != nil {
		s.term.Errorf("Getinfo failed: %v", err)
		return err
	}

	if s.debug {
		s.term.Debugln("DESCRIBE", attackURL, "RTSP/1.0 >", rc)
	}

	// 503 Service Unavailable indicates server is rate-limiting/blocking
	if rc == httpServiceUnavailable {
		return fmt.Errorf("%w: service unavailable (503): server may be rate-limiting", ErrConnection)
	}
	
	// If it's a 401 or 403, it means that the credentials are wrong but the route might be okay.
	// If it's a 200, the stream is accessed successfully.
	if rc == httpOK || rc == httpUnauthorized || rc == httpForbidden {
		return nil
	}
	return fmt.Errorf("unexpected response: %d", rc)
}

func (s *Scanner) credAttack(stream Stream, username string, password string) error {
	c := s.curl.Duphandle()

	route := normalizeRoute(stream.Route())
	attackURL := fmt.Sprintf(
		"rtsp://%s:%d/%s",
		stream.Address,
		stream.Port,
		route,
	)

	s.setCurlOptions(c)

	// Set proper authentication type.
	_ = c.Setopt(curl.OPT_HTTPAUTH, stream.AuthenticationType)
	_ = c.Setopt(curl.OPT_USERNAME, username)
	_ = c.Setopt(curl.OPT_PASSWORD, password)

	// Send a request to the URL of the stream we want to attack.
	_ = c.Setopt(curl.OPT_URL, attackURL)
	// Set the RTSP STREAM URI as the stream URL.
	_ = c.Setopt(curl.OPT_RTSP_STREAM_URI, attackURL)
	// Add Accept header for proper SDP negotiation with strict RTSP servers
	_ = c.Setopt(curl.OPT_HTTPHEADER, []string{"Accept: application/sdp"})
	_ = c.Setopt(curl.OPT_RTSP_REQUEST, rtspDescribe)

	// Perform the request.
	err := c.Perform()
	if err != nil {
		// Check if it's a connection error (reset, timeout, etc.)
		if isConnectionError(err) {
			return fmt.Errorf("connection error: %w", err)
		}
		if !s.verbose {
			s.term.Errorf("Perform failed for %q (auth %d): %v", attackURL, stream.AuthenticationType, err)
		}
		return err
	}

	// Get return code for the request.
	rc, err := c.Getinfo(curl.INFO_RESPONSE_CODE)
	if err != nil {
		s.term.Errorf("Getinfo failed: %v", err)
		return err
	}

	if s.debug {
		s.term.Debugln("DESCRIBE", attackURL, "RTSP/1.0 >", rc)
	}

	// 503 Service Unavailable indicates server is rate-limiting/blocking
	if rc == httpServiceUnavailable {
		return fmt.Errorf("service unavailable (503): server may be rate-limiting")
	}

	// If it's a 404, it means that the route is incorrect but the credentials might be okay.
	// If it's a 200, the stream is accessed successfully.
	if rc == httpOK || rc == httpNotFound {
		return nil
	}
	return fmt.Errorf("unexpected response: %d", rc)
}

func (s *Scanner) validateStream(stream Stream) error {
	c := s.curl.Duphandle()

	route := normalizeRoute(stream.Route())
	attackURL := fmt.Sprintf(
		"rtsp://%s:%d/%s",
		stream.Address,
		stream.Port,
		route,
	)

	s.setCurlOptions(c)

	// Set proper authentication type.
	_ = c.Setopt(curl.OPT_HTTPAUTH, stream.AuthenticationType)
	_ = c.Setopt(curl.OPT_USERNAME, stream.Username)
	_ = c.Setopt(curl.OPT_PASSWORD, stream.Password)

	// Send a request to the URL of the stream we want to attack.
	_ = c.Setopt(curl.OPT_URL, attackURL)
	// Set the RTSP STREAM URI as the stream URL.
	_ = c.Setopt(curl.OPT_RTSP_STREAM_URI, attackURL)
	_ = c.Setopt(curl.OPT_RTSP_REQUEST, rtspSetup)

	_ = c.Setopt(curl.OPT_RTSP_TRANSPORT, "RTP/AVP;unicast;client_port=33332-33333")

	// Perform the request.
	err := c.Perform()
	if err != nil {
		s.term.Errorf("Perform failed for %q (auth %d): %v", attackURL, stream.AuthenticationType, err)
		return err
	}

	// Get return code for the request.
	rc, err := c.Getinfo(curl.INFO_RESPONSE_CODE)
	if err != nil {
		s.term.Errorf("Getinfo failed: %v", err)
		return err
	}

	if s.debug {
		s.term.Debugln("SETUP", attackURL, "RTSP/1.0 >", rc)
	}

	// If it's a 200, the stream is accessed successfully.
	if rc == httpOK {
		return nil
	}
	return fmt.Errorf("unexpected response: %d", rc)
}

func (s *Scanner) setCurlOptions(c Curler) {
	// Do not write sdp in stdout
	_ = c.Setopt(curl.OPT_WRITEFUNCTION, doNotWrite)
	// Do not use signals (would break multithreading).
	_ = c.Setopt(curl.OPT_NOSIGNAL, 1)
	// Do not send a body in the describe request.
	_ = c.Setopt(curl.OPT_NOBODY, 1)
	// Set custom timeout.
	_ = c.Setopt(curl.OPT_TIMEOUT_MS, int(s.timeout/time.Millisecond))

	// Enable verbose logs if verbose mode is on.
	if s.verbose {
		_ = c.Setopt(curl.OPT_VERBOSE, 1)
	} else {
		_ = c.Setopt(curl.OPT_VERBOSE, 0)
	}
}
var ErrConnection = errors.New("connection error")

var connectionErrorSubstrings = []string{
	"connection reset",
	"recv failure",
	"timeout",
	"connection refused",
	"broken pipe",
	"cseq",
	"503",               // service unavailable / possible rate limiting
	"service unavailable",
}

// isConnectionError reports whether err looks like a connection-related error.
func isConnectionError(err error) bool {
	if err == nil {
		return false
	}

	msg := strings.ToLower(err.Error())
	for _, needle := range connectionErrorSubstrings {
		if strings.Contains(msg, needle) {
			return true
		}
	}
	return false
}

// HACK: See https://stackoverflow.com/questions/3572397/lib-curl-in-c-disable-printing
func doNotWrite([]uint8, interface{}) bool {
	return true
}
