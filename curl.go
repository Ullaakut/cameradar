package cmrdr

import (
	curl "github.com/andelf/go-curl"
)

// Curler is an interface that implements the CURL interface of the go-curl library
// Used for mocking
type Curler interface {
	Setopt(opt int, param interface{}) error
	Perform() error
	Getinfo(info curl.CurlInfo) (interface{}, error)
	Duphandle() Curler
}

// Curl is a libcurl wrapper used to make the Curler interface work even though
// golang currently does not support covariance (see https://github.com/golang/go/issues/7512)
type Curl struct {
	*curl.CURL
}

// Duphandle wraps curl.Duphandle
func (c *Curl) Duphandle() Curler {
	return &Curl{c.CURL.Duphandle()}
}
