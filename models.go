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

import "time"

// Stream represents a camera's RTSP stream
type Stream struct {
	Device   string `json:"device"`
	Username string `json:"username"`
	Password string `json:"password"`
	Route    string `json:"route"`
	Address  string `json:"address" validate:"required"`
	Port     uint   `json:"port" validate:"required"`

	CredentialsFound bool `json:"credentials_found"`
	RouteFound       bool `json:"route_found"`
}

// Credentials is a map of credentials
// usernames are keys and passwords are values
// creds['admin'] -> 'secure_password'
type Credentials struct {
	Usernames []string `json:"usernames"`
	Passwords []string `json:"passwords"`
}

// Routes is a slice of Routes
// ['/live.sdp', '/media.amp', ...]
type Routes []string

// Options contains all options needed to launch a complete cameradar scan
type Options struct {
	Target      string        `json:"target" validate:"required"`
	Ports       string        `json:"ports"`
	OutputFile  string        `json:"output_file"`
	Routes      Routes        `json:"routes"`
	Credentials Credentials   `json:"credentials"`
	Speed       int           `json:"speed"`
	Timeout     time.Duration `json:"timeout"`
}
