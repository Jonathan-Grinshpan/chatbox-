let username;
let socket;

function getUsernameFromURL() {
  const params = new URLSearchParams(window.location.search);
  return params.get("Username");
}

async function initializeChat() {
  username = getUsernameFromURL();
  if (username) {
    document.getElementById("usernameLabel").textContent = `${username}:`;

    // Initialize WebSocket connection only if not already initialized
    if (!socket || socket.readyState === WebSocket.CLOSED) {
      socket = new WebSocket(`ws://localhost:8080/ws?Username=${username}`);

      //when opening a new websocket
      socket.onopen = function () {
        console.log("Connected to WebSocket");
        fetchMessages(); // Fetch messages from the backend right on connection
      };

      //triggered whenever a message is recived (not sent) through the WebSocket connection
      //this will update the other user's chatbox when a different user sends a message
      socket.onmessage = function (event) {
        try {
          const message = JSON.parse(event.data);
          displayMessage(message.message);
        } catch (error) {
          console.error("Failed to parse message:", event.data, error);
          displayMessage(event.data); // Optionally display non-JSON messages in chat
        }
      };

      socket.onerror = function (error) {
        console.error("WebSocket error:", error);
      };

      socket.onclose = function (event) {
        console.log("WebSocket connection closed:", event);
      };
    }
  } else {
    alert("Username not found. Redirecting to login page.");
    window.location.href = "./login.html"; // Redirect to the login page
  }

  // Add event listener to send message on "enter" key
  const messageInput = document.getElementById("messageInput");
  messageInput.addEventListener("keypress", function (event) {
    if (event.key === "Enter") {
      event.preventDefault();
      sendMessage();
    }
  });
}

//fetches messages from the backend database using a get API.
async function fetchMessages() {
  try {
    //get the messages
    const response = await fetch("/GetMessages");
    if (!response.ok) {
      throw new Error("Network response was not ok");
    }

    const messages = await response.json();
    if (messages == null) {
      return;
    }
    //add each message from the DB to the chatbox of the user
    if (Array.isArray(messages)) {
      messages.forEach((msg) => {
        displayMessage(`${msg.Username}: ${msg.Text}`);
      });
    } else {
      console.error("Expected an array of messages, but got:", messages);
    }
  } catch (error) {
    console.error("Error fetching messages:", error);
  }
}

//send a message in the chatbox
function sendMessage() {
  const messageInput = document.getElementById("messageInput");
  const message = messageInput.value;
  if (message && socket && socket.readyState === WebSocket.OPEN) {
    const fullMessage = `${username}: ${message}`;
    // Send to backend WebSocket, after sending will send the message to all other users through broadcast
    socket.send(fullMessage);
    displayMessage(fullMessage); // Display in chatbox

    // Send message to the backend to store it
    fetch("/PostMessages", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ username: username, text: message }),
    });

    messageInput.value = "";
  }
}

//add message to the chatbox
function displayMessage(message) {
  const chatbox = document.getElementById("chatbox");
  chatbox.value += `${message}\n`;
  chatbox.scrollTop = chatbox.scrollHeight;
}

//api to fetch the number of clients connected via websocket connection
async function fetchClientCount() {
  try {
    const response = await fetch("/clientCount");
    if (!response.ok) {
      throw new Error("Network response was not ok");
    }
    const data = await response.json();

    //update the number after clicking
    document.getElementById(
      "clientCountDisplay"
    ).textContent = `Users connected: ${data.clientCount}`;
  } catch (error) {
    console.error("Error fetching client count:", error);
  }
}

window.onload = initializeChat;
