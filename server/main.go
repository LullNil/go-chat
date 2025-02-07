package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

const port = ":8081"

// upgrader is used to upgrade HTTP connections to WebSocket connections.
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

// Message represents a message sent over the WebSocket connection.
type Message struct {
	User string `json:"user"`
	Text string `json:"text"`
}

// handleConnections handles WebSocket requests from clients.
func handleConnections(w http.ResponseWriter, r *http.Request) {
	// Upgrade initial HTTP connection to a WebSocket connection.
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	defer ws.Close()

	log.Println("New client connected")

	// Message processing loop
	for {
		// Read JSON from client.
		_, msgBytes, err := ws.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			break
		}

		// Parse JSON.
		var msg Message
		err = json.Unmarshal(msgBytes, &msg)
		if err != nil {
			log.Println("Parse error:", err)
			break
		}

		log.Printf("Received: %s: %s\n", msg.User, msg.Text)

		// Send JSON back to client.
		// response, err := json.Marshal(msg)
		// if err != nil {
		// 	log.Println("Marshal error:", err)
		// 	break
		// }
		response := Message{User: "Server", Text: msg.Text}
		responseBytes, _ := json.Marshal(response)

		err = ws.WriteMessage(websocket.TextMessage, responseBytes)
		if err != nil {
			log.Println("Write error:", err)
			break
		}
		log.Printf("Sent: %s: %s\n", msg.User, msg.Text)
	}

	log.Println("Client disconnected")
}

func main() {
	http.HandleFunc("/ws", handleConnections)

	// Start the server on port 8081.
	log.Printf("Server is running on %s\n", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
