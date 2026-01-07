.PHONY: help docker-up docker-down run start

help:
	@echo "Available commands:"
	@echo "  make start           - Start Docker and run the app (migrations run automatically)"
	@echo "  make docker-up       - Start PostgreSQL in Docker"
	@echo "  make docker-down     - Stop PostgreSQL Docker container"
	@echo "  make run             - Run the application (migrations run automatically)"

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

run:
	go run main.go

start:
	@echo "Starting PostgreSQL..."
	docker-compose up -d
	@echo "Waiting for PostgreSQL to be ready..."
	@sleep 5
	@echo "Starting application (migrations will run automatically)..."
	go run main.go
