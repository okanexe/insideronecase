# Start from the official Golang image
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install git for go mod and build essentials
RUN apk add --no-cache git

# Copy go mod and sum files
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the Go app
RUN go build -o main ./cmd/main.go

# Use a minimal image for running
FROM alpine:latest

WORKDIR /app

# Copy the built binary from the builder
COPY --from=builder /app/main .

# Expose the port
EXPOSE 8080

# Set environment variables (optional, can be overridden by docker-compose)
ENV POSTGRES_URL=postgres://postgres:1234@db:5432/mydb?sslmode=disable
ENV REDIS_ADDR=redis:6379
ENV WEBHOOK_URL=https://webhook.site/9d7cfe72-0001-4753-a5ed-5cd1fef67224

# Run the binary
CMD ["./main"]