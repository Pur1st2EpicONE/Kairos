![kairos banner](assets/banner.png)

<h3 align="center">Robust event booking service using RabbitMQ delayed queues, PostgreSQL storage, JWT authentication, and automatic cancellation of expired reservations.</h3>

<br>

## Table of Contents

- [Architecture](#architecture)
- [Installation](#installation)
- [Configuration](#configuration)
- [Shutting down](#shutting-down)
- [API](#api)
- [Request examples](#request-examples)

<br>

## Architecture

- **App** — the central orchestrator.
  Loads configuration, initializes all components (storage, broker, notifier, service, HTTP server), wires dependencies, and manages startup and graceful shutdown using a shared context.

- **Broker** — messaging layer (RabbitMQ) responsible for delayed cancellation of bookings.
  When a booking is created, a message is published to a per‑booking queue with a TTL equal to the booking’s time‑to‑live. After the TTL expires, the message is dead‑lettered into the main queue, where a consumer picks it up and triggers the cancellation logic.

- **Service** — the application‑level business logic.
  Validates inputs, enforces state transitions, coordinates with storage, broker, and notifier, and exposes a clean API to the HTTP layer. Domain rules live here.

- **Storage** — the persistent data layer (PostgreSQL). 
  Stores users, events, and bookings. Uses transactions to ensure data consistency when creating bookings, confirming payments, or cancelling expired reservations.

- **Notifier** — delivery adapter for outbound notifications. 
  Sends alerts when a booking is created or cancelled (expired). Encapsulates channel‑specific protocols behind a unified interface.

- **Handler** — HTTP layer (Gin).
  Serves both JSON API endpoints and simple HTML pages for user interaction. Handles request parsing, response formatting, and error mapping.

<br>

## Installation
⚠️ Note: This project requires Docker Compose, regardless of how you choose to run it.  

First, clone the repository and enter the project folder:

```bash
git clone https://github.com/Pur1st2EpicONE/Kairos.git
cd Kairos
```

Then you have two options:

#### 1. Run everything in containers
```bash
make
```

This will start the entire project fully containerized using Docker Compose.

#### 2. Run Kairos locally
```bash
make local
```
In this mode, only PostgreSQL and RabbitMQ are started in containers via Docker Compose, while the application itself runs locally.

⚠️ Note: Local mode requires Go 1.25.1 installed on your machine.

<br>

## Configuration

### Runtime configuration

Kairos uses two configuration files, depending on the selected run mode:

[config.full.yaml](./configs/config.full.yaml) — used for the fully containerized setup

[config.dev.yaml](./configs/config.dev.yaml) — used for local development

You may optionally review and adjust the corresponding configuration file to match your preferences. The default values are suitable for most use cases.

### Environment variables and notification credentials

By default, Kairos runs without any external notification credentials. In this mode, Telegram notifications are disabled. If you want to enable Telegram notifications, you must provide the corresponding credentials via environment variable.

Kairos uses a .env file for runtime configuration. You may create your own .env file manually before running the service, or edit [.env.example](.env.example) and let it be copied automatically on startup.
If environment file does not exist, .env.example is copied to create it. If environment file already exists, it is used as-is and will not be overwritten.

⚠️ Note: Keep .env.example for local runs. Some Makefile commands rely on it and may break if it's missing.

<br>

## Shutting down

Stopping Kairos depends on how it was started:

- Local setup — press Ctrl+C to send SIGINT to the application. The service will gracefully close connections and finish any in-progress operations.  
- Full Docker setup — containers run by Docker Compose will be stopped automatically.

In both cases, to stop all services and clean up containers, run:

```bash
make down
```

⚠️ Note: In the full Docker setup, the log folder is created by the container as root and will not be removed automatically. To delete it manually, run:
```bash
sudo rm -rf <log-folder>
```

⚠️ Note: Docker Compose also creates a persistent volume for PostgreSQL data (kairos_postgres_data). This volume is not removed automatically when containers are stopped. To remove it and fully reset the environment, run:
```bash
make reset
```

<br>

## API

All endpoints are mounted under /api/v1. Responses follow a simple wrapper convention:

- Success: **200 OK** with JSON body **{"result": \<value>}**
- Error: appropriate status code with JSON body **{"error": "\<message>"}**

<br>

### Public endpoints

#### Authentication

```bash
POST /api/v1/auth/sign-up      # registration
```

```bash
POST /api/v1/auth/sign-in      # login
```

Request body stays the same for both requests:
```json
{
  "login": "TSoprano",
  "password": "gabagool"
}
```

On success, both endpoints return a JWT token:
```json
{
  "result": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

<br>

#### Get event information

```bash
GET /api/v1/events/{id}
```

Returns event details (public info):
```json
{
  "result": {
    "title": "Go Conference 2026",
    "description": "Annual gathering of Gophers",
    "date": "2026-06-15T10:00:00Z",
    "seats": 100
  }
}
```

<br>

### Protected endpoints (require JWT in Authorization: Bearer <token> header)

#### Create an event

```bash
POST /api/v1/events
```

Request body:
```json
{
  "title": "Bada Bing! Grand Reopening",
  "description": "Free drinks for made guys. No FBI allowed.",
  "date": "2026-08-15T21:00:00Z",
  "seats": 50,
  "booking_ttl": "30m"
}
```

On success, returns the created event ID:
```json
{
  "result": "550e8400-e29b-41d4-a716-446655440000"
}
```

<br>

#### Book a seat

```bash
POST /api/v1/events/{id}/book
```

No request body. On success, returns the booking ID:
```json
{
  "result": 33
}
```

<br>

#### Confirm a booking

```bash
POST /api/v1/events/{id}/confirm
```

No request body. On success, returns:
```json
{
  "result": "confirmed"
}
```

<br>

#### Typical error responses

- **400 Bad Request** — validation errors, missing fields, invalid formats.
- **401 Unauthorized** — missing or invalid JWT, wrong credentials.
- **404 Not Found** — event or booking not found.
- **409 Conflict** — duplicate booking, event full, already confirmed, etc.
- **500 Internal Server Error** — unexpected server failures.

<br>

## Request examples

⚠️ Note: When the service is running, a web-based UI is available at http://localhost:8080. The examples below demonstrate how to interact with the API directly using curl.

### Register a user

```bash
curl -X POST http://localhost:8080/api/v1/auth/sign-up \
  -H "Content-Type: application/json" \
  -d '{
    "login": "TSoprano",
    "password": "gabagool"
  }'
```

### Response

```json
{
  "result": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

<br>

### Login

```bash
curl -X POST http://localhost:8080/api/v1/auth/sign-in \
  -H "Content-Type: application/json" \
  -d '{
    "login": "TSoprano",
    "password": "gabagool"
  }'
```

### Response

```json
{
  "result": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

<br>

### Create an event (authenticated)

```bash
curl -X POST http://localhost:8080/api/v1/events \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <JWT_TOKEN>" \
  -d '{
    "title": "Bada Bing! Grand Reopening",
    "description": "Free drinks for made guys. No FBI allowed.",
    "date": "2026-08-15T21:00:00Z",
    "seats": 50,
    "booking_ttl": "30m"
  }'
```

### Response

```json
{
  "result": "f47ac10b-58cc-4372-a567-0e02b2c3d479"
}
```

<br>

### Get event info (public)

```bash
curl http://localhost:8080/api/v1/events/f47ac10b-58cc-4372-a567-0e02b2c3d479
```

### Response

```json
{
  "result": {
    "title": "Bada Bing! Grand Reopening",
    "description": "Free drinks for made guys. No FBI allowed.",
    "date": "2026-08-15T21:00:00Z",
    "seats": 50
  }
}
```

<br>

### Book a seat (authenticated)

```bash
curl -X POST http://localhost:8080/api/v1/events/f47ac10b-58cc-4372-a567-0e02b2c3d479/book \
  -H "Authorization: Bearer <JWT_TOKEN>"
```

Response:

```json
{
  "result": 101
}
```

<br>

### Confirm the booking (authenticated)

```bash
curl -X POST http://localhost:8080/api/v1/events/f47ac10b-58cc-4372-a567-0e02b2c3d479/confirm \
  -H "Authorization: Bearer <JWT_TOKEN>"
```

Response:

```json
{
  "result": "confirmed"
}
```
