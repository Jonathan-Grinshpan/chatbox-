package main

import (
	"fmt"
	"log"
	"net/http"
	"web_chat_app/handlers/auth"
	"web_chat_app/handlers/messages"
	"web_chat_app/handlers/websocket"
)

// http.ResponseWriter allows you to write data
// (e.g., HTML, JSON, text) back to the client and set HTTP headers or status codes.
func homePage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "homepage.html")
}

func setupRoutes() {

	http.HandleFunc("/", homePage)
	http.HandleFunc("/ws", websocket.WsEndpoint)
	http.HandleFunc("/register", auth.RegisterHandler)
	http.HandleFunc("/login", auth.LoginHandler)
	http.Handle("/chatBox.html", http.FileServer(http.Dir("./")))
	http.HandleFunc("/GetMessages", messages.GetMessages)
	http.HandleFunc("/PostMessages", messages.PostMessage)
	http.HandleFunc("/clientCount", websocket.GetClientCount)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

}

func main() {
	fmt.Println("Starting server on :8080")
	setupRoutes()
	//if the server fails to start it will log errors
	log.Fatal(http.ListenAndServe(":8080", nil))
}
