.DEFAULT_GOAL := run
DB_USER := user
DB_PASS := password
DB_NAME := boonmafarm

# Live reload with air (install: go install github.com/air-verse/air@latest)
run:
	PATH="$(shell go env GOPATH)/bin:$$PATH" air

# One-off run without air
run-once:
	go run src/cmd/api/main.go

gen-mocks:
	PATH="$$(go env GOPATH)/bin:$$PATH" go generate ./...

gen-swag:
	swag init -g src/cmd/api/main.go -o docs --propertyStrategy snakecase

test:
	go test  ./...

migrate-new:
	@if [ -z "$(name)" ]; then \
		echo "Usage: make migrate-new name=<migration_name>"; \
		echo "Example: make migrate-new name=create_users_table"; \
		exit 1; \
	fi
	@TIMESTAMP=$$(date +%Y%m%d%H%M%S); \
	mkdir -p migrations; \
	UP_FILE="migrations/dev/$${TIMESTAMP}_$(name).up.sql"; \
	DOWN_FILE="migrations/dev/$${TIMESTAMP}_$(name).down.sql"; \
	touch "$${UP_FILE}" "$${DOWN_FILE}"; \
	echo "✅ Created migration files:"; \
	echo "  📄 Up:   $${UP_FILE}"; \
	echo "  📄 Down: $${DOWN_FILE}";

migrate-up:
	migrate -path migrations/dev -database "postgres://$(DB_USER):$(DB_PASS)@localhost:5432/$(DB_NAME)?sslmode=disable" up

migrate-down:
	migrate -path migrations/dev -database "postgres://$(DB_USER):$(DB_PASS)@localhost:5432/$(DB_NAME)?sslmode=disable" down

migrate-force:
	migrate -path migrations/dev -database "postgres://$(DB_USER):$(DB_PASS)@localhost:5432/$(DB_NAME)?sslmode=disable" force $(version)

migrate-version:
	migrate -path migrations/dev -database "postgres://$(DB_USER):$(DB_PASS)@localhost:5432/$(DB_NAME)?sslmode=disable" version

# Database commands
db-connect:
	docker-compose exec postgres psql -U $(DB_USER) -d $(DB_NAME)