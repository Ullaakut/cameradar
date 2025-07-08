package attack_test

import (
	"errors"
	"net"
	"net/netip"
	"strings"
	"testing"

	"github.com/bluenviron/gortsplib/v5"
	"github.com/bluenviron/gortsplib/v5/pkg/auth"
	"github.com/bluenviron/gortsplib/v5/pkg/base"
	"github.com/bluenviron/gortsplib/v5/pkg/description"
	"github.com/bluenviron/gortsplib/v5/pkg/format"
	"github.com/bluenviron/gortsplib/v5/pkg/headers"
	"github.com/bluenviron/gortsplib/v5/pkg/liberrors"
	"github.com/stretchr/testify/require"
)

type rtspServerConfig struct {
	allowAll     bool
	allowedRoute string
	requireAuth  bool
	username     string
	password     string
	authMethod   headers.AuthMethod
	authHeader   base.HeaderValue
	failOnAuth   bool
	setupStatus  base.StatusCode
}

type testServerHandler struct {
	stream       *gortsplib.ServerStream
	allowAll     bool
	allowedRoute string
	requireAuth  bool
	username     string
	password     string
	authHeader   base.HeaderValue
	failOnAuth   bool
	setupStatus  base.StatusCode
}

func (h *testServerHandler) OnDescribe(ctx *gortsplib.ServerHandlerOnDescribeCtx) (*base.Response, *gortsplib.ServerStream, error) {
	if !h.routeAllowed(ctx.Path) {
		return &base.Response{StatusCode: base.StatusNotFound}, nil, nil
	}

	if h.failOnAuth && len(ctx.Request.Header["Authorization"]) > 0 {
		return &base.Response{StatusCode: base.StatusBadRequest}, nil, errors.New("forced auth failure")
	}

	if h.requireAuth && !ctx.Conn.VerifyCredentials(ctx.Request, h.username, h.password) {
		return &base.Response{
			StatusCode: base.StatusUnauthorized,
			Header: base.Header{
				"WWW-Authenticate": h.authHeader,
			},
		}, nil, liberrors.ErrServerAuth{}
	}

	return &base.Response{StatusCode: base.StatusOK}, h.stream, nil
}

func (h *testServerHandler) OnSetup(ctx *gortsplib.ServerHandlerOnSetupCtx) (*base.Response, *gortsplib.ServerStream, error) {
	if !h.routeAllowed(ctx.Path) {
		return &base.Response{StatusCode: base.StatusNotFound}, nil, nil
	}

	if h.requireAuth && !ctx.Conn.VerifyCredentials(ctx.Request, h.username, h.password) {
		return &base.Response{
			StatusCode: base.StatusUnauthorized,
			Header: base.Header{
				"WWW-Authenticate": h.authHeader,
			},
		}, nil, liberrors.ErrServerAuth{}
	}

	status := base.StatusOK
	if h.setupStatus != 0 {
		status = h.setupStatus
	}

	return &base.Response{StatusCode: status}, h.stream, nil
}

func (h *testServerHandler) routeAllowed(path string) bool {
	path = strings.TrimLeft(path, "/")
	return h.allowAll || path == h.allowedRoute
}

func startRTSPServer(t *testing.T, cfg rtspServerConfig) (netip.Addr, uint16) {
	t.Helper()

	handler := &testServerHandler{
		allowAll:     cfg.allowAll,
		allowedRoute: cfg.allowedRoute,
		requireAuth:  cfg.requireAuth,
		username:     cfg.username,
		password:     cfg.password,
		failOnAuth:   cfg.failOnAuth,
		setupStatus:  cfg.setupStatus,
	}

	if len(cfg.authHeader) > 0 {
		handler.authHeader = cfg.authHeader
	} else {
		authHeader := headers.Authenticate{
			Method: cfg.authMethod,
			Realm:  "cameradar",
		}
		if cfg.authMethod == headers.AuthMethodDigest {
			authHeader.Nonce = "nonce"
		}
		handler.authHeader = authHeader.Marshal()
	}

	server := &gortsplib.Server{
		Handler:     handler,
		RTSPAddress: "127.0.0.1:0",
		AuthMethods: authMethods(cfg.authMethod),
	}

	err := server.Start()
	require.NoError(t, err)
	t.Cleanup(server.Close)

	desc := &description.Session{
		Medias: []*description.Media{{
			Type: description.MediaTypeVideo,
			Formats: []format.Format{&format.H264{
				PayloadTyp:        96,
				PacketizationMode: 1,
			}},
		}},
	}

	stream := &gortsplib.ServerStream{
		Server: server,
		Desc:   desc,
	}
	err = stream.Initialize()
	require.NoError(t, err)
	t.Cleanup(stream.Close)

	handler.stream = stream

	listener := server.NetListener()
	require.NotNil(t, listener)

	tcpAddr, ok := listener.Addr().(*net.TCPAddr)
	require.True(t, ok)

	return netip.MustParseAddr("127.0.0.1"), uint16(tcpAddr.Port)
}

func authMethods(method headers.AuthMethod) []auth.VerifyMethod {
	switch method {
	case headers.AuthMethodDigest:
		return []auth.VerifyMethod{auth.VerifyMethodDigestMD5}
	case headers.AuthMethodBasic:
		return []auth.VerifyMethod{auth.VerifyMethodBasic}
	default:
		return nil
	}
}
