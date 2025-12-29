# BoonmaFarm Backend

A comprehensive Farm Management API built with Go Fiber, GORM, Dependency Injection, and a clean layered architecture. This backend provides RESTful APIs for managing farms, ponds, workers, merchants, feed collections, and more.

## Features

- **[Fiber](https://gofiber.io/)** - Fast HTTP web framework
- **[GORM](https://gorm.io/)** - Powerful ORM for database interactions
- **Dependency Injection** - Using [Uber Dig](https://github.com/uber-go/dig) for clean dependency management
- **JWT Authentication** - Secure token-based authentication
- **Database Migrations** - Version-controlled schema migrations using `migrate`
- **Docker Support** - PostgreSQL database via Docker Compose
- **Swagger Documentation** - Auto-generated API documentation
- **Comprehensive Testing** - Unit tests for repositories, services, and handlers
- **Snake Case Database** - Consistent snake_case column naming
- **Layered Architecture** - Handler → Service → Repository → Model

## Getting Started

### Prerequisites

- Go 1.24+
- Docker and Docker Compose (for PostgreSQL)
- [migrate](https://github.com/golang-migrate/migrate) CLI tool (for database migrations)
- [swag](https://github.com/swaggo/swag) (for Swagger documentation)

### Installation

1. Clone the repository:

   ```bash
   git clone <repository-url>
   cd backend
   ```

2. Install Go dependencies:

   ```bash
   go mod tidy
   ```

3. Start PostgreSQL database:

   ```bash
   docker-compose up -d
   ```

4. Run database migrations:

   ```bash
   make migrate-up
   ```

5. Configure your environment:
   - Edit `configuration/config.yaml` with your settings
   - Or set environment variables (see Configuration section)

### Running

Start the server:

```bash
make run
# or
go run src/cmd/api/main.go
```

Server runs on `localhost:8080` by default.

## Configuration

Configuration is managed with Viper and supports YAML config files and environment variables.

Example `configuration/config.yaml`:

```yaml
server:
  port: '8080'
  host: 'localhost'

database:
  host: 'localhost'
  port: '5432'
  name: 'boonmafarm'
  user: 'user'
  password: 'password'
  ssl_mode: 'disable'

app:
  environment: 'development'
  log_level: 'info'
  debug: false

authentication:
  jwt_secret: 'FarmSecretKey'
  jwt_expiry: '24h'
```

Environment variables can override config values using uppercase and underscores (e.g., `DATABASE_HOST`, `DATABASE_NAME`).

## Project Structure

```
backend/
├── src/
│   ├── cmd/api/              # Main entry point
│   └── internal/
│       ├── config/           # Configuration loading and DB connection
│       ├── handler/          # HTTP handlers (controllers)
│       ├── service/          # Business logic layer
│       ├── repository/       # Data access layer
│       ├── model/            # Database models
│       ├── dto/              # Data Transfer Objects
│       ├── errors/           # Custom error definitions
│       ├── middleware/       # HTTP middleware (JWT auth, etc.)
│       ├── router/           # Route definitions
│       ├── di/               # Dependency injection container
│       └── utils/            # Utility functions
├── configuration/            # Config files
├── migrations/dev/           # Database migration files
├── docs/                     # Swagger documentation
├── docker-compose.yaml       # PostgreSQL setup
└── Makefile                  # Build and development commands
```

## API Endpoints

### Authentication

- `POST /api/v1/auth/register` - Register a new user
- `POST /api/v1/auth/login` - Login and get JWT token

### Farms

- `POST /api/v1/farm` - Create a new farm
- `GET /api/v1/farm/:id` - Get farm by ID
- `GET /api/v1/farm` - List farms for current client
- `PUT /api/v1/farm` - Update farm

### Ponds

- `POST /api/v1/pond` - Create a new pond
- `POST /api/v1/pond/batch` - Create multiple ponds
- `GET /api/v1/pond/:id` - Get pond by ID
- `GET /api/v1/pond?farmId=X` - List ponds by farm
- `PUT /api/v1/pond` - Update pond
- `DELETE /api/v1/pond/:id` - Delete pond

### Workers

- `POST /api/v1/worker` - Create a new worker
- `GET /api/v1/worker/:id` - Get worker by ID
- `GET /api/v1/worker` - List workers with pagination
- `PUT /api/v1/worker` - Update worker

### Merchants

- `POST /api/v1/merchant` - Create a new merchant
- `GET /api/v1/merchant/:id` - Get merchant by ID
- `GET /api/v1/merchant` - List all merchants
- `PUT /api/v1/merchant` - Update merchant

### Feed Collections

- `POST /api/v1/feedcollection` - Create feed collection with price history
- `GET /api/v1/feedcollection/:id` - Get feed collection by ID
- `GET /api/v1/feedcollection` - List feed collections with pagination
- `PUT /api/v1/feedcollection` - Update feed collection

### Feed Price History

- `POST /api/v1/feedpricehistory` - Create price history entry
- `GET /api/v1/feedpricehistory/:id` - Get price history by ID
- `GET /api/v1/feedpricehistory?feedCollectionId=X` - Get all prices for a feed
- `PUT /api/v1/feedpricehistory` - Update price history

**Note:** Most endpoints require JWT authentication. Include the token in the `Authorization` header: `Bearer <token>`

## API Documentation

Swagger documentation is available at:

```
http://localhost:8080/swagger/index.html
```

To regenerate Swagger docs:

```bash
make gen-swag
```

## Makefile Commands

```bash
make run              # Start the server
make test             # Run all tests
make gen-swag         # Generate Swagger documentation
make gen-mocks        # Generate mock files for testing
make migrate-up       # Run database migrations
make migrate-down     # Rollback last migration
make migrate-version  # Show current migration version
make db-connect       # Connect to PostgreSQL via Docker
```

## Database Migrations

### Create a new migration

```bash
make migrate-new name=add_new_table
```

This creates two files:

- `migrations/dev/TIMESTAMP_add_new_table.up.sql`
- `migrations/dev/TIMESTAMP_add_new_table.down.sql`

### Run migrations

```bash
make migrate-up
```

### Rollback migrations

```bash
make migrate-down
```

## Testing

Run all tests:

```bash
make test
```

Run tests with coverage:

```bash
go test ./... -coverprofile=coverage.out
go tool cover -func=coverage.out
```

Generate HTML coverage report:

```bash
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html
open coverage.html
```

## Docker Setup

The project includes a `docker-compose.yaml` file for easy PostgreSQL setup:

```bash
# Start PostgreSQL
docker-compose up -d

# Stop PostgreSQL
docker-compose down

# View logs
docker-compose logs -f postgres
```

Database credentials (default):

- Database: `boonmafarm`
- User: `user`
- Password: `password`
- Port: `5432`

## Dependency Injection

The project uses [Uber Dig](https://uber-go.github.io/dig/) for dependency injection:

```go
container := di.NewContainer(conf)
```

This automatically wires together:

- Database connection
- Repositories
- Services
- Handlers

## Architecture

### Handler Layer

- Handles HTTP requests/responses
- Validates input
- Extracts JWT context
- Calls service layer

### Service Layer

- Contains business logic
- Validates business rules
- Handles transactions
- Returns DTOs

### Repository Layer

- Database operations
- Query building
- Data mapping
- Error handling

### Model Layer

- Database models
- GORM tags
- Relationships

## Error Handling

Custom error codes are defined in `internal/errors/codes.go`:

- `500010-500019`: Validation errors
- `500020-500029`: Authentication errors
- `500030-500039`: User errors
- `500040-500049`: Farm errors
- `500060-500069`: Merchant errors
- `500070-500079`: Pond errors
- `500080-500089`: Worker errors
- `500090-500099`: FeedCollection errors
- `500100-500109`: FeedPriceHistory errors

## License

MIT

## Credits

- [Fiber](https://gofiber.io/)
- [GORM](https://gorm.io/)
- [Uber Dig](https://github.com/uber-go/dig/)
- [Swagger](https://swaggo.github.io/swaggo/)
