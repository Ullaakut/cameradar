package cameradar

import (
	"fmt"
	"time"

	"github.com/pkg/errors"
	curl "github.com/ullaakut/go-curl"
	v "gopkg.in/go-playground/validator.v9"
)

// HTTP responses
const (
	httpOK           = 200
	httpUnauthorized = 401
	httpForbidden    = 403
	httpNotFound     = 404
)

// CURL RTSP request types
const (
	rtspDescribe = 2
	rtspSetup    = 4
)

// ValidateStreams tries to setup the stream to validate whether or not it is available
func (s *Scanner) ValidateStreams(targets []Stream) ([]Stream, error) {
	for i := range targets {
		targets[i].Available = validateStream(targets[i])
	}

	return targets, nil
}

// AttackCredentials attempts to guess the provided targets' credentials using the given
// dictionary or the default dictionary if none was provided by the user.
func (s *Scanner) AttackCredentials(targets []Stream, credentials Credentials) ([]Stream, error) {
	resChan := make(chan Stream)
	defer close(resChan)

	validate := v.New()
	for i := range targets {
		err := validate.Struct(targets[i])
		if err != nil {
			return targets, errors.Wrap(err, "invalid targets")
		}

		// TODO: Perf Improvement: Skip cameras with no auth type detected, and set their
		// CredentialsFound value to true.
		go attackCameraCredentials(targets[i], credentials, resChan)
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

	return targets, nil
}

// AttackRoute attempts to guess the provided targets' streaming routes using the given
// dictionary or the default dictionary if none was provided by the user.
func (s *Scanner) AttackRoute(targets []Stream, routes Routes) ([]Stream, error) {
	resChan := make(chan Stream)
	defer close(resChan)

	validate := v.New()
	for i := range targets {
		err := validate.Struct(targets[i])
		if err != nil {
			return targets, errors.Wrap(err, "invalid targets")
		}

		go attackCameraRoute(targets[i], routes, resChan)
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

	return targets, nil
}

// DetectAuthMethods attempts to guess the provided targets' authentication types, between
// digest, basic auth or none at all.
func (s *Scanner) DetectAuthMethods(targets []Stream) ([]Stream, error) {
	for i := range targets {
		targets[i].AuthenticationType = detectAuthMethod(targets[i])
	}

	return targets, nil
}

func (s *Scanner) attackCameraCredentials(target Stream, credentials Credentials, resChan chan<- Stream) {
	for _, username := range credentials.Usernames {
		for _, password := range credentials.Passwords {
			ok := credAttack(c.Duphandle(), target, username, password, timeout, log)
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

func (s *Scanner) attackCameraRoute(target Stream, routes Routes, resChan chan<- Stream) {
	for _, route := range routes {
		ok := routeAttack(c.Duphandle(), target, route, timeout, log)
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

func (s *Scanner) detectAuthMethod(stream Stream, timeout time.Duration) int {
	attackURL := fmt.Sprintf(
		"rtsp://%s:%d/%s",
		stream.Address,
		stream.Port,
		stream.Route,
	)

	s.setCurlOptions()

	// Send a request to the URL of the stream we want to attack
	s.curl.Setopt(curl.OPT_URL, attackURL)
	// Set the RTSP STREAM URI as the stream URL
	s.curl.Setopt(curl.OPT_RTSP_STREAM_URI, attackURL)
	// 2 is CURL_RTSPREQ_DESCRIBE
	s.curl.Setopt(curl.OPT_RTSP_REQUEST, 2)
	// Set custom timeout
	s.curl.Setopt(curl.OPT_TIMEOUT_MS, int(timeout/time.Millisecond))

	// Perform the request
	err := c.Perform()
	if err != nil {
		return -1
	}

	authType, err := c.Getinfo(curl.INFO_HTTPAUTH_AVAIL)
	if err != nil {
		return -1
	}

	return authType.(int)
}

func (s *Scanner) routeAttack(stream Stream, route string, timeout time.Duration) bool {
	attackURL := fmt.Sprintf(
		"rtsp://%s:%s@%s:%d/%s",
		stream.Username,
		stream.Password,
		stream.Address,
		stream.Port,
		route,
	)

	s.setCurlOptions()

	// Set proper authentication type.
	s.curl.Setopt(curl.OPT_HTTPAUTH, stream.AuthenticationType)
	s.curl.Setopt(curl.OPT_USERPWD, fmt.Sprint(stream.Username, ":", stream.Password))

	// Send a request to the URL of the stream we want to attack
	s.curl.Setopt(curl.OPT_URL, attackURL)
	// Set the RTSP STREAM URI as the stream URL
	s.curl.Setopt(curl.OPT_RTSP_STREAM_URI, attackURL)
	// 2 is CURL_RTSPREQ_DESCRIBE
	s.curl.Setopt(curl.OPT_RTSP_REQUEST, rtspDescribe)
	// Set custom timeout
	s.curl.Setopt(curl.OPT_TIMEOUT_MS, int(timeout/time.Millisecond))

	// Perform the request
	err := c.Perform()
	if err != nil {
		return false
	}

	// Get return code for the request
	rc, err := c.Getinfo(curl.INFO_RESPONSE_CODE)
	if err != nil {
		return false
	}

	// If it's a 401 or 403, it means that the credentials are wrong but the route might be okay
	// If it's a 200, the stream is accessed successfully
	if rc == httpOK || rc == httpUnauthorized || rc == httpForbidden {
		return true
	}
	return false
}

func (s *Scanner) credAttack(stream Stream, username string, password string, timeout time.Duration) bool {
	attackURL := fmt.Sprintf(
		"rtsp://%s:%s@%s:%d/%s",
		username,
		password,
		stream.Address,
		stream.Port,
		stream.Route,
	)

	s.setCurlOptions()

	// Set proper authentication type.
	s.curl.Setopt(curl.OPT_HTTPAUTH, stream.AuthenticationType)
	s.curl.Setopt(curl.OPT_USERPWD, fmt.Sprint(username, ":", password))

	// Send a request to the URL of the stream we want to attack
	s.curl.Setopt(curl.OPT_URL, attackURL)
	// Set the RTSP STREAM URI as the stream URL
	s.curl.Setopt(curl.OPT_RTSP_STREAM_URI, attackURL)
	// 2 is CURL_RTSPREQ_DESCRIBE
	s.curl.Setopt(curl.OPT_RTSP_REQUEST, 2)
	// Set custom timeout
	s.curl.Setopt(curl.OPT_TIMEOUT_MS, int(timeout/time.Millisecond))

	// Perform the request
	err := c.Perform()
	if err != nil {
		return false
	}

	// Get return code for the request
	rc, err := c.Getinfo(curl.INFO_RESPONSE_CODE)
	if err != nil {
		return false
	}

	// If it's a 404, it means that the route is incorrect but the credentials might be okay
	// If it's a 200, the stream is accessed successfully
	if rc == httpOK || rc == httpNotFound {
		return true
	}
	return false
}

func (s *Scanner) validateStream(stream Stream, timeout time.Duration) bool {
	attackURL := fmt.Sprintf(
		"rtsp://%s:%s@%s:%d/%s",
		stream.Username,
		stream.Password,
		stream.Address,
		stream.Port,
		stream.Route,
	)

	s.setCurlOptions()

	// Set proper authentication type.
	s.curl.Setopt(curl.OPT_HTTPAUTH, stream.AuthenticationType)
	s.curl.Setopt(curl.OPT_USERPWD, fmt.Sprint(stream.Username, ":", stream.Password))

	// Send a request to the URL of the stream we want to attack.
	s.curl.Setopt(curl.OPT_URL, attackURL)
	// Set the RTSP STREAM URI as the stream URL.
	s.curl.Setopt(curl.OPT_RTSP_STREAM_URI, attackURL)
	// 2 is CURL_RTSPREQ_SETUP.
	s.curl.Setopt(curl.OPT_RTSP_REQUEST, rtspSetup)
	// Set custom timeout.
	s.curl.Setopt(curl.OPT_TIMEOUT_MS, int(timeout/time.Millisecond))

	s.curl.Setopt(curl.OPT_RTSP_TRANSPORT, "RTP/AVP;unicast;client_port=33332-33333")

	// Perform the request.
	err := c.Perform()
	if err != nil {
		return false
	}

	// Get return code for the request.
	rc, err := c.Getinfo(curl.INFO_RESPONSE_CODE)
	if err != nil {
		return false
	}

	// If it's a 200, the stream is accessed successfully.
	if rc == httpOK {
		return true
	}
	return false
}

func (s *Scanner) setCurlOptions() {
	if s.debug {
		// Debug logs when logs are enabled
		_ = s.curl.Setopt(curl.OPT_VERBOSE, 1)
	} else {
		// Do not write sdp in stdout
		_ = s.curl.Setopt(curl.OPT_WRITEFUNCTION, doNotWrite)
	}

	// Do not use signals (would break multithreading).
	_ = s.curl.Setopt(curl.OPT_NOSIGNAL, 1)
	// Do not send a body in the describe request.
	_ = s.curl.Setopt(curl.OPT_NOBODY, 1)
}

// HACK: See https://stackoverflow.com/questions/3572397/lib-curl-in-c-disable-printing
func doNotWrite([]uint8, interface{}) bool {
	return true
}
