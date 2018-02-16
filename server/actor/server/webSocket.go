package server

import (
	"fmt"
	"net/http"

	"github.com/Ullaakut/cameradar/server/adaptor"
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
