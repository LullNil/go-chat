package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// Message represents the structure of a chat message.
type Message struct {
	User string `json:"user"`
	Text string `json:"text"`
}

// upgrader is used to upgrade HTTP connections to WebSocket connections.
// It is shared among the handlers.
var upgrader = websocket.Upgrader{
	// Allow connections from any origin.
	CheckOrigin: func(r *http.Request) bool { return true },
}

func main() {
	// Initialize the echo server.
	echoServer := NewEchoServer()
	// Initialize the general chat server.
	generalChat := NewGeneralChat()

	// Register the endpoints for echo and general chat.
	http.HandleFunc("/ws/echo", echoServer.Handle)
	http.HandleFunc("/ws/general", generalChat.Handle)

	// Start the broadcast routine for the general chat.
	go generalChat.Broadcast()

	log.Println("Server started on :8081")
	if err := http.ListenAndServe(":8081", nil); err != nil {
		log.Fatal("ListenAndServe error:", err)
	}
}
