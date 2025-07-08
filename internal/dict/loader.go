package dict

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

// credentials is a map of credentials.
type credentials struct {
	Usernames []string `json:"usernames"`
	Passwords []string `json:"passwords"`
}

// routes is a slice of routes.
type routes []string

// Dictionary groups routes and credentials for attacks.
type Dictionary struct {
	creds  credentials
	routes routes
}

// Usernames returns the usernames list.
func (d Dictionary) Usernames() []string {
	return d.creds.Usernames
}

// Passwords returns the passwords list.
func (d Dictionary) Passwords() []string {
	return d.creds.Passwords
}

// Routes returns the routes list.
func (d Dictionary) Routes() []string {
	return d.routes
}

// New loads a dictionary using the provided configuration.
func New(credentialsPath, routesPath string) (Dictionary, error) {
	creds, err := loadCredentials(credentialsPath)
	if err != nil {
		return Dictionary{}, err
	}

	routes, err := loadRoutes(routesPath)
	if err != nil {
		return Dictionary{}, err
	}

	return Dictionary{
		creds:  creds,
		routes: routes,
	}, nil
}

// loadCredentials loads credentials from a custom path or embedded defaults.
func loadCredentials(credentialsPath string) (credentials, error) {
	if strings.TrimSpace(credentialsPath) != "" {
		content, err := os.ReadFile(credentialsPath)
		if err != nil {
			return credentials{}, fmt.Errorf("reading credentials dictionary %q: %w", credentialsPath, err)
		}

		creds, err := parseCredentials(content)
		if err != nil {
			return credentials{}, err
		}

		return creds, nil
	}

	creds, err := parseCredentials(defaultCredentials)
	if err != nil {
		return credentials{}, err
	}

	return creds, nil
}

// loadRoutes loads routes from a custom path or embedded defaults.
func loadRoutes(routesPath string) (routes, error) {
	if strings.TrimSpace(routesPath) != "" {
		file, err := os.Open(routesPath)
		if err != nil {
			return nil, fmt.Errorf("opening routes dictionary %q: %w", routesPath, err)
		}
		defer file.Close()

		routes, err := parseRoutes(file)
		if err != nil {
			return nil, err
		}
		return routes, nil
	}

	reader := strings.NewReader(defaultRoutes)
	routes, err := parseRoutes(io.NopCloser(reader))
	if err != nil {
		return nil, err
	}

	return routes, nil
}

func parseCredentials(content []byte) (credentials, error) {
	if len(content) == 0 {
		return credentials{}, errors.New("credentials dictionary is empty")
	}

	var creds credentials
	err := json.Unmarshal(content, &creds)
	if err != nil {
		return credentials{}, fmt.Errorf("reading dictionary contents: %w", err)
	}

	return creds, nil
}

func parseRoutes(reader io.ReadCloser) (routes, error) {
	defer reader.Close()

	var routes routes
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		routes = append(routes, scanner.Text())
	}

	return routes, scanner.Err()
}
