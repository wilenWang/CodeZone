# Makefile development workflow design

## Goal

Provide root-level Make targets for the CodeZone local development workflow. `make run` assumes MySQL is already running and orchestrates only the application services.

## Targets

- `help`: print available targets and descriptions.
- `mysql`: run `docker compose up -d mysql`.
- `migrate`: ensure the backend local config exists, then run `go run ./cmd/migrate` from `backend`.
- `backend`: ensure config exists, then run the Go API from `backend` with `ENABLE_DEV_LOGIN=true`.
- `frontend`: run the Vite development server from `frontend`.
- `run`: ensure config exists, verify MySQL is reachable, apply migrations, launch the API in the background, and run Vite in the foreground. An interrupt or Vite exit terminates only the API process started by this invocation.
- `test`: run `go test ./...` in `backend`, then `npm run test` in `frontend`.
- `stop`: run `docker compose stop mysql` only; it does not use broad process matching to kill Go or Node processes.

## Configuration and prerequisites

A reusable Make recipe creates `backend/config.yaml` from `backend/config.local.yaml` only if the ignored local config file does not already exist. Existing local changes are never overwritten.

`run` deliberately does not install npm dependencies and does not start MySQL. It checks `docker compose exec -T mysql mysqladmin ping` before migrations and exits with an instruction to run `make mysql` when the database is unavailable.

## Process handling

`run` captures the API child PID, installs a shell `trap` for normal exit and interrupts, and waits on the foreground frontend process. Cleanup sends a termination signal only to that captured child PID. MySQL remains running.

## Validation

- `make help` renders target descriptions.
- `make -n run` confirms the expected sequence without starting services.
- With MySQL running and frontend dependencies installed, `make run` exposes the API on port 8080 and Vite on port 5173; Ctrl+C stops the API and frontend while MySQL remains available.
- `make test` runs both existing test suites.
