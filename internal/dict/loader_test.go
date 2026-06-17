package dict_test

import (
	"bufio"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Ullaakut/cameradar/v6/internal/dict"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew_LoadsDictionaryFromPaths(t *testing.T) {
	tempDir := t.TempDir()
	credsPath := writeTempFile(t, tempDir, "creds.json", `{"usernames":["alice"],"passwords":["secret"]}`)
	routesPath := writeTempFile(t, tempDir, "routes", "stream\nother\n")

	got, err := dict.New(credsPath, routesPath)
	require.NoError(t, err)

	assert.Equal(t, []string{"alice"}, got.Usernames())
	assert.Equal(t, []string{"secret"}, got.Passwords())
	assert.Equal(t, []string{"stream", "other"}, got.Routes())
}

func TestNew_CustomAndDefaultPaths(t *testing.T) {
	tempDir := t.TempDir()
	customCredsPath := writeTempFile(t, tempDir, "creds.json", `{"usernames":["alice"],"passwords":["secret"]}`)
	customRoutesPath := writeTempFile(t, tempDir, "routes", "stream\nother\n")

	tests := []struct {
		name            string
		credentialsPath string
		routesPath      string
		assertFunc      func(t *testing.T, got dict.Dictionary)
	}{
		{
			name:            "custom credentials and routes",
			credentialsPath: customCredsPath,
			routesPath:      customRoutesPath,
			assertFunc: func(t *testing.T, got dict.Dictionary) {
				assert.Equal(t, []string{"alice"}, got.Usernames())
				assert.Equal(t, []string{"secret"}, got.Passwords())
				assert.Equal(t, []string{"stream", "other"}, got.Routes())
			},
		},
		{
			name:            "custom credentials default routes",
			credentialsPath: customCredsPath,
			assertFunc: func(t *testing.T, got dict.Dictionary) {
				assert.Equal(t, []string{"alice"}, got.Usernames())
				assert.Equal(t, []string{"secret"}, got.Passwords())
				assert.NotEmpty(t, got.Routes())
				assert.Contains(t, got.Routes(), "stream")
			},
		},
		{
			name:       "default credentials custom routes",
			routesPath: customRoutesPath,
			assertFunc: func(t *testing.T, got dict.Dictionary) {
				assert.NotEmpty(t, got.Usernames())
				assert.Contains(t, got.Usernames(), "admin")
				assert.NotEmpty(t, got.Passwords())
				assert.Contains(t, got.Passwords(), "admin")
				assert.Equal(t, []string{"stream", "other"}, got.Routes())
			},
		},
		{
			name:            "whitespace paths use defaults",
			credentialsPath: "  \t\n",
			routesPath:      "\n\t",
			assertFunc: func(t *testing.T, got dict.Dictionary) {
				assert.NotEmpty(t, got.Usernames())
				assert.Contains(t, got.Usernames(), "admin")
				assert.NotEmpty(t, got.Passwords())
				assert.Contains(t, got.Passwords(), "admin")
				assert.NotEmpty(t, got.Routes())
				assert.Contains(t, got.Routes(), "stream")
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := dict.New(test.credentialsPath, test.routesPath)
			require.NoError(t, err)
			test.assertFunc(t, got)
		})
	}
}

func TestNew_Errors(t *testing.T) {
	tempDir := t.TempDir()
	validCredsPath := writeTempFile(t, tempDir, "creds.json", `{"usernames":["alice"],"passwords":["secret"]}`)
	validRoutesPath := writeTempFile(t, tempDir, "routes", "stream\n")
	invalidJSONPath := writeTempFile(t, tempDir, "invalid.json", "{")
	emptyCredsPath := writeTempFile(t, tempDir, "empty.json", "")
	longRoute := strings.Repeat("a", bufio.MaxScanTokenSize+1)
	tooLongRoutesPath := writeTempFile(t, tempDir, "routes-too-long", longRoute)

	tests := []struct {
		name            string
		credentialsPath string
		routesPath      string
		wantErrContains string
		wantErrIs       error
	}{
		{
			name:            "missing credentials file",
			credentialsPath: filepath.Join(tempDir, "missing.json"),
			routesPath:      validRoutesPath,
			wantErrContains: "reading credentials dictionary",
		},
		{
			name:            "invalid credentials json",
			credentialsPath: invalidJSONPath,
			routesPath:      validRoutesPath,
			wantErrContains: "reading dictionary contents",
		},
		{
			name:            "empty credentials file",
			credentialsPath: emptyCredsPath,
			routesPath:      validRoutesPath,
			wantErrContains: "credentials dictionary is empty",
		},
		{
			name:            "missing routes file",
			credentialsPath: validCredsPath,
			routesPath:      filepath.Join(tempDir, "missing-routes"),
			wantErrContains: "opening routes dictionary",
		},
		{
			name:            "routes file too long",
			credentialsPath: validCredsPath,
			routesPath:      tooLongRoutesPath,
			wantErrIs:       bufio.ErrTooLong,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := dict.New(test.credentialsPath, test.routesPath)
			require.Error(t, err)

			if test.wantErrContains != "" {
				assert.ErrorContains(t, err, test.wantErrContains)
			}
			if test.wantErrIs != nil {
				assert.True(t, errors.Is(err, test.wantErrIs))
			}
		})
	}
}

func TestNew_RoutesSanitization(t *testing.T) {
	tempDir := t.TempDir()
	credsPath := writeTempFile(t, tempDir, "creds.json", `{"usernames":["alice"],"passwords":["secret"]}`)

	tests := []struct {
		name    string
		content string
		want    []string
	}{
		{
			name:    "strips BOM from first route",
			content: "\ufeffstream\nother\n",
			want:    []string{"stream", "other"},
		},
		{
			name:    "skips leading, trailing and whitespace-only lines",
			content: "\n\nstream\n   \nother\n\n",
			want:    []string{"stream", "other"},
		},
		{
			name:    "trims surrounding whitespace per line",
			content: "  stream  \n\tother\t\n",
			want:    []string{"stream", "other"},
		},
		{
			name:    "handles CRLF line endings",
			content: "stream\r\nother\r\n",
			want:    []string{"stream", "other"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			routesPath := writeTempFile(t, tempDir, "routes-"+test.name, test.content)

			got, err := dict.New(credsPath, routesPath)
			require.NoError(t, err)

			assert.Equal(t, test.want, got.Routes())
		})
	}
}

func TestNew_CredentialsSanitization(t *testing.T) {
	tempDir := t.TempDir()
	routesPath := writeTempFile(t, tempDir, "routes", "stream\n")

	tests := []struct {
		name            string
		content         string
		wantErr         require.ErrorAssertionFunc
		wantErrContains string
	}{
		{
			name:    "valid credentials parse successfully",
			content: `{"usernames":["alice"],"passwords":["secret"]}`,
			wantErr: require.NoError,
		},
		{
			name:            "rejects CRLF in username",
			content:         `{"usernames":["admin\r\nInjected: header"],"passwords":["secret"]}`,
			wantErr:         require.Error,
			wantErrContains: "control characters are not allowed",
		},
		{
			name:            "rejects control character in password",
			content:         `{"usernames":["alice"],"passwords":["sec\u0001ret"]}`,
			wantErr:         require.Error,
			wantErrContains: "control characters are not allowed",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			credsPath := writeTempFile(t, tempDir, "creds-"+test.name+".json", test.content)

			_, err := dict.New(credsPath, routesPath)
			test.wantErr(t, err)
			if test.wantErrContains != "" {
				assert.ErrorContains(t, err, test.wantErrContains)
			}
		})
	}
}

func writeTempFile(t *testing.T, dir, name, content string) string {
	t.Helper()
	path := filepath.Join(dir, name)
	require.NoError(t, os.WriteFile(path, []byte(content), 0o600))
	return path
}
