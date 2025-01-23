# Use an official Golang image
FROM golang:1.20 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum to the container
COPY go.mod go.sum ./

# Download dependencies
RUN go mod tidy

# Copy the entire project to the container
COPY . .

# Build the application
RUN go build -o task-api main.go

# Use a minimal image for runtime
FROM alpine:latest

# Install required packages (if any)
RUN apk --no-cache add ca-certificates

# Set working directory in runtime container
WORKDIR /root/

# Copy the compiled binary from the builder stage
COPY --from=builder /app/task-api .

# Expose the port the app runs on
EXPOSE 9000

# Command to run the executable
CMD ["./golang-assesment"]
