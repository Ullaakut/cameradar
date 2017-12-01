package adaptor

// WebSocket is an interface that represents an authenticated websocket connection
type WebSocket interface {
	AccessToken() string
	Read() <-chan interface{}
	Write() chan<- interface{}
}
