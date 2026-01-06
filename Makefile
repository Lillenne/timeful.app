# Timeful Docker Makefile
# Convenience commands for managing the Docker deployment

.PHONY: help setup build up down restart logs clean backup restore

help: ## Show this help message
	@echo "Timeful Docker Commands:"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

setup: ## Run initial setup (creates .env from template)
	@./docker-setup.sh

build: ## Build the Docker images
	docker compose build

up: ## Start the application
	docker compose up -d
	@echo ""
	@echo "✅ Timeful is starting!"
	@echo "Access the application at: http://localhost:3002"

down: ## Stop the application
	docker compose down

restart: ## Restart the application
	docker compose restart

logs: ## View application logs
	docker compose logs -f

logs-app: ## View only app logs
	docker compose logs -f app

logs-db: ## View only database logs
	docker compose logs -f mongodb

clean: ## Stop and remove all containers and volumes (⚠️  deletes data!)
	docker compose down -v
	@echo "⚠️  All data has been deleted!"

backup: ## Create a backup of the MongoDB database
	@mkdir -p ./backups
	@echo "Creating backup..."
	docker compose exec mongodb mongodump --db=schej-it --out=/data/backup
	docker cp timeful-mongodb:/data/backup ./backups/backup-$$(date +%Y%m%d-%H%M%S)
	@echo "✅ Backup created in ./backups/"

restore: ## Restore from the latest backup (⚠️  overwrites data!)
	@LATEST=$$(ls -t backups/ | head -1); \
	if [ -z "$$LATEST" ]; then \
		echo "❌ No backups found in ./backups/"; \
		exit 1; \
	fi; \
	echo "Restoring from backup: $$LATEST"; \
	docker cp ./backups/$$LATEST timeful-mongodb:/data/restore; \
	docker compose exec mongodb mongorestore --db=schej-it /data/restore/schej-it --drop
	@echo "✅ Restore completed"

shell-app: ## Open a shell in the app container
	docker compose exec app sh

shell-db: ## Open MongoDB shell
	docker compose exec mongodb mongosh schej-it

dev: ## Start in development mode with live logs
	docker compose up

dev-build: ## Build and start in development mode
	docker compose up --build

pull: ## Pull latest changes and restart
	git pull
	docker compose down
	docker compose build --no-cache
	docker compose up -d
	@echo "✅ Updated to latest version"

status: ## Show status of all containers
	docker compose ps

.DEFAULT_GOAL := help
