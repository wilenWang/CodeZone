# Repository Guidelines

## Project Structure & Module Organization

This repository is a web chat MVP with a Go API, React SPA, and MySQL.

- `backend/cmd/api` starts the HTTP/WebSocket API; `backend/cmd/migrate` applies migrations.
- `backend/internal/*` holds domain packages: `auth`, `users`, `conversations`, `messages`, `realtime`, `agent`, `admin`, `db`, and `httpx`.
- `backend/migrations` contains ordered SQL migrations.
- `frontend/src/app` contains application wiring and routes; `frontend/src/lib` contains shared API, query, and WebSocket clients.
- `frontend/src/features/*` contains feature modules; `frontend/tests` contains Playwright end-to-end tests.
- `docs/superpowers` contains design specs and implementation plans for larger behavior changes.

## Build, Test, and Development Commands

- `docker compose up -d mysql` starts the local MySQL service.
- `cd backend && go run ./cmd/migrate` applies migrations and development seed data.
- `cd backend && ENABLE_DEV_LOGIN=true go run ./cmd/api` starts the API on `:8080`.
- `cd backend && go test ./...` runs all backend tests.
- `cd frontend && npm install` installs frontend dependencies.
- `cd frontend && npm run dev` starts Vite on port `5173`.
- `cd frontend && npm run lint` runs TypeScript checks.
- `cd frontend && npm run test` runs Vitest unit tests.
- `cd frontend && npm run test:e2e` runs Playwright; MySQL, API, and Vite must be running.
- `cd frontend && npm run build` type-checks and builds the frontend.

## Coding Style & Naming Conventions

Use `gofmt` for Go. Keep backend packages small and domain-oriented under `backend/internal`; prefer repository/service/handler separation where established. Go test files use `*_test.go`.

Frontend code uses TypeScript and React. Name components in `PascalCase` (`ChatPage.tsx`) and shared utilities as concise camel-case modules (`api.ts`, `ws.ts`). Keep REST shapes in `frontend/src/lib/api.ts`.

## Testing Guidelines

Backend tests use Go’s standard `testing` package. Frontend unit tests use Vitest with React Testing Library, colocated as `*.test.ts` or `*.test.tsx`. End-to-end coverage lives in `frontend/tests/*.e2e.ts`. Add focused tests for message persistence, WebSocket recovery, auth, and API contracts.

## Commit & Pull Request Guidelines

Recent history uses Conventional Commit prefixes such as `feat:`, `test:`, and `chore:`. Keep subjects short and imperative, for example `feat: add websocket recovery in frontend`.

Pull requests should include a behavior summary, tests run, linked issues or planning docs when relevant, and screenshots for UI changes. Note database migration or environment impacts.

## Security & Configuration Tips

Use `ENABLE_DEV_LOGIN=true` only locally. Do not commit secrets or local database credentials. Document new required variables in README or an example env file.
