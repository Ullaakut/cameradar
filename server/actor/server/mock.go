package server

import (
	"net/http"

	"github.com/stretchr/testify/mock"
)

// Mock mocks a pubsub actor
type Mock struct {
	mock.Mock
}

// Accept mock
func (m *Mock) Accept(w http.ResponseWriter, req *http.Request) {
	m.Called(w, req)
}
