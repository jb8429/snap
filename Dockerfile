# Use Go base image
FROM golang:1.23-alpine

# Install Air for hot reloading
RUN go install github.com/air-verse/air@latest

# Set the working directory to /app/snap-go
WORKDIR /app/snap-go

# Copy the go.mod file (if using Go modules)
COPY snap-go/go.mod ./

# Download dependencies (only if you have go.mod)
RUN go mod download || true

# Copy the entire project directory into the container
COPY snap-go .

# Expose the port for the Go app
EXPOSE 8080

# Start Air for hot reloading
CMD ["air"]