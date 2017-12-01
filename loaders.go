package cmrdr

import (
	"bufio"
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"

	"github.com/pkg/errors"
)

// LoadCredentials opens a dictionary file and returns its contents as a Credentials structure
func LoadCredentials(path string) (Credentials, error) {
	var creds Credentials

	// Open & Read XML file
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return creds, errors.Wrap(err, "could not read credentials dictionary file at "+path+":")
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

// ParseCredentialsFromString parses a dictionary string and returns its contents as a Credentials structure
func ParseCredentialsFromString(content string) (Credentials, error) {
	var creds Credentials

	// Unmarshal content of JSON file into data structure
	err := json.Unmarshal([]byte(content), &creds)
	if err != nil {
		return creds, err
	}

	return creds, nil
}

// ParseRoutesFromString parses a dictionary string and returns its contents as a Routes structure
func ParseRoutesFromString(content string) Routes {
	return strings.Split(content, "\n")
}
