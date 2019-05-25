package cameradar

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
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

// LoadCredentials opens a dictionary file and returns its contents as a Credentials structure.
func (s *Scanner) LoadCredentials() error {
	s.term.Debugf("Loading credentials dictionary from path %q\n", s.credentialDictionaryPath)

	// Open & Read XML file.
	content, err := ioutil.ReadFile(s.credentialDictionaryPath)
	if err != nil {
		return fmt.Errorf("could not read credentials dictionary file at %q: %v", s.credentialDictionaryPath, err)
	}

	// Unmarshal content of JSON file into data structure.
	err = json.Unmarshal(content, &s.credentials)
	if err != nil {
		return err
	}

	s.term.Debugf("Loaded %d usernames and %d passwords\n", len(s.credentials.Usernames), len(s.credentials.Passwords))
	return nil
}

// LoadRoutes opens a dictionary file and returns its contents as a Routes structure.
func (s *Scanner) LoadRoutes() error {
	s.term.Debugf("Loading routes dictionary from path %q\n", s.routeDictionaryPath)

	file, err := os.Open(s.routeDictionaryPath)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		s.routes = append(s.routes, scanner.Text())
	}

	err = scanner.Err()
	if err != nil {
		return err
	}

	s.term.Debugf("Loaded %d routes\n", len(s.routes))
	return nil
}

// ParseCredentialsFromString parses a dictionary string and returns its contents as a Credentials structure.
func ParseCredentialsFromString(content string) (Credentials, error) {
	var creds Credentials

	// Unmarshal content of JSON file into data structure.
	err := json.Unmarshal([]byte(content), &creds)
	if err != nil {
		return creds, err
	}

	return creds, nil
}

// ParseRoutesFromString parses a dictionary string and returns its contents as a Routes structure.
func ParseRoutesFromString(content string) Routes {
	return strings.Split(content, "\n")
}

// LoadTargets parses the file containing hosts to targets, if the targets are
// just set to a file name.
func (s *Scanner) LoadTargets() error {
	if len(s.targets) != 1 {
		return nil
	}

	path := s.targets[0]

	_, err := fs.Stat(path)
	if err != nil {
		return nil
	}

	file, err := fs.Open(path)
	if err != nil {
		return fmt.Errorf("unable to open targets file %q: %v", path, err)
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return fmt.Errorf("unable to read targets file %q: %v", path, err)
	}

	s.targets = strings.Split(string(bytes), "\n")
	return nil
}
