.DEFAULT_GOAL := help

.PHONY: help mysql ensure-config check-mysql migrate backend frontend run test stop

help: ## Show available development commands.
	@awk 'BEGIN {FS = ":.*##"; printf "Usage: make <target>\n\nTargets:\n"} /^[a-zA-Z_-]+:.*##/ {printf "  %-14s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

mysql: ## Start the local MySQL container.
	docker compose up -d mysql

ensure-config:
	@test -f backend/config.yaml || cp backend/config.local.yaml backend/config.yaml

check-mysql:
	@nc -z 127.0.0.1 3306 >/dev/null 2>&1 \
		|| { echo "MySQL is not available on 127.0.0.1:3306. Start it first with: make mysql"; exit 1; }

migrate: ensure-config ## Apply backend database migrations.
	cd backend && go run ./cmd/migrate

backend: ensure-config ## Start the Go API with development login enabled.
	cd backend && ENABLE_DEV_LOGIN=true go run ./cmd/api

frontend: ## Start the Vite development server.
	cd frontend && npm run dev

run: ensure-config check-mysql migrate ## Apply migrations, then run the API and frontend.
	@set -e; \
	cleanup() { \
		if [ -n "$$api_pid" ] && kill -0 "$$api_pid" 2>/dev/null; then \
			echo "Stopping API..."; \
			kill "$$api_pid"; \
			wait "$$api_pid" 2>/dev/null || true; \
		fi; \
	}; \
	trap cleanup EXIT INT TERM; \
	(cd backend && exec env ENABLE_DEV_LOGIN=true go run ./cmd/api) & \
	api_pid=$$!; \
	echo "API started (PID $$api_pid). Starting frontend..."; \
	cd frontend && npm run dev

test: ## Run backend and frontend unit tests.
	cd backend && go test ./...
	cd frontend && npm run test

stop: ## Stop the local MySQL container.
	docker compose stop mysql
