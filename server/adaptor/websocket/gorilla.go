package websocket

import (
	"fmt"
	"time"

	gorilla "github.com/gorilla/websocket"
)

// Gorilla implements WebSocket interface using Gorilla library
type Gorilla struct {
	conn        *gorilla.Conn
	accessToken string

	input  chan interface{}
	output chan interface{}
}

// AccessToken returns the user authentication token
func (g *Gorilla) AccessToken() string {
	return g.accessToken
}

// Write return a chan to retrieve websocket inputs
func (g *Gorilla) Read() <-chan interface{} {
	return g.input
}

// Write returns a chan to write on websocket
func (g *Gorilla) Write() chan<- interface{} {
	return g.output
}

func (g *Gorilla) read(readLimit int64, pongWait time.Duration) {
	defer (func() {
		g.conn.Close()
		close(g.input)
	})()

	// setup read to timeout if no pong is received after `pongWait`
	g.conn.SetReadDeadline(time.Now().Add(pongWait))
	g.conn.SetPongHandler(func(string) error {
		g.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	g.conn.SetReadLimit(readLimit)
	for {
		messageType, message, err := g.conn.ReadMessage()
		if err != nil {
			if _, ok := err.(*gorilla.CloseError); ok {
				fmt.Printf("ws connection closed from %v (%v)\n", g.conn.RemoteAddr(), err)
			} else {
				// most of the time, a read error is not an error (connection closed, ...)
				fmt.Printf("ws read error: %v\n", err)
			}
			break
		}

		switch messageType {
		case gorilla.TextMessage:
			g.input <- string(message)
		case gorilla.BinaryMessage:
			g.input <- message
		default:
			fmt.Printf("received invalid message type: %v\n", messageType)
		}

	}
}

func (g *Gorilla) pingAndWrite(pingInterval time.Duration) {
	defer g.conn.Close()

	pinger := time.NewTicker(pingInterval)

	for {
		select {
		case <-pinger.C:
			if err := g.conn.WriteMessage(gorilla.PingMessage, []byte{}); err != nil {
				fmt.Printf("ping write error: %v\n", err)
				return
			}
		case message, ok := <-g.output:
			if !ok {
				// chan closed, stop write routine
				return
			}

			var err error

			switch msg := message.(type) {
			case []byte:
				err = g.conn.WriteMessage(gorilla.BinaryMessage, msg)
			case string:
				err = g.conn.WriteMessage(gorilla.TextMessage, []byte(msg))
			default:
				err = fmt.Errorf("invalid message type: %T", msg)
			}

			if err != nil {
				fmt.Printf("write error: %v\n", err)
				return
			}
		}
	}
}
