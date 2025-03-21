package echo

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// upgrader is used to upgrade HTTP connections to WebSocket connections.
// It is shared among the handlers.
var upgrader = websocket.Upgrader{
	// Allow connections from any origin.
	CheckOrigin: func(r *http.Request) bool { return true },
}

// Message represents the structure of a chat message.
type Message struct {
	User string `json:"user"`
	Text string `json:"text"`
}

// EchoServer represents a simple echo server.
type EchoServer struct{}

// NewEchoServer creates and returns a new EchoServer instance.
func NewEchoServer() *EchoServer {
	return &EchoServer{}
}

// Handle upgrades the HTTP connection to a WebSocket and echoes received
// messages back to the client.
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
