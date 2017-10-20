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

package actor

// Sub/Unsub event type
const (
	SubscribeEvent   = "subscribe"
	UnsubscribeEvent = "unsubscribe"
)

// Subscription contains a sub/unsub event
type Subscription struct {
	Command          string
	Channel          string
	SubscribersCount uint
}

// Publication contains a publish event
type Publication struct {
	Channel string
	Data    string
}

// PubSub is a generic interface for publishing data to subscribers using channels
// It exposes subscriptions events so the controller can create/delete
// data sources depending on the channels users subscribe to.
// ex: launch a camera stream only when users subscribe to it
type PubSub interface {
	Run()
	Sub() <-chan *Subscription
	Pub() chan<- *Publication
}

// ChannelAccessChecker allows to check for accesses on a given channel
type ChannelAccessChecker interface {
	CheckAccess(channel, accessToken string) bool
	ClearCache(accessToken string)
}
