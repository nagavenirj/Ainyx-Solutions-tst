# User Management API

A RESTful API built with Go, Fiber, SQLC, and PostgreSQL to manage users with dynamic age calculation.

## 🚀 Tech Stack

- **GoFiber**: Web framework
- **SQLC**: Type-safe database code generation
- **PostgreSQL**: Database
- **Uber Zap**: Structured logging
- **Validator**: Input validation
- **Docker**: Containerization

## 🛠️ Project Structure

```
/cmd/server/      # Main application entry point
/config/          # Configuration loading
/db/              # Database migrations and queries
/internal/
  /handler/       # HTTP handlers
  /middleware/    # Request ID and logging middleware
  /models/        # Domain models and DTOs
  /repository/    # Data access layer
  /routes/        # Route definitions
  /service/       # Business logic (age calculation)
```

## ⚙️ Setup

### Prerequisites

- Go 1.21 or higher
- Docker and Docker Compose (optional)
- PostgreSQL (if running locally without Docker)

### Running with Docker (Recommended)

```bash
docker-compose up --build
```

The API will be available at `http://localhost:8080`.

### Running Locally

1. **Install dependencies**
   ```bash
   go mod tidy
   go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
   ```

2. **Start PostgreSQL**
   Ensure PostgreSQL is running and create a database named `userdb`.

3. **Run Migrations**
   Connect to your database and execute the SQL in `db/migrations/001_create_users_table.sql`.

4. **Generate SQLC Code**
   ```bash
   sqlc generate
   ```

5. **Start the Server**
   ```bash
   go run cmd/server/main.go
   ```

## 🧪 API Endpoints

| Method | Endpoint | Description |
| --- | --- | --- |
| POST | `/users` | Create a new user |
| GET | `/users/:id` | Get user details with age |
| GET | `/users` | List users (paginated) |
| PUT | `/users/:id` | Update user details |
| DELETE | `/users/:id` | Delete a user |

### Example Request (Create User)

```bash
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{"name": "Alice", "dob": "1990-05-10"}'
```

## ✅ Features

- **Dynamic Age Calculation**: Age is calculated on-the-fly based on DOB.
- **Pagination**: List users with `?page=1&page_size=10`.
- **Validation**: Strict input validation using `go-playground/validator`.
- **Logging**: Structured logs with request ID and duration.
- **Docker**: Ready-to-use Docker configuration.
