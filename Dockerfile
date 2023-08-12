# Start with a base Golang image
FROM golang:latest

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files
COPY go.mod go.sum ./

# Download the dependencies
RUN go mod download

# Copy the source code files
COPY . .

# Build the Go application
RUN go build -o main .

# Expose the port the application listens on
EXPOSE 8080

# Command to run the executable
CMD ["./main"]
