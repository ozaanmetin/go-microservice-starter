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