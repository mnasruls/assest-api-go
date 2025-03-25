# Use the official Golang image with build tools
FROM golang:1.24-alpine as builder

# Install build dependencies
RUN apk add --no-cache gcc musl-dev

# Set the working directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application with CGO
RUN CGO_ENABLED=1 GOOS=linux go build -o assets-api ./cmd/main.go

# Use a minimal alpine image for the final stage
FROM alpine:latest

# Install SQLite and CA certificates
RUN apk --no-cache add sqlite ca-certificates

# Set the working directory
WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/assets-api .

# Copy the .env file
COPY .env .

# Create empty SQLite database
RUN touch /app/assets.db

# Expose the application port
EXPOSE 9123

# Command to run the application
CMD ["./assets-api"]