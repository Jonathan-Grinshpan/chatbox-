package client

import (
	"fmt"
)

type Client struct {
	Username string
	Password string
	loggedIn bool
}

// SetIsLoggedIn sets the loggedIn status for a Client instance.
func (c *Client) SetIsLoggedIn(logIn bool) {
	c.loggedIn = logIn
}

// Check if user is logged in
func (c *Client) IsLoggedIn() bool {
	return c.loggedIn
}

// Global map to hold registered and login users
var RegisteredUsers = make(map[string]*Client)
var LoggedInUsers = make(map[string]*Client)

func RegisterUser(username, password string) {
	RegisteredUsers[username] = &Client{Username: username, Password: password}
	printAllRegisterdUsers()
}

func LogInUser(username, password string) {
	// Check if the user is registered
	user, exists := RegisteredUsers[username]
	if !exists {
		fmt.Println("User not registered")
		return
	}

	// Check password (this is a basic check; you might want to hash passwords in a real application)
	if user.Password != password {
		fmt.Println("Invalid password")
		return
	}

	user.SetIsLoggedIn(true)
	LoggedInUsers[username] = user
	PrintAllLoggedInUsers()
}

func printAllRegisterdUsers() {
	for _, user := range RegisteredUsers {
		fmt.Printf("username: %s, password: %s, isLoggedIn: %t\n", user.Username, user.Password, user.loggedIn)

	}
	fmt.Println("--------------")

}

func PrintAllLoggedInUsers() {
	for _, user := range LoggedInUsers {
		fmt.Printf("username: %s, password: %s, isLoggedIn: %t\n", user.Username, user.Password, user.loggedIn)

	}
	fmt.Println("--------------")

}

func GetRegisteredUser(username string) (*Client, bool) {
	user, exists := RegisteredUsers[username]
	return user, exists
}

// GetLoggedIndUser retrieves a logged-in user safely.
func GetLoggedIndUser(username string) (*Client, bool) {

	user, exists := LoggedInUsers[username]
	return user, exists
}

// GetAllUserNames returns a slice of usernames for all logged-in users
func GetAllUserNames() []string {
	usernames := make([]string, 0, len(LoggedInUsers))
	for username := range LoggedInUsers {
		usernames = append(usernames, username)
	}
	return usernames
}

func RemoveFromLoggedIn(username string) {
	delete(LoggedInUsers, username)
}
