package attack_test

import (
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/Ullaakut/cameradar/v6"
	"github.com/Ullaakut/cameradar/v6/internal/attack"
	"github.com/Ullaakut/cameradar/v6/internal/ui"
	"github.com/bluenviron/gortsplib/v5/pkg/base"
	"github.com/bluenviron/gortsplib/v5/pkg/headers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name    string
		dict    attack.Dictionary
		wantErr require.ErrorAssertionFunc
	}{
		{
			name:    "rejects nil dictionary",
			dict:    nil,
			wantErr: require.Error,
		},
		{
			name: "accepts dictionary",
			dict: testDictionary{
				routes:    []string{"stream"},
				usernames: []string{"user"},
				passwords: []string{"pass"},
			},
			wantErr: require.NoError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			attacker, err := attack.New(test.dict, 10*time.Millisecond, time.Second, ui.NopReporter{})
			test.wantErr(t, err)
			if err != nil {
				assert.NotNil(t, attacker)
			}
		})
	}
}

func TestAttacker_Attack_BasicAuth(t *testing.T) {
	addr, port := startRTSPServer(t, rtspServerConfig{
		allowedRoute: "stream",
		requireAuth:  true,
		username:     "user",
		password:     "pass",
		authMethod:   headers.AuthMethodBasic,
	})

	dict := testDictionary{
		routes:    []string{"stream"},
		usernames: []string{"user", "other"},
		passwords: []string{"pass", "bad"},
	}

	testInterval := time.Millisecond
	testRequestTimeout := time.Second
	attacker, err := attack.New(dict, testInterval, testRequestTimeout, ui.NopReporter{})
	require.NoError(t, err)

	streams := []cameradar.Stream{{
		Address: addr,
		Port:    port,
	}}

	got, err := attacker.Attack(t.Context(), streams)
	require.NoError(t, err)
	require.Len(t, got, 1)

	assert.True(t, got[0].RouteFound)
	assert.True(t, got[0].CredentialsFound)
	assert.True(t, got[0].Available)
	assert.Equal(t, cameradar.AuthBasic, got[0].AuthenticationType)
	assert.Equal(t, "user", got[0].Username)
	assert.Equal(t, "pass", got[0].Password)
	assert.Contains(t, got[0].Routes, "stream")
}

func TestAttacker_Attack_AuthVariants(t *testing.T) {
	tests := []struct {
		name         string
		config       rtspServerConfig
		dict         testDictionary
		wantAuthType cameradar.AuthType
		wantRoute    bool
		wantCreds    bool
		wantAvail    bool
		wantErr      require.ErrorAssertionFunc
		errContains  string
	}{
		{
			name: "no authentication",
			config: rtspServerConfig{
				allowedRoute: "stream",
				requireAuth:  false,
				authMethod:   headers.AuthMethodBasic,
			},
			dict: testDictionary{
				routes: []string{"stream"},
			},
			wantAuthType: cameradar.AuthNone,
			wantRoute:    true,
			wantCreds:    false,
			wantAvail:    true,
			wantErr:      require.NoError,
		},
		{
			name: "digest authentication",
			config: rtspServerConfig{
				allowedRoute: "stream",
				requireAuth:  true,
				username:     "user",
				password:     "pass",
				authMethod:   headers.AuthMethodDigest,
			},
			dict: testDictionary{
				routes:    []string{"stream"},
				usernames: []string{"user"},
				passwords: []string{"pass"},
			},
			wantAuthType: cameradar.AuthDigest,
			wantRoute:    true,
			wantCreds:    true,
			wantAvail:    true,
			wantErr:      require.NoError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			addr, port := startRTSPServer(t, test.config)

			attacker, err := attack.New(test.dict, 0, time.Second, ui.NopReporter{})
			require.NoError(t, err)

			streams := []cameradar.Stream{{
				Address: addr,
				Port:    port,
			}}

			got, err := attacker.Attack(t.Context(), streams)
			test.wantErr(t, err)

			if test.errContains != "" {
				assert.ErrorContains(t, err, test.errContains)
			}

			require.Len(t, got, 1)
			assert.Equal(t, test.wantAuthType, got[0].AuthenticationType)
			assert.Equal(t, test.wantRoute, got[0].RouteFound)
			assert.Equal(t, test.wantCreds, got[0].CredentialsFound)
			assert.Equal(t, test.wantAvail, got[0].Available)
		})
	}
}

func TestAttacker_Attack_ValidationErrors(t *testing.T) {
	attacker, err := attack.New(testDictionary{routes: []string{"stream"}}, 0, time.Second, ui.NopReporter{})
	require.NoError(t, err)

	tests := []struct {
		name     string
		attacker attack.Attacker
		targets  []cameradar.Stream
		wantErr  string
	}{
		{
			name:     "fails with no targets",
			attacker: attacker,
			targets:  nil,
			wantErr:  "no stream found",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := test.attacker.Attack(t.Context(), test.targets)
			require.Error(t, err)
			assert.ErrorContains(t, err, test.wantErr)
		})
	}
}

func TestAttacker_Attack_ReturnsErrorWhenRouteMissing(t *testing.T) {
	addr, port := startRTSPServer(t, rtspServerConfig{
		allowedRoute: "stream",
		requireAuth:  false,
		authMethod:   headers.AuthMethodBasic,
	})

	dict := testDictionary{
		routes:    []string{"missing"},
		usernames: []string{"user"},
		passwords: []string{"pass"},
	}

	attacker, err := attack.New(dict, 0, time.Second, ui.NopReporter{})
	require.NoError(t, err)

	streams := []cameradar.Stream{{
		Address: addr,
		Port:    port,
	}}

	got, err := attacker.Attack(t.Context(), streams)
	require.Error(t, err)
	assert.ErrorContains(t, err, "validating streams")
	require.Len(t, got, 1)
	assert.False(t, got[0].RouteFound)
}

func TestAttacker_Attack_ReturnsErrorWhenCredentialsMissing(t *testing.T) {
	addr, port := startRTSPServer(t, rtspServerConfig{
		allowedRoute: "stream",
		requireAuth:  true,
		username:     "user",
		password:     "pass",
		authMethod:   headers.AuthMethodBasic,
	})

	dict := testDictionary{
		routes:    []string{"stream"},
		usernames: []string{"user"},
		passwords: []string{"wrong"},
	}

	attacker, err := attack.New(dict, 0, time.Second, ui.NopReporter{})
	require.NoError(t, err)

	streams := []cameradar.Stream{{
		Address: addr,
		Port:    port,
	}}

	got, err := attacker.Attack(t.Context(), streams)
	require.Error(t, err)
	assert.ErrorContains(t, err, "validating streams")
	require.Len(t, got, 1)
	assert.Equal(t, cameradar.AuthBasic, got[0].AuthenticationType)
	assert.False(t, got[0].CredentialsFound)
}

func TestAttacker_Attack_CredentialAttemptFails(t *testing.T) {
	reporter := &recordingReporter{}

	addr, port := startRTSPServer(t, rtspServerConfig{
		allowedRoute: "stream",
		requireAuth:  true,
		username:     "user",
		password:     "pass",
		authMethod:   headers.AuthMethodBasic,
		failOnAuth:   true,
	})

	dict := testDictionary{
		routes:    []string{"stream"},
		usernames: []string{"user"},
		passwords: []string{"pass"},
	}

	attacker, err := attack.New(dict, 0, time.Second, reporter)
	require.NoError(t, err)

	streams := []cameradar.Stream{{
		Address: addr,
		Port:    port,
	}}

	got, err := attacker.Attack(t.Context(), streams)
	require.Error(t, err)
	assert.ErrorContains(t, err, "validating streams")
	require.Len(t, got, 1)
	assert.False(t, got[0].CredentialsFound)
}

func TestAttacker_Attack_AllowsDummyRoute(t *testing.T) {
	addr, port := startRTSPServer(t, rtspServerConfig{
		allowAll:    true,
		requireAuth: false,
		authMethod:  headers.AuthMethodBasic,
	})

	dict := testDictionary{}

	attacker, err := attack.New(dict, 0, time.Second, ui.NopReporter{})
	require.NoError(t, err)

	streams := []cameradar.Stream{{
		Address: addr,
		Port:    port,
	}}

	got, err := attacker.Attack(t.Context(), streams)
	require.NoError(t, err)
	require.Len(t, got, 1)
	assert.True(t, got[0].RouteFound)
	assert.Equal(t, []string{""}, got[0].Routes)
	assert.True(t, got[0].Available)
}

func TestAttacker_Attack_ValidationFailsWhenSetupErrors(t *testing.T) {
	addr, port := startRTSPServer(t, rtspServerConfig{
		allowedRoute: "stream",
		requireAuth:  false,
		authMethod:   headers.AuthMethodBasic,
		setupStatus:  base.StatusUnsupportedTransport,
	})

	dict := testDictionary{
		routes: []string{"stream"},
	}

	attacker, err := attack.New(dict, 0, time.Second, ui.NopReporter{})
	require.NoError(t, err)

	streams := []cameradar.Stream{{
		Address: addr,
		Port:    port,
	}}

	got, err := attacker.Attack(t.Context(), streams)
	require.NoError(t, err)
	require.Len(t, got, 1)
	assert.False(t, got[0].Available)
	assert.True(t, got[0].RouteFound)
}

type testDictionary struct {
	routes    []string
	usernames []string
	passwords []string
}

func (d testDictionary) Routes() []string {
	return d.routes
}

func (d testDictionary) Usernames() []string {
	return d.usernames
}

func (d testDictionary) Passwords() []string {
	return d.passwords
}

type recordingReporter struct {
	mu            sync.Mutex
	debugMessages []string
}

func (r *recordingReporter) Start(cameradar.Step, string) {}

func (r *recordingReporter) Done(cameradar.Step, string) {}

func (r *recordingReporter) Progress(cameradar.Step, string) {}

func (r *recordingReporter) Debug(_ cameradar.Step, message string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.debugMessages = append(r.debugMessages, message)
}

func (r *recordingReporter) Error(cameradar.Step, error) {}

func (r *recordingReporter) Summary([]cameradar.Stream, error) {}

func (r *recordingReporter) Close() {}

func (r *recordingReporter) HasDebugContaining(value string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, message := range r.debugMessages {
		if strings.Contains(message, value) {
			return true
		}
	}
	return false
}
