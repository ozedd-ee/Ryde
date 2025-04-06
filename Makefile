# Default target
.PHONY: up down restart logs build

# Start all services
default: up

# Create and start all service containers
up:
	docker compose up -d

# Stop and remove containers, networks
down:
	docker compose down

# Restart service containers
restart:
	docker compose restart -d

# View output from containers
logs:
	docker compose logs -f

# Build or rebuild services
build:
	docker compose build

# List containers
ps:
	docker compose ps

# Stop and remove containers, networks
clean:
	docker compose down -v
