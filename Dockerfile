# Stage 1: Build
FROM golang:1.23.4 AS builder
WORKDIR /shitshat

# Copy dependencies and source code
COPY go.mod go.sum ./
RUN go mod download
COPY . .

# Build the Go binary for Linux
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o shitshat .

# Stage 2: Minimal Runtime
FROM alpine:latest
WORKDIR /shitshat

# Copy the binary from the builder stage
COPY --from=builder /shitshat/shitshat .

# Ensure the binary is executable
RUN chmod +x shitshat

# Set default arguments
CMD ["./shitshat", "--server"]