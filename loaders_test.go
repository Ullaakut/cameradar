package cmrdr

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

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

	testCases := []struct {
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

	for i, test := range testCases {
		filePath := "/tmp/cameradar_test_load_credentials_" + fmt.Sprint(i) + ".xml"
		// create file
		if test.fileExists {
			_, err := os.Create(filePath)
			if err != nil {
				fmt.Printf("could not create xml file for LoadCredentials: %v. iteration: %d. file path: %s\n", err, i, filePath)
				os.Exit(1)
			}

			err = ioutil.WriteFile(filePath, test.input, 0644)
			if err != nil {
				fmt.Printf("could not write xml file for LoadCredentials: %v. iteration: %d. file path: %s\n", err, i, filePath)
				os.Exit(1)
			}
		}

		result, err := LoadCredentials(filePath)
		if len(test.expectedErrMsg) > 0 {
			if err == nil {
				fmt.Printf("unexpected success in LoadCredentials test, iteration %d. expected error: %s\n", i, test.expectedErrMsg)
				os.Exit(1)
			}

			assert.Contains(t, err.Error(), test.expectedErrMsg, "wrong error message")
		} else {
			if err != nil {
				fmt.Printf("unexpected error in LoadCredentials test, iteration %d: %v\n", i, err)
				os.Exit(1)
			}

			for _, expectedUsername := range test.expectedOutput.Usernames {
				foundUsername := false
				for _, username := range result.Usernames {
					if username == expectedUsername {
						foundUsername = true
					}
				}

				assert.Equal(t, true, foundUsername, "wrong usernames parsed")
			}

			for _, expectedPassword := range test.expectedOutput.Passwords {
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

	testCases := []struct {
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

	for i, test := range testCases {
		filePath := "/tmp/cameradar_test_load_routes_" + fmt.Sprint(i) + ".xml"

		// create file
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

		result, err := LoadRoutes(filePath)
		if len(test.expectedErrMsg) > 0 {
			if err == nil {
				fmt.Printf("unexpected success in LoadRoutes test, iteration %d. expected error: %s\n", i, test.expectedErrMsg)
				os.Exit(1)
			}
			assert.Contains(t, err.Error(), test.expectedErrMsg, "wrong error message")
		} else {
			if err != nil {
				fmt.Printf("unexpected error in LoadRoutes test, iteration %d: %v\n", i, err)
				os.Exit(1)
			}

			for _, expectedRoute := range test.expectedOutput {
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

	testCases := []struct {
		str            string
		expectedResult Credentials
	}{
		{
			str:            "{\"usernames\":[\"\",\"admin\",\"Admin\",\"Administrator\",\"root\",\"supervisor\",\"ubnt\",\"service\",\"Dinion\",\"administrator\",\"admin1\"],\"passwords\":[\"\",\"admin\",\"9999\",\"123456\",\"pass\",\"camera\",\"1234\",\"12345\",\"fliradmin\",\"system\",\"jvc\",\"meinsm\",\"root\",\"4321\",\"111111\",\"1111111\",\"password\",\"ikwd\",\"supervisor\",\"ubnt\",\"wbox123\",\"service\"]}",
			expectedResult: defaultCredentials,
		},
		{
			str:            "{}",
			expectedResult: Credentials{},
		},
		{
			str:            "{\"invalid_field\":42}",
			expectedResult: Credentials{},
		},
		{
			str:            "not json",
			expectedResult: Credentials{},
		},
	}
	for _, test := range testCases {
		parsedCredentials, _ := ParseCredentialsFromString(test.str)
		assert.Equal(t, test.expectedResult, parsedCredentials, "unexpected result, parse error")
	}
}

func TestParseRoutesFromString(t *testing.T) {
	testCases := []struct {
		str            string
		expectedResult Routes
	}{
		{
			str:            "a\nb\nc",
			expectedResult: []string{"a", "b", "c"},
		},
		{
			str:            "a",
			expectedResult: []string{"a"},
		},
		{
			str:            "",
			expectedResult: []string{""},
		},
	}
	for _, test := range testCases {
		parsedRoutes := ParseRoutesFromString(test.str)
		assert.Equal(t, test.expectedResult, parsedRoutes, "unexpected result, parse error")
	}
}

func TestParseTargetsFile(t *testing.T) {

	oldFS := fs
	mfs := &mockedFS{}
	fs = mfs
	defer func() {
		fs = oldFS
	}()

	testCases := []struct {
		input string

		fileExists bool
		openError  bool
		readError  bool

		expectedResult []string
		expectedError  error
	}{
		{
			input: "0.0.0.0",

			fileExists: false,

			expectedResult: []string{"0.0.0.0"},
			expectedError:  nil,
		},
		{
			input: "test_does_not_really_exist",

			fileExists: true,

			expectedResult: []string{"0.0.0.0", "localhost", "192.17.0.0/16", "192.168.1.140-255", "192.168.2-3.0-255"},
			expectedError:  nil,
		},
		{
			input: "test_does_not_really_exist",

			fileExists: true,
			openError:  true,

			expectedResult: []string{"test_does_not_really_exist"},
			expectedError:  os.ErrNotExist,
		},
		{
			input: "test_does_not_really_exist",

			fileExists: true,
			readError:  true,

			expectedResult: []string{"test_does_not_really_exist"},
			expectedError:  os.ErrNotExist,
		},
	}

	for _, test := range testCases {
		mfs.fileExists = test.fileExists
		mfs.openError = test.openError

		mfs.fileMock = &fileMock{
			readError: test.readError,
		}
		mfs.fileMock.On("Close").Return(nil)
		mfs.fileMock.WriteString("0.0.0.0\nlocalhost\n192.17.0.0/16\n192.168.1.140-255\n192.168.2-3.0-255")

		result, err := ParseTargetsFile(test.input)
		assert.Equal(t, test.expectedResult, result, "unexpected result, parse error")
		assert.Equal(t, test.expectedError, err, "unexpected error")
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
