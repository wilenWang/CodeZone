# Web Chat App Design

Date: 2026-07-03

## Summary

Build a maintainable MVP for a web chat application. The first version supports seed-user login, direct chats, basic group chats, Markdown text messages, lightweight realtime updates, conversation unread counts, a responsive React interface, and a simple development admin page.

The system should reserve clean extension points for future Agent users. In the MVP, one Mock Agent is seeded as an `agent` user and can auto-reply with simple deterministic behavior. Real task execution, tools, memory, and advanced agent orchestration are out of scope for the first release.

## Goals

- Provide a working web chat experience with users, conversations, messages, history, and realtime delivery.
- Use React for the frontend, Go for the backend, and MySQL for persistence.
- Keep the backend as a modular monolith: one deployable Go service with clear internal package boundaries.
- Support both direct chats and basic group chats through one conversation model.
- Store all core records under a default workspace so future organization and tenant support can be added without rewriting the schema.
- Treat Agent accounts as a user type from day one.
- Prioritize local development with Docker Compose.

## Non-Goals

- Public registration, SMS login, OAuth, or production account recovery.
- Image and file attachments.
- Online presence, typing indicators, delivery receipts, or per-message read receipts.
- Multi-workspace switching in the UI.
- Real Agent task execution, tool calls, long-term memory, or task lifecycle management.
- PWA install, push notifications, or offline support.
- Production-grade operations admin, moderation, export, or audit tooling.

## Product Scope

### Authentication

The MVP uses seed users. The login page supports:

- One-click development login by choosing a seed user.
- Username/password login shape, backed by seed credentials, so the backend auth boundary is not throwaway.

Sessions are stored server-side as token hashes with expiration.

### Conversations

Direct chats and group chats share the same model:

- A direct chat is a conversation with two members.
- A group chat is a conversation with more than two members and an optional title.
- Conversation list items show title, recent message, update time, and unread count.

### Messages

Messages support text with safe Markdown rendering:

- Paragraphs, links, lists, inline code, and code blocks.
- Raw HTML is not allowed.
- The database stores Markdown source and a plain-text summary.

Message UI states:

- Sending.
- Failed.
- Sent.

The MVP does not show per-message read receipts.

### Unread Counts

Unread state is tracked per conversation member, not on each message. `conversation_members` stores a reading position and unread count. Opening a conversation marks it read through a REST endpoint.

### Mock Agent

The seed data includes one Agent user:

- `users.user_type = "agent"`.
- `agent_profiles.kind = "mock"`.
- The Mock Agent can be added to a conversation.
- When a human sends a message in a conversation containing the enabled Mock Agent, the backend asynchronously creates a simple Agent reply and broadcasts it.

This validates the future design where human users and Agent users both participate as conversation members.

### Admin Page

The admin page is for development diagnostics, not operations:

- View seed users.
- View conversations.
- View recent messages.

No formal permissions, moderation, user bans, or audit export are included in the MVP.

## Data Model

### `workspaces`

- `id`
- `name`
- `created_at`

The MVP seeds one default workspace.

### `users`

- `id`
- `workspace_id`
- `username`
- `display_name`
- `avatar_url`
- `user_type`: `human` or `agent`
- `password_hash`, nullable where appropriate for development login
- `created_at`

Both humans and Agents live in this table so conversations, permissions, and membership logic stay unified.

### `sessions`

- `id`
- `user_id`
- `token_hash`
- `expires_at`
- `created_at`

### `conversations`

- `id`
- `workspace_id`
- `type`: `direct` or `group`
- `title`, nullable for direct chats
- `created_by`
- `last_message_id`
- `last_message_at`
- `created_at`

`last_message_id` and `last_message_at` make conversation list queries simple and fast.

### `conversation_members`

- `conversation_id`
- `user_id`
- `role`
- `last_read_message_id`, nullable
- `unread_count`
- `joined_at`

This table owns unread state.

### `messages`

- `id`
- `conversation_id`
- `sender_id`
- `content_markdown`
- `content_plain`
- `created_at`
- `edited_at`, nullable

Server-persisted messages are considered sent. Sending and failed states are client-side UI states for optimistic messages until the server confirms or rejects the request. The MVP should not persist failed messages as normal chat records.

### `agent_profiles`

- `user_id`
- `kind`: initially `mock`
- `config_json`
- `enabled`
- `updated_at`

### Deferred Tables

Do not build these for the MVP unless implementation reveals a direct need:

- Attachments.
- OAuth identities.
- Full audit logs.
- Agent tasks.
- Tool call logs.
- Message reactions.
- Per-user device cursors.

Future Agent task and tool events can be added later through dedicated task tables or a `message_events` style append-only log.

## Backend Design

The backend is one Go service that exposes REST APIs and a WebSocket endpoint.

### Package Layout

- `cmd/api`: service entrypoint, configuration load, router setup, lifecycle.
- `internal/config`: environment parsing and typed configuration.
- `internal/db`: MySQL connection, migrations, transaction helpers.
- `internal/auth`: session token handling, seed-user login, username/password login boundary.
- `internal/users`: user queries and human/agent user handling.
- `internal/conversations`: direct/group conversation creation, membership, lists.
- `internal/messages`: send message, history pagination, read markers, unread updates.
- `internal/realtime`: WebSocket authentication, connection registry, user/conversation broadcasts.
- `internal/agent`: Mock Agent detection and reply generation.
- `internal/admin`: development admin APIs.
- `internal/httpx`: request parsing, response writing, middleware, error format.

### REST API

- `POST /api/auth/login`
- `POST /api/auth/dev-login`
- `POST /api/auth/logout`
- `GET /api/me`
- `GET /api/users`
- `GET /api/conversations`
- `POST /api/conversations`
- `GET /api/conversations/:id/messages?before=&limit=`
- `POST /api/conversations/:id/messages`
- `POST /api/conversations/:id/read`
- `GET /api/admin/users`
- `GET /api/admin/conversations`
- `GET /api/admin/messages`

### WebSocket API

- `GET /api/ws`

The WebSocket connection is authenticated with the same session token model as REST.

Server-pushed event types:

- `message.created`
- `conversation.updated`
- `message.failed`, only when needed for server-originated async work such as Mock Agent failure

The MVP sends messages through REST, not WebSocket. This keeps validation, persistence, transactions, and error handling simpler. After a successful REST write, the backend broadcasts realtime events.

### Mock Agent Flow

1. A human sends a message through `POST /api/conversations/:id/messages`.
2. `messages` persists the message and updates conversation metadata and unread counts.
3. `realtime` broadcasts `message.created` and `conversation.updated`.
4. `agent` checks whether the conversation has an enabled Mock Agent member.
5. If yes, `agent` schedules an asynchronous mock reply.
6. The mock reply is persisted as a normal message from the Agent user.
7. `realtime` broadcasts the Agent message.

The `agent` module should expose an interface that can later be backed by a real Agent runner without changing the message API.

## Frontend Design

Use React with Vite and TypeScript.

### Routes

- `/login`: seed-user login and username/password login form.
- `/app`: main chat interface.
- `/admin`: development admin page.

### Main Layout

Desktop:

- Left column: conversation list.
- Right column: active conversation.
- Conversation details and members appear in a drawer or modal.

Mobile:

- Conversation list as the first view.
- Tapping a conversation navigates to the chat view.
- Chat header includes a back button.

### Core Components

- `ConversationList`: search, conversation rows, unread counts, recent message preview.
- `ChatHeader`: conversation title and member/detail entry point.
- `MessageList`: historical pagination, new-message insertion, scroll behavior.
- `MessageBubble`: Markdown rendering and human/agent visual distinction.
- `MessageComposer`: input, send action, retry failed message.
- `ConversationDrawer`: members and conversation metadata.
- `NewConversationDialog`: create direct or group conversations.
- `AdminDashboard`: users, conversations, and messages.

### State Management

- Use React Query for REST data: current user, users, conversations, messages, admin data.
- Use a small WebSocket client to receive server events and update React Query caches.
- Keep UI state local where possible: current conversation, drawer, dialogs, pending messages.
- For sending messages, create an optimistic local message, replace it with the server message on success, and mark it failed on error.

### Markdown Rendering

Use a safe Markdown renderer:

- Raw HTML disabled.
- Links sanitized.
- Code blocks supported.

Human and Agent messages use the same renderer.

## Realtime and Recovery

- REST is the source of truth for writes.
- WebSocket is used for event delivery only.
- If the WebSocket disconnects, the UI shows a lightweight connection status.
- On reconnect, the client refetches the conversation list and the active conversation's recent messages.
- This strategy avoids complex multi-device sync cursors in the MVP while still recovering from missed events.

## Error Handling

REST errors use one response shape:

```json
{
  "code": "message_send_failed",
  "message": "Could not send message",
  "details": {}
}
```

Frontend behavior:

- Failed sends remain visible as local failed messages.
- Users can retry failed sends.
- WebSocket connection issues do not block REST sending.
- Mock Agent async failures are logged server-side; the MVP may optionally surface a lightweight system event later.

## Development Environment

Use Docker Compose for local dependencies:

- MySQL database.
- Optional API and frontend containers if desired.

The first implementation should also support running Go and React directly on the host:

- Go API reads `.env`.
- React dev server proxies `/api` to Go.
- Database migration command is separate from seed command.

Expected developer commands:

- Start MySQL through Docker Compose.
- Run migrations.
- Seed development users, conversations, and Mock Agent.
- Start Go API.
- Start React dev server.

Provide `.env.example` with:

- MySQL DSN parts.
- Session secret.
- CORS origin.
- Development seed toggle.

## Testing Strategy

### Backend

Unit tests:

- Auth/session handling.
- Conversation creation rules.
- Message validation and plain-text extraction.
- Unread count updates.
- Mock Agent reply generation.

Integration tests:

- Login.
- Create direct chat.
- Create group chat.
- Send message.
- Load paginated history.
- Mark conversation read.
- Mock Agent auto-reply.

WebSocket tests:

- Reject unauthenticated connections.
- Broadcast new messages to conversation members.
- Client can recover missed messages through REST after reconnect.

### Frontend

Component tests:

- Conversation list renders unread state.
- Message list renders Markdown safely.
- Message composer handles sending, success, failure, and retry.
- Mobile layout switches between list and chat views.

End-to-end test:

- User A logs in.
- User A sends a message to User B.
- User B receives it.
- User A chats with Mock Agent and receives an automatic reply.

## Implementation Sequence

1. Scaffold repository structure for frontend, backend, and Docker Compose.
2. Implement MySQL migrations and seed data.
3. Implement backend auth, users, conversations, and messages.
4. Implement WebSocket connection manager and event broadcast.
5. Implement Mock Agent as an asynchronous module.
6. Implement React login and main chat UI.
7. Implement admin page.
8. Add tests around core flows.
9. Verify local development instructions from a clean checkout.

## Open Decisions

No product decisions remain open for the MVP. Specific library choices, such as Go router, migration tool, React Query version, and Markdown renderer, can be selected during implementation planning based on current ecosystem fit and project simplicity.
