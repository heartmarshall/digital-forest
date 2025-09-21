include .env
export

# --- Docker Compose команды ---
up:
	@echo "Starting application and database via Docker Compose..."
	@docker-compose up -d

down:
	@echo "Stopping and removing containers..."
	@docker-compose down

logs:
	@echo "Following application logs..."
	@docker-compose logs -f app

build:
	@echo "Building Docker images..."
	@docker-compose build

# --- Миграции БД (выполняются ВНУТРИ контейнера) ---
db/create-migration:
	@read -p "Enter migration name: " name; \
	docker-compose run --rm app goose -dir ./migrations create $$name sql

db/migrate-up:
	@echo "Applying database migrations..."
	@docker-compose run --rm app goose -dir ./migrations postgres "$(POSTGRES_URL)" up

db/migrate-down:
	@echo "Rolling back the last database migration..."
	@docker-compose run --rm app goose -dir ./migrations postgres "$(POSTGRES_URL)" down

db/status:
	@echo "Checking migration status..."
	@docker-compose run --rm app goose -dir ./migrations postgres "$(POSTGRES_URL)" status

.PHONY: up down logs build db/create-migration db/migrate-up db/migrate-down db/status