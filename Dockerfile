# Start from the official Go image
FROM docker.arvancloud.ir/golang:1.24-alpine

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the Go app
RUN go build -o app .

# Expose the port your app runs on
EXPOSE 15000 15001 15002 15003

# Command to run the app
CMD ["./app"] 