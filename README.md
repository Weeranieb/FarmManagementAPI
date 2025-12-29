# Go Kit Base

A simple Go starter kit with Fiber, GORM, Dependency Injection, and a layered architecture. This template helps you quickly get started with REST APIs using best practices for configuration, repository layer, service layer, and handlers.

## Features

- [Fiber](https://gofiber.io/) web framework for fast HTTP services
- [GORM](https://gorm.io/) ORM for database interactions
- Dependency Injection with [Uber Dig](https://github.com/uber-go/dig)
- YAML config with environment variable overrides
- Layered architecture: Config, Handler, Service, Repository, Model
- Simple User example to get started

## Getting Started

### Prerequisites

- Go 1.18+
- PostgreSQL (or update for your DB)

### Installation

1. Clone the repository:

   ```bash
   git clone https://github.com/weeranieb/go-kit-base.git
   cd go-kit-base
   ```

2. Install Go dependencies:

   ```bash
   go mod tidy
   ```

3. Configure your environment:
   - Copy `config/config.example.yaml` to `config/config.yaml` and edit your DB creds
   - OR set environment variables, e.g. `DATABASE_HOST`, `APP_LOG_LEVEL` (see config section below)

### Running

Start the server:

```bash
go run src/cmd/api/main.go
```

Server runs on the port/host set in config (`localhost:8080` by default).

## Configuration

Configuration is managed with Viper and supports YAML config files and environment variables.

Example `config.yaml`:

```yaml
server:
  port: '8080'
  host: 'localhost'

database:
  host: 'localhost'
  port: '5432'
  name: 'example_db'
  user: 'user'
  password: 'your-password'
  ssl_mode: 'disable'

app:
  environment: 'development'
  log_level: 'info'
  debug: true
```

Environment variables can override config, using uppercase and underscores (e.g. `DATABASE_HOST`).

## Project Structure

```
src/
  cmd/api/          # Main entry point
  internal/
    config/         # Config loading, DB connect
    handler/        # HTTP handlers
    model/          # Structs for database/models
    repository/     # Data layer
    service/        # Business logic
    di/             # Dependency injection setup
  config/           # Config files
```

## Example API

- User CRUD routes are scaffolded (see `internal/handler/user_handler.go`)

## Dependency Injection

The [Uber Dig](https://uber-go.github.io/dig/) container wires dependencies:

```go
container := di.NewContainer(conf)
```

Which wires together DB connection, repositories, services, and handlers.

## Database

- Connects on start via config in `internal/config/database.go`
- Uses GORM for migrations and queries

## License

MIT

## Credits

- [Fiber](https://gofiber.io/)
- [GORM](https://gorm.io/)
- [Uber Dig](https://github.com/uber-go/dig/)

---

```

```
