# Repository Guidelines

## Project Overview

CodeZone is a web chat MVP built as a modular monolith: a single Go HTTP/WebSocket API, a React single-page app, and MySQL. Conversations support direct and group chats between human users and a built-in Mock Agent that echoes messages back. REST is the write path; WebSocket (`/api/ws`) is delivery-only. Both authenticate against the same server-side session stored in `sessions.token_hash`.

Authoritative design docs live in `docs/superpowers/`: `specs/` holds product/behavior specs and `plans/` holds task-by-task implementation plans. For larger behavior changes, check the matching spec first; when spec and plan disagree, the spec wins for behavior and the plan wins for structure. `CLAUDE.md` contains additional architecture notes, but its "pre-implementation" status section is stale — the implementation described there now exists on `main`.

## Project Structure & Module Organization

- `backend/cmd/api` — HTTP/WebSocket API entrypoint (wiring, router, auth middleware).
- `backend/cmd/migrate` — applies SQL migrations (and dev seed data when `dev_seed: true`).
- `backend/internal/*` — domain packages:
  - `auth` — dev-login and session tokens (tokens are random 32 bytes, base64url; stored as `sha256(token + session_secret)`).
  - `users` — profiles, avatar upload; repository/service/handler triad.
  - `conversations` — direct/group conversations and membership; repository/service/handler triad.
  - `messages` — message validation, persistence, history pagination, read markers, unread counts. **All message writes go through this package.**
  - `realtime` — in-process WebSocket hub and event fan-out only; no SQL beyond broadcast targeting.
  - `agent` — Mock Agent reply policy. Replies asynchronously via the `messages` service (`MessageSender` interface), never writes SQL directly. This is the extension point for real agents later.
  - `admin` — read-only admin listings (users, conversations, messages).
  - `db` — connection setup and the migration runner (splits files on `;`, runs them in filename order).
  - `httpx` — shared JSON/error helpers and user-id context middleware (`httpx.WithUserID` / `httpx.UserID`).
  - `config` — YAML config loading and validation.
  - `storage` — Aliyun OSS avatar storage (via the private `codeup.aliyun.com/.../golib` module); returns nil when OSS config is incomplete, so avatar upload degrades gracefully.
- `backend/migrations` — ordered SQL files (`0001_schema.sql`, ...). Files containing `_seed_dev` are skipped unless `dev_seed: true` in config.
- `frontend/src/app` — application wiring (`App.tsx` switches between login/chat/admin by token state and path).
- `frontend/src/lib` — shared clients: `api.ts` (all REST request/response shapes), `ws.ts` (WebSocket reconnect + event dispatch), `query.ts` (React Query client).
- `frontend/src/features/*` — feature modules: `auth` (LoginPage), `chat` (ChatPage, ConversationList, MessageList, MessageComposer, MarkdownMessage, ProfileDrawer, etc.), `admin` (AdminPage).
- `frontend/tests` — Playwright end-to-end tests (`*.e2e.ts`).

## Build, Test, and Development Commands

The root `Makefile` wraps the common flows (`make help` lists targets):

- `make mysql` — start local MySQL 8.4 (docker compose; user/password/db all `chat`, port 3306).
- `make migrate` — apply migrations (auto-copies `backend/config.local.yaml` to `backend/config.yaml` if missing).
- `make backend` — start the API on `:8080` with `ENABLE_DEV_LOGIN=true`.
- `make frontend` — start Vite on `:5173` (proxies `/api`, including WebSocket, to `:8080`).
- `make run` — migrate, then start API and frontend together.
- `make test` — backend and frontend unit tests.
- `make stop` — stop the MySQL container.

Equivalent raw commands:

```bash
docker compose up -d mysql
cd backend && go run ./cmd/migrate
cd backend && ENABLE_DEV_LOGIN=true go run ./cmd/api
cd backend && go test ./...                                     # all backend tests
cd backend && go test ./internal/messages -run TestX -count=1   # single test
cd frontend && npm install
cd frontend && npm run dev
cd frontend && npm run lint          # tsc -b --noEmit (this is the only "linter")
cd frontend && npm run test          # Vitest unit tests
cd frontend && npm run build         # tsc -b && vite build
cd frontend && npm run test:e2e      # Playwright; requires MySQL + API + Vite running
```

## Configuration

- The backend reads a YAML file: `./config.yaml` by default (git-ignored), or the path in the `CODEZONE_CONFIG` env var. Versioned templates: `config.local.yaml`, `config.test.yaml`, `config.gray.yaml`, `config.prod.yaml`. Copy the local template to get started: `cp backend/config.local.yaml backend/config.yaml`.
- Required keys (validated at startup): `port`, `mysql_dsn`, `session_secret`, `cors_origin`. Optional flags: `dev_seed` (apply seed migration), `enable_dev_login` (expose `POST /api/auth/dev-login`). The `oss` block configures avatar storage; all values are `REPLACE_ME` placeholders outside of real deployments.
- Production/gray/test templates keep `dev_seed: false` and `enable_dev_login: false`; their secrets must be injected by the deployment pipeline, never committed.
- Frontend env vars go in `frontend/.env` (see `.env.example`) and must be prefixed with `VITE_`. They are not needed for local dev because Vite proxies `/api`.

## Runtime Architecture

- Auth: dev-login issues a bearer token valid 24h. Protected routes accept `Authorization: Bearer <token>` or an `access_token` query parameter (used by the WebSocket client). The dev-login route is gated solely by the `enable_dev_login` YAML key (`config.local.yaml` sets it true). Note the `ENABLE_DEV_LOGIN=true` env var that the Makefile/README pass is not read by any Go code — it is vestigial.
- Messaging flow: `POST /api/conversations/{id}/messages` persists via the `messages` service, which then publishes `message.created` / `conversation.updated` events through the realtime hub (per-user subscriptions, bounded backlog of 256 events) and triggers `agent.Orchestrator.MaybeReply` for conversations containing an enabled mock agent (async, 5s timeout).
- Frontend recovery: on WebSocket disconnect the client shows a status pill, reconnects after ~1.2s, and refetches the conversation list plus recent messages. Do not add sync cursors or per-device state.

Key invariants to preserve:

- Persisted messages have **no** server-side `failed` state; sending/failed are client-only UI states for optimistic messages.
- Unread state lives on `conversation_members` (`unread_count`, `last_read_message_id`), not per-message.
- Agents are ordinary `users` rows with `user_type='agent'` plus a matching `agent_profiles` row.
- All records carry `workspace_id`; the MVP seeds one default workspace and does not expose switching.

## Coding Style & Naming Conventions

- Go: use `gofmt`. Keep packages small and domain-oriented under `backend/internal`; follow the established repository/service/handler separation. Test files use `*_test.go`.
- TypeScript/React: strict mode (`tsconfig` has `"strict": true`). Components in `PascalCase` (`ChatPage.tsx`); shared utilities as concise camel-case modules (`api.ts`, `ws.ts`). Keep all REST shapes in `frontend/src/lib/api.ts`. Feature view components stay presentational; data fetching goes through React Query.
- Migrations: add a new sequentially-numbered `NNNN_description.sql` file; never edit applied migrations. Statements must be idempotent (`CREATE TABLE IF NOT EXISTS`, `ON DUPLICATE KEY UPDATE`, or information-schema guards as in `0001_schema.sql`) because migrations re-run on every migrate invocation.

## Testing Guidelines

- Backend tests use Go's standard `testing` package: `go test ./...`. The MySQL-backed integration test in `backend/internal/messages/repository_integration_test.go` skips itself unless `CHAT_TEST_MYSQL_DSN` (or `MYSQL_DSN`) is set — point it at a scratch database, e.g. the docker-compose MySQL.
- Frontend unit tests use Vitest + React Testing Library with jsdom, colocated as `*.test.ts` / `*.test.tsx`. Run with `npm run test`; a single file with `npm run test -- --run src/lib/api.test.ts`.
- E2E tests live in `frontend/tests/*.e2e.ts` (Playwright, chromium + mobile projects, `workers: 1`). They require MySQL, the API with dev login enabled, and the Vite dev server all running.
- When adding behavior, add focused tests for message persistence, WebSocket recovery, auth, and API contracts.

## Commit & Pull Request Guidelines

- Recent history uses Conventional Commit prefixes (`feat:`, `test:`, `chore:`). Keep subjects short and imperative, e.g. `feat: add websocket recovery in frontend`.
- PRs should include a behavior summary, tests run, linked issues or planning docs when relevant, and screenshots for UI changes. Note database migration or environment impacts.

## Security Considerations

- `ENABLE_DEV_LOGIN=true` / `enable_dev_login: true` is for local development only; it must stay off in test/gray/prod configs.
- Never commit real secrets: `backend/config.yaml` is git-ignored for this reason, and `config.*.yaml` templates keep `REPLACE_ME` placeholders. Do not commit local database credentials or OSS keys.
- Session tokens are only stored hashed (`sha256(token + session_secret)`); keep it that way if touching auth code.
- Document new required environment variables or config keys in `README.md` or an example env file.
