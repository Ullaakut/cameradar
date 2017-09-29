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
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// Again, since these tests use the curl library, I don't want to spend ages trying to mock
// the lib right now.

func TestAttackCredentials(t *testing.T) {
	validStream1 := Stream{
		Device:  "fakeDevice",
		Address: "fakeAddress",
		Port:    1337,
	}

	validStream2 := Stream{
		Device:  "fakeDevice",
		Address: "differentFakeAddress",
		Port:    1337,
	}

	fakeTargets := []Stream{validStream1, validStream2}
	fakeCredentials := Credentials{
		Usernames: []string{"admin", "root"},
		Passwords: []string{"12345", "root"},
	}

	vectors := []struct {
		targets     []Stream
		credentials Credentials
		timeout     time.Duration
		log         bool

		expectedStreams []Stream
		expectedErrMsg  string
	}{
		// Valid baseline
		{
			targets:     fakeTargets,
			credentials: fakeCredentials,
			timeout:     1 * time.Millisecond,
			log:         true,

			expectedStreams: fakeTargets,
		},
		// Valid baseline without logs
		{
			targets:     fakeTargets,
			credentials: fakeCredentials,
			timeout:     0 * time.Millisecond,
			log:         false,

			expectedStreams: fakeTargets,
		},
	}
	for i, vector := range vectors {
		results, err := AttackCredentials(vector.targets, vector.credentials, vector.timeout, vector.log)

		if len(vector.expectedErrMsg) > 0 {
			if err == nil {
				fmt.Printf("unexpected success in AttackCredentials test, iteration %d. expected error: %s\n", i, vector.expectedErrMsg)
				os.Exit(1)
			}
			assert.Contains(t, err.Error(), vector.expectedErrMsg, "wrong error message")
		} else {
			if err != nil {
				fmt.Printf("unexpected error in AttackCredentials test, iteration %d: %v\n", i, err)
				os.Exit(1)
			}
			for _, stream := range vector.expectedStreams {
				foundStream := false
				for _, result := range results {
					if result.Address == stream.Address && result.Device == stream.Device && result.Port == stream.Port {
						foundStream = true
					}
				}
				assert.Equal(t, true, foundStream, "wrong streams parsed")
			}
		}
		assert.Equal(t, len(vector.expectedStreams), len(results), "wrong streams parsed")

	}
}

func TestAttackRoute(t *testing.T) {
	validStream1 := Stream{
		Device:  "fakeDevice",
		Address: "fakeAddress",
		Port:    1337,
	}

	validStream2 := Stream{
		Device:  "fakeDevice",
		Address: "differentFakeAddress",
		Port:    1337,
	}

	fakeTargets := []Stream{validStream1, validStream2}
	fakeRoutes := Routes{"live.sdp", "media.amp"}

	vectors := []struct {
		targets []Stream
		routes  Routes
		timeout time.Duration
		log     bool

		expectedStreams []Stream
		expectedErrMsg  string
	}{
		// Valid baseline
		{
			targets: fakeTargets,
			routes:  fakeRoutes,
			timeout: 1 * time.Millisecond,
			log:     true,

			expectedStreams: fakeTargets,
		},
		// Valid baseline without logs
		{
			targets: fakeTargets,
			routes:  fakeRoutes,
			timeout: 0 * time.Millisecond,
			log:     false,

			expectedStreams: fakeTargets,
		},
	}
	for i, vector := range vectors {
		results, err := AttackRoute(vector.targets, vector.routes, vector.timeout, vector.log)

		if len(vector.expectedErrMsg) > 0 {
			if err == nil {
				fmt.Printf("unexpected success in AttackRoute test, iteration %d. expected error: %s\n", i, vector.expectedErrMsg)
				os.Exit(1)
			}
			assert.Contains(t, err.Error(), vector.expectedErrMsg, "wrong error message")
		} else {
			if err != nil {
				fmt.Printf("unexpected error in AttackRoute test, iteration %d: %v\n", i, err)
				os.Exit(1)
			}
			for _, stream := range vector.expectedStreams {
				foundStream := false
				for _, result := range results {
					if result.Address == stream.Address && result.Device == stream.Device && result.Port == stream.Port {
						foundStream = true
					}
				}
				assert.Equal(t, true, foundStream, "wrong streams parsed")
			}
		}
		assert.Equal(t, len(vector.expectedStreams), len(results), "wrong streams parsed")
	}
}
