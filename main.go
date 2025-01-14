package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// Server-Side Code
type Message struct {
	Content   string
	Username  string
	Timestamp time.Time
}

type Client struct {
	conn     *websocket.Conn
	send     chan []byte
	username string
}

type Hub struct {
	clients    map[*Client]bool
	broadcast  chan Message
	register   chan *Client
	unregister chan *Client
	messages   []Message
	mu         sync.Mutex
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan Message),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		messages:   []Message{},
	}
}

func (hub *Hub) Run() {
	for {
		select {
		case client := <-hub.register:
			hub.clients[client] = true
			hub.sendRecentMessages(client)
			fmt.Printf("New client connected: %s\n", client.username)
		case client := <-hub.unregister:
			if _, ok := hub.clients[client]; ok {
				delete(hub.clients, client)
				close(client.send)
				fmt.Printf("Client disconnected: %s\n", client.username)
			}
		case message := <-hub.broadcast:
			hub.storeMessage(message)
			hub.broadcastToClients(message)
		}
	}
}

func (hub *Hub) sendRecentMessages(client *Client) {
	hub.mu.Lock()
	defer hub.mu.Unlock()

	cutoff := time.Now().Add(-1 * time.Hour)
	for _, msg := range hub.messages {
		if msg.Timestamp.After(cutoff) {
			client.send <- []byte(fmt.Sprintf("[%s] %s: %s", msg.Timestamp.Format("15:04:05"), msg.Username, msg.Content))
		}
	}
}

func (hub *Hub) storeMessage(message Message) {
	hub.mu.Lock()
	defer hub.mu.Unlock()

	hub.messages = append(hub.messages, message)

	cutoff := time.Now().Add(-1 * time.Hour)
	var filteredMessages []Message
	for _, msg := range hub.messages {
		if msg.Timestamp.After(cutoff) {
			filteredMessages = append(filteredMessages, msg)
		}
	}
	hub.messages = filteredMessages
}

func (hub *Hub) broadcastToClients(message Message) {
	for client := range hub.clients {
		select {
		case client.send <- []byte(fmt.Sprintf("[%s] %s: %s", message.Timestamp.Format("15:04:05"), message.Username, message.Content)):
		default:
			delete(hub.clients, client)
			close(client.send)
		}
	}
}

func (client *Client) Read(hub *Hub) {
	defer func() {
		hub.unregister <- client
		client.conn.Close()
	}()

	for {
		_, msg, err := client.conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			break
		}
		hub.broadcast <- Message{Content: string(msg), Username: client.username, Timestamp: time.Now()}
	}
}

func (client *Client) Write() {
	defer client.conn.Close()
	for msg := range client.send {
		err := client.conn.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			log.Println("Error writing message:", err)
			break
		}
	}
}

func handleConnections(hub *Hub, w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	if username == "" {
		http.Error(w, "Username is required", http.StatusBadRequest)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading connection:", err)
		return
	}

	client := &Client{conn: conn, send: make(chan []byte), username: username}
	hub.register <- client

	go client.Read(hub)
	go client.Write()
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Client-Side Code
func startClient() {
	fmt.Print("Enter your username: ")
	var username string
	fmt.Scanln(&username)

	serverAddr := "ws://localhost:8080/ws?username=" + username
	conn, _, err := websocket.DefaultDialer.Dial(serverAddr, nil)
	if err != nil {
		log.Fatal("Failed to connect to server:", err)
	}
	defer conn.Close()

	fmt.Printf("Connected as '%s'. Type your messages below:\n", username)

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt)

	go func() {
		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				log.Println("Error reading message:", err)
				return
			}
			fmt.Println(string(msg))
		}
	}()

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		if scanner.Scan() {
			text := scanner.Text()
			err := conn.WriteMessage(websocket.TextMessage, []byte(text))
			if err != nil {
				log.Println("Error sending message:", err)
				break
			}
		}

		select {
		case <-done:
			fmt.Println("\nDisconnecting...")
			return
		default:
		}
	}
}

func main() {
	fmt.Println("Choose mode:")
	fmt.Println("1. Run as server")
	fmt.Println("2. Run as client")
	fmt.Print("Enter choice (1/2): ")

	var choice int
	fmt.Scanln(&choice)

	if choice == 1 {
		hub := NewHub()
		go hub.Run()

		http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
			handleConnections(hub, w, r)
		})

		fmt.Println("Server started on :8080")
		log.Fatal(http.ListenAndServe(":8080", nil))
	} else if choice == 2 {
		startClient()
	} else {
		fmt.Println("Invalid choice. Exiting.")
	}
}
