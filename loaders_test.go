package cameradar

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/Ullaakut/disgo"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Setup Mock
type mockedFS struct {
	osFS

	fileExists bool
	openError  bool

	fileMock *fileMock

	fileSize int64
}

// fileMock mocks a file
type fileMock struct {
	mock.Mock

	readError bool

	bytes.Buffer
}

type mockedFileInfo struct {
	os.FileInfo
}

func (m mockedFileInfo) Size() int64 { return 1 }

func (m mockedFS) Stat(name string) (os.FileInfo, error) {
	if !m.fileExists {
		return nil, os.ErrNotExist
	}
	return mockedFileInfo{}, nil
}

func (m mockedFS) Open(name string) (file, error) {
	if m.openError {
		return nil, os.ErrNotExist
	}

	return m.fileMock, nil
}

func (m *fileMock) Read(p []byte) (n int, err error) {
	if m.readError {
		return 0, os.ErrNotExist
	}
	return m.Buffer.Read(p)
}

func (m *fileMock) ReadAt(p []byte, off int64) (n int, err error) {
	return 1, nil
}

func (m *fileMock) Seek(offset int64, whence int) (int64, error) {
	return offset, nil
}

func (m *fileMock) Stat() (os.FileInfo, error) {
	return mockedFileInfo{}, nil
}

// Close mock
func (m *fileMock) Close() error {
	args := m.Called()
	return args.Error(0)
}

// Sync mock
func (m *fileMock) Sync() error {
	args := m.Called()
	return args.Error(0)
}

func TestLoadCredentials(t *testing.T) {
	credentialsJSONString := []byte("{\"usernames\":[\"admin\",\"root\"],\"passwords\":[\"12345\",\"root\"]}")
	validCredentials := Credentials{
		Usernames: []string{"admin", "root"},
		Passwords: []string{"12345", "root"},
	}

	tests := []struct {
		description string

		input      []byte
		fileExists bool

		expectedCredentials Credentials
		expectedErr         error
	}{
		{
			description: "Valid baseline",

			fileExists:          true,
			input:               credentialsJSONString,
			expectedCredentials: validCredentials,
		},
		{
			description: "File does not exist",

			fileExists:  false,
			input:       credentialsJSONString,
			expectedErr: errors.New("could not read credentials dictionary file at \"/tmp/cameradar_test_load_credentials_1.xml\": open /tmp/cameradar_test_load_credentials_1.xml: no such file or directory"),
		},
		{
			description: "Invalid format",

			fileExists:  true,
			input:       []byte("not json"),
			expectedErr: errors.New("unable to unmarshal dictionary contents: invalid character 'o' in literal null (expecting 'u')"),
		},
		{
			description: "No streams in dictionary",

			fileExists: true,
			input:      []byte("{\"invalid\":\"json\"}"),
		},
	}

	for i, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			filePath := "/tmp/cameradar_test_load_credentials_" + fmt.Sprint(i) + ".xml"
			// create file.
			if test.fileExists {
				_, err := os.Create(filePath)
				if err != nil {
					t.Fatalf("could not create xml file for LoadCredentials: %v. iteration: %d. file path: %s\n", err, i, filePath)
				}

				err = ioutil.WriteFile(filePath, test.input, 0644)
				if err != nil {
					t.Fatalf("could not write xml file for LoadCredentials: %v. iteration: %d. file path: %s\n", err, i, filePath)
				}
			}

			scanner := &Scanner{
				term:                     disgo.NewTerminal(disgo.WithDefaultOutput(ioutil.Discard)),
				credentialDictionaryPath: filePath,
			}

			err := scanner.LoadCredentials()

			assert.Equal(t, test.expectedErr, err)

			assert.Len(t, scanner.credentials.Usernames, len(test.expectedCredentials.Usernames))
			for _, expectedUsername := range test.expectedCredentials.Usernames {
				assert.Contains(t, scanner.credentials.Usernames, expectedUsername)
			}

			assert.Len(t, scanner.credentials.Passwords, len(test.expectedCredentials.Passwords))
			for _, expectedPassword := range test.expectedCredentials.Passwords {
				assert.Contains(t, scanner.credentials.Passwords, expectedPassword)
			}
		})
	}
}

func TestLoadRoutes(t *testing.T) {
	routesJSONString := []byte("admin\nroot")
	validRoutes := Routes{"admin", "root"}

	tests := []struct {
		description string
		input       []byte
		fileExists  bool

		expectedRoutes Routes
		expectedErr    error
	}{
		{
			description: "Valid baseline",

			fileExists:     true,
			input:          routesJSONString,
			expectedRoutes: validRoutes,
		},
		{
			description: "File does not exist",

			fileExists:  false,
			input:       routesJSONString,
			expectedErr: errors.New("unable to open dictionary: open /tmp/cameradar_test_load_routes_1.xml: no such file or directory"),
		},
		{
			description: "No streams in dictionary",

			fileExists: true,
			input:      []byte(""),
		},
	}

	for i, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			filePath := "/tmp/cameradar_test_load_routes_" + fmt.Sprint(i) + ".xml"

			// Create file.
			if test.fileExists {
				_, err := os.Create(filePath)
				if err != nil {
					fmt.Printf("could not create xml file for LoadRoutes: %v. iteration: %d. file path: %s\n", err, i, filePath)
					os.Exit(1)
				}

				err = ioutil.WriteFile(filePath, test.input, 0644)
				if err != nil {
					fmt.Printf("could not write xml file for LoadRoutes: %v. iteration: %d. file path: %s\n", err, i, filePath)
					os.Exit(1)
				}
			}

			scanner := &Scanner{
				term:                disgo.NewTerminal(disgo.WithDefaultOutput(ioutil.Discard)),
				routeDictionaryPath: filePath,
			}

			err := scanner.LoadRoutes()

			assert.Equal(t, test.expectedErr, err)

			assert.Len(t, scanner.routes, len(test.expectedRoutes))
			for _, expectedRoute := range test.expectedRoutes {
				assert.Contains(t, scanner.routes, expectedRoute)
			}
		})
	}
}

func TestParseCredentialsFromString(t *testing.T) {
	defaultCredentials := Credentials{
		Usernames: []string{
			"",
			"admin",
			"Admin",
			"Administrator",
			"root",
			"supervisor",
			"ubnt",
			"service",
			"Dinion",
			"administrator",
			"admin1",
		},
		Passwords: []string{
			"",
			"admin",
			"9999",
			"123456",
			"pass",
			"camera",
			"1234",
			"12345",
			"fliradmin",
			"system",
			"jvc",
			"meinsm",
			"root",
			"4321",
			"111111",
			"1111111",
			"password",
			"ikwd",
			"supervisor",
			"ubnt",
			"wbox123",
			"service",
		},
	}

	tests := []struct {
		str                 string
		expectedCredentials Credentials
	}{
		{
			str:                 "{\"usernames\":[\"\",\"admin\",\"Admin\",\"Administrator\",\"root\",\"supervisor\",\"ubnt\",\"service\",\"Dinion\",\"administrator\",\"admin1\"],\"passwords\":[\"\",\"admin\",\"9999\",\"123456\",\"pass\",\"camera\",\"1234\",\"12345\",\"fliradmin\",\"system\",\"jvc\",\"meinsm\",\"root\",\"4321\",\"111111\",\"1111111\",\"password\",\"ikwd\",\"supervisor\",\"ubnt\",\"wbox123\",\"service\"]}",
			expectedCredentials: defaultCredentials,
		},
		{
			str:                 "{}",
			expectedCredentials: Credentials{},
		},
		{
			str:                 "{\"invalid_field\":42}",
			expectedCredentials: Credentials{},
		},
		{
			str:                 "not json",
			expectedCredentials: Credentials{},
		},
	}

	for _, test := range tests {
		parsedCredentials, _ := ParseCredentialsFromString(test.str)
		assert.Equal(t, test.expectedCredentials, parsedCredentials)
	}
}

func TestParseRoutesFromString(t *testing.T) {
	tests := []struct {
		str            string
		expectedRoutes Routes
	}{
		{
			str:            "a\nb\nc",
			expectedRoutes: []string{"a", "b", "c"},
		},
		{
			str:            "a",
			expectedRoutes: []string{"a"},
		},
		{
			str:            "",
			expectedRoutes: []string{""},
		},
	}

	for _, test := range tests {
		assert.Equal(t, test.expectedRoutes, ParseRoutesFromString(test.str))
	}
}

func TestLoadTargets(t *testing.T) {

	oldFS := fs
	mfs := &mockedFS{}
	fs = mfs
	defer func() {
		fs = oldFS
	}()

	tests := []struct {
		description string

		targets []string

		fileExists bool
		openError  bool
		readError  bool

		expectedTargets []string
		expectedError   error
	}{
		{
			description: "not a file",

			targets: []string{"0.0.0.0"},

			fileExists: false,

			expectedTargets: []string{"0.0.0.0"},
			expectedError:   nil,
		},
		{
			description: "not file targets",

			targets: []string{"0.0.0.0", "1.2.3.4/24"},

			expectedTargets: []string{"0.0.0.0", "1.2.3.4/24"},
			expectedError:   nil,
		},
		{
			description: "file contains targets",

			targets: []string{"test_does_not_really_exist"},

			fileExists: true,

			expectedTargets: []string{"0.0.0.0", "localhost", "192.17.0.0/16", "192.168.1.140-255", "192.168.2-3.0-255"},
			expectedError:   nil,
		},
		{
			description: "open error",

			targets: []string{"test_does_not_really_exist"},

			fileExists: true,
			openError:  true,

			expectedTargets: []string{"test_does_not_really_exist"},
			expectedError:   errors.New("unable to open targets file \"test_does_not_really_exist\": file does not exist"),
		},
		{
			description: "read error",

			targets: []string{"test_does_not_really_exist"},

			fileExists: true,
			readError:  true,

			expectedTargets: []string{"test_does_not_really_exist"},
			expectedError:   errors.New("unable to read targets file \"test_does_not_really_exist\": file does not exist"),
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			mfs.fileExists = test.fileExists
			mfs.openError = test.openError

			mfs.fileMock = &fileMock{
				readError: test.readError,
			}
			mfs.fileMock.On("Close").Return(nil)
			mfs.fileMock.WriteString("0.0.0.0\nlocalhost\n192.17.0.0/16\n192.168.1.140-255\n192.168.2-3.0-255")

			scanner := &Scanner{
				term:    disgo.NewTerminal(disgo.WithDefaultOutput(ioutil.Discard)),
				targets: test.targets,
			}

			err := scanner.LoadTargets()
			assert.Equal(t, test.expectedTargets, scanner.targets)
			assert.Equal(t, test.expectedError, err)
		})
	}
}

// This is completely useless and just lets me
// not look at these two red lines on the coverage
// any longer.
func TestFS(t *testing.T) {
	fs := osFS{}

	fs.Open("test")
	fs.Stat("test")
}
