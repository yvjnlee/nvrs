# Use an official Go image with Debian as the base
FROM golang:1.23.1-bullseye

# Set the working directory inside the container
WORKDIR /app

# Install SQLite dependencies
RUN apt-get update && apt-get install -y sqlite3 libsqlite3-dev

# Copy the Go mod and sum files, then download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application files
COPY . .

# Set environment variables to enable CGO
ENV CGO_ENABLED=1

# Build the Go application binary
RUN go build -o nvrs-gateway cmd/nvrs-gateway/main.go

# Expose the port your application listens on
EXPOSE 8080

# Command to run the executable
CMD ["./nvrs-gateway"]
