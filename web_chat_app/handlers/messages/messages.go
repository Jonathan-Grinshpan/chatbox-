package messages

import (
	"encoding/json"
	"net/http"
)

type Message struct {
	Username string
	Text     string
}

var messages_arr []Message

// api to add messages to the messages_arr
func PostMessage(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	//if sending a message, add the message from the front end to the backend
	var msg Message
	err := json.NewDecoder(r.Body).Decode(&msg)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	messages_arr = append(messages_arr, msg)

	w.WriteHeader(http.StatusOK)

}

// api to get messages from the database
func GetMessages(w http.ResponseWriter, r *http.Request) {
	// Check if the request method is GET
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	// Encode and send the messages array as JSON
	if err := json.NewEncoder(w).Encode(messages_arr); err != nil {
		http.Error(w, "Failed to encode messages", http.StatusInternalServerError)
	}

}
