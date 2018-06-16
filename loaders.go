package cmrdr

import (
	"bufio"
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/pkg/errors"
)

var fs fileSystem = osFS{}

type fileSystem interface {
	Open(name string) (file, error)
	Stat(name string) (os.FileInfo, error)
}

type file interface {
	io.Closer
	io.Reader
	io.ReaderAt
	io.Seeker
	Stat() (os.FileInfo, error)
}

// osFS implements fileSystem using the local disk.
type osFS struct{}

func (osFS) Open(name string) (file, error)        { return os.Open(name) }
func (osFS) Stat(name string) (os.FileInfo, error) { return os.Stat(name) }

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

// ParseTargetsFile parses an input file containing hosts to targets
func ParseTargetsFile(path string) (string, error) {
	_, err := fs.Stat(path)
	if err != nil {
		return path, nil
	}

	file, err := fs.Open(path)
	if err != nil {
		return path, err
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return path, err
	}

	return strings.Replace(string(bytes), "\n", " ", -1), nil
}
