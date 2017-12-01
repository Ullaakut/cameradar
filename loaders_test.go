package cmrdr

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadCredentials(t *testing.T) {
	credentialsJSONString := []byte("{\"usernames\":[\"admin\",\"root\"],\"passwords\":[\"12345\",\"root\"]}")
	validCredentials := Credentials{
		Usernames: []string{"admin", "root"},
		Passwords: []string{"12345", "root"},
	}

	vectors := []struct {
		input      []byte
		fileExists bool

		expectedOutput Credentials
		expectedErrMsg string
	}{
		// Valid baseline
		{
			fileExists:     true,
			input:          credentialsJSONString,
			expectedOutput: validCredentials,
		},
		// File does not exist
		{
			fileExists:     false,
			input:          credentialsJSONString,
			expectedErrMsg: "could not read credentials dictionary file at",
		},
		// Invalid format
		{
			fileExists:     true,
			input:          []byte("not json"),
			expectedErrMsg: "invalid character",
		},
		// No streams in dictionary
		{
			fileExists: true,
			input:      []byte("{\"invalid\":\"json\"}"),
		},
	}
	for i, vector := range vectors {
		filePath := "/tmp/cameradar_test_load_credentials_" + fmt.Sprint(i) + ".xml"
		// create file
		if vector.fileExists {
			_, err := os.Create(filePath)
			if err != nil {
				fmt.Printf("could not create xml file for LoadCredentials: %v. iteration: %d. file path: %s\n", err, i, filePath)
				os.Exit(1)
			}

			err = ioutil.WriteFile(filePath, vector.input, 0644)
			if err != nil {
				fmt.Printf("could not write xml file for LoadCredentials: %v. iteration: %d. file path: %s\n", err, i, filePath)
				os.Exit(1)
			}
		}

		result, err := LoadCredentials(filePath)
		if len(vector.expectedErrMsg) > 0 {
			if err == nil {
				fmt.Printf("unexpected success in LoadCredentials test, iteration %d. expected error: %s\n", i, vector.expectedErrMsg)
				os.Exit(1)
			}
			assert.Contains(t, err.Error(), vector.expectedErrMsg, "wrong error message")
		} else {
			if err != nil {
				fmt.Printf("unexpected error in LoadCredentials test, iteration %d: %v\n", i, err)
				os.Exit(1)
			}
			for _, expectedUsername := range vector.expectedOutput.Usernames {
				foundUsername := false
				for _, username := range result.Usernames {
					if username == expectedUsername {
						foundUsername = true
					}
				}
				assert.Equal(t, true, foundUsername, "wrong usernames parsed")
			}
			for _, expectedPassword := range vector.expectedOutput.Passwords {
				foundPassword := false
				for _, password := range result.Passwords {
					if password == expectedPassword {
						foundPassword = true
					}
				}
				assert.Equal(t, true, foundPassword, "wrong passwords parsed")
			}
		}
	}
}

func TestLoadRoutes(t *testing.T) {
	routesJSONString := []byte("admin\nroot")
	validRoutes := Routes{"admin", "root"}

	vectors := []struct {
		input      []byte
		fileExists bool

		expectedOutput Routes
		expectedErrMsg string
	}{
		// Valid baseline
		{
			fileExists:     true,
			input:          routesJSONString,
			expectedOutput: validRoutes,
		},
		// File does not exist
		{
			fileExists:     false,
			input:          routesJSONString,
			expectedErrMsg: "no such file or directory",
		},
		// No streams in dictionary
		{
			fileExists: true,
			input:      []byte(""),
		},
	}
	for i, vector := range vectors {
		filePath := "/tmp/cameradar_test_load_routes_" + fmt.Sprint(i) + ".xml"
		// create file
		if vector.fileExists {
			_, err := os.Create(filePath)
			if err != nil {
				fmt.Printf("could not create xml file for LoadRoutes: %v. iteration: %d. file path: %s\n", err, i, filePath)
				os.Exit(1)
			}

			err = ioutil.WriteFile(filePath, vector.input, 0644)
			if err != nil {
				fmt.Printf("could not write xml file for LoadRoutes: %v. iteration: %d. file path: %s\n", err, i, filePath)
				os.Exit(1)
			}
		}

		result, err := LoadRoutes(filePath)
		if len(vector.expectedErrMsg) > 0 {
			if err == nil {
				fmt.Printf("unexpected success in LoadRoutes test, iteration %d. expected error: %s\n", i, vector.expectedErrMsg)
				os.Exit(1)
			}
			assert.Contains(t, err.Error(), vector.expectedErrMsg, "wrong error message")
		} else {
			if err != nil {
				fmt.Printf("unexpected error in LoadRoutes test, iteration %d: %v\n", i, err)
				os.Exit(1)
			}
			for _, expectedRoute := range vector.expectedOutput {
				foundRoute := false
				for _, route := range result {
					if route == expectedRoute {
						foundRoute = true
					}
				}
				assert.Equal(t, true, foundRoute, "wrong routes parsed")
			}
		}
	}
}
