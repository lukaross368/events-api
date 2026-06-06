# Events API

A RESTful API for managing events and user registrations, built with Go, Gin, and SQLite.

## Features

- JWT authentication
- Custom in-memory rate limiter — per-IP token bucket built from scratch using `sync.Mutex` and `golang.org/x/time/rate`, with a background goroutine that evicts stale clients to keep memory bounded
- Event CRUD (owner-only writes)
- User registration for events

## Requirements

- Go 1.21+
- GCC (required by the SQLite driver)

## Getting Started

### Local

```bash
go mod download
JWT_SECRET=your-secret go run main.go
```

The server starts on `http://localhost:8080`.

### Docker (recommended)

```bash
docker compose -f docker-compose-dev.yaml up
```

## Environment Variables

| Variable    | Default              | Description                              |
|-------------|----------------------|------------------------------------------|
| `JWT_SECRET` | *(insecure default)* | Secret key for signing JWTs — **required in production** |
| `DB_PATH`   | `/app/api.db`        | Path to the SQLite database file         |

## API Reference

### Authentication

Authenticated endpoints require a `Authorization` header with the JWT returned from `/login`.

```
Authorization: <token>
```

---

### Users

#### `POST /signup`

Create a new user account.

**Body**
```json
{ "email": "user@example.com", "password": "secret" }
```

**Response** `201 Created`
```json
{ "message": "User created" }
```

---

#### `POST /login`

Authenticate and receive a JWT.

**Body**
```json
{ "email": "user@example.com", "password": "secret" }
```

**Response** `200 OK`
```json
{ "message": "Login successful", "token": "<jwt>" }
```

---

### Events

#### `GET /events`

List all events.

**Response** `200 OK` — array of event objects.

---

#### `GET /events/:id`

Get a single event by ID.

**Response** `200 OK`
```json
{
  "ID": 1,
  "Name": "Go Meetup",
  "Description": "Monthly Go meetup",
  "Location": "London",
  "DateTime": "2026-07-01T18:00:00Z",
  "UserId": 3
}
```

---

#### `GET /totalEvents`

Get the total count of events.

**Response** `200 OK`
```json
{ "totalEvents": 42 }
```

---

#### `POST /events` *(auth required)*

Create a new event.

**Body**
```json
{
  "Name": "Go Meetup",
  "Description": "Monthly Go meetup",
  "Location": "London",
  "DateTime": "2026-07-01T18:00:00Z"
}
```

**Response** `201 Created`
```json
{ "message": "Event created", "event": { ... } }
```

---

#### `PUT /events/:id` *(auth required, owner only)*

Update an event. Only the user who created the event can update it.

**Body** — same shape as `POST /events`.

**Response** `200 OK`
```json
{ "message": "Event Updated" }
```

---

#### `DELETE /events/:id` *(auth required, owner only)*

Delete an event.

**Response** `200 OK`
```json
{ "message": "Event Deleted" }
```

---

### Registrations

#### `POST /events/:id/register` *(auth required)*

Register the authenticated user for an event.

**Response** `201 Created`
```json
{ "message": "Event Registered" }
```

---

#### `DELETE /events/:id/register` *(auth required)*

Cancel the authenticated user's registration for an event.

**Response** `200 OK`
```json
{ "message": "Event Canceled" }
```

---

## Rate Limiting

The API uses a **custom in-memory rate limiter** rather than an off-the-shelf middleware. Each unique IP gets its own token bucket (5 req/s, burst of 1), stored in a mutex-guarded map. A background goroutine runs every minute and evicts any client entry that hasn't been seen in 3 minutes, keeping memory usage proportional to active traffic rather than all-time unique IPs.

Exceeding the limit returns `429 Too Many Requests`.

## Project Structure

```
.
├── db/           # Database init and connection
├── middlewares/  # Auth and rate-limiting middleware
├── models/       # Data models and DB queries
├── processes/    # Background workers (limiter cleanup)
├── routes/       # HTTP handlers
├── utils/        # JWT and password hashing helpers
└── main.go
```
