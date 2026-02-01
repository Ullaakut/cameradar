package attack

import (
	"context"
	"errors"
	"fmt"

	"github.com/Ullaakut/cameradar/v6"
	"github.com/bluenviron/gortsplib/v5/pkg/base"
	"github.com/bluenviron/gortsplib/v5/pkg/liberrors"
)

func (a Attacker) detectAuthMethods(ctx context.Context, targets []cameradar.Stream) ([]cameradar.Stream, error) {
	streams, err := runParallel(ctx, targets, a.detectAuthMethod)
	if err != nil {
		return streams, err
	}

	for i := range streams {
		a.reporter.Progress(cameradar.StepDetectAuth, cameradar.ProgressTickMessage())

		var authMethod string
		switch streams[i].AuthenticationType {
		case cameradar.AuthNone:
			authMethod = "no"
		case cameradar.AuthBasic:
			authMethod = "basic"
		case cameradar.AuthDigest:
			authMethod = "digest"
		case cameradar.AuthUnknown:
			authMethod = "unknown"
		default:
			authMethod = fmt.Sprintf("unknown (%d)", streams[i].AuthenticationType)
		}

		a.reporter.Progress(cameradar.StepDetectAuth, fmt.Sprintf("Detected %s authentication for %s:%d", authMethod, streams[i].Address.String(), streams[i].Port))
	}

	return streams, nil
}

func (a Attacker) detectAuthMethod(ctx context.Context, stream cameradar.Stream) (cameradar.Stream, error) {
	if ctx.Err() != nil {
		return stream, ctx.Err()
	}
	u, urlStr, err := buildRTSPURL(stream, stream.Route(), "", "")
	if err != nil {
		return stream, fmt.Errorf("building rtsp url: %w", err)
	}

	client, err := a.newRTSPClient(u)
	if err != nil {
		return stream, fmt.Errorf("starting rtsp client: %w", err)
	}
	defer client.Close()

	res, err := describeRTSP(client, u)
	if err == nil {
		if res != nil {
			a.reporter.Debug(cameradar.StepDetectAuth, fmt.Sprintf("DESCRIBE %s RTSP/1.0 > %d", urlStr, res.StatusCode))
		}

		stream.AuthenticationType = cameradar.AuthNone
		return stream, nil
	}

	return a.handleDetectAuthError(stream, urlStr, res, err)
}

func (a Attacker) handleDetectAuthError(stream cameradar.Stream, urlStr string, res *base.Response, err error) (cameradar.Stream, error) {
	var badStatus liberrors.ErrClientBadStatusCode
	if !errors.As(err, &badStatus) || badStatus.Code != base.StatusUnauthorized {
		a.reporter.Debug(cameradar.StepDetectAuth, fmt.Sprintf("DESCRIBE %s RTSP/1.0 > error: %v", urlStr, err))
		stream.AuthenticationType = cameradar.AuthUnknown
		return stream, fmt.Errorf("performing describe request at %q: %w", urlStr, err)
	}

	stream.AuthenticationType = cameradar.AuthUnknown
	if res == nil {
		a.reporter.Debug(cameradar.StepDetectAuth, fmt.Sprintf("DESCRIBE %s RTSP/1.0 > %d", urlStr, badStatus.Code))
		return stream, nil
	}

	stream.AuthenticationType = authTypeFromHeaders(res.Header["WWW-Authenticate"])
	a.reporter.Debug(cameradar.StepDetectAuth, fmt.Sprintf("DESCRIBE %s RTSP/1.0 > %d", urlStr, badStatus.Code))
	a.reporter.Debug(cameradar.StepDetectAuth, "WWW-Authenticate header value is "+fmt.Sprint(res.Header["WWW-Authenticate"]))

	return stream, nil
}
