package attack

import (
	"context"
	"fmt"

	"github.com/Ullaakut/cameradar/v6"
	"github.com/bluenviron/gortsplib/v5/pkg/base"
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

	statusCode, headers, err := a.probeDescribeHeaders(ctx, u, urlStr)
	if err != nil {
		a.reporter.Debug(cameradar.StepDetectAuth, fmt.Sprintf("DESCRIBE %s RTSP/1.0 > error: %v", urlStr, err))
		stream.AuthenticationType = cameradar.AuthUnknown
		return stream, fmt.Errorf("performing describe request at %q: %w", urlStr, err)
	}

	a.reporter.Debug(cameradar.StepDetectAuth, fmt.Sprintf("DESCRIBE %s RTSP/1.0 > %d", urlStr, statusCode))
	values := headerValues(headers, "WWW-Authenticate")
	switch statusCode {
	case base.StatusOK:
		stream.AuthenticationType = cameradar.AuthNone
	case base.StatusUnauthorized:
		stream.AuthenticationType = authTypeFromHeaders(values)
	default:
		stream.AuthenticationType = cameradar.AuthUnknown
	}

	return stream, nil
}
