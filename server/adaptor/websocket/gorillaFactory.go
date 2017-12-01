package websocket

import (
	"fmt"
	"net/http"
	"time"

	"github.com/EtixLabs/cameradar/server/adaptor"

	gorilla "github.com/gorilla/websocket"
	"github.com/pkg/errors"
)

// GorillaFactory is a websocket Factory using Gorilla websocket client
type GorillaFactory struct {
	readLimit           int64
	pingInterval        time.Duration
	pongWait            time.Duration
	writeChanBufferSize uint
	upgrader            gorilla.Upgrader
}

// NewGorillaFactory instantiates and returns a Gorilla Factory
func NewGorillaFactory(options ...func(*GorillaFactory)) *GorillaFactory {
	gf := &GorillaFactory{
		// readLimit: default to 1MB
		readLimit:    1024 * 1024,
		pingInterval: 5 * time.Second,
		pongWait:     10 * time.Second,

		// allow about 1 frame per grid cell to be buffered (grids contain about ~16 cameras)
		// NOTE: this should be about the same size as the number of subcriptions the client has
		writeChanBufferSize: 20,

		// default upgrader: don't check requests origin
		upgrader: gorilla.Upgrader{
			CheckOrigin: func(r *http.Request) bool { return true },
		},
	}
	// apply the options to the struct
	for _, option := range options {
		option(gf)
	}
	return gf
}

// NewIncomingWebSocket instantiates a Gorilla websocket from an http incoming connection
func (gf *GorillaFactory) NewIncomingWebSocket(w http.ResponseWriter, req *http.Request) (adaptor.WebSocket, error) {
	fmt.Printf("new ws connection from %v\n", req.RemoteAddr)

	// create WS connection
	conn, err := gf.upgrader.Upgrade(w, req, nil)
	if err != nil {
		return nil, errors.Wrap(err, "cannot upgrade ws connection")
	}

	g := &Gorilla{
		conn:        conn,
		accessToken: req.Header.Get("Sec-WebSocket-Protocol"),

		input:  make(chan interface{}),
		output: make(chan interface{}, gf.writeChanBufferSize),
	}

	go g.pingAndWrite(gf.pingInterval)
	go g.read(gf.readLimit, gf.pongWait)

	return g, nil
}

// NewWebSocket attemps to connect to a ws server using Gorilla library
func (gf *GorillaFactory) NewWebSocket(url string) (adaptor.WebSocket, error) {
	fmt.Printf("opening new ws connection to %v\n", url)

	// create WS connection
	conn, _, err := gorilla.DefaultDialer.Dial(url, nil)
	if err != nil {
		return nil, errors.Wrap(err, "cannot open ws connection")
	}

	g := &Gorilla{
		conn: conn,

		input:  make(chan interface{}),
		output: make(chan interface{}, gf.writeChanBufferSize),
	}

	go g.pingAndWrite(gf.pingInterval)
	go g.read(gf.readLimit, gf.pongWait)

	return g, nil
}

// SetReadLimit sets the maximum size read from an incoming message
func SetReadLimit(limit int64) func(gf *GorillaFactory) {
	return func(gf *GorillaFactory) {
		gf.readLimit = limit
	}
}

// SetPingInterval sets the interval between pings
func SetPingInterval(interval time.Duration) func(gf *GorillaFactory) {
	return func(gf *GorillaFactory) {
		gf.pingInterval = interval
	}
}

// SetPongWait sets the time before read should timeout if no pong is received
func SetPongWait(pongWait time.Duration) func(gf *GorillaFactory) {
	return func(gf *GorillaFactory) {
		gf.pongWait = pongWait
	}
}

// SetWriteChanBufferSize sets the buffer size of the write channel
func SetWriteChanBufferSize(size uint) func(gf *GorillaFactory) {
	return func(gf *GorillaFactory) {
		gf.writeChanBufferSize = size
	}
}
