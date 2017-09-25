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
	"bufio"
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/pkg/errors"
)

// LoadCredentials opens a dictionary file and returns its contents as a Credentials structure
func LoadCredentials(path string) (Credentials, error) {
	var creds Credentials

	// Open & Read XML file
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return creds, errors.Wrap(err, "Could not read credentials dictionary file at "+path+":")
	}

	// Unmarshal content of JSON file into data structure
	err = json.Unmarshal(content, &creds)
	if err != nil {
		return creds, err
	}

	return creds, nil
}

// LoadRoutes opens a dictionary file and returns its contents as a Routes structure
func LoadRoutes(path string) (Routes, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var routes Routes
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		routes = append(routes, scanner.Text())
	}

	return routes, scanner.Err()
}
