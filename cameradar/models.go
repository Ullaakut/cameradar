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

// Stream represents a camera's RTSP stream
type Stream struct {
	device   string
	username string
	password string
	route    string
	address  string `validate:"required"`
	port     uint   `validate:"required"`
}

// Attack represents the state of the attack on a stream
type Attack struct {
	stream           Stream
	credentialsFound bool
	routeFound       bool
}

// Credentials is a map of credentials
// usernames are keys and passwords are values
// creds['admin'] -> 'secure_password'
type Credentials struct {
	Usernames []string
	Passwords []string
}

// Routes is a slice of Routes
// ['/live.sdp', '/media.amp', ...]
type Routes []string
