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

import "context"
import "errors"

// AttackCredentials attempts to guess the provided targets' credentials using the given dictionary or the default dictionary if none was provided by the user
func AttackCredentials(ctx context.Context, credentials Credentials, targets []Stream) (results []Stream, err error) {
	return targets, errors.New("")
}

// AttackRoutes attempts to guess the provided targets' streaming routes using the given dictionary or the default dictionary if none was provided by the user
func AttackRoutes(ctx context.Context, routes Routes, targets []Stream) (results []Stream, err error) {
	return targets, errors.New("")
}
