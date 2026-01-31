package attack

import (
	"net/netip"
	"testing"
	"time"

	"github.com/Ullaakut/cameradar/v6"
	"github.com/Ullaakut/cameradar/v6/internal/ui"
	"github.com/bluenviron/gortsplib/v5"
	"github.com/bluenviron/gortsplib/v5/pkg/base"
	"github.com/bluenviron/gortsplib/v5/pkg/headers"
	"github.com/bluenviron/gortsplib/v5/pkg/liberrors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDetectAuthMethod_UnauthorizedWithoutResponseDoesNotError(t *testing.T) {
	originalDescribe := describeRTSP
	t.Cleanup(func() {
		describeRTSP = originalDescribe
	})
	describeRTSP = func(*gortsplib.Client, *base.URL) (*base.Response, error) {
		return nil, liberrors.ErrClientBadStatusCode{Code: base.StatusUnauthorized}
	}

	attacker, err := New(stubDictionary{}, 0, time.Second, ui.NopReporter{})
	require.NoError(t, err)

	stream := cameradar.Stream{
		Address: netip.MustParseAddr("127.0.0.1"),
		Port:    8554,
		Routes:  []string{"stream"},
	}

	got, err := attacker.detectAuthMethod(t.Context(), stream)
	assert.Equal(t, cameradar.AuthUnknown, got.AuthenticationType)
	assert.Equal(t, cameradar.AuthUnknown, got.AuthenticationType)
}

func TestDetectAuthMethod_SetsAuthNoneOnSuccess(t *testing.T) {
	originalDescribe := describeRTSP
	t.Cleanup(func() {
		describeRTSP = originalDescribe
	})
	describeRTSP = func(*gortsplib.Client, *base.URL) (*base.Response, error) {
		return &base.Response{StatusCode: base.StatusOK}, nil
	}

	attacker, err := New(stubDictionary{}, 0, time.Second, ui.NopReporter{})
	require.NoError(t, err)

	stream := cameradar.Stream{
		Address: netip.MustParseAddr("127.0.0.1"),
		Port:    8554,
		Routes:  []string{"stream"},
	}

	got, err := attacker.detectAuthMethod(t.Context(), stream)
	require.NoError(t, err)
	assert.Equal(t, cameradar.AuthNone, got.AuthenticationType)
}

func TestDetectAuthMethod_SetsAuthFromWWWAuthenticate(t *testing.T) {
	originalDescribe := describeRTSP
	t.Cleanup(func() {
		describeRTSP = originalDescribe
	})

	attacker, err := New(stubDictionary{}, 0, time.Second, ui.NopReporter{})
	require.NoError(t, err)

	stream := cameradar.Stream{
		Address: netip.MustParseAddr("127.0.0.1"),
		Port:    8554,
		Routes:  []string{"stream"},
	}

	res := &base.Response{StatusCode: base.StatusUnauthorized}

	res.Header = base.Header{
		"WWW-Authenticate": headers.Authenticate{Method: headers.AuthMethodBasic, Realm: "cameradar"}.Marshal(),
	}
	describeRTSP = func(*gortsplib.Client, *base.URL) (*base.Response, error) {
		return res, liberrors.ErrClientBadStatusCode{Code: base.StatusUnauthorized}
	}

	got, err := attacker.detectAuthMethod(t.Context(), stream)
	require.NoError(t, err)
	assert.Equal(t, cameradar.AuthBasic, got.AuthenticationType)

	res.Header = base.Header{
		"WWW-Authenticate": headers.Authenticate{Method: headers.AuthMethodDigest, Realm: "cameradar", Nonce: "nonce"}.Marshal(),
	}

	got, err = attacker.detectAuthMethod(t.Context(), stream)
	require.NoError(t, err)
	assert.Equal(t, cameradar.AuthDigest, got.AuthenticationType)
}

type stubDictionary struct{}

func (stubDictionary) Routes() []string { return nil }

func (stubDictionary) Usernames() []string { return nil }

func (stubDictionary) Passwords() []string { return nil }
