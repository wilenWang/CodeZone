# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Status

This repo is a **pre-implementation** MVP for a web chat app. Only design docs and empty package scaffolding (`backend/internal/{conversations,messages}/`) exist on `main`. Active implementation happens on the `chat-app-implementation` worktree under `.worktrees/`.

The authoritative documents are:

- `docs/superpowers/specs/2026-07-03-chat-app-design.md` — product scope, data model, REST/WS API surface, non-goals.
- `docs/superpowers/plans/2026-07-03-chat-app-implementation.md` — task-by-task build order with concrete file paths, code snippets, and shell commands. Follow this plan when implementing.

Before adding new features, cross-check both. If the plan and spec disagree, the spec wins for behavior; the plan wins for structure.

## Architecture

Modular monolith: single Go service + React SPA + MySQL. REST is the write path; WebSocket is delivery-only. Both authenticate against the same server-side session (`sessions.token_hash`).

Package boundaries (backend, `backend/internal/*`) are load-bearing — the plan calls them out explicitly:

- `messages` — owns validation, persistence orchestration, history pagination, read markers, unread updates. All message writes go through here.
- `realtime` — WebSocket connections and event fan-out only. No SQL beyond looking up conversation members for broadcast targeting.
- `agent` — Mock Agent reply policy. Must call the `messages` service (via `MessageSender` interface) rather than writing SQL directly. This is the extension point for real agents later.
- `conversations`, `users`, `auth`, `admin` — normal repo/service/handler triads.
- `httpx` — shared JSON/error helpers plus context-based user-id middleware (`httpx.WithUserID` / `httpx.UserID`).

Frontend boundaries (`frontend/src/`):

- `lib/api.ts` — all REST request/response shape.
- `lib/ws.ts` — WebSocket reconnect + event dispatch.
- Feature dirs (`features/chat`, `features/auth`, `features/admin`) host page components and their hooks; view components stay presentational.

Key invariants to preserve:

- Persisted messages have **no** server-side `failed` state. Sending/failed are client-only UI states for optimistic messages.
- Unread state lives on `conversation_members`, not per-message.
- Agents are ordinary rows in `users` with `user_type='agent'` + a matching `agent_profiles` row.
- All records carry `workspace_id`. MVP seeds one default workspace and does not expose switching.

## Realtime and recovery

On WebSocket disconnect, the client shows a status pill and refetches the conversation list plus the active conversation's recent messages on reconnect. Do not add sync cursors or per-device state.

## Development commands

Full command reference lives in the plan; the common ones once implementation lands:

```bash
# MySQL (from repo root)
docker compose up -d mysql

# Backend
cd backend
go run ./cmd/migrate        # apply migrations + seed
go run ./cmd/api             # start API on :8080
go test ./...                # all backend tests
go test ./internal/messages -run TestPlainTextFromMarkdown -count=1   # single test

# Frontend
cd frontend
npm run dev                  # Vite on :5173, proxies /api → :8080
npm run lint                 # tsc -b --noEmit
npm run test                 # vitest
npm run test -- --run src/lib/api.test.ts  # single test file
npm run test:e2e             # Playwright (needs API + dev server running)
npm run build
```

Env config is loaded via `internal/config` from environment variables with dev-friendly defaults; see `.env.example` (created in plan Task 1) for the full list.

## Working with the plan

The plan uses `- [ ]` checkboxes and expects TDD: each task writes a failing test, then implements to pass. When picking up work, check the `chat-app-implementation` worktree for progress before re-doing a task from scratch.
