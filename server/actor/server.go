package actor

// Server is a generic interface for creating a bidirectional
// communication server through websocket.
type Server interface {
	Run()
}
