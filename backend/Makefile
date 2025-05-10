include .env
MIGRATION_PATH=cmd/migration
DB_URL=postgresql://${DB_USERNAME}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_DATABASE}?sslmode=disable

# Create container using postgres 16 image
postgres:
	@docker run --name postgres -p ${DB_PORT}:5432 -e POSTGRES_USER=${DB_USERNAME} -e POSTGRES_PASSWORD=${DB_PASSWORD} -d postgres:16-alpine

# Create database
db/create:
	@docker exec -it postgres createdb --username=${DB_USERNAME} --owner=${DB_USERNAME} ${DB_DATABASE}

# Drop database
db/drop:
	@docker exec -it postgres dropdb ${DB_DATABASE}

# Create a new migration
migrate/new:
	@migrate create -ext sql -dir "${MIGRATION_PATH}" -seq ${name}

# Up all migrations
migrate/up:
	@migrate -path ${MIGRATION_PATH} -database "$(DB_URL)" -verbose up

# Up the latest migration
migrate/up1:
	@migrate -path ${MIGRATION_PATH} -database "$(DB_URL)" -verbose up 1

# Down all migrations
migrate/down:
	@migrate -path ${MIGRATION_PATH} -database "$(DB_URL)" -verbose down

# Down the latest migration
migrate/down1:
	@migrate -path ${MIGRATION_PATH} -database "$(DB_URL)" -verbose down 1

# Generate Go code using sqlc
sqlc:
	@sqlc generate

all: build test

# Build the application
build:
	@echo "Building..."
	@go build -o main cmd/api/main.go

# Run the application
run:
	@go run cmd/api/main.go

# Create DB container
docker/run:
	@if docker compose up --build 2>/dev/null; then \
		: ; \
	else \
		echo "Falling back to Docker Compose V1"; \
		docker-compose up --build; \
	fi

# Shutdown DB container
docker/down:
	@if docker compose down 2>/dev/null; then \
		: ; \
	else \
		echo "Falling back to Docker Compose V1"; \
		docker-compose down; \
	fi

# Test the application
test:
	@echo "Testing..."
	@go test ./... -v

# Integrations Tests for the application
itest:
	@echo "Running integration tests..."
	@go test ./internal/database -v

# Clean the binary
clean:
	@echo "Cleaning..."
	@rm -f main

# Live Reload
watch:
	@if command -v air > /dev/null; then \
            air; \
            echo "Watching...";\
        else \
            read -p "Go's 'air' is not installed on your machine. Do you want to install it? [Y/n] " choice; \
            if [ "$$choice" != "n" ] && [ "$$choice" != "N" ]; then \
                go install github.com/air-verse/air@latest; \
                air; \
                echo "Watching...";\
            else \
                echo "You chose not to install air. Exiting..."; \
                exit 1; \
            fi; \
        fi

.PHONY: postgres db/create db/drop migrate/new migrate/up migrate/up1 migrate/down migrate/down1 all build run test clean watch docker-run docker-down itest
