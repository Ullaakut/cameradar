package adaptor

import "net/http"

// WebSocketFactory is an interface for creating generic websocket connections
type WebSocketFactory interface {
	NewIncomingWebSocket(w http.ResponseWriter, req *http.Request) (WebSocket, error)
	NewWebSocket(url string) (WebSocket, error)
}
