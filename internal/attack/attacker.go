package attack

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Ullaakut/cameradar/v6"
	"github.com/bluenviron/gortsplib/v5"
	"github.com/bluenviron/gortsplib/v5/pkg/base"
	"github.com/bluenviron/gortsplib/v5/pkg/description"
	"github.com/bluenviron/gortsplib/v5/pkg/format"
	"github.com/bluenviron/gortsplib/v5/pkg/liberrors"
	"github.com/pion/rtp"
)

// Route that should never be a constructor default.
const dummyRoute = "0x8b6c42"

var errFrameProbeNoMedias = errors.New("describe succeeded but no media tracks found")

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
	framecheck     bool
}

// New builds an Attacker with the provided dependencies.
func New(dict Dictionary, attackInterval, timeout time.Duration, framecheck bool, reporter Reporter) (Attacker, error) {
	if dict == nil {
		return Attacker{}, errors.New("dictionary is required")
	}

	return Attacker{
		dictionary:     dict,
		attackInterval: attackInterval,
		timeout:        timeout,
		framecheck:     framecheck,
		reporter:       reporter,
	}, nil
}

// Attack attacks the given targets and returns the accessed streams.
func (a Attacker) Attack(ctx context.Context, targets []cameradar.Stream) ([]cameradar.Stream, error) {
	if len(targets) == 0 {
		return nil, errors.New("no stream found")
	}

	// Each phase processes every target even when one camera errors, so a
	// single unreachable host cannot drop results for the rest of the batch.
	// The first non-cancellation error is remembered and surfaced after all
	// phases have run, while healthy cameras still progress to validation.
	var firstErr error
	record := func(err error) {
		if err != nil && firstErr == nil {
			firstErr = err
		}
	}

	streams, err := a.attackRoutesPhase(ctx, targets)
	record(err)
	if ctx.Err() != nil {
		return streams, ctx.Err()
	}

	streams, err = a.detectAuthPhase(ctx, streams)
	record(err)
	if ctx.Err() != nil {
		return streams, ctx.Err()
	}

	streams, err = a.attackCredentialsPhase(ctx, streams)
	record(err)
	if ctx.Err() != nil {
		return streams, ctx.Err()
	}

	streams, err = a.validateStreamsPhase(ctx, streams)
	record(err)
	if ctx.Err() != nil {
		return streams, ctx.Err()
	}

	// Some cameras run an inaccurate version of the RTSP protocol which prioritizes 401 over 404.
	// For these cameras, running another route attack solves the problem.
	if needsReattack(streams) {
		streams, err = a.reattackRoutes(ctx, streams)
		record(err)
	}

	return streams, firstErr
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
			// This stream is fully discovered, no need to re-attack.
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
			ok, err := a.credAttack(ctx, target, username, password)
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

				return target, nil
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
	ok, err := a.routeAttack(ctx, target, dummyRoute)
	if err != nil {
		a.reporter.Debug(cameradar.StepAttackRoutes, fmt.Sprintf("route probe failed for %s:%d: %v", target.Address.String(), target.Port, err))
		return target, nil
	}
	if ok {
		target.RouteFound = true
		target.Routes = append(target.Routes, "") // Add empty route for default.
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
		ok, err := a.routeAttack(ctx, target, route)
		if err != nil {
			a.reporter.Debug(cameradar.StepAttackRoutes, fmt.Sprintf("route attempt failed for %s:%d (%s): %v", target.Address.String(), target.Port, route, err))
			return target, nil
		}
		if ok {
			target.RouteFound = true
			target.Routes = append(target.Routes, route)
			a.reporter.Progress(cameradar.StepAttackRoutes, fmt.Sprintf("Route found for %s:%d -> %s", target.Address.String(), target.Port, route))
		}
	}

	return target, nil
}

func (a Attacker) routeAttack(ctx context.Context, stream cameradar.Stream, route string) (bool, error) {
	stream.Routes = []string{route}
	code, err := a.describeStatus(stream)
	if err != nil {
		return false, fmt.Errorf("performing describe request at %q: %w", stream, err)
	}

	a.reporter.Debug(cameradar.StepAttackRoutes, fmt.Sprintf("DESCRIBE %s RTSP/1.0 > %d", stream, code))
	if code == base.StatusOK && !a.acceptStatusOK(ctx, stream, cameradar.StepAttackRoutes) {
		return false, nil
	}

	return code == base.StatusOK || code == base.StatusUnauthorized || code == base.StatusForbidden, nil
}

func (a Attacker) credAttack(ctx context.Context, stream cameradar.Stream, username, password string) (bool, error) {
	stream.Username = username
	stream.Password = password
	code, err := a.describeStatus(stream)
	if err != nil {
		return false, fmt.Errorf("performing describe request at %q: %w", stream, err)
	}

	a.reporter.Debug(cameradar.StepAttackCredentials, fmt.Sprintf("DESCRIBE %s RTSP/1.0 > %d", stream, code))
	if code == base.StatusOK && !a.acceptStatusOK(ctx, stream, cameradar.StepAttackCredentials) {
		return false, nil
	}

	return code == base.StatusOK || code == base.StatusNotFound, nil
}

func (a Attacker) acceptStatusOK(ctx context.Context, stream cameradar.Stream, step cameradar.Step) bool {
	if !a.framecheck {
		return true
	}

	ok, statusCode, err := a.probeFrameGeneration(ctx, stream, step)
	if err != nil {
		if errors.Is(err, errFrameProbeNoMedias) {
			a.reporter.Debug(step, fmt.Sprintf("Ignoring RTSP 200 for %s because DESCRIBE returned no media tracks", stream))
			return false
		}
		a.reporter.Debug(step, fmt.Sprintf("Frame probe failed for %s: %v", stream, err))
		return false
	}
	if !ok && step == cameradar.StepAttackRoutes && (statusCode == base.StatusUnauthorized || statusCode == base.StatusForbidden) {
		a.reporter.Debug(step, fmt.Sprintf("Keeping RTSP 200 route for %s because frame probe returned RTSP %d", stream, statusCode))
		return true
	}
	if ok {
		a.reporter.Debug(step, fmt.Sprintf("Frame probe succeeded for %s", stream))
	}
	if !ok {
		a.reporter.Debug(step, fmt.Sprintf("Ignoring RTSP 200 for %s because no RTP packet was received", stream))
	}

	return ok
}

func (a Attacker) probeFrameGeneration(ctx context.Context, stream cameradar.Stream, step cameradar.Step) (bool, base.StatusCode, error) {
	ok, statusCode, err := a.probeFrameGenerationWithProtocol(ctx, stream, nil, step)
	if ok || statusCode != 0 || err != nil {
		return ok, statusCode, err
	}

	// When UDP packets are blocked or not delivered, retry over interleaved TCP.
	tcpProtocol := gortsplib.ProtocolTCP
	return a.probeFrameGenerationWithProtocol(ctx, stream, &tcpProtocol, step)
}

//nolint:cyclop // Splitting this function does not make it clearer.
func (a Attacker) probeFrameGenerationWithProtocol(
	ctx context.Context,
	stream cameradar.Stream,
	protocol *gortsplib.Protocol,
	step cameradar.Step,
) (bool, base.StatusCode, error) {
	if ctx.Err() != nil {
		return false, 0, ctx.Err()
	}

	client, err := a.newRTSPClient(stream)
	if err != nil {
		return false, 0, fmt.Errorf("starting rtsp client: %w", err)
	}
	defer client.Close()

	if protocol != nil {
		client.Protocol = protocol
	}

	desc, _, err := a.describeWithRetry(ctx, client, stream, step)
	if err != nil {
		if noMediaSDPError(err) {
			return false, 0, errFrameProbeNoMedias
		}
		if code, ok := badStatusCode(err); ok {
			return false, code, nil
		}
		return false, 0, fmt.Errorf("performing describe request at %q: %w", stream, err)
	}

	if desc == nil || len(desc.Medias) == 0 {
		return false, 0, errFrameProbeNoMedias
	}

	err = client.SetupAll(desc.BaseURL, desc.Medias)
	if err != nil {
		if code, ok := badStatusCode(err); ok {
			return false, code, nil
		}
		return false, 0, fmt.Errorf("performing setup requests at %q: %w", stream, err)
	}

	rtpPacketReceived := make(chan struct{}, 1)
	client.OnPacketRTPAny(func(_ *description.Media, _ format.Format, _ *rtp.Packet) {
		select {
		case rtpPacketReceived <- struct{}{}:
		default:
		}
	})

	_, err = client.Play(nil)
	if err != nil {
		if code, ok := badStatusCode(err); ok {
			return false, code, nil
		}
		return false, 0, fmt.Errorf("performing play request at %q: %w", stream, err)
	}

	timer := time.NewTimer(frameProbeTimeout(a.timeout))
	defer timer.Stop()

	select {
	case <-ctx.Done():
		return false, 0, ctx.Err()
	case <-timer.C:
		return false, 0, nil
	case <-rtpPacketReceived:
		return true, 0, nil
	}
}

func badStatusCode(err error) (base.StatusCode, bool) {
	var badStatus liberrors.ErrClientBadStatusCode
	if !errors.As(err, &badStatus) {
		return 0, false
	}

	return badStatus.Code, true
}

func noMediaSDPError(err error) bool {
	var sdpErr liberrors.ErrClientSDPInvalid
	if !errors.As(err, &sdpErr) {
		return false
	}

	if sdpErr.Err == nil {
		return false
	}

	return strings.Contains(strings.ToLower(sdpErr.Err.Error()), "no media streams")
}

func frameProbeTimeout(timeout time.Duration) time.Duration {
	if timeout > 0 {
		return timeout
	}

	return 2 * time.Second
}

func (a Attacker) validateStream(ctx context.Context, stream cameradar.Stream, emitProgress bool) (cameradar.Stream, error) {
	if emitProgress {
		defer a.reporter.Progress(cameradar.StepValidateStreams, cameradar.ProgressTickMessage())
	}

	if ctx.Err() != nil {
		return stream, ctx.Err()
	}

	client, err := a.newRTSPClient(stream)
	if err != nil {
		return stream, fmt.Errorf("starting rtsp client: %w", err)
	}
	defer client.Close()

	desc, res, err := a.describeWithRetry(ctx, client, stream, cameradar.StepValidateStreams)
	if err != nil {
		return a.handleDescribeError(stream, err)
	}
	a.logDescribeResponse(stream.String(), res)

	if desc == nil || len(desc.Medias) == 0 {
		return stream, fmt.Errorf("no media tracks found for %q", stream)
	}

	res, err = client.Setup(desc.BaseURL, desc.Medias[0], 0, 0)
	if err != nil {
		return a.handleSetupError(stream, err)
	}
	a.logSetupResponse(stream.String(), res)

	stream.Available = res != nil && res.StatusCode == base.StatusOK
	if stream.Available {
		a.reporter.Progress(cameradar.StepValidateStreams, fmt.Sprintf("Stream validated for %s:%d", stream.Address.String(), stream.Port))
	}

	return stream, nil
}

func (a Attacker) describeWithRetry(
	ctx context.Context,
	client *gortsplib.Client,
	stream cameradar.Stream,
	step cameradar.Step,
) (*description.Session, *base.Response, error) {
	u, err := stream.URL()
	if err != nil {
		return nil, nil, fmt.Errorf("building rtsp url: %w", err)
	}

	var (
		desc *description.Session
		res  *base.Response
	)
	for range 5 {
		desc, res, err = client.Describe(u)
		if err == nil {
			return desc, res, nil
		}

		var badStatus liberrors.ErrClientBadStatusCode
		if errors.As(err, &badStatus) && badStatus.Code == base.StatusServiceUnavailable {
			a.reporter.Debug(step, fmt.Sprintf("DESCRIBE %s RTSP/1.0 > %d (retrying)", stream, badStatus.Code))
			select {
			case <-ctx.Done():
				return nil, nil, ctx.Err()
			case <-time.After(time.Second):
			}
			continue
		}

		return nil, nil, err
	}

	return nil, nil, fmt.Errorf("describe retries exhausted for %q: %w", stream, err)
}

func (a Attacker) handleDescribeError(stream cameradar.Stream, err error) (cameradar.Stream, error) {
	var badStatus liberrors.ErrClientBadStatusCode
	if errors.As(err, &badStatus) && badStatus.Code == base.StatusServiceUnavailable {
		a.reporter.Debug(cameradar.StepValidateStreams, fmt.Sprintf("DESCRIBE %s RTSP/1.0 > %d", stream, badStatus.Code))
		a.reporter.Progress(cameradar.StepValidateStreams, fmt.Sprintf("Stream unavailable for %s:%d (RTSP %d)",
			stream.Address.String(),
			stream.Port,
			badStatus.Code,
		))
		stream.Available = false
		return stream, nil
	}

	a.reporter.Debug(cameradar.StepValidateStreams, fmt.Sprintf("DESCRIBE %s RTSP/1.0 > error: %v", stream, err))

	return stream, fmt.Errorf("performing describe request at %q: %w", stream, err)
}

func (a Attacker) handleSetupError(stream cameradar.Stream, err error) (cameradar.Stream, error) {
	var badStatus liberrors.ErrClientBadStatusCode
	if errors.As(err, &badStatus) {
		a.reporter.Debug(cameradar.StepValidateStreams, fmt.Sprintf("SETUP %s RTSP/1.0 > %d", stream, badStatus.Code))
		stream.Available = badStatus.Code == base.StatusOK
		return stream, nil
	}

	return stream, fmt.Errorf("performing setup request at %q: %w", stream, err)
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
