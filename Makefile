# Variables
API_IMAGE_NAME := my-api-image
SCRAPPER_IMAGE_NAME := my-scrapper-image

# Build API Docker image
build-api:# docker build -t $(API_IMAGE_NAME) -f api/Dockerfile .
	docker-compose build api
	
# Enter the API Docker container
enter-api:
	docker-compose exec api sh	

# Build Scrapper Docker image
build-scrapper:
	docker build -t $(SCRAPPER_IMAGE_NAME) ./scrapper

# Run the gRPC servers
run-grpc:
	@echo "Starting gRPC servers..."
	# Add command to start the gRPC servers here

# Start the API service using Docker Compose
run-api:
	@echo "Starting API service..."
	docker-compose up --build api

# Stop and remove the containers
stop:
	docker-compose down

# Build the Docker images and start the services
start: build-api build-scrapper run-grpc run-api

# Help command to display available commands
help:
	@echo "Available commands:"
	@echo "  build-api        : Build the API Docker image"
	@echo "  build-scrapper   : Build the Scrapper Docker image"
	@echo "  run-grpc         : Run the gRPC servers"
	@echo "  run-api          : Start the API service using Docker Compose"
	@echo "  stop             : Stop and remove the running containers"
	@echo "  start            : Build the Docker images and start the gRPC servers and API service"
	@echo "  help             : Display available commands"

.PHONY: build-api build-scrapper run-grpc run-api stop start help
