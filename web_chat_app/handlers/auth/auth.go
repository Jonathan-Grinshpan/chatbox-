package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"web_chat_app/client"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var newClient client.Client
	if err := json.NewDecoder(r.Body).Decode(&newClient); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	if _, exists := client.GetRegisteredUser(newClient.Username); exists {
		w.Header().Set("Content-Type", "application/json")
		message := fmt.Sprintf("User %s already registered", newClient.Username)
		json.NewEncoder(w).Encode(map[string]string{"message": message})
		return
	}

	// Call the RegisterUser function from the client package
	client.RegisterUser(newClient.Username, newClient.Password)

	// Send a response back to the client
	w.Header().Set("Content-Type", "application/json")
	message := fmt.Sprintf("User %s registered successfully!", newClient.Username)
	json.NewEncoder(w).Encode(map[string]string{"message": message})
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var newClient client.Client
	if err := json.NewDecoder(r.Body).Decode(&newClient); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	// Check if the user exists and the password matches
	user, exists := client.GetRegisteredUser(newClient.Username)
	if exists && user.Password == newClient.Password {

		client.LogInUser(user.Username, user.Password)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"message":  "Login successful!",
			"username": newClient.Username,
		})
		return
	}

	// If credentials are invalid
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	json.NewEncoder(w).Encode(map[string]string{"message": "Invalid username or password"})
}
