package attack_test

import (
	"errors"
	"net"
	"net/netip"
	"strings"
	"testing"
	"time"

	"github.com/bluenviron/gortsplib/v5"
	"github.com/bluenviron/gortsplib/v5/pkg/auth"
	"github.com/bluenviron/gortsplib/v5/pkg/base"
	"github.com/bluenviron/gortsplib/v5/pkg/description"
	"github.com/bluenviron/gortsplib/v5/pkg/format"
	"github.com/bluenviron/gortsplib/v5/pkg/headers"
	"github.com/bluenviron/gortsplib/v5/pkg/liberrors"
	"github.com/pion/rtp"
	"github.com/stretchr/testify/require"
)

type rtspServerConfig struct {
	allowAll                  bool
	describeAllowAll          bool
	allowedRoute              string
	requireAuth               bool
	describeIgnoreAuth        bool
	describeAcceptInvalidAuth bool
	username                  string
	password                  string
	authMethod                headers.AuthMethod
	authHeader                base.HeaderValue
	failOnAuth                bool
	setupStatus               base.StatusCode
	playStatus                base.StatusCode
	sendFrames                bool
}

type testServerHandler struct {
	stream                    *gortsplib.ServerStream
	allowAll                  bool
	describeAllowAll          bool
	allowedRoute              string
	requireAuth               bool
	describeIgnoreAuth        bool
	describeAcceptInvalidAuth bool
	username                  string
	password                  string
	authHeader                base.HeaderValue
	failOnAuth                bool
	setupStatus               base.StatusCode
	playStatus                base.StatusCode
	sendFrames                bool
}

func (h *testServerHandler) OnDescribe(ctx *gortsplib.ServerHandlerOnDescribeCtx) (*base.Response, *gortsplib.ServerStream, error) {
	if !h.describeRouteAllowed(ctx.Path) {
		return &base.Response{StatusCode: base.StatusNotFound}, nil, nil
	}

	if h.failOnAuth && len(ctx.Request.Header["Authorization"]) > 0 {
		return &base.Response{StatusCode: base.StatusBadRequest}, nil, errors.New("forced auth failure")
	}

	if h.requireAuth && !ctx.Conn.VerifyCredentials(ctx.Request, h.username, h.password) {
		authorization := ctx.Request.Header["Authorization"]
		if h.describeIgnoreAuth || (h.describeAcceptInvalidAuth && len(authorization) > 0) {
			return &base.Response{StatusCode: base.StatusOK}, h.stream, nil
		}

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

func (h *testServerHandler) OnPlay(ctx *gortsplib.ServerHandlerOnPlayCtx) (*base.Response, error) {
	if !h.routeAllowed(ctx.Path) {
		return &base.Response{StatusCode: base.StatusNotFound}, nil
	}

	if h.requireAuth && !ctx.Conn.VerifyCredentials(ctx.Request, h.username, h.password) {
		return &base.Response{
			StatusCode: base.StatusUnauthorized,
			Header: base.Header{
				"WWW-Authenticate": h.authHeader,
			},
		}, liberrors.ErrServerAuth{}
	}

	status := base.StatusOK
	if h.playStatus != 0 {
		status = h.playStatus
	}

	if status == base.StatusOK && h.sendFrames {
		h.emitFrame()
	}

	return &base.Response{StatusCode: status}, nil
}

func (h *testServerHandler) routeAllowed(path string) bool {
	path = strings.TrimLeft(path, "/")
	return h.allowAll || path == h.allowedRoute
}

func (h *testServerHandler) describeRouteAllowed(path string) bool {
	if h.describeAllowAll {
		return true
	}

	return h.routeAllowed(path)
}

func (h *testServerHandler) emitFrame() {
	if h.stream == nil || h.stream.Desc == nil || len(h.stream.Desc.Medias) == 0 {
		return
	}

	media := h.stream.Desc.Medias[0]
	go func() {
		time.Sleep(20 * time.Millisecond)
		_ = h.stream.WritePacketRTP(media, &rtp.Packet{
			Header: rtp.Header{
				Version:        2,
				PayloadType:    96,
				SequenceNumber: 1,
				Timestamp:      90_000,
				SSRC:           1,
			},
			Payload: []byte{0x05, 0x01},
		})
	}()
}

func startRTSPServer(t *testing.T, cfg rtspServerConfig) (netip.Addr, uint16) {
	t.Helper()

	handler := &testServerHandler{
		allowAll:                  cfg.allowAll,
		describeAllowAll:          cfg.describeAllowAll,
		allowedRoute:              cfg.allowedRoute,
		requireAuth:               cfg.requireAuth,
		describeIgnoreAuth:        cfg.describeIgnoreAuth,
		describeAcceptInvalidAuth: cfg.describeAcceptInvalidAuth,
		username:                  cfg.username,
		password:                  cfg.password,
		failOnAuth:                cfg.failOnAuth,
		setupStatus:               cfg.setupStatus,
		playStatus:                cfg.playStatus,
		sendFrames:                cfg.sendFrames,
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
