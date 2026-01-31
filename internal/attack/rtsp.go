package attack

import (
	"errors"
	"net"
	"net/url"
	"strconv"

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

// describeRTSP is a variable to allow mocking in tests.
var describeRTSP = func(client *gortsplib.Client, u *base.URL) (*base.Response, error) {
	_, res, err := client.Describe(u)
	return res, err
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
		if errors.As(err, &badStatus) && res != nil {
			return badStatus.Code, nil
		}
		return 0, err
	}
	if res == nil {
		return 0, errors.New("no response received")
	}

	return res.StatusCode, nil
}

func authTypeFromHeaders(values base.HeaderValue) cameradar.AuthType {
	if len(values) == 0 {
		return cameradar.AuthNone
	}

	var hasBasic bool
	var hasDigest bool

	for _, value := range values {
		var authHeader headers.Authenticate
		err := authHeader.Unmarshal(base.HeaderValue{value})
		if err != nil {
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
	return cameradar.AuthType(-1)
}

func buildRTSPURL(stream cameradar.Stream, route, username, password string) (*base.URL, string, error) {
	host := net.JoinHostPort(stream.Address.String(), strconv.Itoa(int(stream.Port)))
	path := "/" + route
	if route == "" {
		path = "/"
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
