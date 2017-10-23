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

package server

import (
	"fmt"
	"net/http"

	"github.com/EtixLabs/cameradar/server/adaptor"
)

// WebSocket manages server communication using a websocket adaptor
type WebSocket struct {
	wsf adaptor.WebSocketFactory

	client adaptor.WebSocket

	fromClient chan<- string
	toClient   <-chan string
	disconnect chan interface{}
}

// New creates a Server actor that uses a WebSocketFactory
func New(
	wsf adaptor.WebSocketFactory,
	fromClient chan string,
	toClient chan string,
) *WebSocket {
	wsServer := &WebSocket{
		wsf: wsf,

		fromClient: fromClient,
		toClient:   toClient,
	}
	return wsServer
}

// Accept a new incoming connection and create a websocket using the factory
func (ws *WebSocket) Accept(w http.ResponseWriter, req *http.Request) {
	client, err := ws.wsf.NewIncomingWebSocket(w, req)
	if err != nil {
		fmt.Printf("cannot accept incoming connection: %v\n", err)
		return
	}

	go ws.readClient(client)
}

func (ws *WebSocket) readClient(client adaptor.WebSocket) {
	for {
		select {
		case message, ok := <-client.Read():
			if !ok {
				// connection channel closed, disconnect client (in the main routine)
				ws.disconnect <- struct{}{}
				println("client disconnected")
				return
			}

			// expect text message
			msg, ok := message.(string)
			if !ok {
				fmt.Printf("invalid non-text message: %v\n", message)
				return
			}
			ws.fromClient <- msg
		case msg := <-ws.toClient:
			client.Write() <- msg
		}
	}
}
