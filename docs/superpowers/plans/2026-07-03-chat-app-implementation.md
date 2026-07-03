# Web Chat App Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build the approved React + Go + MySQL web chat MVP with seed-user login, direct/group conversations, Markdown messages, lightweight realtime updates, unread counts, a Mock Agent, and a development admin page.

**Architecture:** Use a modular monolith Go backend for REST, WebSocket, persistence, and Mock Agent behavior. Use a React SPA for login, chat, responsive layout, and the development admin page. MySQL is the source of truth; REST handles writes and WebSocket delivers server events.

**Tech Stack:** Go, `net/http`, `github.com/go-chi/chi/v5`, `github.com/go-sql-driver/mysql`, `github.com/gorilla/websocket`, `golang.org/x/crypto/bcrypt`, React, Vite, TypeScript, TanStack Query, React Testing Library, Playwright, MySQL, Docker Compose.

---

## File Structure

Create this repository structure:

```text
.
├── backend/
│   ├── cmd/
│   │   ├── api/main.go
│   │   └── migrate/main.go
│   ├── internal/
│   │   ├── admin/
│   │   │   ├── handler.go
│   │   │   └── handler_test.go
│   │   ├── agent/
│   │   │   ├── mock.go
│   │   │   └── mock_test.go
│   │   ├── auth/
│   │   │   ├── handler.go
│   │   │   ├── service.go
│   │   │   └── service_test.go
│   │   ├── config/config.go
│   │   ├── conversations/
│   │   │   ├── handler.go
│   │   │   ├── repository.go
│   │   │   ├── service.go
│   │   │   └── service_test.go
│   │   ├── db/
│   │   │   ├── db.go
│   │   │   ├── migrate.go
│   │   │   └── migrate_test.go
│   │   ├── httpx/
│   │   │   ├── auth.go
│   │   │   ├── json.go
│   │   │   └── errors.go
│   │   ├── messages/
│   │   │   ├── handler.go
│   │   │   ├── markdown.go
│   │   │   ├── repository.go
│   │   │   ├── service.go
│   │   │   └── service_test.go
│   │   ├── realtime/
│   │   │   ├── event.go
│   │   │   ├── hub.go
│   │   │   ├── hub_test.go
│   │   │   └── handler.go
│   │   └── users/
│   │       ├── handler.go
│   │       ├── repository.go
│   │       └── service.go
│   ├── migrations/
│   │   ├── 0001_schema.sql
│   │   └── 0002_seed_dev.sql
│   ├── go.mod
│   └── go.sum
├── frontend/
│   ├── index.html
│   ├── package.json
│   ├── playwright.config.ts
│   ├── src/
│   │   ├── app/App.tsx
│   │   ├── app/routes.tsx
│   │   ├── features/admin/AdminPage.tsx
│   │   ├── features/auth/LoginPage.tsx
│   │   ├── features/chat/ChatPage.tsx
│   │   ├── features/chat/ConversationList.tsx
│   │   ├── features/chat/MessageComposer.tsx
│   │   ├── features/chat/MessageList.tsx
│   │   ├── features/chat/MarkdownMessage.tsx
│   │   ├── lib/api.ts
│   │   ├── lib/query.ts
│   │   ├── lib/ws.ts
│   │   ├── main.tsx
│   │   └── styles.css
│   ├── tests/chat.e2e.ts
│   ├── tsconfig.json
│   └── vite.config.ts
├── docker-compose.yml
├── .env.example
└── README.md
```

Boundary rules:

- `backend/internal/messages` owns message validation, persistence orchestration, history pagination, read markers, and unread updates.
- `backend/internal/realtime` owns WebSocket connections and event broadcast only.
- `backend/internal/agent` owns Mock Agent reply policy and uses message service interfaces instead of writing SQL directly.
- `frontend/src/lib/api.ts` owns REST request shape and response parsing.
- `frontend/src/lib/ws.ts` owns reconnect and event dispatch.
- React components stay view-focused; data fetching lives in page components and hooks created inside feature directories.

---

## Task 1: Repository Scaffold and Developer Commands

**Files:**
- Create: `docker-compose.yml`
- Create: `.env.example`
- Create: `README.md`
- Create: `backend/go.mod`
- Create: `backend/cmd/api/main.go`
- Create: `backend/internal/config/config.go`
- Create: `frontend/package.json`
- Create: `frontend/index.html`
- Create: `frontend/src/main.tsx`
- Create: `frontend/src/app/App.tsx`
- Create: `frontend/src/styles.css`
- Create: `frontend/vite.config.ts`
- Create: `frontend/tsconfig.json`

- [ ] **Step 1: Create backend module skeleton**

Create `backend/go.mod`:

```go
module vibework-chat/backend

go 1.23
```

Create `backend/cmd/api/main.go`:

```go
package main

import (
	"fmt"
	"log"
	"net/http"

	"vibework-chat/backend/internal/config"
)

func main() {
	cfg := config.Load()
	addr := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("api listening on %s", addr)
	if err := http.ListenAndServe(addr, http.NewServeMux()); err != nil {
		log.Fatal(err)
	}
}
```

Create `backend/internal/config/config.go`:

```go
package config

import "os"

type Config struct {
	Port          string
	MySQLDSN      string
	SessionSecret string
	CORSOrigin    string
	DevSeed       bool
}

func Load() Config {
	return Config{
		Port:          env("PORT", "8080"),
		MySQLDSN:      env("MYSQL_DSN", "chat:chat@tcp(127.0.0.1:3306)/chat?parseTime=true&multiStatements=true"),
		SessionSecret: env("SESSION_SECRET", "dev-session-secret-change-me"),
		CORSOrigin:    env("CORS_ORIGIN", "http://localhost:5173"),
		DevSeed:       env("DEV_SEED", "true") == "true",
	}
}

func env(key string, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
```

- [ ] **Step 2: Add frontend skeleton**

Create `frontend/package.json`:

```json
{
  "scripts": {
    "dev": "vite",
    "build": "tsc -b && vite build",
    "test": "vitest run",
    "test:e2e": "playwright test",
    "lint": "tsc -b --noEmit"
  },
  "dependencies": {
    "@tanstack/react-query": "^5.0.0",
    "@vitejs/plugin-react": "^4.0.0",
    "react": "^18.2.0",
    "react-dom": "^18.2.0",
    "react-markdown": "^9.0.0",
    "remark-gfm": "^4.0.0",
    "vite": "^5.0.0"
  },
  "devDependencies": {
    "@playwright/test": "^1.45.0",
    "@testing-library/jest-dom": "^6.4.0",
    "@testing-library/react": "^15.0.0",
    "@types/react": "^18.2.0",
    "@types/react-dom": "^18.2.0",
    "typescript": "^5.4.0",
    "vitest": "^1.6.0"
  }
}
```

Create `frontend/index.html`:

```html
<!doctype html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <title>Vibework Chat</title>
  </head>
  <body>
    <div id="root"></div>
    <script type="module" src="/src/main.tsx"></script>
  </body>
</html>
```

Create `frontend/src/main.tsx`:

```tsx
import React from "react";
import ReactDOM from "react-dom/client";
import { App } from "./app/App";
import "./styles.css";

ReactDOM.createRoot(document.getElementById("root")!).render(
  <React.StrictMode>
    <App />
  </React.StrictMode>
);
```

Create `frontend/src/app/App.tsx`:

```tsx
export function App() {
  return <div className="app-shell">Vibework Chat</div>;
}
```

Create `frontend/src/styles.css`:

```css
:root {
  color: #1d2329;
  background: #f5f7f8;
  font-family: Inter, ui-sans-serif, system-ui, -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif;
}

* {
  box-sizing: border-box;
}

body {
  margin: 0;
}

.app-shell {
  min-height: 100vh;
  display: grid;
  place-items: center;
}
```

Create `frontend/vite.config.ts`:

```ts
import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";

export default defineConfig({
  plugins: [react()],
  server: {
    port: 5173,
    proxy: {
      "/api": "http://localhost:8080"
    }
  }
});
```

Create `frontend/tsconfig.json`:

```json
{
  "compilerOptions": {
    "target": "ES2020",
    "useDefineForClassFields": true,
    "lib": ["DOM", "DOM.Iterable", "ES2020"],
    "allowJs": false,
    "skipLibCheck": true,
    "esModuleInterop": true,
    "allowSyntheticDefaultImports": true,
    "strict": true,
    "forceConsistentCasingInFileNames": true,
    "module": "ESNext",
    "moduleResolution": "Node",
    "resolveJsonModule": true,
    "isolatedModules": true,
    "noEmit": true,
    "jsx": "react-jsx"
  },
  "include": ["src", "tests", "vite.config.ts", "playwright.config.ts"]
}
```

- [ ] **Step 3: Add Docker Compose and environment documentation**

Create `docker-compose.yml`:

```yaml
services:
  mysql:
    image: mysql:8.4
    ports:
      - "3306:3306"
    environment:
      MYSQL_ROOT_PASSWORD: root
      MYSQL_DATABASE: chat
      MYSQL_USER: chat
      MYSQL_PASSWORD: chat
    volumes:
      - mysql_data:/var/lib/mysql
    healthcheck:
      test: ["CMD", "mysqladmin", "ping", "-h", "localhost", "-uchat", "-pchat"]
      interval: 5s
      timeout: 3s
      retries: 20

volumes:
  mysql_data:
```

Create `.env.example`:

```dotenv
PORT=8080
MYSQL_DSN=chat:chat@tcp(127.0.0.1:3306)/chat?parseTime=true&multiStatements=true
SESSION_SECRET=dev-session-secret-change-me
CORS_ORIGIN=http://localhost:5173
DEV_SEED=true
```

Create `README.md`:

```markdown
# Vibework Chat

Web chat MVP with React, Go, and MySQL.

## Development

1. Start MySQL:

   ```bash
   docker compose up -d mysql
   ```

2. Run backend migrations after Task 2 is complete:

   ```bash
   cd backend
   go run ./cmd/migrate
   ```

3. Start the API:

   ```bash
   cd backend
   go run ./cmd/api
   ```

4. Start the frontend:

   ```bash
   cd frontend
   npm install
   npm run dev
   ```
```

- [ ] **Step 4: Verify scaffold builds enough to start**

Run:

```bash
cd backend && go mod tidy && go test ./...
```

Expected: PASS with no packages failing.

Run:

```bash
cd frontend && npm install && npm run lint
```

Expected: TypeScript passes.

- [ ] **Step 5: Commit scaffold**

```bash
git add docker-compose.yml .env.example README.md backend frontend
git commit -m "chore: scaffold chat app workspace"
```

---

## Task 2: MySQL Schema, Migration Runner, and Seed Data

**Files:**
- Create: `backend/internal/db/db.go`
- Create: `backend/internal/db/migrate.go`
- Create: `backend/internal/db/migrate_test.go`
- Create: `backend/cmd/migrate/main.go`
- Create: `backend/migrations/0001_schema.sql`
- Create: `backend/migrations/0002_seed_dev.sql`
- Modify: `backend/go.mod`

- [ ] **Step 1: Write migration ordering test**

Create `backend/internal/db/migrate_test.go`:

```go
package db

import "testing"

func TestMigrationFilesAreSorted(t *testing.T) {
	files := []string{"0002_seed_dev.sql", "0001_schema.sql"}
	got := SortMigrationFiles(files)
	want := []string{"0001_schema.sql", "0002_seed_dev.sql"}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("index %d: got %q want %q", i, got[i], want[i])
		}
	}
}
```

- [ ] **Step 2: Run migration test to verify it fails**

Run:

```bash
cd backend && go test ./internal/db -run TestMigrationFilesAreSorted -count=1
```

Expected: FAIL with `undefined: SortMigrationFiles`.

- [ ] **Step 3: Implement DB connection and migration helpers**

Create `backend/internal/db/db.go`:

```go
package db

import (
	"database/sql"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func Open(dsn string) (*sql.DB, error) {
	conn, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	conn.SetMaxOpenConns(10)
	conn.SetMaxIdleConns(10)
	conn.SetConnMaxLifetime(30 * time.Minute)
	if err := conn.Ping(); err != nil {
		_ = conn.Close()
		return nil, err
	}
	return conn, nil
}
```

Create `backend/internal/db/migrate.go`:

```go
package db

import (
	"database/sql"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func SortMigrationFiles(files []string) []string {
	out := append([]string(nil), files...)
	sort.Strings(out)
	return out
}

func RunMigrations(conn *sql.DB, dir string) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}
	var files []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".sql") {
			files = append(files, entry.Name())
		}
	}
	for _, name := range SortMigrationFiles(files) {
		body, err := os.ReadFile(filepath.Join(dir, name))
		if err != nil {
			return err
		}
		if _, err := conn.Exec(string(body)); err != nil {
			return err
		}
	}
	return nil
}
```

- [ ] **Step 4: Add schema migration**

Create `backend/migrations/0001_schema.sql`:

```sql
CREATE TABLE IF NOT EXISTS workspaces (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  name VARCHAR(120) NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS users (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  workspace_id BIGINT NOT NULL,
  username VARCHAR(80) NOT NULL,
  display_name VARCHAR(120) NOT NULL,
  avatar_url VARCHAR(500) NULL,
  user_type ENUM('human', 'agent') NOT NULL,
  password_hash VARCHAR(255) NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  UNIQUE KEY users_workspace_username_uq (workspace_id, username),
  CONSTRAINT users_workspace_fk FOREIGN KEY (workspace_id) REFERENCES workspaces(id)
);

CREATE TABLE IF NOT EXISTS sessions (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  user_id BIGINT NOT NULL,
  token_hash CHAR(64) NOT NULL,
  expires_at TIMESTAMP NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  UNIQUE KEY sessions_token_hash_uq (token_hash),
  KEY sessions_user_id_idx (user_id),
  CONSTRAINT sessions_user_fk FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE TABLE IF NOT EXISTS conversations (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  workspace_id BIGINT NOT NULL,
  type ENUM('direct', 'group') NOT NULL,
  title VARCHAR(160) NULL,
  created_by BIGINT NOT NULL,
  last_message_id BIGINT NULL,
  last_message_at TIMESTAMP NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  KEY conversations_workspace_updated_idx (workspace_id, last_message_at),
  CONSTRAINT conversations_workspace_fk FOREIGN KEY (workspace_id) REFERENCES workspaces(id),
  CONSTRAINT conversations_created_by_fk FOREIGN KEY (created_by) REFERENCES users(id)
);

CREATE TABLE IF NOT EXISTS conversation_members (
  conversation_id BIGINT NOT NULL,
  user_id BIGINT NOT NULL,
  role ENUM('owner', 'member') NOT NULL DEFAULT 'member',
  last_read_message_id BIGINT NULL,
  unread_count INT NOT NULL DEFAULT 0,
  joined_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (conversation_id, user_id),
  KEY conversation_members_user_idx (user_id),
  CONSTRAINT conversation_members_conversation_fk FOREIGN KEY (conversation_id) REFERENCES conversations(id),
  CONSTRAINT conversation_members_user_fk FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE TABLE IF NOT EXISTS messages (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  conversation_id BIGINT NOT NULL,
  sender_id BIGINT NOT NULL,
  content_markdown TEXT NOT NULL,
  content_plain TEXT NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  edited_at TIMESTAMP NULL,
  KEY messages_conversation_id_id_idx (conversation_id, id),
  CONSTRAINT messages_conversation_fk FOREIGN KEY (conversation_id) REFERENCES conversations(id),
  CONSTRAINT messages_sender_fk FOREIGN KEY (sender_id) REFERENCES users(id)
);

CREATE TABLE IF NOT EXISTS agent_profiles (
  user_id BIGINT PRIMARY KEY,
  kind ENUM('mock') NOT NULL,
  config_json JSON NOT NULL,
  enabled BOOLEAN NOT NULL DEFAULT TRUE,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  CONSTRAINT agent_profiles_user_fk FOREIGN KEY (user_id) REFERENCES users(id)
);
```

- [ ] **Step 5: Add dev seed migration**

Create `backend/migrations/0002_seed_dev.sql`:

```sql
INSERT INTO workspaces (id, name)
VALUES (1, 'Default')
ON DUPLICATE KEY UPDATE name = VALUES(name);

INSERT INTO users (id, workspace_id, username, display_name, avatar_url, user_type, password_hash)
VALUES
  (1, 1, 'alice', 'Alice Chen', NULL, 'human', NULL),
  (2, 1, 'bob', 'Bob Lin', NULL, 'human', NULL),
  (3, 1, 'carol', 'Carol Wu', NULL, 'human', NULL),
  (10, 1, 'mock-agent', 'Mock Agent', NULL, 'agent', NULL)
ON DUPLICATE KEY UPDATE
  display_name = VALUES(display_name),
  user_type = VALUES(user_type);

INSERT INTO agent_profiles (user_id, kind, config_json, enabled)
VALUES (10, 'mock', JSON_OBJECT('replyPrefix', 'Mock Agent received:'), TRUE)
ON DUPLICATE KEY UPDATE
  config_json = VALUES(config_json),
  enabled = VALUES(enabled);

INSERT INTO conversations (id, workspace_id, type, title, created_by)
VALUES
  (1, 1, 'direct', NULL, 1),
  (2, 1, 'group', 'Project Room', 1),
  (3, 1, 'direct', NULL, 1)
ON DUPLICATE KEY UPDATE title = VALUES(title);

INSERT INTO conversation_members (conversation_id, user_id, role)
VALUES
  (1, 1, 'owner'),
  (1, 2, 'member'),
  (2, 1, 'owner'),
  (2, 2, 'member'),
  (2, 3, 'member'),
  (3, 1, 'owner'),
  (3, 10, 'member')
ON DUPLICATE KEY UPDATE role = VALUES(role);
```

- [ ] **Step 6: Add migration command**

Create `backend/cmd/migrate/main.go`:

```go
package main

import (
	"log"
	"path/filepath"

	"vibework-chat/backend/internal/config"
	"vibework-chat/backend/internal/db"
)

func main() {
	cfg := config.Load()
	conn, err := db.Open(cfg.MySQLDSN)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	if err := db.RunMigrations(conn, filepath.Join("migrations")); err != nil {
		log.Fatal(err)
	}
	log.Println("migrations applied")
}
```

- [ ] **Step 7: Install backend dependencies and verify migration helper test**

Run:

```bash
cd backend && go get github.com/go-sql-driver/mysql && go mod tidy && go test ./internal/db -count=1
```

Expected: PASS.

- [ ] **Step 8: Verify migration against local MySQL**

Run:

```bash
docker compose up -d mysql
cd backend && go run ./cmd/migrate
```

Expected: command prints `migrations applied`.

- [ ] **Step 9: Commit database foundation**

```bash
git add backend docker-compose.yml .env.example
git commit -m "feat: add mysql schema and seed data"
```

---

## Task 3: Auth, Users, and Session Middleware

**Files:**
- Create: `backend/internal/httpx/json.go`
- Create: `backend/internal/httpx/errors.go`
- Create: `backend/internal/httpx/auth.go`
- Create: `backend/internal/users/repository.go`
- Create: `backend/internal/users/service.go`
- Create: `backend/internal/users/handler.go`
- Create: `backend/internal/auth/service.go`
- Create: `backend/internal/auth/service_test.go`
- Create: `backend/internal/auth/handler.go`
- Modify: `backend/cmd/api/main.go`
- Modify: `backend/go.mod`

- [ ] **Step 1: Write auth service tests**

Create `backend/internal/auth/service_test.go`:

```go
package auth

import (
	"context"
	"testing"
	"time"
)

type fakeUsers struct {
	users map[string]User
}

func (f fakeUsers) FindByUsername(ctx context.Context, username string) (User, error) {
	user, ok := f.users[username]
	if !ok {
		return User{}, ErrInvalidCredentials
	}
	return user, nil
}

type fakeSessions struct {
	createdFor int64
}

func (f *fakeSessions) Create(ctx context.Context, userID int64, tokenHash string, expiresAt time.Time) error {
	f.createdFor = userID
	return nil
}

func TestDevLoginCreatesSessionForSeedUser(t *testing.T) {
	sessionRepo := &fakeSessions{}
	service := NewService(fakeUsers{users: map[string]User{
		"alice": {ID: 1, Username: "alice", DisplayName: "Alice Chen", UserType: "human"},
	}}, sessionRepo, "secret")

	result, err := service.DevLogin(context.Background(), "alice")
	if err != nil {
		t.Fatal(err)
	}
	if result.User.ID != 1 {
		t.Fatalf("got user %d want 1", result.User.ID)
	}
	if result.Token == "" {
		t.Fatal("expected token")
	}
	if sessionRepo.createdFor != 1 {
		t.Fatalf("session user = %d want 1", sessionRepo.createdFor)
	}
}
```

- [ ] **Step 2: Run auth test to verify it fails**

Run:

```bash
cd backend && go test ./internal/auth -run TestDevLoginCreatesSessionForSeedUser -count=1
```

Expected: FAIL with undefined auth types and `NewService`.

- [ ] **Step 3: Implement shared HTTP helpers**

Create `backend/internal/httpx/json.go`:

```go
package httpx

import (
	"encoding/json"
	"net/http"
)

func ReadJSON(r *http.Request, dst any) error {
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(dst)
}

func WriteJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}
```

Create `backend/internal/httpx/errors.go`:

```go
package httpx

import "net/http"

type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details any    `json:"details"`
}

func WriteError(w http.ResponseWriter, status int, code string, message string) {
	WriteJSON(w, status, ErrorResponse{Code: code, Message: message, Details: map[string]any{}})
}
```

Create `backend/internal/httpx/auth.go`:

```go
package httpx

import "context"

type contextKey string

const userIDKey contextKey = "userID"

func WithUserID(ctx context.Context, userID int64) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}

func UserID(ctx context.Context) (int64, bool) {
	value, ok := ctx.Value(userIDKey).(int64)
	return value, ok
}
```

- [ ] **Step 4: Implement user repository and service**

Create `backend/internal/users/repository.go`:

```go
package users

import (
	"context"
	"database/sql"
)

type User struct {
	ID          int64   `json:"id"`
	WorkspaceID int64  `json:"workspaceId"`
	Username    string `json:"username"`
	DisplayName string `json:"displayName"`
	AvatarURL   *string `json:"avatarUrl"`
	UserType    string `json:"userType"`
}

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) FindByUsername(ctx context.Context, username string) (User, error) {
	var user User
	err := r.db.QueryRowContext(ctx, `
		SELECT id, workspace_id, username, display_name, avatar_url, user_type
		FROM users
		WHERE username = ?
	`, username).Scan(&user.ID, &user.WorkspaceID, &user.Username, &user.DisplayName, &user.AvatarURL, &user.UserType)
	return user, err
}

func (r *Repository) FindByID(ctx context.Context, id int64) (User, error) {
	var user User
	err := r.db.QueryRowContext(ctx, `
		SELECT id, workspace_id, username, display_name, avatar_url, user_type
		FROM users
		WHERE id = ?
	`, id).Scan(&user.ID, &user.WorkspaceID, &user.Username, &user.DisplayName, &user.AvatarURL, &user.UserType)
	return user, err
}

func (r *Repository) List(ctx context.Context, workspaceID int64) ([]User, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, workspace_id, username, display_name, avatar_url, user_type
		FROM users
		WHERE workspace_id = ?
		ORDER BY user_type, display_name
	`, workspaceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.WorkspaceID, &user.Username, &user.DisplayName, &user.AvatarURL, &user.UserType); err != nil {
			return nil, err
		}
		out = append(out, user)
	}
	return out, rows.Err()
}
```

Create `backend/internal/users/service.go`:

```go
package users

import "context"

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) List(ctx context.Context, workspaceID int64) ([]User, error) {
	return s.repo.List(ctx, workspaceID)
}

func (s *Service) FindByID(ctx context.Context, id int64) (User, error) {
	return s.repo.FindByID(ctx, id)
}
```

- [ ] **Step 5: Implement auth service**

Create `backend/internal/auth/service.go`:

```go
package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"time"
)

var ErrInvalidCredentials = errors.New("invalid credentials")

type User struct {
	ID          int64  `json:"id"`
	Username    string `json:"username"`
	DisplayName string `json:"displayName"`
	UserType    string `json:"userType"`
}

type UserFinder interface {
	FindByUsername(ctx context.Context, username string) (User, error)
}

type SessionStore interface {
	Create(ctx context.Context, userID int64, tokenHash string, expiresAt time.Time) error
}

type Service struct {
	users    UserFinder
	sessions SessionStore
	secret   string
}

type LoginResult struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

func NewService(users UserFinder, sessions SessionStore, secret string) *Service {
	return &Service{users: users, sessions: sessions, secret: secret}
}

func (s *Service) DevLogin(ctx context.Context, username string) (LoginResult, error) {
	user, err := s.users.FindByUsername(ctx, username)
	if err != nil {
		return LoginResult{}, ErrInvalidCredentials
	}
	token, hash, err := s.newToken()
	if err != nil {
		return LoginResult{}, err
	}
	expiresAt := time.Now().Add(24 * time.Hour)
	if err := s.sessions.Create(ctx, user.ID, hash, expiresAt); err != nil {
		return LoginResult{}, err
	}
	return LoginResult{Token: token, User: user}, nil
}

func (s *Service) newToken() (string, string, error) {
	raw := make([]byte, 32)
	if _, err := rand.Read(raw); err != nil {
		return "", "", err
	}
	token := hex.EncodeToString(raw)
	sum := sha256.Sum256([]byte(token + s.secret))
	return token, hex.EncodeToString(sum[:]), nil
}

func HashToken(token string, secret string) string {
	sum := sha256.Sum256([]byte(token + secret))
	return hex.EncodeToString(sum[:])
}
```

- [ ] **Step 6: Bridge user repository to auth service**

Add adapter methods in `backend/internal/auth/handler.go` and session repository in the same file for this task:

```go
package auth

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	"vibework-chat/backend/internal/httpx"
)

type SessionRepository struct {
	db *sql.DB
}

func NewSessionRepository(db *sql.DB) *SessionRepository {
	return &SessionRepository{db: db}
}

func (r *SessionRepository) Create(ctx context.Context, userID int64, tokenHash string, expiresAt time.Time) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO sessions (user_id, token_hash, expires_at)
		VALUES (?, ?, ?)
	`, userID, tokenHash, expiresAt)
	return err
}

func (r *SessionRepository) UserIDByToken(ctx context.Context, tokenHash string) (int64, error) {
	var userID int64
	err := r.db.QueryRowContext(ctx, `
		SELECT user_id
		FROM sessions
		WHERE token_hash = ? AND expires_at > NOW()
	`, tokenHash).Scan(&userID)
	return userID, err
}

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) DevLogin(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username"`
	}
	if err := httpx.ReadJSON(r, &req); err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "bad_json", "Invalid JSON body")
		return
	}
	result, err := h.service.DevLogin(r.Context(), req.Username)
	if err != nil {
		httpx.WriteError(w, http.StatusUnauthorized, "invalid_credentials", "Invalid credentials")
		return
	}
	httpx.WriteJSON(w, http.StatusOK, result)
}
```

- [ ] **Step 7: Add user handler**

Create `backend/internal/users/handler.go`:

```go
package users

import (
	"net/http"

	"vibework-chat/backend/internal/httpx"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	users, err := h.service.List(r.Context(), 1)
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, "users_failed", "Could not load users")
		return
	}
	httpx.WriteJSON(w, http.StatusOK, map[string]any{"users": users})
}
```

- [ ] **Step 8: Run auth test and tidy dependencies**

Run:

```bash
cd backend && go mod tidy && go test ./internal/auth -count=1
```

Expected: PASS.

- [ ] **Step 9: Wire auth and users into API router**

Modify `backend/cmd/api/main.go` to use chi and DB:

```go
package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"

	"vibework-chat/backend/internal/auth"
	"vibework-chat/backend/internal/config"
	"vibework-chat/backend/internal/db"
	"vibework-chat/backend/internal/users"
)

func main() {
	cfg := config.Load()
	conn, err := db.Open(cfg.MySQLDSN)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	userRepo := users.NewRepository(conn)
	userService := users.NewService(userRepo)
	sessionRepo := auth.NewSessionRepository(conn)
	authService := auth.NewService(authUserFinder{repo: userRepo}, sessionRepo, cfg.SessionSecret)

	r := chi.NewRouter()
	r.Post("/api/auth/dev-login", auth.NewHandler(authService).DevLogin)
	r.Get("/api/users", users.NewHandler(userService).List)

	addr := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("api listening on %s", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatal(err)
	}
}

type authUserFinder struct {
	repo *users.Repository
}

func (f authUserFinder) FindByUsername(ctx context.Context, username string) (auth.User, error) {
	user, err := f.repo.FindByUsername(ctx, username)
	if err != nil {
		return auth.User{}, err
	}
	return auth.User{ID: user.ID, Username: user.Username, DisplayName: user.DisplayName, UserType: user.UserType}, nil
}
```

Add `context` to the import list in `backend/cmd/api/main.go`.

- [ ] **Step 10: Verify API compiles**

Run:

```bash
cd backend && go get github.com/go-chi/chi/v5 && go mod tidy && go test ./...
```

Expected: PASS.

- [ ] **Step 11: Commit auth and users**

```bash
git add backend
git commit -m "feat: add seed auth and user APIs"
```

---

## Task 4: Conversations and Messages Core

**Files:**
- Create: `backend/internal/conversations/repository.go`
- Create: `backend/internal/conversations/service.go`
- Create: `backend/internal/conversations/service_test.go`
- Create: `backend/internal/conversations/handler.go`
- Create: `backend/internal/messages/markdown.go`
- Create: `backend/internal/messages/repository.go`
- Create: `backend/internal/messages/service.go`
- Create: `backend/internal/messages/service_test.go`
- Create: `backend/internal/messages/handler.go`
- Modify: `backend/cmd/api/main.go`

- [ ] **Step 1: Write message plain-text extraction test**

Create `backend/internal/messages/service_test.go`:

```go
package messages

import "testing"

func TestPlainTextFromMarkdown(t *testing.T) {
	input := "Hello **Bob**\n\n```go\nfmt.Println(\"x\")\n```"
	got := PlainTextFromMarkdown(input)
	want := "Hello Bob fmt.Println(\"x\")"
	if got != want {
		t.Fatalf("got %q want %q", got, want)
	}
}
```

- [ ] **Step 2: Run message test to verify it fails**

Run:

```bash
cd backend && go test ./internal/messages -run TestPlainTextFromMarkdown -count=1
```

Expected: FAIL with `undefined: PlainTextFromMarkdown`.

- [ ] **Step 3: Implement Markdown plain-text helper**

Create `backend/internal/messages/markdown.go`:

```go
package messages

import (
	"regexp"
	"strings"
)

var markdownMarks = regexp.MustCompile("[*_`#>\\[\\]()]+")
var whitespace = regexp.MustCompile("\\s+")

func PlainTextFromMarkdown(input string) string {
	stripped := markdownMarks.ReplaceAllString(input, "")
	return strings.TrimSpace(whitespace.ReplaceAllString(stripped, " "))
}
```

- [ ] **Step 4: Run message helper test**

Run:

```bash
cd backend && go test ./internal/messages -run TestPlainTextFromMarkdown -count=1
```

Expected: PASS.

- [ ] **Step 5: Write conversation service test**

Create `backend/internal/conversations/service_test.go`:

```go
package conversations

import (
	"context"
	"testing"
)

type fakeRepo struct {
	createdType string
	memberIDs   []int64
}

func (f *fakeRepo) Create(ctx context.Context, input CreateConversationInput) (Conversation, error) {
	f.createdType = input.Type
	f.memberIDs = input.MemberIDs
	return Conversation{ID: 42, Type: input.Type, Title: input.Title}, nil
}

func TestCreateGroupRequiresAtLeastThreeMembers(t *testing.T) {
	service := NewService(&fakeRepo{})
	_, err := service.Create(context.Background(), CreateConversationInput{
		WorkspaceID: 1,
		CreatedBy:   1,
		Type:        "group",
		Title:       "Room",
		MemberIDs:   []int64{1, 2},
	})
	if err != ErrGroupNeedsThreeMembers {
		t.Fatalf("got %v want %v", err, ErrGroupNeedsThreeMembers)
	}
}
```

- [ ] **Step 6: Run conversation test to verify it fails**

Run:

```bash
cd backend && go test ./internal/conversations -run TestCreateGroupRequiresAtLeastThreeMembers -count=1
```

Expected: FAIL with undefined conversation service types.

- [ ] **Step 7: Implement conversation service and repository interfaces**

Create `backend/internal/conversations/service.go`:

```go
package conversations

import (
	"context"
	"errors"
)

var ErrGroupNeedsThreeMembers = errors.New("group conversations need at least three members")

type Conversation struct {
	ID            int64   `json:"id"`
	WorkspaceID   int64   `json:"workspaceId"`
	Type          string  `json:"type"`
	Title         *string `json:"title"`
	LastMessageID *int64  `json:"lastMessageId"`
	LastMessageAt *string `json:"lastMessageAt"`
	UnreadCount   int     `json:"unreadCount"`
}

type CreateConversationInput struct {
	WorkspaceID int64
	CreatedBy   int64
	Type        string
	Title       *string
	MemberIDs   []int64
}

type Repository interface {
	Create(ctx context.Context, input CreateConversationInput) (Conversation, error)
}

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, input CreateConversationInput) (Conversation, error) {
	if input.Type == "group" && len(input.MemberIDs) < 3 {
		return Conversation{}, ErrGroupNeedsThreeMembers
	}
	return s.repo.Create(ctx, input)
}
```

Create `backend/internal/conversations/repository.go`:

```go
package conversations

import (
	"context"
	"database/sql"
)

type SQLRepository struct {
	db *sql.DB
}

func NewSQLRepository(db *sql.DB) *SQLRepository {
	return &SQLRepository{db: db}
}

func (r *SQLRepository) Create(ctx context.Context, input CreateConversationInput) (Conversation, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return Conversation{}, err
	}
	defer tx.Rollback()

	result, err := tx.ExecContext(ctx, `
		INSERT INTO conversations (workspace_id, type, title, created_by)
		VALUES (?, ?, ?, ?)
	`, input.WorkspaceID, input.Type, input.Title, input.CreatedBy)
	if err != nil {
		return Conversation{}, err
	}
	conversationID, err := result.LastInsertId()
	if err != nil {
		return Conversation{}, err
	}
	for _, memberID := range input.MemberIDs {
		role := "member"
		if memberID == input.CreatedBy {
			role = "owner"
		}
		if _, err := tx.ExecContext(ctx, `
			INSERT INTO conversation_members (conversation_id, user_id, role)
			VALUES (?, ?, ?)
		`, conversationID, memberID, role); err != nil {
			return Conversation{}, err
		}
	}
	if err := tx.Commit(); err != nil {
		return Conversation{}, err
	}
	return Conversation{ID: conversationID, WorkspaceID: input.WorkspaceID, Type: input.Type, Title: input.Title}, nil
}

func (r *SQLRepository) ListForUser(ctx context.Context, workspaceID int64, userID int64) ([]Conversation, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT
		  c.id,
		  c.workspace_id,
		  c.type,
		  COALESCE(c.title, GROUP_CONCAT(CASE WHEN u.id <> ? THEN u.display_name END SEPARATOR ', ')) AS title,
		  c.last_message_id,
		  CAST(c.last_message_at AS CHAR),
		  cm.unread_count
		FROM conversations c
		JOIN conversation_members cm ON cm.conversation_id = c.id AND cm.user_id = ?
		JOIN conversation_members all_cm ON all_cm.conversation_id = c.id
		JOIN users u ON u.id = all_cm.user_id
		WHERE c.workspace_id = ?
		GROUP BY c.id, c.workspace_id, c.type, c.title, c.last_message_id, c.last_message_at, cm.unread_count
		ORDER BY COALESCE(c.last_message_at, c.created_at) DESC
	`, userID, userID, workspaceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []Conversation
	for rows.Next() {
		var conversation Conversation
		var title sql.NullString
		var lastMessageID sql.NullInt64
		var lastMessageAt sql.NullString
		if err := rows.Scan(
			&conversation.ID,
			&conversation.WorkspaceID,
			&conversation.Type,
			&title,
			&lastMessageID,
			&lastMessageAt,
			&conversation.UnreadCount,
		); err != nil {
			return nil, err
		}
		if title.Valid {
			conversation.Title = &title.String
		}
		if lastMessageID.Valid {
			conversation.LastMessageID = &lastMessageID.Int64
		}
		if lastMessageAt.Valid {
			conversation.LastMessageAt = &lastMessageAt.String
		}
		out = append(out, conversation)
	}
	return out, rows.Err()
}
```

- [ ] **Step 8: Implement message repository and service**

Create `backend/internal/messages/service.go`:

```go
package messages

import (
	"context"
	"errors"
	"strings"
)

var ErrEmptyMessage = errors.New("message content is empty")

type Message struct {
	ID              int64  `json:"id"`
	ConversationID  int64  `json:"conversationId"`
	SenderID        int64  `json:"senderId"`
	ContentMarkdown string `json:"contentMarkdown"`
	ContentPlain    string `json:"contentPlain"`
	CreatedAt       string `json:"createdAt"`
}

type SendInput struct {
	ConversationID  int64
	SenderID        int64
	ContentMarkdown string
}

type Repository interface {
	Create(ctx context.Context, input SendInput, contentPlain string) (Message, error)
	ListBefore(ctx context.Context, conversationID int64, beforeID int64, limit int) ([]Message, error)
	MarkRead(ctx context.Context, conversationID int64, userID int64) error
}

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Send(ctx context.Context, input SendInput) (Message, error) {
	if strings.TrimSpace(input.ContentMarkdown) == "" {
		return Message{}, ErrEmptyMessage
	}
	return s.repo.Create(ctx, input, PlainTextFromMarkdown(input.ContentMarkdown))
}

func (s *Service) ListBefore(ctx context.Context, conversationID int64, beforeID int64, limit int) ([]Message, error) {
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	return s.repo.ListBefore(ctx, conversationID, beforeID, limit)
}

func (s *Service) MarkRead(ctx context.Context, conversationID int64, userID int64) error {
	return s.repo.MarkRead(ctx, conversationID, userID)
}
```

Create `backend/internal/messages/repository.go`:

```go
package messages

import (
	"context"
	"database/sql"
)

type SQLRepository struct {
	db *sql.DB
}

func NewSQLRepository(db *sql.DB) *SQLRepository {
	return &SQLRepository{db: db}
}

func (r *SQLRepository) Create(ctx context.Context, input SendInput, contentPlain string) (Message, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return Message{}, err
	}
	defer tx.Rollback()

	result, err := tx.ExecContext(ctx, `
		INSERT INTO messages (conversation_id, sender_id, content_markdown, content_plain)
		VALUES (?, ?, ?, ?)
	`, input.ConversationID, input.SenderID, input.ContentMarkdown, contentPlain)
	if err != nil {
		return Message{}, err
	}
	messageID, err := result.LastInsertId()
	if err != nil {
		return Message{}, err
	}
	if _, err := tx.ExecContext(ctx, `
		UPDATE conversations
		SET last_message_id = ?, last_message_at = NOW()
		WHERE id = ?
	`, messageID, input.ConversationID); err != nil {
		return Message{}, err
	}
	if _, err := tx.ExecContext(ctx, `
		UPDATE conversation_members
		SET unread_count = unread_count + 1
		WHERE conversation_id = ? AND user_id <> ?
	`, input.ConversationID, input.SenderID); err != nil {
		return Message{}, err
	}
	if err := tx.Commit(); err != nil {
		return Message{}, err
	}
	return Message{
		ID:              messageID,
		ConversationID:  input.ConversationID,
		SenderID:        input.SenderID,
		ContentMarkdown: input.ContentMarkdown,
		ContentPlain:    contentPlain,
	}, nil
}

func (r *SQLRepository) ListBefore(ctx context.Context, conversationID int64, beforeID int64, limit int) ([]Message, error) {
	query := `
		SELECT id, conversation_id, sender_id, content_markdown, content_plain, created_at
		FROM messages
		WHERE conversation_id = ? AND (? = 0 OR id < ?)
		ORDER BY id DESC
		LIMIT ?
	`
	rows, err := r.db.QueryContext(ctx, query, conversationID, beforeID, beforeID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []Message
	for rows.Next() {
		var message Message
		if err := rows.Scan(&message.ID, &message.ConversationID, &message.SenderID, &message.ContentMarkdown, &message.ContentPlain, &message.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, message)
	}
	return out, rows.Err()
}

func (r *SQLRepository) MarkRead(ctx context.Context, conversationID int64, userID int64) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE conversation_members
		SET unread_count = 0,
		    last_read_message_id = (SELECT last_message_id FROM conversations WHERE id = ?)
		WHERE conversation_id = ? AND user_id = ?
	`, conversationID, conversationID, userID)
	return err
}
```

- [ ] **Step 9: Run service tests**

Run:

```bash
cd backend && go test ./internal/messages ./internal/conversations -count=1
```

Expected: PASS.

- [ ] **Step 10: Add conversation and message handlers**

Create `backend/internal/conversations/handler.go`:

```go
package conversations

import (
	"net/http"

	"vibework-chat/backend/internal/httpx"
)

type Handler struct {
	service *Service
	lister  interface {
		ListForUser(ctx context.Context, workspaceID int64, userID int64) ([]Conversation, error)
	}
}

func NewHandler(service *Service, lister interface {
	ListForUser(ctx context.Context, workspaceID int64, userID int64) ([]Conversation, error)
}) *Handler {
	return &Handler{service: service, lister: lister}
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	userID, ok := httpx.UserID(r.Context())
	if !ok {
		httpx.WriteError(w, http.StatusUnauthorized, "unauthorized", "Login required")
		return
	}
	conversations, err := h.lister.ListForUser(r.Context(), 1, userID)
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, "conversations_failed", "Could not load conversations")
		return
	}
	httpx.WriteJSON(w, http.StatusOK, map[string]any{"conversations": conversations})
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := httpx.UserID(r.Context())
	if !ok {
		httpx.WriteError(w, http.StatusUnauthorized, "unauthorized", "Login required")
		return
	}
	var req struct {
		Type      string  `json:"type"`
		Title     *string `json:"title"`
		MemberIDs []int64 `json:"memberIds"`
	}
	if err := httpx.ReadJSON(r, &req); err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "bad_json", "Invalid JSON body")
		return
	}
	input := CreateConversationInput{
		WorkspaceID: 1,
		CreatedBy:   userID,
		Type:        req.Type,
		Title:       req.Title,
		MemberIDs:   append(req.MemberIDs, userID),
	}
	conversation, err := h.service.Create(r.Context(), input)
	if err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "conversation_invalid", err.Error())
		return
	}
	httpx.WriteJSON(w, http.StatusCreated, conversation)
}
```

Add `context` to the import list in `backend/internal/conversations/handler.go`.

Create `backend/internal/messages/handler.go`:

```go
package messages

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"vibework-chat/backend/internal/httpx"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) ListMessages(w http.ResponseWriter, r *http.Request) {
	conversationID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "bad_conversation_id", "Invalid conversation id")
		return
	}
	beforeID, _ := strconv.ParseInt(r.URL.Query().Get("before"), 10, 64)
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	messages, err := h.service.ListBefore(r.Context(), conversationID, beforeID, limit)
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, "messages_failed", "Could not load messages")
		return
	}
	httpx.WriteJSON(w, http.StatusOK, map[string]any{"messages": messages})
}

func (h *Handler) SendMessage(w http.ResponseWriter, r *http.Request) {
	userID, ok := httpx.UserID(r.Context())
	if !ok {
		httpx.WriteError(w, http.StatusUnauthorized, "unauthorized", "Login required")
		return
	}
	conversationID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "bad_conversation_id", "Invalid conversation id")
		return
	}
	var req struct {
		ContentMarkdown string `json:"contentMarkdown"`
	}
	if err := httpx.ReadJSON(r, &req); err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "bad_json", "Invalid JSON body")
		return
	}
	message, err := h.service.Send(r.Context(), SendInput{
		ConversationID:  conversationID,
		SenderID:        userID,
		ContentMarkdown: req.ContentMarkdown,
	})
	if err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "message_invalid", err.Error())
		return
	}
	httpx.WriteJSON(w, http.StatusCreated, message)
}

func (h *Handler) MarkRead(w http.ResponseWriter, r *http.Request) {
	userID, ok := httpx.UserID(r.Context())
	if !ok {
		httpx.WriteError(w, http.StatusUnauthorized, "unauthorized", "Login required")
		return
	}
	conversationID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "bad_conversation_id", "Invalid conversation id")
		return
	}
	if err := h.service.MarkRead(r.Context(), conversationID, userID); err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, "read_failed", "Could not mark conversation read")
		return
	}
	httpx.WriteJSON(w, http.StatusOK, map[string]any{"ok": true})
}
```

- [ ] **Step 11: Wire routes into API**

Modify `backend/cmd/api/main.go` to add:

```go
r.Get("/api/conversations", conversationHandler.List)
r.Post("/api/conversations", conversationHandler.Create)
r.Get("/api/conversations/{id}/messages", messageHandler.ListMessages)
r.Post("/api/conversations/{id}/messages", messageHandler.SendMessage)
r.Post("/api/conversations/{id}/read", messageHandler.MarkRead)
```

Keep auth middleware minimal for this task by reading `Authorization: Bearer <token>` and attaching `userID` through `httpx.WithUserID`.

- [ ] **Step 12: Verify backend packages**

Run:

```bash
cd backend && go test ./...
```

Expected: PASS.

- [ ] **Step 13: Commit conversations and messages**

```bash
git add backend
git commit -m "feat: add conversations and messages"
```

---

## Task 5: Realtime WebSocket Events

**Files:**
- Create: `backend/internal/realtime/event.go`
- Create: `backend/internal/realtime/hub.go`
- Create: `backend/internal/realtime/hub_test.go`
- Create: `backend/internal/realtime/handler.go`
- Create: `backend/internal/realtime/notifier.go`
- Modify: `backend/internal/messages/service.go`
- Modify: `backend/cmd/api/main.go`
- Modify: `backend/go.mod`

- [ ] **Step 1: Write hub broadcast test**

Create `backend/internal/realtime/hub_test.go`:

```go
package realtime

import "testing"

func TestHubStoresConnectionsByUser(t *testing.T) {
	hub := NewHub()
	ch := make(chan Event, 1)
	hub.Register(1, ch)
	hub.SendToUser(1, Event{Type: "message.created", Payload: map[string]any{"id": 5}})
	got := <-ch
	if got.Type != "message.created" {
		t.Fatalf("got %q", got.Type)
	}
	hub.Unregister(1, ch)
	if len(hub.users[1]) != 0 {
		t.Fatal("expected user connections removed")
	}
}
```

- [ ] **Step 2: Run hub test to verify it fails**

Run:

```bash
cd backend && go test ./internal/realtime -run TestHubStoresConnectionsByUser -count=1
```

Expected: FAIL with undefined realtime types.

- [ ] **Step 3: Implement realtime event and hub**

Create `backend/internal/realtime/event.go`:

```go
package realtime

type Event struct {
	Type    string `json:"type"`
	Payload any   `json:"payload"`
}
```

Create `backend/internal/realtime/hub.go`:

```go
package realtime

import "sync"

type Hub struct {
	mu    sync.RWMutex
	users map[int64]map[chan Event]struct{}
}

func NewHub() *Hub {
	return &Hub{users: map[int64]map[chan Event]struct{}{}}
}

func (h *Hub) Register(userID int64, ch chan Event) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.users[userID] == nil {
		h.users[userID] = map[chan Event]struct{}{}
	}
	h.users[userID][ch] = struct{}{}
}

func (h *Hub) Unregister(userID int64, ch chan Event) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.users[userID] == nil {
		return
	}
	delete(h.users[userID], ch)
	close(ch)
}

func (h *Hub) SendToUser(userID int64, event Event) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	for ch := range h.users[userID] {
		select {
		case ch <- event:
		default:
		}
	}
}
```

- [ ] **Step 4: Run hub test**

Run:

```bash
cd backend && go test ./internal/realtime -count=1
```

Expected: PASS.

- [ ] **Step 5: Add WebSocket handler**

Create `backend/internal/realtime/handler.go`:

```go
package realtime

import (
	"net/http"

	"github.com/gorilla/websocket"

	"vibework-chat/backend/internal/httpx"
)

type Handler struct {
	hub      *Hub
	upgrader websocket.Upgrader
}

func NewHandler(hub *Hub, allowedOrigin string) *Handler {
	return &Handler{
		hub: hub,
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return allowedOrigin == "" || r.Header.Get("Origin") == allowedOrigin
			},
		},
	}
}

func (h *Handler) ServeWS(w http.ResponseWriter, r *http.Request) {
	userID, ok := httpx.UserID(r.Context())
	if !ok {
		httpx.WriteError(w, http.StatusUnauthorized, "unauthorized", "Login required")
		return
	}
	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	out := make(chan Event, 16)
	h.hub.Register(userID, out)
	defer h.hub.Unregister(userID, out)

	for event := range out {
		if err := conn.WriteJSON(event); err != nil {
			return
		}
	}
}
```

- [ ] **Step 6: Broadcast from message creation**

Modify `backend/internal/messages/service.go` so `Service` accepts a notifier:

```go
type Notifier interface {
	MessageCreated(ctx context.Context, message Message) error
	ConversationUpdated(ctx context.Context, conversationID int64) error
}
```

After `repo.Create` succeeds in `Send`, call:

```go
if s.notifier != nil {
	_ = s.notifier.MessageCreated(ctx, message)
	_ = s.notifier.ConversationUpdated(ctx, message.ConversationID)
}
```

Create `backend/internal/realtime/notifier.go`:

```go
package realtime

import (
	"context"
	"database/sql"

	"vibework-chat/backend/internal/messages"
)

type SQLNotifier struct {
	db  *sql.DB
	hub *Hub
}

func NewSQLNotifier(db *sql.DB, hub *Hub) *SQLNotifier {
	return &SQLNotifier{db: db, hub: hub}
}

func (n *SQLNotifier) MessageCreated(ctx context.Context, message messages.Message) error {
	userIDs, err := n.memberIDs(ctx, message.ConversationID)
	if err != nil {
		return err
	}
	for _, userID := range userIDs {
		n.hub.SendToUser(userID, Event{Type: "message.created", Payload: message})
	}
	return nil
}

func (n *SQLNotifier) ConversationUpdated(ctx context.Context, conversationID int64) error {
	userIDs, err := n.memberIDs(ctx, conversationID)
	if err != nil {
		return err
	}
	for _, userID := range userIDs {
		n.hub.SendToUser(userID, Event{
			Type:    "conversation.updated",
			Payload: map[string]any{"conversationId": conversationID},
		})
	}
	return nil
}

func (n *SQLNotifier) memberIDs(ctx context.Context, conversationID int64) ([]int64, error) {
	rows, err := n.db.QueryContext(ctx, `
		SELECT user_id
		FROM conversation_members
		WHERE conversation_id = ?
	`, conversationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []int64
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, rows.Err()
}
```

- [ ] **Step 7: Wire `/api/ws`**

Modify `backend/cmd/api/main.go`:

```go
hub := realtime.NewHub()
realtimeHandler := realtime.NewHandler(hub, cfg.CORSOrigin)
r.Get("/api/ws", realtimeHandler.ServeWS)
```

Ensure the auth middleware also wraps `/api/ws`.

- [ ] **Step 8: Install WebSocket dependency and test**

Run:

```bash
cd backend && go get github.com/gorilla/websocket && go mod tidy && go test ./...
```

Expected: PASS.

- [ ] **Step 9: Commit realtime**

```bash
git add backend
git commit -m "feat: add websocket realtime events"
```

---

## Task 6: Mock Agent Reply Flow

**Files:**
- Create: `backend/internal/agent/mock.go`
- Create: `backend/internal/agent/mock_test.go`
- Modify: `backend/internal/messages/service.go`
- Modify: `backend/cmd/api/main.go`

- [ ] **Step 1: Write Mock Agent reply test**

Create `backend/internal/agent/mock_test.go`:

```go
package agent

import (
	"context"
	"testing"
)

func TestMockRunnerBuildsReply(t *testing.T) {
	runner := NewMockRunner("Mock Agent received:")
	reply := runner.BuildReply(context.Background(), "hello")
	want := "Mock Agent received: hello"
	if reply != want {
		t.Fatalf("got %q want %q", reply, want)
	}
}
```

- [ ] **Step 2: Run Mock Agent test to verify it fails**

Run:

```bash
cd backend && go test ./internal/agent -run TestMockRunnerBuildsReply -count=1
```

Expected: FAIL with undefined `NewMockRunner`.

- [ ] **Step 3: Implement Mock Agent runner**

Create `backend/internal/agent/mock.go`:

```go
package agent

import (
	"context"
	"strings"
)

type MockRunner struct {
	prefix string
}

func NewMockRunner(prefix string) *MockRunner {
	return &MockRunner{prefix: prefix}
}

func (r *MockRunner) BuildReply(ctx context.Context, input string) string {
	trimmed := strings.TrimSpace(input)
	if trimmed == "" {
		return r.prefix + " empty message"
	}
	return r.prefix + " " + trimmed
}
```

- [ ] **Step 4: Add agent orchestration interface**

In `backend/internal/agent/mock.go`, add:

```go
type MessageSender interface {
	SendFromAgent(ctx context.Context, conversationID int64, agentUserID int64, contentMarkdown string) error
}

type ConversationAgentFinder interface {
	EnabledMockAgentForConversation(ctx context.Context, conversationID int64) (agentUserID int64, enabled bool, err error)
}

type Orchestrator struct {
	finder ConversationAgentFinder
	sender MessageSender
	runner *MockRunner
}

func NewOrchestrator(finder ConversationAgentFinder, sender MessageSender, runner *MockRunner) *Orchestrator {
	return &Orchestrator{finder: finder, sender: sender, runner: runner}
}

func (o *Orchestrator) MaybeReply(ctx context.Context, conversationID int64, humanMessage string) {
	go func() {
		agentUserID, enabled, err := o.finder.EnabledMockAgentForConversation(context.Background(), conversationID)
		if err != nil || !enabled {
			return
		}
		reply := o.runner.BuildReply(context.Background(), humanMessage)
		_ = o.sender.SendFromAgent(context.Background(), conversationID, agentUserID, reply)
	}()
}
```

- [ ] **Step 5: Trigger Mock Agent after human messages**

Modify `backend/internal/messages/service.go`:

```go
type AgentResponder interface {
	MaybeReply(ctx context.Context, conversationID int64, humanMessage string)
}
```

In `Send`, after broadcasting the human message:

```go
if s.agentResponder != nil {
	s.agentResponder.MaybeReply(ctx, message.ConversationID, message.ContentMarkdown)
}
```

Add `SendFromAgent` to the message service so the orchestrator can create the reply without triggering another agent reply:

```go
func (s *Service) SendFromAgent(ctx context.Context, conversationID int64, agentUserID int64, contentMarkdown string) error {
	_, err := s.repo.Create(ctx, SendInput{
		ConversationID: conversationID,
		SenderID:       agentUserID,
		ContentMarkdown: contentMarkdown,
	}, PlainTextFromMarkdown(contentMarkdown))
	return err
}
```

- [ ] **Step 6: Implement enabled agent finder**

In the SQL-backed conversation or agent repository, implement:

```sql
SELECT u.id
FROM conversation_members cm
JOIN users u ON u.id = cm.user_id
JOIN agent_profiles ap ON ap.user_id = u.id
WHERE cm.conversation_id = ?
  AND u.user_type = 'agent'
  AND ap.kind = 'mock'
  AND ap.enabled = TRUE
LIMIT 1
```

Return `(0, false, nil)` when no row exists.

- [ ] **Step 7: Verify agent tests and backend**

Run:

```bash
cd backend && go test ./...
```

Expected: PASS.

- [ ] **Step 8: Commit Mock Agent**

```bash
git add backend
git commit -m "feat: add mock agent replies"
```

---

## Task 7: Development Admin API

**Files:**
- Create: `backend/internal/admin/handler.go`
- Create: `backend/internal/admin/handler_test.go`
- Modify: `backend/cmd/api/main.go`

- [ ] **Step 1: Write admin response shape test**

Create `backend/internal/admin/handler_test.go`:

```go
package admin

import "testing"

func TestLimitRecentMessages(t *testing.T) {
	if got := normalizeLimit(0); got != 50 {
		t.Fatalf("got %d want 50", got)
	}
	if got := normalizeLimit(500); got != 100 {
		t.Fatalf("got %d want 100", got)
	}
	if got := normalizeLimit(20); got != 20 {
		t.Fatalf("got %d want 20", got)
	}
}
```

- [ ] **Step 2: Run admin test to verify it fails**

Run:

```bash
cd backend && go test ./internal/admin -run TestLimitRecentMessages -count=1
```

Expected: FAIL with `undefined: normalizeLimit`.

- [ ] **Step 3: Implement admin handler**

Create `backend/internal/admin/handler.go`:

```go
package admin

import (
	"database/sql"
	"net/http"
	"strconv"

	"vibework-chat/backend/internal/httpx"
)

type Handler struct {
	db *sql.DB
}

func NewHandler(db *sql.DB) *Handler {
	return &Handler{db: db}
}

func normalizeLimit(limit int) int {
	if limit <= 0 {
		return 50
	}
	if limit > 100 {
		return 100
	}
	return limit
}

func (h *Handler) Users(w http.ResponseWriter, r *http.Request) {
	rows, err := h.db.QueryContext(r.Context(), `
		SELECT id, username, display_name, user_type, created_at
		FROM users
		ORDER BY id
	`)
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, "admin_users_failed", "Could not load users")
		return
	}
	defer rows.Close()
	httpx.WriteJSON(w, http.StatusOK, map[string]any{"users": rowsToMaps(rows)})
}

func (h *Handler) Conversations(w http.ResponseWriter, r *http.Request) {
	rows, err := h.db.QueryContext(r.Context(), `
		SELECT id, type, title, last_message_at, created_at
		FROM conversations
		ORDER BY id DESC
		LIMIT 100
	`)
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, "admin_conversations_failed", "Could not load conversations")
		return
	}
	defer rows.Close()
	httpx.WriteJSON(w, http.StatusOK, map[string]any{"conversations": rowsToMaps(rows)})
}

func (h *Handler) Messages(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	rows, err := h.db.QueryContext(r.Context(), `
		SELECT id, conversation_id, sender_id, content_plain, created_at
		FROM messages
		ORDER BY id DESC
		LIMIT ?
	`, normalizeLimit(limit))
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, "admin_messages_failed", "Could not load messages")
		return
	}
	defer rows.Close()
	httpx.WriteJSON(w, http.StatusOK, map[string]any{"messages": rowsToMaps(rows)})
}
```

Append this helper to `backend/internal/admin/handler.go`:

```go
func rowsToMaps(rows *sql.Rows) []map[string]any {
	columns, err := rows.Columns()
	if err != nil {
		return []map[string]any{}
	}
	var result []map[string]any
	for rows.Next() {
		values := make([]any, len(columns))
		pointers := make([]any, len(columns))
		for i := range values {
			pointers[i] = &values[i]
		}
		if err := rows.Scan(pointers...); err != nil {
			return result
		}
		row := map[string]any{}
		for i, column := range columns {
			switch value := values[i].(type) {
			case []byte:
				row[column] = string(value)
			default:
				row[column] = value
			}
		}
		result = append(result, row)
	}
	return result
}
```

- [ ] **Step 4: Wire admin routes**

Modify `backend/cmd/api/main.go`:

```go
adminHandler := admin.NewHandler(conn)
r.Get("/api/admin/users", adminHandler.Users)
r.Get("/api/admin/conversations", adminHandler.Conversations)
r.Get("/api/admin/messages", adminHandler.Messages)
```

- [ ] **Step 5: Verify backend**

Run:

```bash
cd backend && go test ./...
```

Expected: PASS.

- [ ] **Step 6: Commit admin API**

```bash
git add backend
git commit -m "feat: add development admin api"
```

---

## Task 8: Frontend API Client, Query Setup, and Login

**Files:**
- Create: `frontend/src/lib/api.ts`
- Create: `frontend/src/lib/query.ts`
- Create: `frontend/src/app/routes.tsx`
- Create: `frontend/src/features/auth/LoginPage.tsx`
- Modify: `frontend/src/app/App.tsx`
- Modify: `frontend/src/styles.css`
- Modify: `frontend/package.json`

- [ ] **Step 1: Write API client unit test**

Create `frontend/src/lib/api.test.ts`:

```ts
import { describe, expect, it } from "vitest";
import { authHeader } from "./api";

describe("authHeader", () => {
  it("returns bearer header when token exists", () => {
    expect(authHeader("abc")).toEqual({ Authorization: "Bearer abc" });
  });

  it("returns empty object without token", () => {
    expect(authHeader(null)).toEqual({});
  });
});
```

- [ ] **Step 2: Run frontend test to verify it fails**

Run:

```bash
cd frontend && npm run test -- --run src/lib/api.test.ts
```

Expected: FAIL because `src/lib/api.ts` does not exist.

- [ ] **Step 3: Implement API client**

Create `frontend/src/lib/api.ts`:

```ts
export type User = {
  id: number;
  username: string;
  displayName: string;
  userType: "human" | "agent";
};

export type LoginResult = {
  token: string;
  user: User;
};

export function authHeader(token: string | null): Record<string, string> {
  return token ? { Authorization: `Bearer ${token}` } : {};
}

export async function apiGet<T>(path: string, token: string | null): Promise<T> {
  const response = await fetch(path, { headers: authHeader(token) });
  return parseResponse<T>(response);
}

export async function apiPost<T>(path: string, token: string | null, body: unknown): Promise<T> {
  const response = await fetch(path, {
    method: "POST",
    headers: { "Content-Type": "application/json", ...authHeader(token) },
    body: JSON.stringify(body)
  });
  return parseResponse<T>(response);
}

async function parseResponse<T>(response: Response): Promise<T> {
  const data = await response.json();
  if (!response.ok) {
    throw new Error(data.message ?? "Request failed");
  }
  return data as T;
}

export function devLogin(username: string): Promise<LoginResult> {
  return apiPost<LoginResult>("/api/auth/dev-login", null, { username });
}
```

- [ ] **Step 4: Add React Query provider**

Create `frontend/src/lib/query.ts`:

```ts
import { QueryClient } from "@tanstack/react-query";

export const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      refetchOnWindowFocus: false,
      retry: 1
    }
  }
});
```

- [ ] **Step 5: Implement login page**

Create `frontend/src/features/auth/LoginPage.tsx`:

```tsx
import { useState } from "react";
import { devLogin, type User } from "../../lib/api";

const seedUsers = ["alice", "bob", "carol"];

type Props = {
  onLogin: (token: string, user: User) => void;
};

export function LoginPage({ onLogin }: Props) {
  const [username, setUsername] = useState("alice");
  const [error, setError] = useState<string | null>(null);

  async function loginAs(name: string) {
    setError(null);
    try {
      const result = await devLogin(name);
      onLogin(result.token, result.user);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Login failed");
    }
  }

  return (
    <main className="login-page">
      <section className="login-panel">
        <h1>Vibework Chat</h1>
        <div className="seed-list">
          {seedUsers.map((name) => (
            <button key={name} onClick={() => loginAs(name)}>
              Continue as {name}
            </button>
          ))}
        </div>
        <form
          onSubmit={(event) => {
            event.preventDefault();
            void loginAs(username);
          }}
        >
          <input value={username} onChange={(event) => setUsername(event.target.value)} />
          <button type="submit">Login</button>
        </form>
        {error ? <p className="error-text">{error}</p> : null}
      </section>
    </main>
  );
}
```

- [ ] **Step 6: Wire app state**

Modify `frontend/src/app/App.tsx`:

```tsx
import { QueryClientProvider } from "@tanstack/react-query";
import { useState } from "react";
import { LoginPage } from "../features/auth/LoginPage";
import { queryClient } from "../lib/query";
import type { User } from "../lib/api";

export function App() {
  const [token, setToken] = useState(() => localStorage.getItem("chat.token"));
  const [user, setUser] = useState<User | null>(null);

  return (
    <QueryClientProvider client={queryClient}>
      {token && user ? (
        <div className="app-shell">Logged in as {user.displayName}</div>
      ) : (
        <LoginPage
          onLogin={(nextToken, nextUser) => {
            localStorage.setItem("chat.token", nextToken);
            setToken(nextToken);
            setUser(nextUser);
          }}
        />
      )}
    </QueryClientProvider>
  );
}
```

- [ ] **Step 7: Add login styles**

Append to `frontend/src/styles.css`:

```css
.login-page {
  min-height: 100vh;
  display: grid;
  place-items: center;
  padding: 24px;
}

.login-panel {
  width: min(420px, 100%);
  border: 1px solid #d8dee4;
  border-radius: 8px;
  background: #ffffff;
  padding: 24px;
}

.seed-list {
  display: grid;
  gap: 8px;
  margin: 18px 0;
}

button,
input {
  min-height: 40px;
  border-radius: 6px;
  border: 1px solid #c7d0d9;
  padding: 0 12px;
}

button {
  background: #0f766e;
  color: #ffffff;
  cursor: pointer;
}

.error-text {
  color: #b42318;
}
```

- [ ] **Step 8: Verify frontend tests and build**

Run:

```bash
cd frontend && npm run test -- --run src/lib/api.test.ts && npm run build
```

Expected: PASS.

- [ ] **Step 9: Commit login frontend**

```bash
git add frontend
git commit -m "feat: add frontend login flow"
```

---

## Task 9: Frontend Chat Data and Responsive Layout

**Files:**
- Create: `frontend/src/features/chat/ChatPage.tsx`
- Create: `frontend/src/features/chat/ConversationList.tsx`
- Create: `frontend/src/features/chat/MessageList.tsx`
- Create: `frontend/src/features/chat/MarkdownMessage.tsx`
- Create: `frontend/src/features/chat/MessageComposer.tsx`
- Create: `frontend/src/features/chat/MessageComposer.test.tsx`
- Modify: `frontend/src/lib/api.ts`
- Modify: `frontend/src/app/App.tsx`
- Modify: `frontend/src/styles.css`

- [ ] **Step 1: Extend API types**

Add to `frontend/src/lib/api.ts`:

```ts
export type Conversation = {
  id: number;
  type: "direct" | "group";
  title: string | null;
  lastMessageId: number | null;
  lastMessageAt: string | null;
  unreadCount: number;
};

export type Message = {
  id: number;
  conversationId: number;
  senderId: number;
  contentMarkdown: string;
  contentPlain: string;
  createdAt: string;
};

export function listConversations(token: string): Promise<{ conversations: Conversation[] }> {
  return apiGet("/api/conversations", token);
}

export function listMessages(token: string, conversationId: number): Promise<{ messages: Message[] }> {
  return apiGet(`/api/conversations/${conversationId}/messages?limit=50`, token);
}

export function sendMessage(token: string, conversationId: number, contentMarkdown: string): Promise<Message> {
  return apiPost(`/api/conversations/${conversationId}/messages`, token, { contentMarkdown });
}
```

- [ ] **Step 2: Write composer test**

Create `frontend/src/features/chat/MessageComposer.test.tsx`:

```tsx
import { fireEvent, render, screen } from "@testing-library/react";
import { describe, expect, it, vi } from "vitest";
import { MessageComposer } from "./MessageComposer";

describe("MessageComposer", () => {
  it("submits trimmed message text", () => {
    const onSend = vi.fn();
    render(<MessageComposer disabled={false} onSend={onSend} />);
    fireEvent.change(screen.getByRole("textbox"), { target: { value: " hello " } });
    fireEvent.click(screen.getByRole("button", { name: /send/i }));
    expect(onSend).toHaveBeenCalledWith("hello");
  });
});
```

- [ ] **Step 3: Run composer test to verify it fails**

Run:

```bash
cd frontend && npm run test -- --run src/features/chat/MessageComposer.test.tsx
```

Expected: FAIL because `MessageComposer` does not exist.

- [ ] **Step 4: Implement message composer**

Create `frontend/src/features/chat/MessageComposer.tsx`:

```tsx
import { useState } from "react";

type Props = {
  disabled: boolean;
  onSend: (text: string) => void;
};

export function MessageComposer({ disabled, onSend }: Props) {
  const [value, setValue] = useState("");
  return (
    <form
      className="composer"
      onSubmit={(event) => {
        event.preventDefault();
        const trimmed = value.trim();
        if (!trimmed) return;
        onSend(trimmed);
        setValue("");
      }}
    >
      <textarea
        aria-label="Message"
        value={value}
        disabled={disabled}
        onChange={(event) => setValue(event.target.value)}
      />
      <button type="submit" disabled={disabled}>Send</button>
    </form>
  );
}
```

- [ ] **Step 5: Implement Markdown renderer**

Create `frontend/src/features/chat/MarkdownMessage.tsx`:

```tsx
import ReactMarkdown from "react-markdown";
import remarkGfm from "remark-gfm";

type Props = {
  content: string;
};

export function MarkdownMessage({ content }: Props) {
  return (
    <ReactMarkdown remarkPlugins={[remarkGfm]} skipHtml>
      {content}
    </ReactMarkdown>
  );
}
```

- [ ] **Step 6: Implement list and message views**

Create `frontend/src/features/chat/ConversationList.tsx`:

```tsx
import type { Conversation } from "../../lib/api";

type Props = {
  conversations: Conversation[];
  selectedId: number | null;
  onSelect: (id: number) => void;
};

export function ConversationList({ conversations, selectedId, onSelect }: Props) {
  return (
    <aside className="conversation-list">
      <div className="list-title">Chats</div>
      {conversations.map((conversation) => (
        <button
          key={conversation.id}
          className={conversation.id === selectedId ? "conversation-row selected" : "conversation-row"}
          onClick={() => onSelect(conversation.id)}
        >
          <span>{conversation.title ?? `Conversation ${conversation.id}`}</span>
          {conversation.unreadCount > 0 ? <strong>{conversation.unreadCount}</strong> : null}
        </button>
      ))}
    </aside>
  );
}
```

Create `frontend/src/features/chat/MessageList.tsx`:

```tsx
import type { Message } from "../../lib/api";
import { MarkdownMessage } from "./MarkdownMessage";

type Props = {
  currentUserId: number;
  messages: Message[];
};

export function MessageList({ currentUserId, messages }: Props) {
  return (
    <div className="message-list">
      {[...messages].reverse().map((message) => (
        <article
          key={message.id}
          className={message.senderId === currentUserId ? "message mine" : "message"}
        >
          <MarkdownMessage content={message.contentMarkdown} />
        </article>
      ))}
    </div>
  );
}
```

- [ ] **Step 7: Implement ChatPage**

Create `frontend/src/features/chat/ChatPage.tsx`:

```tsx
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useState } from "react";
import {
  listConversations,
  listMessages,
  sendMessage,
  type Conversation,
  type Message,
  type User
} from "../../lib/api";
import { ConversationList } from "./ConversationList";
import { MessageComposer } from "./MessageComposer";
import { MessageList } from "./MessageList";

type Props = {
  token: string;
  user: User;
};

export function ChatPage({ token, user }: Props) {
  const queryClient = useQueryClient();
  const [selectedId, setSelectedId] = useState<number | null>(null);
  const conversations = useQuery({
    queryKey: ["conversations"],
    queryFn: () => listConversations(token)
  });
  const activeId = selectedId ?? conversations.data?.conversations[0]?.id ?? null;
  const messages = useQuery({
    queryKey: ["messages", activeId],
    queryFn: () => listMessages(token, activeId!),
    enabled: activeId !== null
  });
  const send = useMutation({
    mutationFn: (text: string) => sendMessage(token, activeId!, text),
    onSuccess: (message: Message) => {
      queryClient.setQueryData<{ messages: Message[] }>(["messages", activeId], (old) => ({
        messages: [message, ...(old?.messages ?? [])]
      }));
      void queryClient.invalidateQueries({ queryKey: ["conversations"] });
    }
  });

  return (
    <main className="chat-layout">
      <ConversationList
        conversations={conversations.data?.conversations ?? []}
        selectedId={activeId}
        onSelect={setSelectedId}
      />
      <section className="chat-panel">
        {activeId ? (
          <>
            <header className="chat-header">
              {(conversations.data?.conversations ?? []).find((item: Conversation) => item.id === activeId)?.title ?? `Conversation ${activeId}`}
            </header>
            <MessageList currentUserId={user.id} messages={messages.data?.messages ?? []} />
            <MessageComposer disabled={send.isPending} onSend={(text) => send.mutate(text)} />
          </>
        ) : (
          <div className="empty-state">Select a conversation</div>
        )}
      </section>
    </main>
  );
}
```

- [ ] **Step 8: Wire App to ChatPage**

Modify `frontend/src/app/App.tsx` to render:

```tsx
{token && user ? (
  <ChatPage token={token} user={user} />
) : (
  <LoginPage onLogin={...} />
)}
```

Import `ChatPage`.

- [ ] **Step 9: Add responsive chat styles**

Append to `frontend/src/styles.css`:

```css
.chat-layout {
  min-height: 100vh;
  display: grid;
  grid-template-columns: 320px minmax(0, 1fr);
  background: #eef3f4;
}

.conversation-list {
  border-right: 1px solid #d7e0e3;
  background: #ffffff;
  padding: 12px;
  overflow: auto;
}

.list-title {
  font-weight: 700;
  margin: 8px 8px 14px;
}

.conversation-row {
  width: 100%;
  display: flex;
  justify-content: space-between;
  align-items: center;
  min-height: 48px;
  margin-bottom: 6px;
  background: transparent;
  color: #1d2329;
}

.conversation-row.selected {
  background: #dbeafe;
}

.chat-panel {
  min-width: 0;
  display: grid;
  grid-template-rows: 56px minmax(0, 1fr) auto;
}

.chat-header {
  display: flex;
  align-items: center;
  padding: 0 18px;
  border-bottom: 1px solid #d7e0e3;
  background: #ffffff;
  font-weight: 700;
}

.message-list {
  padding: 18px;
  overflow: auto;
}

.message {
  max-width: min(680px, 78%);
  margin-bottom: 12px;
  padding: 10px 12px;
  border-radius: 8px;
  background: #ffffff;
}

.message.mine {
  margin-left: auto;
  background: #cce7ff;
}

.composer {
  display: grid;
  grid-template-columns: minmax(0, 1fr) 96px;
  gap: 8px;
  padding: 12px;
  background: #ffffff;
  border-top: 1px solid #d7e0e3;
}

.composer textarea {
  min-height: 44px;
  resize: vertical;
  border-radius: 6px;
  border: 1px solid #c7d0d9;
  padding: 10px;
}

@media (max-width: 720px) {
  .chat-layout {
    grid-template-columns: 1fr;
  }

  .conversation-list {
    min-height: 42vh;
    border-right: 0;
    border-bottom: 1px solid #d7e0e3;
  }
}
```

- [ ] **Step 10: Verify frontend**

Run:

```bash
cd frontend && npm run test -- --run src/features/chat/MessageComposer.test.tsx && npm run build
```

Expected: PASS.

- [ ] **Step 11: Commit chat UI**

```bash
git add frontend
git commit -m "feat: add responsive chat ui"
```

---

## Task 10: Frontend WebSocket Recovery and Send Failure State

**Files:**
- Create: `frontend/src/lib/ws.ts`
- Create: `frontend/src/lib/ws.test.ts`
- Modify: `frontend/src/features/chat/ChatPage.tsx`
- Modify: `frontend/src/features/chat/MessageList.tsx`
- Modify: `frontend/src/styles.css`

- [ ] **Step 1: Write WebSocket URL test**

Create `frontend/src/lib/ws.test.ts`:

```ts
import { describe, expect, it } from "vitest";
import { wsUrl } from "./ws";

describe("wsUrl", () => {
  it("converts http origin to ws url", () => {
    expect(wsUrl("http://localhost:5173", "/api/ws")).toBe("ws://localhost:5173/api/ws");
  });

  it("converts https origin to wss url", () => {
    expect(wsUrl("https://example.com", "/api/ws")).toBe("wss://example.com/api/ws");
  });
});
```

- [ ] **Step 2: Run WebSocket test to verify it fails**

Run:

```bash
cd frontend && npm run test -- --run src/lib/ws.test.ts
```

Expected: FAIL because `src/lib/ws.ts` does not exist.

- [ ] **Step 3: Implement WebSocket client helper**

Create `frontend/src/lib/ws.ts`:

```ts
export type ServerEvent = {
  type: "message.created" | "conversation.updated" | "message.failed";
  payload: unknown;
};

export function wsUrl(origin: string, path: string): string {
  const url = new URL(path, origin);
  url.protocol = url.protocol === "https:" ? "wss:" : "ws:";
  return url.toString();
}

export function connectEvents(token: string, onEvent: (event: ServerEvent) => void, onStatus: (connected: boolean) => void) {
  let closed = false;
  let socket: WebSocket | null = null;

  function connect() {
    socket = new WebSocket(`${wsUrl(window.location.origin, "/api/ws")}?token=${encodeURIComponent(token)}`);
    socket.onopen = () => onStatus(true);
    socket.onclose = () => {
      onStatus(false);
      if (!closed) window.setTimeout(connect, 1200);
    };
    socket.onmessage = (message) => {
      onEvent(JSON.parse(message.data) as ServerEvent);
    };
  }

  connect();

  return () => {
    closed = true;
    socket?.close();
  };
}
```

- [ ] **Step 4: Add event handling in ChatPage**

Modify `frontend/src/features/chat/ChatPage.tsx`:

```tsx
import { useEffect, useState } from "react";
import { connectEvents } from "../../lib/ws";
```

Inside `ChatPage`, add:

```tsx
const [connected, setConnected] = useState(false);

useEffect(() => {
  return connectEvents(
    token,
    (event) => {
      if (event.type === "message.created") {
        void queryClient.invalidateQueries({ queryKey: ["messages"] });
      }
      if (event.type === "conversation.updated") {
        void queryClient.invalidateQueries({ queryKey: ["conversations"] });
      }
    },
    (nextConnected) => {
      setConnected(nextConnected);
      if (nextConnected) {
        void queryClient.invalidateQueries({ queryKey: ["conversations"] });
        if (activeId) void queryClient.invalidateQueries({ queryKey: ["messages", activeId] });
      }
    }
  );
}, [token, queryClient, activeId]);
```

Render connection status in the chat header:

```tsx
<span className={connected ? "ws-status connected" : "ws-status"}>{connected ? "Live" : "Reconnecting"}</span>
```

- [ ] **Step 5: Add local failed send tracking**

In `ChatPage`, add local failed state:

```tsx
const [failedMessages, setFailedMessages] = useState<{ id: string; contentMarkdown: string }[]>([]);
```

In send mutation:

```tsx
onError: (_error, text) => {
  setFailedMessages((items) => [...items, { id: crypto.randomUUID(), contentMarkdown: text }]);
}
```

Pass failed messages to `MessageList` and render them with a retry button.

- [ ] **Step 6: Verify frontend**

Run:

```bash
cd frontend && npm run test -- --run src/lib/ws.test.ts && npm run build
```

Expected: PASS.

- [ ] **Step 7: Commit realtime frontend**

```bash
git add frontend
git commit -m "feat: add websocket recovery in frontend"
```

---

## Task 11: Development Admin Page

**Files:**
- Create: `frontend/src/features/admin/AdminPage.tsx`
- Modify: `frontend/src/lib/api.ts`
- Modify: `frontend/src/app/App.tsx`
- Modify: `frontend/src/styles.css`

- [ ] **Step 1: Extend API client with admin calls**

Add to `frontend/src/lib/api.ts`:

```ts
export type AdminRow = Record<string, string | number | null>;

export function adminUsers(token: string): Promise<{ users: AdminRow[] }> {
  return apiGet("/api/admin/users", token);
}

export function adminConversations(token: string): Promise<{ conversations: AdminRow[] }> {
  return apiGet("/api/admin/conversations", token);
}

export function adminMessages(token: string): Promise<{ messages: AdminRow[] }> {
  return apiGet("/api/admin/messages", token);
}
```

- [ ] **Step 2: Implement admin page**

Create `frontend/src/features/admin/AdminPage.tsx`:

```tsx
import { useQuery } from "@tanstack/react-query";
import { adminConversations, adminMessages, adminUsers } from "../../lib/api";

type Props = {
  token: string;
};

export function AdminPage({ token }: Props) {
  const users = useQuery({ queryKey: ["admin", "users"], queryFn: () => adminUsers(token) });
  const conversations = useQuery({ queryKey: ["admin", "conversations"], queryFn: () => adminConversations(token) });
  const messages = useQuery({ queryKey: ["admin", "messages"], queryFn: () => adminMessages(token) });

  return (
    <main className="admin-page">
      <h1>Admin</h1>
      <AdminTable title="Users" rows={users.data?.users ?? []} />
      <AdminTable title="Conversations" rows={conversations.data?.conversations ?? []} />
      <AdminTable title="Recent Messages" rows={messages.data?.messages ?? []} />
    </main>
  );
}

function AdminTable({ title, rows }: { title: string; rows: Record<string, unknown>[] }) {
  const columns = Object.keys(rows[0] ?? {});
  return (
    <section className="admin-section">
      <h2>{title}</h2>
      <table>
        <thead>
          <tr>{columns.map((column) => <th key={column}>{column}</th>)}</tr>
        </thead>
        <tbody>
          {rows.map((row, index) => (
            <tr key={index}>
              {columns.map((column) => <td key={column}>{String(row[column] ?? "")}</td>)}
            </tr>
          ))}
        </tbody>
      </table>
    </section>
  );
}
```

- [ ] **Step 3: Add simple route switch**

Modify `frontend/src/app/App.tsx` to render admin when the path is `/admin`:

```tsx
const path = window.location.pathname;
if (token && user && path === "/admin") {
  return <QueryClientProvider client={queryClient}><AdminPage token={token} /></QueryClientProvider>;
}
```

Import `AdminPage`.

- [ ] **Step 4: Add admin styles**

Append to `frontend/src/styles.css`:

```css
.admin-page {
  padding: 24px;
}

.admin-section {
  margin-bottom: 28px;
  overflow-x: auto;
}

table {
  width: 100%;
  border-collapse: collapse;
  background: #ffffff;
}

th,
td {
  border: 1px solid #d7e0e3;
  padding: 8px 10px;
  text-align: left;
  white-space: nowrap;
}

th {
  background: #eef3f4;
}
```

- [ ] **Step 5: Verify frontend**

Run:

```bash
cd frontend && npm run build
```

Expected: PASS.

- [ ] **Step 6: Commit admin frontend**

```bash
git add frontend
git commit -m "feat: add development admin page"
```

---

## Task 12: End-to-End Verification and Documentation

**Files:**
- Create: `frontend/playwright.config.ts`
- Create: `frontend/tests/chat.e2e.ts`
- Modify: `README.md`
- Modify: `backend/cmd/api/main.go`
- Modify: `frontend/src/features/chat/ChatPage.tsx`

- [ ] **Step 1: Add Playwright config**

Create `frontend/playwright.config.ts`:

```ts
import { defineConfig, devices } from "@playwright/test";

export default defineConfig({
  testDir: "./tests",
  use: {
    baseURL: "http://localhost:5173",
    trace: "on-first-retry"
  },
  projects: [
    { name: "chromium", use: { ...devices["Desktop Chrome"] } },
    { name: "mobile", use: { ...devices["Pixel 7"] } }
  ]
});
```

- [ ] **Step 2: Add E2E smoke test**

Create `frontend/tests/chat.e2e.ts`:

```ts
import { expect, test } from "@playwright/test";

test("user can chat with mock agent", async ({ page }) => {
  await page.goto("/");
  await page.getByRole("button", { name: /continue as alice/i }).click();
  await expect(page.getByText(/Chats/i)).toBeVisible();
  await page.getByRole("button", { name: /Mock Agent/i }).click();
  await page.getByRole("textbox", { name: /message/i }).fill("hello agent");
  await page.getByRole("button", { name: /send/i }).click();
  await expect(page.getByText(/hello agent/i)).toBeVisible();
  await expect(page.getByText(/Mock Agent received: hello agent/i)).toBeVisible({ timeout: 5000 });
});
```

- [ ] **Step 3: Verify seeded conversation titles in the E2E test**

Modify `frontend/tests/chat.e2e.ts` so the Mock Agent conversation is selected by its computed title from `ListForUser`:

```ts
await page.getByRole("button", { name: /Mock Agent/i }).click();
```

The final test should include these assertions:

```ts
await expect(page.getByText(/Chats/i)).toBeVisible();
await page.getByRole("button", { name: /Mock Agent/i }).click();
await page.getByRole("textbox", { name: /message/i }).fill("hello agent");
await page.getByRole("button", { name: /send/i }).click();
await expect(page.getByText(/hello agent/i)).toBeVisible();
await expect(page.getByText(/Mock Agent received: hello agent/i)).toBeVisible({ timeout: 5000 });
```

- [ ] **Step 4: Update README with verified commands**

Modify `README.md`:

```markdown
## Verification

Backend:

```bash
cd backend
go test ./...
```

Frontend:

```bash
cd frontend
npm run test
npm run build
```

End-to-end:

```bash
docker compose up -d mysql
cd backend && go run ./cmd/migrate
cd backend && go run ./cmd/api
cd frontend && npm run dev
cd frontend && npm run test:e2e
```
```

- [ ] **Step 5: Run full backend verification**

Run:

```bash
cd backend && go test ./...
```

Expected: PASS.

- [ ] **Step 6: Run full frontend verification**

Run:

```bash
cd frontend && npm run test && npm run build
```

Expected: PASS.

- [ ] **Step 7: Run local app smoke test**

Start services in separate terminals:

```bash
docker compose up -d mysql
cd backend && go run ./cmd/migrate
cd backend && go run ./cmd/api
cd frontend && npm run dev
```

Then run:

```bash
cd frontend && npm run test:e2e
```

Expected: Playwright passes for desktop and mobile projects.

- [ ] **Step 8: Commit verification and docs**

```bash
git add README.md backend frontend
git commit -m "test: add e2e chat verification"
```

---

## Self-Review Results

Spec coverage:

- Seed-user login is covered by Task 3 and Task 8.
- Direct/group conversations are covered by Task 4 and Task 9.
- Markdown messages are covered by Task 4 and Task 9.
- Lightweight realtime delivery is covered by Task 5 and Task 10.
- Unread counts are covered by Task 4.
- Mock Agent is covered by Task 6 and Task 12.
- Development admin page is covered by Task 7 and Task 11.
- Docker Compose and local development are covered by Task 1, Task 2, and Task 12.
- Tests are covered incrementally in every implementation task and fully in Task 12.

Consistency check:

- REST remains the write path for messages.
- WebSocket remains event delivery only.
- Persisted messages have no server-side failed status.
- Agent users remain ordinary users with `user_type = "agent"`.
- Workspace support stays as a default `workspace_id` without UI switching.

Implementation risk:

- Auth middleware and route wiring must keep `go test ./...` green after each backend task because Task 3 introduces adapters before all route modules exist.
- The WebSocket query-token path in frontend must match backend auth middleware behavior.
- Admin `rowsToMaps` should convert `[]byte` database values to strings before JSON encoding.
