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

import "fmt"

func replace(streams []Stream, new Stream) []Stream {
	updatedSlice := streams[:0]

	for _, old := range streams {
		if old.Address == new.Address && old.Port == new.Port {
			updatedSlice = append(updatedSlice, new)
		} else {
			updatedSlice = append(updatedSlice, old)
		}
	}
	return updatedSlice
}

// RTSPURL generates a stream's RTSP URL
func RTSPURL(stream Stream) string {
	return "rtsp://" + stream.Username + ":" + stream.Password + "@" + stream.Address + ":" + fmt.Sprint(stream.Port) + "/" + stream.Route
}

// AdminPanelURL returns the URL to the camera's admin panel
func AdminPanelURL(stream Stream) string {
	return "http://" + stream.Address + "/"
}
