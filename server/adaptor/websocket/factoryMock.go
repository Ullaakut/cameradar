package websocket

import (
	"net/http"

	"github.com/Ullaakut/cameradar/server/adaptor"

	"github.com/stretchr/testify/mock"
)

// FactoryMock mocks a websocket factory
type FactoryMock struct {
	mock.Mock
}

// NewIncomingWebSocket mocks the creation of a websocket adaptor
func (m *FactoryMock) NewIncomingWebSocket(
	w http.ResponseWriter,
	req *http.Request,
) (adaptor.WebSocket, error) {
	args := m.Called(w, req)
	return args.Get(0).(adaptor.WebSocket), args.Error(1)
}

// NewWebSocket mocks the creation of a websocket adaptor
func (m *FactoryMock) NewWebSocket(url string) (adaptor.WebSocket, error) {
	args := m.Called(url)
	ws := args.Get(0)
	if ws != nil {
		return ws.(adaptor.WebSocket), args.Error(1)
	}
	return nil, args.Error(1)
}
