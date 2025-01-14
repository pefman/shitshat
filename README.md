# Real-Time Chat Application in Go

A simple real-time chat application built with Go, featuring WebSocket-based communication. Users can connect to the server via a command-line interface (CLI) and exchange messages in real-time.

---

## Features

- Real-time communication using WebSockets.
- User-friendly CLI for connecting and chatting.
- Message timestamps for better context.
- Usernames for personalized communication.
- Temporary message storage (up to 1 hour) for new clients to see recent chat history.

---

## Requirements

- [Go](https://golang.org/doc/install) (1.18 or higher)

---

## Installation

1. Clone the repository:

   ```bash
   git clone https://github.com/your-username/real-time-chat.git
   cd real-time-chat

2. Build the project:

   - **Server**:
     ```bash
     cd server
     go build -o chat-server
     ```

   - **Client**:
     ```bash
     cd ../client
     go build -o chat-client
     ```

3. Run the application:

   - **Start the server**:
     ```bash
     cd server
     ./chat-server
     ```

   - **Start the client**:
     Open a new terminal and run:
     ```bash
     cd client
     ./chat-client
     ```

   - Enter your username when prompted and start chatting!

   - Repeat the client step in multiple terminals to test real-time chat functionality.

4. Test the application:

   - Open multiple terminals and start the client in each one.
   - Use different usernames for each client.
   - Send messages from one client and see them appear in real-time on all connected clients.
   - New clients will receive the last hour's chat history upon connecting.