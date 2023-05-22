.PHONY: build run clean proto lint test coverage help

# Define the executable names
API_EXECUTABLE := api
SCRAPPER_EXECUTABLE := scrapper
DATABASE_EXECUTABLE := database

# Define the directories
CMD_API_DIR := cmd/api
CMD_SCRAPPER_DIR := cmd/scrapper
CMD_DATABASE_DIR := cmd/database
PROTO_DIR := proto
INTERNAL_API_DIR := internal/api
INTERNAL_SCRAPPER_DIR := internal/scrapper
INTERNAL_DATABASE_DIR := internal/database

# Generate protobuf files
proto:
	protoc --go_out=. --go_opt=module=github.com/brochadoluis/temperature-exercise \
		--go-grpc_out=. --go-grpc_opt=module=github.com/brochadoluis/temperature-exercise \
		--proto_path=./proto ./proto/temperature.proto

# Default target
build-all: build-api build-scrapper build-database build-mongodb

# Run all containers target
run-all:
	docker-compose up -d

# Run the API component
run-api:
	docker-compose up api

# Run the Scrapper component
run-scrapper:
	docker-compose up scrapper

# Run the Database component
run-database:
	docker-compose up database

# Clean up the project
clean:
	docker-compose down -v
	rm -f $(CMD_API_DIR)/$(API_EXECUTABLE)
	rm -f $(CMD_SCRAPPER_DIR)/$(SCRAPPER_EXECUTABLE)
	rm -f $(CMD_DATABASE_DIR)/$(DATABASE_EXECUTABLE)

# Run linter
lint:
	docker-compose run --rm lint

# Run tests
test:
	docker-compose run --rm test

# Run tests and generate coverage report
coverage:
	docker-compose run --rm coverage


# Build the API container
build-api:
	@docker-compose build api

# Build the Scrapper container
build-scrapper:
	@docker-compose build scrapper

# Build the Database container
build-database:
	@docker-compose build database

# Run the API container
run-api:
	@docker-compose up api

# Run the Scrapper container
run-scrapper:
	@docker-compose up scrapper

# Run the Database container
run-database:
	@docker-compose up database

# Run linter
lint:
	docker-compose run --rm temperature-exercise golangci-lint run

# Run tests
test:
	docker-compose run --rm temperature-exercise go test -v ./...

# Run tests and generate coverage report
coverage:
	docker-compose run --rm temperature-exercise sh -c "go test -v -coverpkg=./... -coverprofile=coverage.out ./... && go tool cover -html=coverage.out -o coverage.html"

# Display help message
help:
	@echo "Available commands:"
	@echo "  all             : Build all containers (default target)"
	@echo "  clean           : Clean up the project"
	@echo "  proto           : Generate protobuf files"
	@echo "  run-api         : Run the API container"
	@echo "  run-scrapper    : Run the Scrapper container"
	@echo "  run-database    : Run the Database container"
	@echo "  lint            : Run linter"
	@echo "  test            : Run tests"
	@echo "  coverage        : Run tests and generate coverage report"
	@echo "  help            : Show this help message"
