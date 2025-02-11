package main

import (
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/gorilla/websocket"
)

var (
	// generalClients stores the WebSocket connections with their assigned client names.
	generalClients = make(map[*websocket.Conn]string)
	// gMutex protects access to generalClients.
	gMutex sync.Mutex
	// broadcast is the channel used to broadcast messages to all connected clients.
	broadcast = make(chan Message)
	// nextClientID is used to assign a unique name to each client.
	nextClientID = 1
)

// GeneralChat represents a general chat server.
type GeneralChat struct{}

// NewGeneralChat creates and returns a new GeneralChat instance.
func NewGeneralChat() *GeneralChat {
	return &GeneralChat{}
}

// Handle upgrades the HTTP connection to a WebSocket, assigns a client name, and listens for messages from the client.
func (gc *GeneralChat) Handle(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("General: upgrade error:", err)
		return
	}
	// Automatically assign a client name.
	gMutex.Lock()
	clientName := "Client " + strconv.Itoa(nextClientID)
	nextClientID++
	generalClients[ws] = clientName
	gMutex.Unlock()

	// Send a welcome message with the assigned client name.
	welcomeMsg := Message{
		User: "Server",
		Text: "Welcome, you are " + clientName,
	}
	ws.WriteJSON(welcomeMsg)
	log.Println("General: new client:", clientName)

	// Read messages from the client.
	for {
		var msg Message
		if err := ws.ReadJSON(&msg); err != nil {
			log.Println("General: read error:", err)
			break
		}
		// If the user did not specify a name, use the assigned name.
		if msg.User == "" {
			msg.User = clientName
		}
		log.Printf("General: message from %s: %s\n", msg.User, msg.Text)
		broadcast <- msg
	}

	// Remove the client upon disconnection.
	gMutex.Lock()
	delete(generalClients, ws)
	gMutex.Unlock()
	ws.Close()
	log.Println("General: client disconnected:", clientName)
}

// Broadcast listens for messages on the broadcast channel and sends each received message to all connected clients.
func (gc *GeneralChat) Broadcast() {
	for {
		msg := <-broadcast
		gMutex.Lock()
		for client := range generalClients {
			if err := client.WriteJSON(msg); err != nil {
				log.Println("General: write error:", err)
				client.Close()
				delete(generalClients, client)
			}
		}
		gMutex.Unlock()
	}
}
