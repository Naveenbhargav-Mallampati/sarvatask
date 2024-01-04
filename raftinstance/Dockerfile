# Use a minimal base image for Go applications
FROM golang:1.21.4-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the source code into the container
COPY . .

# Build the Go application
RUN go build -o raftinstance.exe .

# Create a new image from scratch (no underlying OS)
FROM alpine:latest

# Set the working directory inside the container
WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/raftinstance.exe /app/

# Copy the entry script into the container
COPY entry.sh /app/

# Make the entry script executable
RUN chmod +x /app/entry.sh

# Set the entry point for the container
CMD ["/app/entry.sh"]
