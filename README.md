# Shitshat a Chat Application in Go

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

- [Go](https://golang.org/doc/install) (1.23.4 or higher)

---

## Installation

1. Clone the repository:

   ```bash
   git clone https://github.com/pefman/shitshat.git
   cd shitshat
   ```

2. Build the project:

   ```bash
   go build -o shitshat main.go
   ```

3. Run the application:

   - **Start the server**:
     ```bash
     ./shitshat --server --address <address> --port <port>
     ```

   - **Start the client**:
     Open a new terminal and run:
     ```bash
     ./shitshat --client --address <address> --port <port>
     ```

   - Enter your username when prompted and start chatting!

   - Repeat the client step in multiple terminals to test real-time chat functionality.

4. Test the application:

   - Open multiple terminals and start the client in each one.
   - Use different usernames for each client.
   - Send messages from one client and see them appear in real-time on all connected clients.
   - New clients will receive the last hour's chat history upon connecting.

---

## Docker

You can also build and run the application using Docker.

1. Build the Docker image:

   ```bash
   ./build_docker_image.sh
   ```

2. Run the Docker container:

   ```bash
   docker run -p <port>:<port> shitshat --server --address 0.0.0.0 --port <port>
   ```

3. Connect clients to the server using the same steps as above, but ensure the address and port match the Docker container settings.

---

## License

This project is licensed under the Apache License 2.0. See the [LICENSE](LICENSE) file for details.
