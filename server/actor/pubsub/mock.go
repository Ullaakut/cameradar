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

package pubsub

import (
	"github.com/stretchr/testify/mock"

	"github.com/EtixLabs/cameradar/server/actor"
)

// Mock mocks a pubsub actor
type Mock struct {
	mock.Mock
}

// Sub mock
func (m *Mock) Sub() <-chan *actor.Subscription {
	args := m.Called()
	return args.Get(0).(<-chan *actor.Subscription)
}

// Pub mock
func (m *Mock) Pub() chan<- *actor.Publication {
	args := m.Called()
	return args.Get(0).(chan<- *actor.Publication)
}

// Run mock
func (m *Mock) Run() {
	m.Called()
}

// AccessCheckerMock mocks a channel access checker
type AccessCheckerMock struct {
	mock.Mock
}

// CheckAccess mocks an access check
func (m *AccessCheckerMock) CheckAccess(channel, accessToken string) bool {
	args := m.Called(channel, accessToken)
	return args.Bool(0)
}

// ClearCache mocks a cache clear
func (m *AccessCheckerMock) ClearCache(accessToken string) {
	m.Called(accessToken)
}
