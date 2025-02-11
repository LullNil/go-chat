package main

import (
	"log"
	"net/http"
)

// EchoServer represents a simple echo server.
type EchoServer struct{}

// NewEchoServer creates and returns a new EchoServer instance.
func NewEchoServer() *EchoServer {
	return &EchoServer{}
}

// Handle upgrades the HTTP connection to a WebSocket and echoes received messages back to the client.
func (es *EchoServer) Handle(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Echo: upgrade error:", err)
		return
	}
	defer ws.Close()
	log.Println("Echo: client connected")

	for {
		var msg Message
		// Read a JSON-encoded message from the client.
		if err := ws.ReadJSON(&msg); err != nil {
			log.Println("Echo: read error:", err)
			break
		}
		log.Printf("Echo: received from %s: %s\n", msg.User, msg.Text)

		// Prepare the echo response.
		response := Message{
			User: "Server",
			Text: msg.Text,
		}

		// Write the JSON-encoded response back to the client.
		if err := ws.WriteJSON(response); err != nil {
			log.Println("Echo: write error:", err)
			break
		}
	}
	log.Println("Echo: client disconnected")
}
