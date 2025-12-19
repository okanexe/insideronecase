docker compose up --build

migrate -path . -database "postgres://postgres:1234@localhost:5432/mydb?sslmode=disable" up

curl --location --request POST 'localhost:8080/stop' --data ''

curl --location --request POST 'localhost:8080/start' --data ''

curl --location 'localhost:8080/messages'

swag init --dir cmd,api --output docs

# Insider One Message Sender Service

This project is a message sending service built with **Go**, **PostgreSQL**, and **Redis**.  
It supports background workers for message delivery, API-based control (start/stop), sent message history retrieval, and Swagger-based API documentation.

---

## Architecture Overview

- **Go** – Core application & HTTP API
- **PostgreSQL** – Persistent message storage
- **Redis** – caching layer
- **Worker** – Background process for sending messages
- **Swagger (swaggo)** – API documentation

---

## Requirements

- Docker & Docker Compose
- Go 1.20+
- `migrate` CLI
- `swag` CLI

---

## Running the Project

### Start Redis, PostgreSQL, and Go services

```bash
docker compose up --build
```

### Database Migrations

Run migrations manually using:
```bash
migrate -path . -database "postgres://postgres:1234@localhost:5432/mydb?sslmode=disable" up
```

### API Usage

Start Message Sending service
```bash
curl --location --request POST 'localhost:8080/start' --data ''
```

Stop Message Sending Worker
```bash
curl --location --request POST 'localhost:8080/stop' --data ''
```

Get Sent Messages
```bash
curl --location 'localhost:8080/messages'
```

### Swagger API Documentation

Generate Swagger Docs
```bash
swag init --dir cmd,api,internal/repository --output docs
```

### Access Swagger UI

Once the server is running, open:
```bash
http://localhost:8080/swagger/index.html
```

### Project Structure
.
├── api/                # HTTP handlers & routes
├── cmd/                # Application entrypoint
├── docs/               # Swagger generated files
├── internal/
│   ├── cache/          # Redis client wrapper
│   ├── migrations/     # SQL Migrations
│   ├── repository/     # PostgreSQL access layer
│   ├── service/        # Business logic
│   ├── webhook/        # External webhook client
│   └── worker/         # Background worker
├── docker-compose.yml
├── go.mod
└── README.md

