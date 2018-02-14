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
}
