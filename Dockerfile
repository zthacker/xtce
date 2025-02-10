# Stage 1: Build the Go application
FROM golang:1.23 AS builder
WORKDIR /app

# Copy go modules and dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire project
COPY . .

# Build the entire Go module (not just main.go)
RUN go build -o mcs-service .

# Stage 2: Create a lightweight container
FROM debian:latest
WORKDIR /root/

# Install dependencies
RUN apt-get update && apt-get install -y ca-certificates

# Copy built binary
COPY --from=builder /app/mcs-service .

# Run the service
CMD ["./mcs-service"]
