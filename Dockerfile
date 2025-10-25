# Step 1: Use official Golang image as builder
FROM golang:1.21-alpine AS builder

# Set working directory inside container
WORKDIR /app

# Copy go.mod and go.sum first for dependency caching
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the project
COPY . .

# Build the Go binary
RUN go build -o main ./cmd/server/main.go

# Step 2: Create a lightweight image
FROM alpine:latest

# Install CA certificates
RUN apk --no-cache add ca-certificates

# Set working directory
WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/main .

# Expose port
EXPOSE 8080

# Command to run the binary
CMD ["./main"]
