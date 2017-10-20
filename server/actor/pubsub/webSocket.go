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
	"fmt"
	"net/http"
	"strings"

	"github.com/EtixLabs/cameradar/server/actor"
	"github.com/EtixLabs/cameradar/server/adaptor"
)

// events allow to serialize in one go routine all events from the subscribers
const (
	eventSubscribe   = "subscribe"
	eventUnsubscribe = "unsubscribe"
	eventDisconnect  = "disconnect"
)

type pubSubEvent struct {
	name    string
	channel string
	client  adaptor.WebSocket
}

// WebSocket manage pubsub communication using a websocket adaptor
type WebSocket struct {
	wsf adaptor.WebSocketFactory

	subscriptions map[string][]adaptor.WebSocket
	sub           chan *actor.Subscription
	pub           chan *actor.Publication
	events        chan *pubSubEvent
}

// NewWebSocket creates a PubSub actor that uses a websockets factory
func NewWebSocket(
	wsf adaptor.WebSocketFactory,
) *WebSocket {
	wsPubSub := &WebSocket{
		wsf: wsf,

		subscriptions: make(map[string][]adaptor.WebSocket),
		sub:           make(chan *actor.Subscription),
		pub:           make(chan *actor.Publication),
		events:        make(chan *pubSubEvent),
	}
	go wsPubSub.Run()
	return wsPubSub
}

// Sub return the chan where websocket event gonna be pushed
func (b *WebSocket) Sub() <-chan *actor.Subscription {
	return b.sub
}

// Pub return the chan where we consider publishement will be asked
func (b *WebSocket) Pub() chan<- *actor.Publication {
	return b.pub
}

// Run start to listen on pubsub events
func (b *WebSocket) Run() {
	for {
		select {
		case event := <-b.events:

			client := event.client
			channel := event.channel

			switch event.name {
			case eventSubscribe:
				b.handleSubscribe(client, channel)
			case eventUnsubscribe:
				b.handleUnsubscribe(client, channel)
			case eventDisconnect:
				b.handleDisconnect(client)
			}
		case publication := <-b.pub:

			subscribers := b.subscriptions[publication.Channel]

			// prepend channel name to message, so client knows from which channel
			// the message comes from
			message := fmt.Sprintf("%s/%s", publication.Channel, publication.Data)

			// broadcast message to subscribers
			for _, client := range subscribers {
				select {
				case client.Write() <- message:
				default:
					// drop frame
				}
			}
		}
	}
}

// Accept a new incoming connection and create a websocket using the factory
func (b *WebSocket) Accept(w http.ResponseWriter, req *http.Request) {
	client, err := b.wsf.NewIncomingWebSocket(w, req)
	if err != nil {
		fmt.Printf("cannot accept incoming connection: %v\n", err)
		return
	}

	go b.readClient(client)
}

func (b *WebSocket) readClient(client adaptor.WebSocket) {
	for {
		message, ok := <-client.Read()
		if !ok {
			// connection channel closed, disconnect client (in the main routine)
			b.events <- &pubSubEvent{
				name:   eventDisconnect,
				client: client,
			}
			return
		}

		// expect text message
		command, ok := message.(string)
		if !ok {
			fmt.Printf("invalid non-text message: %v\n", message)
			return
		}

		// process command
		// NOTE: if another protocol is needed, extract this behavior
		if strings.HasPrefix(command, "s/") {
			channel := strings.TrimPrefix(command, "s/")

			// process in main routine
			b.events <- &pubSubEvent{
				name:    eventSubscribe,
				client:  client,
				channel: channel,
			}
		} else if strings.HasPrefix(command, "u/") {
			channel := strings.TrimPrefix(command, "u/")

			// process in main routine
			b.events <- &pubSubEvent{
				name:    eventUnsubscribe,
				client:  client,
				channel: channel,
			}
		} else {
			fmt.Printf("invalid message '%s', should be [s|u]/{channel}\n", command)
		}
	}
}

func (b *WebSocket) handleSubscribe(client adaptor.WebSocket, channel string) {
	// if client is already subscribed, ignore
	if b.alreadySubscribed(client, channel) {
		return
	}

	// add to subscribers map
	subscribersCount := b.addSubscription(client, channel)

	// notify external world
	b.sub <- &actor.Subscription{
		Command:          actor.SubscribeEvent,
		Channel:          channel,
		SubscribersCount: subscribersCount,
	}
}

func (b *WebSocket) handleUnsubscribe(client adaptor.WebSocket, channel string) {
	// if client didn't subscribe, ignore
	if !b.alreadySubscribed(client, channel) {
		return
	}

	// remove from map
	subscribersCount := b.removeSubscription(client, channel)

	// notify external world
	b.sub <- &actor.Subscription{
		Command:          actor.UnsubscribeEvent,
		Channel:          channel,
		SubscribersCount: subscribersCount,
	}
}

func (b *WebSocket) handleDisconnect(client adaptor.WebSocket) {

	// unsubscribe client from all its channels
	for channel := range b.subscriptions {
		b.handleUnsubscribe(client, channel)
	}

	// close client write channel
	close(client.Write())
}

func (b *WebSocket) alreadySubscribed(client adaptor.WebSocket, channel string) bool {
	clients := b.subscriptions[channel]
	for _, c := range clients {
		if c == client {
			return true
		}
	}
	return false
}

// addSubscription adds a subscription of a client to a channel and returns
// the new number of subscriptions on the given channel
func (b *WebSocket) addSubscription(client adaptor.WebSocket, channel string) uint {
	b.subscriptions[channel] = append(b.subscriptions[channel], client)
	return uint(len(b.subscriptions[channel]))
}

func (b *WebSocket) removeSubscription(client adaptor.WebSocket, channel string) uint {
	clients := b.subscriptions[channel]
	for i, current := range clients {
		if current == client {
			clients = append(clients[:i], clients[i+1:]...)
			break
		}
	}
	b.subscriptions[channel] = clients
	return uint(len(clients))
}
