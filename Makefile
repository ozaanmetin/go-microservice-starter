# Run the Go microservice server
.PHONY: runserver
runserver:
		go run cmd/server/main.go


# Run all tests
.PHONY: test
test:
		go test -v ./...


# Run tests with coverage report
.PHONY: test-coverage
test-coverage:
		go test -v -coverprofile=coverage.out ./...
		go tool cover -html=coverage.out -o coverage.html
		@echo "Coverage report generated at coverage.html"


# Build the Go microservice server
.PHONY: build
build:
		go build -o bin/server cmd/server/main.go


# Clean build artifacts and coverage reports
.PHONY: clean
clean:
		rm -rf bin/
		rm -f coverage.out coverage.html


# Database migration commands
.PHONY: migrate-up
migrate-up:
		@echo "Running migrations..."
		goose -dir db/migrations postgres "host=127.0.0.1 port=5433 user=postgres password=postgres dbname=microservice_db sslmode=disable" up

.PHONY: migrate-down
migrate-down:
		@echo "Rolling back migration..."
		goose -dir db/migrations postgres "host=127.0.0.1 port=5433 user=postgres password=postgres dbname=microservice_db sslmode=disable" down

.PHONY: migrate-status
migrate-status:
		@echo "Checking migration status..."
		goose -dir db/migrations postgres "host=127.0.0.1 port=5433 user=postgres password=postgres dbname=microservice_db sslmode=disable" status

.PHONY: migrate-reset
migrate-reset:
		@echo "Resetting database..."
		goose -dir db/migrations postgres "host=127.0.0.1 port=5433 user=postgres password=postgres dbname=microservice_db sslmode=disable" reset