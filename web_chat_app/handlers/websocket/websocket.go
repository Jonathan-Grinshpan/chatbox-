package websocket

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"web_chat_app/client"
	"web_chat_app/handlers/messages"

	"github.com/gorilla/websocket"
)

// upgrade HTTP connections to WebSocket connections.
var upgrader = websocket.Upgrader{
	//data is read or written in chunks up to the specified buffer size
	ReadBufferSize:  1024,
	WriteBufferSize: 1024, //For a simple chat application, 1024 bytes (1 KB) is typically sufficient
	CheckOrigin:     func(r *http.Request) bool { return true },
}

// userConnections to store active connections, with a mutex for concurrent access.
var (
	userConnections = make(map[string]*websocket.Conn)
	connMutex       = &sync.Mutex{}
	messageChannel  = make(chan messages.Message) // Channel for broadcasting messages
)

// Start a goroutine for handling messages, the init func is called automatically in Go
func init() {

	//The messageHandler function runs in its own goroutine. It continuously listens for messages received on messageChannel.
	go messageHandler()
}

// initizalized a websocket connection between the server and a user
func WsEndpoint(w http.ResponseWriter, r *http.Request) {

	username := r.URL.Query().Get("Username")
	fmt.Println("Received username:", username) // Debugging line

	// Check if the user is logged in
	if _, exists := client.GetLoggedIndUser(username); !exists {
		fmt.Println("User not found in logged-in users:", client.GetAllUserNames()) //debug
		http.Error(w, "Unauthorized or not logged in", http.StatusUnauthorized)
		return
	}

	// Upgrade the connection to WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Failed to upgrade connection:", err)
		return
	}
	
	connMutex.Lock()
	userConnections[username] = conn
	connMutex.Unlock()

	fmt.Println("Client connected:", conn.RemoteAddr(), "Username:", username)
	// when the websocket connection is closed
	defer func() {
		conn.Close()
		connMutex.Lock()
		delete(userConnections, username) // Remove user on disconnect
		connMutex.Unlock()
	}()

	// Listen for incoming messages
	for {
		// The server listens for messages sent by the client
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			break
		}
		fmt.Println("Message from user", username, ":", string(msg))
		// Send the message to the message channel for broadcasting
		messageChannel <- messages.Message{Username: username, Text: string(msg)}
	}
}

// messageHandler listens for incoming messages from the messageChannel.

func messageHandler() {

	// When a new message is received, it triggers the broadcastMessage function
	for message := range messageChannel {
		// Broadcast the message to all active connections
		broadcastMessage(message.Username, message.Text)
	}
}

// Broadcast message to all active connections
func broadcastMessage(sender string, message string) {
	// Wrap message in JSON format
	jsonMessage := fmt.Sprintf(`{"type": "chatMessage", "sender": "%s", "message": "%s"}`, sender, string(message))

	connMutex.Lock()
	defer connMutex.Unlock()

	//for every user that has a websocket connection
	for username, conn := range userConnections {
		if username == sender {
			continue // Skip sending the message back to the sender to avoid duplicates
		}

		// Launch a goroutine to send the message concurrently
		go sendMessage(username, conn, jsonMessage)
	}
}

// sendMessage handles sending a message to a specific WebSocket connection.
func sendMessage(username string, conn *websocket.Conn, jsonMessage string) {

	//attempt to send the message over the websocket, will trigger the reciver (socket.onmessage)
	err := conn.WriteMessage(websocket.TextMessage, []byte(jsonMessage))
	if err != nil {
		fmt.Println("Error broadcasting message to", username, ":", err)

		// Safely close the connection and remove it from the map
		connMutex.Lock()
		conn.Close()
		delete(userConnections, username)
		connMutex.Unlock()
	}
}

func GetClientCount(w http.ResponseWriter, r *http.Request) {
	connMutex.Lock()
	userCount := len(userConnections)
	connMutex.Unlock()

	response := map[string]int{"clientCount": userCount}
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
	}
}
