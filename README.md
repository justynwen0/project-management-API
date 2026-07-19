# Project Management API

A RESTful API for Kanban-style project management, built with Go and Fiber. Supports boards, lists, cards, team collaboration, and dashboard analytics.

## Table of Contents

- [Project Overview](#project-overview)
- [Tech Stack](#tech-stack)
- [Prerequisites](#prerequisites)
- [Getting Started](#getting-started)
- [Environment Configuration](#environment-configuration)
- [Project Structure](#project-structure)
- [Architecture](#architecture)
- [API Reference](#api-reference)
- [Database Schema](#database-schema)

---

## Project Overview

This API implements a Trello-like project management system with the following core features:

- **Boards** — Workspaces owned by a user, shared with team members
- **Lists** — Columns within a board (e.g., To Do, In Progress, Done)
- **Cards** — Tasks within a list with titles, descriptions, due dates, and positions
- **Members** — Users can be added to boards and assigned to individual cards
- **Labels & Attachments** — Metadata and files on cards
- **Dashboard** — Workload distribution and task status percentage analytics
- **Authentication** — JWT-based auth with access and refresh tokens

---

## Tech Stack

| Package | Version | Purpose |
|---------|---------|---------|
| [Go](https://go.dev) | 1.21+ | Language |
| [gofiber/fiber](https://github.com/gofiber/fiber) | v2.52.12 | HTTP web framework |
| [gorm.io/gorm](https://gorm.io) | v1.31.1 | ORM |
| [gorm.io/driver/postgres](https://gorm.io/docs/connecting_to_the_database.html) | v1.6.0 | PostgreSQL driver |
| [golang-jwt/jwt](https://github.com/golang-jwt/jwt) | v5.3.1 | JWT token signing |
| [gofiber/jwt](https://github.com/gofiber/jwt) | v3.3.10 | JWT middleware for Fiber |
| [jackc/pgx](https://github.com/jackc/pgx) | v5.6.0 | PostgreSQL client |
| [google/uuid](https://github.com/google/uuid) | v1.6.0 | UUID generation |
| [joho/godotenv](https://github.com/joho/godotenv) | v1.5.1 | `.env` file loader |
| bcrypt (Go stdlib) | — | Password hashing |
| [golang-migrate/migrate](https://github.com/golang-migrate/migrate) | latest | Database migrations (CLI tool) |

---

## Prerequisites

- **Go** 1.21 or higher
- **PostgreSQL** 13 or higher
- **golang-migrate CLI** — download from [github.com/golang-migrate/migrate/releases](https://github.com/golang-migrate/migrate/releases) and add to your PATH

---

## Getting Started

### 1. Clone the repository

```bash
git clone <repo-url>
cd project-management-API
```

### 2. Configure environment

Copy and edit the environment file:

```bash
cp .env.example .env
```

Fill in your database credentials, JWT secret, and other values. See [Environment Configuration](#environment-configuration) for all variables.

### 3. Create the database

```sql
CREATE DATABASE project_management;
```

> Note: The database name in the connection string uses an underscore (`project_management`), not a hyphen.

### 4. Run migrations

```bash
migrate -path database/migrations \
  -database "postgres://USER:PASSWORD@localhost:5432/project_management?sslmode=disable" \
  up
```

Replace `USER` and `PASSWORD` with your PostgreSQL credentials.

### 5. Run the server

```bash
go run main.go
```

The server starts on port `3030` by default (configurable via `APP_PORT`).

### Build binary

```bash
go build -o project-management-API
./project-management-API
```

---

## Environment Configuration

Create a `.env` file at the project root with the following variables:

| Variable | Default | Required | Description |
|----------|---------|----------|-------------|
| `APP_PORT` | `3030` | No | HTTP server port |
| `DB_HOST` | `localhost` | Yes | PostgreSQL host |
| `DB_PORT` | `5432` | Yes | PostgreSQL port |
| `DB_USER` | — | Yes | PostgreSQL username |
| `DB_PASSWORD` | — | Yes | PostgreSQL password |
| `DB_NAME` | — | Yes | PostgreSQL database name |
| `APP_URL` | `http://localhost:3030` | No | Base URL of the API |
| `CORS_ORIGIN` | — | Yes | Allowed frontend origin (e.g. `http://localhost:5173`) |
| `JWT_SECRET` | — | Yes | Secret key for JWT signing |
| `JWT_EXPIRED` | `6h` | No | Access token TTL |
| `REFRESH_TOKEN_EXPIRED` | `24h` | No | Refresh token TTL |
| `ADMIN_EMAIL` | — | Yes | Email for the seeded admin account |
| `ADMIN_PASSWORD` | — | Yes | Password for the seeded admin account |
| `ADMIN_ROLE` | `admin` | No | Role assigned to the seeded admin |

Example `.env`:

```env
APP_PORT=3030
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=yourpassword
DB_NAME=project_management
APP_URL=http://localhost:3030
CORS_ORIGIN=http://localhost:5173
JWT_SECRET=your-secret-key
JWT_EXPIRED=6h
REFRESH_TOKEN_EXPIRED=24h
ADMIN_EMAIL=admin@example.com
ADMIN_PASSWORD=admin123
ADMIN_ROLE=admin
```

---

## Project Structure

```
project-management-API/
├── main.go                     # Entry point — initializes app, wires DI, starts server
├── go.mod                      # Module definition and dependencies
├── go.sum                      # Dependency checksums
├── .env                        # Environment configuration (not committed)
│
├── config/
│   └── config.go               # Loads .env, establishes DB connection, exposes config values
│
├── routes/
│   └── route.go                # Registers all routes, applies JWT middleware to protected group
│
├── controllers/                # HTTP request handlers — parse input, call service, return JSON
│   ├── user_controller.go
│   ├── board_controller.go
│   ├── list_controller.go
│   ├── card_controller.go
│   └── dashboard_controller.go
│
├── services/                   # Business logic — validation, orchestration across repositories
│   ├── user_services.go
│   ├── board_services.go
│   ├── list_services.go
│   ├── card_services.go
│   └── dashboard_services.go
│
├── repositories/               # Data access layer — GORM queries against PostgreSQL
│   ├── user_repository.go
│   ├── board_repository.go
│   ├── board_member_repository.go
│   ├── list_repository.go
│   ├── card_repository.go
│   └── dashboard_repository.go
│
├── models/                     # GORM struct definitions for all entities
│   ├── user.go
│   ├── board.go
│   ├── board_member.go
│   ├── list.go
│   ├── list_position.go
│   ├── card.go
│   ├── card_position.go
│   ├── card_assignees.go
│   ├── card_label.go
│   ├── card_attachment.go
│   ├── label.go
│   ├── comment.go
│   └── dashboard.go
│
├── utils/
│   ├── jwt.go                  # Token generation (access + refresh)
│   ├── password.go             # Bcrypt hashing and verification
│   └── response.go             # Standardized JSON response helpers
│
└── database/
    ├── migrations/             # 15 migration pairs (.up.sql / .down.sql)
    └── seed/
        └── seed_admin.go       # Creates default admin user on startup
```

---

## Architecture

### Request Flow

```
HTTP Request
    │
    ▼
[Fiber Router]
    │
    ├── Public routes (no auth)
    │       POST /v1/auth/register
    │       POST /v1/auth/login
    │
    └── Protected routes (/api/v1/*)
            │
            ▼
        [JWT Middleware]  — validates Bearer token, injects claims into context
            │
            ▼
        [Controller]  — parses request (body, params, query), calls service
            │
            ▼
        [Service]     — enforces business rules, validates references, orchestrates repos
            │
            ▼
        [Repository]  — executes GORM queries against PostgreSQL
            │
            ▼
        [PostgreSQL]
```

### Key Design Patterns

**Layered Architecture (MVC + Repository + Service)**  
Each layer has a single responsibility. Controllers never touch the DB directly; repositories never contain business logic.

**Constructor-based Dependency Injection**  
All dependencies are wired in `main.go`:
```go
userRepo    := repositories.NewUserRepository(config.DB)
userService := services.NewUserService(userRepo)
userCtrl    := controllers.NewUserController(userService)
```

**Dual ID Strategy**  
Every entity carries two IDs:
- `internal_id` — auto-increment integer, used for DB foreign keys (never exposed in API responses)
- `public_id` — UUID v4, exposed in all API responses and used in URL parameters

**Soft Deletes**  
User records use GORM's `DeletedAt` field. Deleted users are hidden from queries but remain in the database.

**UUID Array Ordering**  
`list_positions` and `card_positions` tables store ordering as `UUID[]` (PostgreSQL native array). This enables drag-and-drop reordering without updating every row's position index.

**Standardized Response Envelope**  
All responses follow a consistent shape:

```json
{
  "status": "Success",
  "response_code": 200,
  "message": "Human-readable message",
  "data": { }
}
```

Paginated responses include a `meta` field:
```json
{
  "status": "Success",
  "response_code": 200,
  "message": "Data ditemukan",
  "data": [],
  "meta": {
    "page": 1,
    "limit": 10,
    "total": 100,
    "total_pages": 10,
    "filter": "",
    "sort": "-id"
  }
}
```

**Authentication**  
JWT HS256 tokens. Access token (`JWT_EXPIRED`, default 6h) + refresh token (`REFRESH_TOKEN_EXPIRED`, default 24h). The JWT middleware injects `user_id`, `pub_id`, `role`, and `email` into the Fiber context for downstream handlers.

---

## API Reference

### Authentication

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| POST | `/v1/auth/register` | No | Register new user |
| POST | `/v1/auth/login` | No | Login, returns JWT access + refresh tokens |

### Users

All routes require `Authorization: Bearer <token>`.

| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/v1/users/:page` | List users with pagination (`?limit=10&filter=name&sort=-id`) |
| GET | `/api/v1/users/:id` | Get user by `public_id` |
| PUT | `/api/v1/users/:id` | Update user by `public_id` |
| DELETE | `/api/v1/users/:id` | Delete user by `public_id` |

### Boards

| Method | Path | Description |
|--------|------|-------------|
| POST | `/api/v1/boards` | Create board (owner = authenticated user) |
| GET | `/api/v1/boards/my` | List boards belonging to authenticated user (paginated) |
| GET | `/api/v1/boards/:id` | Get board details by `public_id` |
| PUT | `/api/v1/boards/:id` | Update board by `public_id` |
| DELETE | `/api/v1/boards/:id` | Delete board and cascade-delete all lists, cards |
| POST | `/api/v1/boards/:id/members` | Add members to board (`body: { user_ids: [...] }`) |
| DELETE | `/api/v1/boards/:id/members` | Remove members from board |
| GET | `/api/v1/boards/:id/members` | Get all members of a board |

### Lists

| Method | Path | Description |
|--------|------|-------------|
| POST | `/api/v1/lists` | Create list (requires `board_id`) |
| PUT | `/api/v1/lists/:id` | Update list by `public_id` |
| DELETE | `/api/v1/lists/:id` | Delete list by `public_id` |
| GET | `/api/v1/boards/:board_id/lists` | Get all lists in a board |

### Cards

| Method | Path | Description |
|--------|------|-------------|
| POST | `/api/v1/cards` | Create card (requires `list_id`, `title`, `position`) |
| GET | `/api/v1/cards/:id` | Get card details (includes assignees, labels, attachments) |
| PUT | `/api/v1/cards/:id` | Update card by `public_id` |
| DELETE | `/api/v1/cards/:id` | Delete card by `public_id` |
| GET | `/api/v1/lists/:id/cards` | Get all cards in a list |
| POST | `/api/v1/cards/:id/assignees` | Add assignees to card (`body: { user_ids: [...] }`) |
| DELETE | `/api/v1/cards/:id/assignees` | Remove assignees from card |

### Dashboard

| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/v1/dashboard/workload` | Task count grouped by assignee |
| GET | `/api/v1/dashboard/task-percentage` | Task completion status as percentages |

---

## Database Schema

### Entity Relationships

```
users
 ├── owns → boards (owner_internal_id)
 ├── member of → board_members (many-to-many with boards)
 └── assigned to → card_assignees (many-to-many with cards)

boards
 ├── has many → lists
 └── has one → list_positions (UUID[] ordering)

lists
 ├── has many → cards
 └── has one → card_positions (UUID[] ordering)

cards
 ├── has many → card_assignees
 ├── has many → card_labels (many-to-many with labels)
 ├── has many → card_attachments
 └── has many → comments
```

### Tables

| Table | Key Columns | Notes |
|-------|------------|-------|
| `users` | `internal_id`, `public_id`, `name`, `email`, `password`, `role` | Soft delete via `deleted_at` |
| `boards` | `internal_id`, `public_id`, `title`, `description`, `owner_internal_id`, `due_date` | — |
| `board_members` | `board_internal_id`, `user_internal_id`, `joined_at` | Composite PK |
| `lists` | `internal_id`, `public_id`, `board_internal_id`, `title`, `position` | — |
| `list_positions` | `board_internal_id` (unique), `list_order UUID[]` | One row per board |
| `cards` | `internal_id`, `public_id`, `list_internal_id`, `title`, `description`, `due_date`, `position` | — |
| `card_positions` | `list_internal_id` (unique), `card_order UUID[]` | One row per list |
| `card_assignees` | `card_internal_id`, `user_internal_id` | Composite PK |
| `card_labels` | `card_internal_id`, `label_internal_id` | Composite PK |
| `card_attachment` | `internal_id`, `public_id`, `file`, `user_internal_id`, `card_internal_id` | — |
| `labels` | `internal_id`, `public_id`, `name`, `color` | — |
| `comments` | `internal_id`, `public_id`, `card_internal_id`, `user_internal_id`, `message` | — |

### ID Strategy

Every entity uses a **dual ID system**:
- `internal_id` — auto-increment integer primary key, used in all foreign key relationships (never returned by the API)
- `public_id` — UUID v4, exposed in all API responses and used as the identifier in URL parameters

This decouples the public-facing API from the internal database structure and prevents enumeration attacks.
