# Set default environment variables
COMPOSE_FILE=docker-compose.yml

# Run everything
.PHONY: up
up:
	@echo "Starting MCS services..."
	docker-compose -f $(COMPOSE_FILE) up --build -d

# Stop and remove containers
.PHONY: down
down:
	@echo "Stopping MCS services..."
	docker-compose -f $(COMPOSE_FILE) down

# Restart services
.PHONY: restart
restart:
	@echo "Restarting MCS services..."
	docker-compose -f $(COMPOSE_FILE) down
	docker-compose -f $(COMPOSE_FILE) up --build -d

# View running containers
.PHONY: ps
ps:
	@echo "Listing running MCS containers..."
	docker ps

# View logs from all services
.PHONY: logs
logs:
	@echo "Fetching logs from all MCS services..."
	docker-compose -f $(COMPOSE_FILE) logs -f

# View logs from a specific service (use: make logs SERVICE=service_name)
.PHONY: service-logs
service-logs:
	@echo "Fetching logs for $(SERVICE)..."
	docker-compose -f $(COMPOSE_FILE) logs -f $(SERVICE)

# Connect to PostgreSQL inside the container
.PHONY: psql
psql:
	@echo "Connecting to PostgreSQL..."
	docker exec -it postgres_mcs psql -U mcs_user -d mcs_db

# Connect to MongoDB shell
.PHONY: mongo
mongo:
	@echo "Connecting to MongoDB..."
	docker exec -it mongo_mcs mongosh

# Clean up everything (containers, images, volumes, networks)
.PHONY: clean
clean:
	@echo "Removing all containers, volumes, and networks..."
	docker-compose -f $(COMPOSE_FILE) down -v
	docker system prune -a -f
	docker volume prune -f

# Check Docker system info
.PHONY: info
info:
	@echo "Checking Docker system information..."
	docker system df
