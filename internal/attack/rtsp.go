package attack

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"net"
	"net/textproto"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/Ullaakut/cameradar/v6"
	"github.com/bluenviron/gortsplib/v5"
	"github.com/bluenviron/gortsplib/v5/pkg/base"
	"github.com/bluenviron/gortsplib/v5/pkg/headers"
	"github.com/bluenviron/gortsplib/v5/pkg/liberrors"
)

func (a Attacker) newRTSPClient(u *base.URL) (*gortsplib.Client, error) {
	client := &gortsplib.Client{
		ReadTimeout:  a.timeout,
		WriteTimeout: a.timeout,
	}
	client.Scheme = u.Scheme
	client.Host = u.Host

	err := client.Start()
	if err != nil {
		return nil, err
	}

	return client, nil
}

func (a Attacker) describeStatus(u *base.URL) (base.StatusCode, error) {
	client, err := a.newRTSPClient(u)
	if err != nil {
		return 0, err
	}
	defer client.Close()

	_, res, err := client.Describe(u)
	if err != nil {
		var badStatus liberrors.ErrClientBadStatusCode
		if errors.As(err, &badStatus) {
			return badStatus.Code, nil
		}
		return 0, err
	}
	if res == nil {
		return 0, errors.New("no response received")
	}

	return res.StatusCode, nil
}

// probeDescribeHeaders performs a manual DESCRIBE request and returns the status code and headers.
//
// NOTE: We do not use gortsplib here because it does not expose response headers when the status code is 401 Unauthorized,
// which is exactly what we need in order to detect authentication methods.
func (a Attacker) probeDescribeHeaders(ctx context.Context, u *base.URL, urlStr string) (base.StatusCode, base.Header, error) {
	dialer := &net.Dialer{Timeout: a.timeout}
	conn, err := dialer.DialContext(ctx, "tcp", u.Host)
	if err != nil {
		return 0, nil, err
	}
	defer conn.Close()

	deadline, ok := ctx.Deadline()
	if !ok {
		deadline = time.Now().Add(a.timeout)
	}

	err = conn.SetDeadline(deadline)
	if err != nil {
		return 0, nil, err
	}

	request := fmt.Sprintf(
		"DESCRIBE %s RTSP/1.0\r\nCSeq: 1\r\nUser-Agent: cameradar\r\nAccept: application/sdp\r\nHost: %s\r\n\r\n",
		urlStr,
		u.Host,
	)
	_, err = conn.Write([]byte(request))
	if err != nil {
		return 0, nil, err
	}

	reader := textproto.NewReader(bufio.NewReader(conn))
	statusLine, err := reader.ReadLine()
	if err != nil {
		return 0, nil, err
	}
	fields := strings.Fields(statusLine)
	if len(fields) < 2 {
		return 0, nil, fmt.Errorf("invalid RTSP status line: %q", statusLine)
	}

	code, err := strconv.Atoi(fields[1])
	if err != nil {
		return 0, nil, fmt.Errorf("parsing RTSP status code %q: %w", fields[1], err)
	}

	mimeHeader, err := reader.ReadMIMEHeader()
	if err != nil {
		return 0, nil, err
	}

	headers := make(base.Header)
	for key, values := range mimeHeader {
		headers[key] = append(base.HeaderValue(nil), values...)
	}

	return base.StatusCode(code), headers, nil
}

func authTypeFromHeaders(values base.HeaderValue) cameradar.AuthType {
	if len(values) == 0 {
		return cameradar.AuthUnknown
	}

	var hasBasic bool
	var hasDigest bool

	for _, value := range values {
		var authHeader headers.Authenticate
		err := authHeader.Unmarshal(base.HeaderValue{value})
		if err != nil {
			lower := strings.ToLower(value)
			hasDigest = hasDigest || strings.Contains(lower, "digest")
			hasBasic = hasBasic || strings.Contains(lower, "basic")
			continue
		}

		switch authHeader.Method {
		case headers.AuthMethodDigest:
			hasDigest = true
		case headers.AuthMethodBasic:
			hasBasic = true
		}
	}

	if hasDigest {
		return cameradar.AuthDigest
	}
	if hasBasic {
		return cameradar.AuthBasic
	}
	return cameradar.AuthUnknown
}

func headerValues(header base.Header, name string) base.HeaderValue {
	if header == nil {
		return nil
	}
	for key, values := range header {
		if strings.EqualFold(key, name) {
			return values
		}
	}
	return nil
}

func buildRTSPURL(stream cameradar.Stream, route, username, password string) (*base.URL, string, error) {
	host := net.JoinHostPort(stream.Address.String(), strconv.Itoa(int(stream.Port)))
	path := strings.TrimSpace(route)
	if path != "" && !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	u := &url.URL{
		Scheme: "rtsp",
		Host:   host,
		Path:   path,
	}
	if username != "" || password != "" {
		u.User = url.UserPassword(username, password)
	}

	urlStr := u.String()
	parsed, err := base.ParseURL(urlStr)
	if err != nil {
		return nil, "", err
	}

	return parsed, urlStr, nil
}
