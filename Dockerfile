# Use the official golang image as the base image
FROM golang:latest

# Set the working directory inside the container
WORKDIR /app

# Copy the Go server files into the container
COPY . .

# Build the Go application
RUN go mod init github.com/Ash7540/go-sorting-server
RUN go mod tidy
RUN go build -o mygoserver

# Expose the port the server listens on
EXPOSE 8000

# Command to run the Go application
CMD ["./mygoserver"]
