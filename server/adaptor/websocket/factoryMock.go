// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package websocket

import (
	"net/http"

	"github.com/EtixLabs/cameradar/server/adaptor"

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
