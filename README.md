# Insider One Message Sender Service

This project is a message sending service built with **Go**, **PostgreSQL**, and **Redis**.  
It supports background workers for message delivery, API-based control (start/stop), sent message history retrieval, and Swagger-based API documentation.

### Project Structure
```text
.
├── api/                # HTTP handlers & routes
├── cmd/                # Application entrypoint
├── docs/               # Swagger generated files
├── internal/
│   ├── cache/          # Redis client wrapper
│   ├── migrations/     # SQL migrations
│   ├── repository/     # PostgreSQL access layer
│   ├── service/        # Business logic
│   ├── webhook/        # External webhook client
│   └── worker/         # Background worker
├── docker-compose.yml
├── go.mod
└── README.md
```

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

# Running the Project

### Start Redis, PostgreSQL, and Go services
All components (Golang Server, PostgreSQL, Redis, and database migrations) are fully containerized.
You only need to run Docker Compose; the database and migrations will be initialized automatically.

```bash
docker compose up --build
```

### Swagger API Documentation
You can start and stop the message service and retrieve sent messages using the Swagger UI.

### Access Swagger UI

Once the server is running, open:
```bash
http://localhost:8080/swagger/index.html
```

---

## Notes

### Generate Swagger Docs
```bash
swag init --dir cmd,api,internal/repository --output docs
```

### Database Migrations
If you want to manage migrations manually (outside Docker), use the following commands:

Run migrations manually using
Apply migrations:
```bash
migrate -path . -database "postgres://postgres:1234@localhost:5432/mydb?sslmode=disable" up
```

Rollback migrations:
```bash
migrate -path . -database "postgres://postgres:1234@localhost:5432/mydb?sslmode=disable" down
```

### API Usage
You can interact with the service either via curl or the Swagger UI.

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
