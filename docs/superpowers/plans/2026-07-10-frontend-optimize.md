# Frontend Visual & UX Optimization Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Apply the approved frontend design spec to upgrade the existing React + Vite chat interface with a cohesive token system, Geist typography, full UI states, and polished micro-interactions without migrating frameworks.

**Architecture:** Keep the existing vanilla CSS + React stack. Centralize design tokens in `frontend/src/styles.css`. Update each page/component to consume the tokens and implement the required states (loading, empty, error, hover, active, focus). Add minimal font dependencies `@fontsource-variable/geist` and `@fontsource-variable/geist-mono` for typography.

**Tech Stack:** React, Vite, TypeScript, vanilla CSS, TanStack Query, `@fontsource-variable/geist` and `@fontsource-variable/geist-mono` font packages, Vitest, React Testing Library.

---

## File Structure

Files to create: none.

Files to modify:

- `frontend/package.json` — add `@fontsource-variable/geist` and `@fontsource-variable/geist-mono` dependencies.
- `frontend/package-lock.json` — updated by npm install.
- `frontend/index.html` — preload Geist fonts, add meta tags.
- `frontend/src/main.tsx` — import Geist font CSS.
- `frontend/src/styles.css` — full rewrite with design tokens and component styles.
- `frontend/src/features/auth/LoginPage.tsx` — new layout, loading/error states.
- `frontend/src/features/chat/ChatPage.tsx` — loading/error/empty states, mobile back navigation, status pill.
- `frontend/src/features/chat/ConversationList.tsx` — new row styling, unread badge.
- `frontend/src/features/chat/MessageList.tsx` — bubble styling, metadata, failed state, entrance animation.
- `frontend/src/features/chat/MessageComposer.tsx` — auto-resize, keyboard handling, spinner state.
- `frontend/src/features/chat/MarkdownMessage.tsx` — no JSX changes; code block styling handled in `styles.css`.
- `frontend/src/features/admin/AdminPage.tsx` — card/table styling, loading state.
- `frontend/src/features/chat/MessageComposer.test.tsx` — update for new behavior.
- `frontend/src/features/chat/ConversationList.test.tsx` — new test file.
- `frontend/src/features/chat/MessageList.test.tsx` — new test file.

---

## Task 1: Add Geist Font Dependency and Preload

**Files:**
- Modify: `frontend/package.json`
- Modify: `frontend/index.html`
- Modify: `frontend/src/main.tsx`

- [ ] **Step 1: Install the `@fontsource-variable/geist` and `@fontsource-variable/geist-mono` packages**

Run:
```bash
cd frontend && npm install @fontsource-variable/geist @fontsource-variable/geist-mono
```

Expected: `package.json` and `package-lock.json` are updated. No errors.

- [ ] **Step 2: Import Geist font CSS in main.tsx**

Modify `frontend/src/main.tsx` to import the font CSS at the top:

```tsx
import "@fontsource-variable/geist";
import "@fontsource-variable/geist-mono";
import React from "react";
import ReactDOM from "react-dom/client";
import { App } from "./app/App";
import "./styles.css";

ReactDOM.createRoot(document.getElementById("root")!).render(
  <React.StrictMode>
    <App />
  </React.StrictMode>,
);
```

If the exact import path differs in the installed version, use the path documented by the `geist` package.

- [ ] **Step 3: Preload fonts and add meta tags in index.html**

Modify `frontend/index.html`:

```html
<!doctype html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <meta name="description" content="Vibework Chat — a lightweight web chat for teams." />
    <meta name="theme-color" content="#f9fafb" />
    <link rel="preconnect" href="https://fonts.googleapis.com" crossorigin />
    <title>Vibework Chat</title>
  </head>
  <body>
    <div id="root"></div>
    <script type="module" src="/src/main.tsx"></script>
  </body>
</html>
```

Note: `@fontsource-variable/geist` and `@fontsource-variable/geist-mono` are imported via CSS in `main.tsx`; the preconnect is optional but keeps the head tidy.

- [ ] **Step 4: Verify build still compiles**

Run:
```bash
cd frontend && npm run lint
```

Expected: TypeScript passes with no errors.

- [ ] **Step 5: Commit**

```bash
git add frontend/package.json frontend/package-lock.json frontend/index.html frontend/src/main.tsx
git commit -m "chore: add geist font dependency and preload"
```

---

## Task 2: Rewrite Global Styles with Design Tokens

**Files:**
- Modify: `frontend/src/styles.css`

- [ ] **Step 1: Replace styles.css with the new token-based stylesheet**

Write `frontend/src/styles.css`:

```css
:root {
  --font-sans: "Geist", ui-sans-serif, system-ui, -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif;
  --font-mono: "Geist Mono", ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;

  --color-bg: #f9fafb;
  --color-surface: #ffffff;
  --color-surface-raised: #ffffff;
  --color-border: #e4e4e7;
  --color-border-subtle: #f4f4f5;

  --color-text: #18181b;
  --color-text-secondary: #52525b;
  --color-text-tertiary: #71717a;

  --color-accent: #0d9488;
  --color-accent-hover: #0f766e;
  --color-accent-subtle: #ccfbf1;
  --color-accent-text: #134e4a;

  --color-error: #dc2626;
  --color-error-bg: #fef2f2;
  --color-error-border: #fecaca;

  --color-success: #16a34a;
  --color-warning: #ca8a04;

  --radius-sm: 6px;
  --radius-md: 10px;
  --radius-lg: 16px;
  --radius-xl: 24px;

  --shadow-sm: 0 1px 2px 0 rgba(0, 0, 0, 0.04);
  --shadow-md: 0 4px 12px -2px rgba(0, 0, 0, 0.06);
  --shadow-lg: 0 12px 24px -4px rgba(0, 0, 0, 0.08);

  --ease-out: cubic-bezier(0.16, 1, 0.3, 1);

  font-family: var(--font-sans);
  color: var(--color-text);
  background: var(--color-bg);
}

* {
  box-sizing: border-box;
}

html {
  scroll-behavior: smooth;
}

body {
  margin: 0;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
}

button,
input,
textarea {
  font: inherit;
}

:focus-visible {
  outline: 2px solid var(--color-accent);
  outline-offset: 2px;
}

@media (prefers-reduced-motion: reduce) {
  *,
  *::before,
  *::after {
    animation-duration: 0.01ms !important;
    animation-iteration-count: 1 !important;
    transition-duration: 0.01ms !important;
  }

  html {
    scroll-behavior: auto;
  }
}

/* Utility helpers */
.sr-only {
  position: absolute;
  width: 1px;
  height: 1px;
  padding: 0;
  margin: -1px;
  overflow: hidden;
  clip: rect(0, 0, 0, 0);
  white-space: nowrap;
  border: 0;
}

@keyframes message-in {
  from {
    opacity: 0;
    transform: translateY(8px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

@keyframes pulse-dot {
  0%, 100% {
    opacity: 1;
  }
  50% {
    opacity: 0.4;
  }
}

@keyframes skeleton-shimmer {
  0% {
    background-position: 200% 0;
  }
  100% {
    background-position: -200% 0;
  }
}

.skeleton {
  border-radius: var(--radius-md);
  background: linear-gradient(90deg, var(--color-border-subtle) 25%, var(--color-border) 50%, var(--color-border-subtle) 75%);
  background-size: 200% 100%;
  animation: skeleton-shimmer 1.5s infinite linear;
}

/* Spinner */
.spinner {
  width: 16px;
  height: 16px;
  border: 2px solid currentColor;
  border-right-color: transparent;
  border-radius: 50%;
  animation: spin 0.6s linear infinite;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

/* App shell */
.app-shell {
  min-height: 100dvh;
  display: grid;
  place-items: center;
  padding: 24px;
}

/* Login page */
.login-page {
  min-height: 100dvh;
  display: grid;
  place-items: center;
  padding: 24px;
}

.login-panel {
  width: min(420px, 100%);
  background: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-md);
  padding: 32px;
}

.login-panel h1 {
  font-size: 1.5rem;
  font-weight: 600;
  letter-spacing: -0.025em;
  margin: 0 0 24px;
}

.seed-list {
  display: grid;
  gap: 10px;
  margin: 0 0 24px;
}

.seed-button {
  display: flex;
  align-items: center;
  gap: 12px;
  width: 100%;
  min-height: 48px;
  padding: 10px 14px;
  background: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  color: var(--color-text);
  cursor: pointer;
  transition: all 0.2s var(--ease-out);
}

.seed-button:hover {
  background: var(--color-bg);
  border-color: var(--color-accent);
  transform: translateY(-1px);
}

.seed-button:active {
  transform: scale(0.99);
}

.seed-button:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.seed-avatar {
  display: grid;
  place-items: center;
  width: 32px;
  height: 32px;
  border-radius: 50%;
  background: var(--color-border-subtle);
  color: var(--color-text-secondary);
  font-size: 0.875rem;
  font-weight: 600;
  flex-shrink: 0;
}

.login-panel form {
  display: grid;
  gap: 16px;
}

.field-label {
  display: grid;
  gap: 6px;
  font-size: 0.875rem;
  font-weight: 500;
}

.text-input {
  min-height: 44px;
  padding: 0 12px;
  background: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  color: var(--color-text);
  transition: all 0.2s var(--ease-out);
}

.text-input:hover {
  border-color: var(--color-text-tertiary);
}

.text-input:focus {
  border-color: var(--color-accent);
  box-shadow: 0 0 0 3px var(--color-accent-subtle);
  outline: none;
}

.text-input--error,
.text-input--error:focus {
  border-color: var(--color-error);
  box-shadow: 0 0 0 3px rgba(220, 38, 38, 0.12);
}

.primary-button {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  min-height: 44px;
  padding: 0 18px;
  background: var(--color-accent);
  border: none;
  border-radius: var(--radius-md);
  color: #ffffff;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s var(--ease-out);
}

.primary-button:hover {
  background: var(--color-accent-hover);
  transform: translateY(-1px);
}

.primary-button:active {
  transform: scale(0.98);
}

.primary-button:disabled {
  opacity: 0.6;
  cursor: not-allowed;
  transform: none;
}

.error-banner {
  margin-top: 16px;
  padding: 12px;
  background: var(--color-error-bg);
  border: 1px solid var(--color-error-border);
  border-radius: var(--radius-md);
  color: var(--color-error);
  font-size: 0.875rem;
}

.error-text {
  color: var(--color-error);
  font-size: 0.875rem;
  margin: 4px 0 0;
}

/* Chat layout */
.chat-layout {
  min-height: 100dvh;
  display: grid;
  grid-template-columns: 320px minmax(0, 1fr);
  background: var(--color-bg);
}

.conversation-list {
  border-right: 1px solid var(--color-border);
  background: var(--color-surface);
  padding: 16px 12px;
  overflow: auto;
}

.list-title {
  font-size: 0.75rem;
  font-weight: 600;
  letter-spacing: 0.05em;
  text-transform: uppercase;
  color: var(--color-text-tertiary);
  margin: 0 8px 12px;
}

.conversation-row {
  position: relative;
  width: 100%;
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 10px;
  min-height: 44px;
  padding: 10px 12px;
  margin-bottom: 6px;
  background: transparent;
  border: none;
  border-radius: var(--radius-md);
  color: var(--color-text);
  font-size: 0.9375rem;
  font-weight: 500;
  text-align: left;
  cursor: pointer;
  transition: all 0.2s var(--ease-out);
}

.conversation-row:hover {
  background: var(--color-border-subtle);
}

.conversation-row:active {
  transform: scale(0.99);
}

.conversation-row.selected {
  background: var(--color-accent-subtle);
  color: var(--color-accent-text);
}

.conversation-row.selected::before {
  content: "";
  position: absolute;
  left: 0;
  top: 50%;
  transform: translateY(-50%);
  width: 3px;
  height: 20px;
  background: var(--color-accent);
  border-radius: 0 3px 3px 0;
}

.unread-badge {
  display: grid;
  place-items: center;
  min-width: 18px;
  height: 18px;
  padding: 0 5px;
  background: var(--color-accent);
  color: #ffffff;
  font-size: 0.6875rem;
  font-weight: 600;
  border-radius: 9999px;
}

.chat-panel {
  min-width: 0;
  display: grid;
  grid-template-rows: 56px minmax(0, 1fr) auto;
  background: var(--color-bg);
}

.chat-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  padding: 0 18px;
  background: var(--color-surface);
  border-bottom: 1px solid var(--color-border);
}

.chat-header-title {
  font-weight: 600;
  font-size: 1rem;
}

.chat-header-start {
  display: flex;
  align-items: center;
  gap: 10px;
  min-width: 0;
}

.back-button {
  display: none;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  padding: 0;
  background: transparent;
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  color: var(--color-text-secondary);
  cursor: pointer;
  transition: all 0.2s var(--ease-out);
}

.back-button:hover {
  background: var(--color-border-subtle);
  color: var(--color-text);
}

.ws-status {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 4px 10px;
  border-radius: 9999px;
  background: var(--color-border-subtle);
  color: var(--color-text-secondary);
  font-size: 0.75rem;
  font-weight: 600;
}

.ws-status.connected {
  background: var(--color-accent-subtle);
  color: var(--color-accent-text);
}

.ws-status-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: currentColor;
}

.ws-status:not(.connected) .ws-status-dot {
  animation: pulse-dot 2s infinite;
}

.message-list {
  padding: 18px;
  overflow: auto;
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.message {
  max-width: min(680px, 78%);
  padding: 12px 14px;
  border-radius: var(--radius-lg);
  background: var(--color-surface);
  border: 1px solid var(--color-border);
  box-shadow: var(--shadow-sm);
  animation: message-in 0.25s var(--ease-out);
}

.message.mine {
  margin-left: auto;
  background: var(--color-accent-subtle);
  border-color: rgba(13, 148, 136, 0.16);
  border-bottom-right-radius: var(--radius-sm);
}

.message:not(.mine) {
  border-bottom-left-radius: var(--radius-sm);
}

.message.failed {
  background: var(--color-error-bg);
  border-color: var(--color-error-border);
}

.message-content {
  font-size: 0.9375rem;
  line-height: 1.5;
  word-wrap: break-word;
}

.message-content p:first-child {
  margin-top: 0;
}

.message-content p:last-child {
  margin-bottom: 0;
}

.message-content pre {
  background: var(--color-bg);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-sm);
  padding: 12px;
  overflow-x: auto;
}

.message-content code {
  font-family: var(--font-mono);
  font-size: 0.875rem;
}

.message-content :not(pre) > code {
  background: var(--color-bg);
  padding: 2px 5px;
  border-radius: 4px;
}

.message-meta {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-top: 8px;
  font-size: 0.75rem;
  color: var(--color-text-tertiary);
}

.message-failed-actions {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 10px;
  margin-top: 10px;
  padding-top: 10px;
  border-top: 1px solid var(--color-error-border);
  color: var(--color-error);
  font-size: 0.8125rem;
}

.message-failed-actions button {
  min-height: 28px;
  padding: 0 10px;
  background: var(--color-error);
  border: none;
  border-radius: var(--radius-sm);
  color: #ffffff;
  font-size: 0.75rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s var(--ease-out);
}

.message-failed-actions button:hover {
  background: #b91c1c;
}

.composer {
  display: grid;
  grid-template-columns: minmax(0, 1fr) auto;
  gap: 10px;
  align-items: end;
  padding: 12px;
  background: var(--color-surface);
  border-top: 1px solid var(--color-border);
}

.composer textarea {
  min-height: 44px;
  max-height: 180px;
  resize: none;
  overflow-y: auto;
  border-radius: var(--radius-md);
  border: 1px solid var(--color-border);
  padding: 10px 12px;
  background: var(--color-surface);
  color: var(--color-text);
  line-height: 1.5;
  transition: all 0.2s var(--ease-out);
}

.composer textarea:hover {
  border-color: var(--color-text-tertiary);
}

.composer textarea:focus {
  border-color: var(--color-accent);
  box-shadow: 0 0 0 3px var(--color-accent-subtle);
  outline: none;
}

.composer textarea:disabled {
  background: var(--color-bg);
  opacity: 0.7;
  cursor: not-allowed;
}

.empty-state {
  display: grid;
  place-items: center;
  padding: 32px;
  color: var(--color-text-secondary);
  text-align: center;
}

.empty-state-title {
  font-weight: 600;
  color: var(--color-text);
  margin: 0 0 6px;
}

.empty-state-hint {
  font-size: 0.875rem;
  color: var(--color-text-tertiary);
  margin: 0;
}

.error-state {
  display: grid;
  place-items: center;
  gap: 12px;
  padding: 32px;
  text-align: center;
}

.error-state-message {
  color: var(--color-error);
  font-weight: 500;
}

.loading-state {
  display: grid;
  gap: 12px;
  padding: 18px;
}

/* Admin page */
.admin-page {
  min-height: 100dvh;
  padding: 24px;
}

.admin-page-inner {
  max-width: 1200px;
  margin: 0 auto;
}

.admin-page h1 {
  font-size: 1.5rem;
  font-weight: 600;
  letter-spacing: -0.025em;
  margin: 0 0 24px;
}

.admin-section {
  margin-bottom: 28px;
  background: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-sm);
  overflow: hidden;
}

.admin-section h2 {
  font-size: 0.875rem;
  font-weight: 600;
  letter-spacing: 0.03em;
  text-transform: uppercase;
  color: var(--color-text-secondary);
  margin: 0;
  padding: 16px 20px;
  background: var(--color-bg);
  border-bottom: 1px solid var(--color-border);
}

.admin-table-wrap {
  overflow-x: auto;
}

table {
  width: 100%;
  border-collapse: collapse;
  background: var(--color-surface);
}

th,
td {
  padding: 12px 16px;
  text-align: left;
  border-bottom: 1px solid var(--color-border);
  white-space: nowrap;
}

th {
  font-size: 0.6875rem;
  font-weight: 600;
  letter-spacing: 0.05em;
  text-transform: uppercase;
  color: var(--color-text-tertiary);
  background: var(--color-bg);
}

tbody tr:hover {
  background: var(--color-bg);
}

tbody tr:last-child td {
  border-bottom: none;
}

.tabular-nums {
  font-variant-numeric: tabular-nums;
}

/* Mobile */
@media (max-width: 767px) {
  .chat-layout {
    grid-template-columns: 1fr;
  }

  .conversation-list.mobile-hidden {
    display: none;
  }

  .chat-panel {
    display: none;
  }

  .chat-panel.mobile-visible {
    display: grid;
  }

  .back-button {
    display: inline-flex;
  }
}
```

- [ ] **Step 2: Verify lint passes**

Run:
```bash
cd frontend && npm run lint
```

Expected: TypeScript passes. CSS does not fail the build.

- [ ] **Step 3: Commit**

```bash
git add frontend/src/styles.css
git commit -m "feat: add design tokens and global styles"
```

---

## Task 3: Refactor Login Page

**Files:**
- Modify: `frontend/src/features/auth/LoginPage.tsx`

- [ ] **Step 1: Update LoginPage.tsx with new layout and states**

Replace the component with:

```tsx
import { useState } from "react";
import { devLogin, type User } from "../../lib/api";

const seedUsers = ["alice", "bob", "carol"];

const avatarColors: Record<string, string> = {
  alice: "#e4e4e7",
  bob: "#d4d4d8",
  carol: "#e4e4e7",
};

type Props = {
  onLogin: (token: string, user: User) => void;
};

export function LoginPage({ onLogin }: Props) {
  const [username, setUsername] = useState("alice");
  const [error, setError] = useState<string | null>(null);
  const [inputError, setInputError] = useState(false);
  const [loadingSeed, setLoadingSeed] = useState<string | null>(null);

  async function loginAs(name: string) {
    setError(null);
    setInputError(false);
    setLoadingSeed(name);
    try {
      const result = await devLogin(name);
      onLogin(result.token, result.user);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Login failed");
    } finally {
      setLoadingSeed(null);
    }
  }

  return (
    <main className="login-page">
      <section className="login-panel">
        <h1>Vibework Chat</h1>
        <div className="seed-list" role="list">
          {seedUsers.map((name) => (
            <button
              key={name}
              className="seed-button"
              disabled={loadingSeed !== null}
              onClick={() => void loginAs(name)}
              role="listitem"
              type="button"
            >
              <span
                className="seed-avatar"
                style={{ background: avatarColors[name] ?? "#e4e4e7" }}
                aria-hidden="true"
              >
                {name[0]?.toUpperCase()}
              </span>
              <span>Continue as {name}</span>
              {loadingSeed === name ? <span className="spinner" aria-hidden="true" /> : null}
            </button>
          ))}
        </div>
        <form
          onSubmit={(event) => {
            event.preventDefault();
            const trimmed = username.trim();
            if (!trimmed) {
              setInputError(true);
              return;
            }
            void loginAs(trimmed);
          }}
        >
          <label className="field-label">
            Username
            <input
              className={inputError ? "text-input text-input--error" : "text-input"}
              value={username}
              onChange={(event) => {
                setUsername(event.target.value);
                setInputError(false);
              }}
              aria-invalid={inputError}
              aria-describedby={inputError ? "username-error" : undefined}
            />
            {inputError ? (
              <p id="username-error" className="error-text">
                Please enter a username
              </p>
            ) : null}
          </label>
          <button className="primary-button" type="submit" disabled={loadingSeed !== null}>
            {loadingSeed === username ? <span className="spinner" aria-hidden="true" /> : null}
            Login
          </button>
        </form>
        {error ? <div className="error-banner">{error}</div> : null}
      </section>
    </main>
  );
}
```

- [ ] **Step 2: Verify lint passes**

Run:
```bash
cd frontend && npm run lint
```

Expected: PASS.

- [ ] **Step 3: Commit**

```bash
git add frontend/src/features/auth/LoginPage.tsx
git commit -m "feat: redesign login page with loading and error states"
```

---

## Task 4: Add ConversationList Tests

**Files:**
- Create: `frontend/src/features/chat/ConversationList.test.tsx`

- [ ] **Step 1: Create the test file**

Write `frontend/src/features/chat/ConversationList.test.tsx`:

```tsx
import { fireEvent, render, screen } from "@testing-library/react";
import { describe, expect, it, vi } from "vitest";
import { ConversationList } from "./ConversationList";

describe("ConversationList", () => {
  const conversations = [
    { id: 1, title: "Project Room", unreadCount: 3 },
    { id: 2, title: "Direct with Bob", unreadCount: 0 },
  ];

  it("renders conversation titles", () => {
    render(<ConversationList conversations={conversations} selectedId={null} onSelect={() => {}} />);
    expect(screen.getByText("Project Room")).toBeInTheDocument();
    expect(screen.getByText("Direct with Bob")).toBeInTheDocument();
  });

  it("shows unread badge only for conversations with unread messages", () => {
    render(<ConversationList conversations={conversations} selectedId={null} onSelect={() => {}} />);
    expect(screen.getByText("3")).toBeInTheDocument();
    expect(screen.queryByText("0")).not.toBeInTheDocument();
  });

  it("calls onSelect when a conversation is clicked", () => {
    const onSelect = vi.fn();
    render(<ConversationList conversations={conversations} selectedId={null} onSelect={onSelect} />);
    fireEvent.click(screen.getByText("Project Room"));
    expect(onSelect).toHaveBeenCalledWith(1);
  });

  it("marks selected conversation", () => {
    const { container } = render(
      <ConversationList conversations={conversations} selectedId={2} onSelect={() => {}} />,
    );
    const selected = container.querySelector(".selected");
    expect(selected).toHaveTextContent("Direct with Bob");
  });
});
```

- [ ] **Step 2: Run the test to verify it fails**

Run:
```bash
cd frontend && npm run test -- src/features/chat/ConversationList.test.tsx
```

Expected: FAIL because the file is new and may reference missing helpers, or because the badge text "0" is not rendered (which is expected behavior).

- [ ] **Step 3: Update ConversationList.tsx to match design**

Modify `frontend/src/features/chat/ConversationList.tsx`:

```tsx
import type { Conversation } from "../../lib/api";

type Props = {
  conversations: Conversation[];
  selectedId: number | null;
  onSelect: (id: number) => void;
  className?: string;
};

export function ConversationList({ conversations, selectedId, onSelect, className }: Props) {
  return (
    <aside className={`conversation-list ${className ?? ""}`.trim()} aria-label="Conversations">
      <div className="list-title">Chats</div>
      {conversations.map((conversation) => (
        <button
          key={conversation.id}
          className={conversation.id === selectedId ? "conversation-row selected" : "conversation-row"}
          onClick={() => onSelect(conversation.id)}
          type="button"
          aria-current={conversation.id === selectedId ? "true" : undefined}
        >
          <span>{conversation.title ?? `Conversation ${conversation.id}`}</span>
          {conversation.unreadCount > 0 ? (
            <span className="unread-badge" aria-label={`${conversation.unreadCount} unread`}>
              {conversation.unreadCount}
            </span>
          ) : null}
        </button>
      ))}
    </aside>
  );
}
```

- [ ] **Step 4: Run tests to verify they pass**

Run:
```bash
cd frontend && npm run test -- src/features/chat/ConversationList.test.tsx
```

Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add frontend/src/features/chat/ConversationList.tsx frontend/src/features/chat/ConversationList.test.tsx
git commit -m "feat: restyle conversation list and add tests"
```

---

## Task 5: Refactor MessageList and Add Tests

**Files:**
- Modify: `frontend/src/features/chat/MessageList.tsx`
- Create: `frontend/src/features/chat/MessageList.test.tsx`

- [ ] **Step 1: Update MessageList.tsx**

Modify `frontend/src/features/chat/MessageList.tsx`:

```tsx
import type { Message } from "../../lib/api";
import { MarkdownMessage } from "./MarkdownMessage";

type Props = {
  currentUserId: number;
  messages: Message[];
  failedMessages: { id: string; contentMarkdown: string }[];
  onRetry: (id: string, contentMarkdown: string) => void;
};

function formatTime(dateString: string): string {
  const date = new Date(dateString);
  return date.toLocaleTimeString([], { hour: "2-digit", minute: "2-digit" });
}

export function MessageList({ currentUserId, messages, failedMessages, onRetry }: Props) {
  return (
    <div className="message-list" role="log" aria-live="polite" aria-label="Messages">
      {[...messages].reverse().map((message) => {
        const isMine = message.senderId === currentUserId;
        return (
          <article
            key={message.id}
            className={isMine ? "message mine" : "message"}
          >
            <div className="message-content">
              <MarkdownMessage content={message.contentMarkdown} />
            </div>
            <div className="message-meta">
              <span>{isMine ? "You" : `User ${message.senderId}`}</span>
              <span aria-hidden="true">·</span>
              <time dateTime={message.createdAt}>{formatTime(message.createdAt)}</time>
            </div>
          </article>
        );
      })}
      {failedMessages.map((message) => (
        <article key={message.id} className="message mine failed">
          <div className="message-content">
            <MarkdownMessage content={message.contentMarkdown} />
          </div>
          <div className="message-failed-actions">
            <span>Failed to send</span>
            <button type="button" onClick={() => onRetry(message.id, message.contentMarkdown)}>
              Retry
            </button>
          </div>
        </article>
      ))}
    </div>
  );
}
```

- [ ] **Step 2: Create MessageList.test.tsx**

Write `frontend/src/features/chat/MessageList.test.tsx`:

```tsx
import { fireEvent, render, screen } from "@testing-library/react";
import { describe, expect, it, vi } from "vitest";
import { MessageList } from "./MessageList";

describe("MessageList", () => {
  const messages = [
    {
      id: 1,
      senderId: 1,
      contentMarkdown: "Hello",
      createdAt: "2026-07-10T10:00:00Z",
    },
    {
      id: 2,
      senderId: 2,
      contentMarkdown: "Hi there",
      createdAt: "2026-07-10T10:01:00Z",
    },
  ];

  it("renders messages", () => {
    render(
      <MessageList
        currentUserId={1}
        messages={messages}
        failedMessages={[]}
        onRetry={() => {}}
      />,
    );
    expect(screen.getByText("Hello")).toBeInTheDocument();
    expect(screen.getByText("Hi there")).toBeInTheDocument();
  });

  it("renders failed messages with retry button", () => {
    const onRetry = vi.fn();
    render(
      <MessageList
        currentUserId={1}
        messages={[]}
        failedMessages={[{ id: "failed-1", contentMarkdown: "Oops" }]}
        onRetry={onRetry}
      />,
    );
    expect(screen.getByText("Oops")).toBeInTheDocument();
    expect(screen.getByText("Failed to send")).toBeInTheDocument();
    fireEvent.click(screen.getByRole("button", { name: "Retry" }));
    expect(onRetry).toHaveBeenCalledWith("failed-1", "Oops");
  });
});
```

- [ ] **Step 3: Run tests to verify they pass**

Run:
```bash
cd frontend && npm run test -- src/features/chat/MessageList.test.tsx
```

Expected: PASS.

- [ ] **Step 4: Commit**

```bash
git add frontend/src/features/chat/MessageList.tsx frontend/src/features/chat/MessageList.test.tsx
git commit -m "feat: restyle message list with metadata and failed state"
```

---

## Task 6: Refactor MessageComposer and Update Tests

**Files:**
- Modify: `frontend/src/features/chat/MessageComposer.tsx`
- Modify: `frontend/src/features/chat/MessageComposer.test.tsx`

- [ ] **Step 1: Update MessageComposer.tsx with auto-resize and keyboard handling**

Modify `frontend/src/features/chat/MessageComposer.tsx`:

```tsx
import { useEffect, useRef, useState } from "react";

type Props = {
  disabled: boolean;
  onSend: (text: string) => void;
};

export function MessageComposer({ disabled, onSend }: Props) {
  const [value, setValue] = useState("");
  const textareaRef = useRef<HTMLTextAreaElement>(null);

  useEffect(() => {
    const el = textareaRef.current;
    if (!el) return;
    el.style.height = "auto";
    el.style.height = `${Math.min(el.scrollHeight, 180)}px`;
  }, [value]);

  function submit() {
    const trimmed = value.trim();
    if (!trimmed || disabled) {
      return;
    }
    onSend(trimmed);
    setValue("");
    const el = textareaRef.current;
    if (el) {
      el.style.height = "auto";
    }
  }

  return (
    <form
      className="composer"
      onSubmit={(event) => {
        event.preventDefault();
        submit();
      }}
    >
      <textarea
        ref={textareaRef}
        aria-label="Message"
        rows={1}
        value={value}
        disabled={disabled}
        onChange={(event) => setValue(event.target.value)}
        onKeyDown={(event) => {
          if (event.key === "Enter" && !event.shiftKey) {
            event.preventDefault();
            submit();
          }
        }}
        placeholder="Type a message..."
      />
      <button className="primary-button" type="submit" disabled={disabled || !value.trim()}>
        {disabled ? <span className="spinner" aria-hidden="true" /> : "Send"}
      </button>
    </form>
  );
}
```

- [ ] **Step 2: Update MessageComposer.test.tsx**

Replace `frontend/src/features/chat/MessageComposer.test.tsx` with:

```tsx
import { fireEvent, render, screen } from "@testing-library/react";
import { describe, expect, it, vi } from "vitest";
import { MessageComposer } from "./MessageComposer";

describe("MessageComposer", () => {
  it("sends message on submit", () => {
    const onSend = vi.fn();
    render(<MessageComposer disabled={false} onSend={onSend} />);
    const input = screen.getByLabelText("Message");
    fireEvent.change(input, { target: { value: "Hello" } });
    fireEvent.click(screen.getByRole("button", { name: "Send" }));
    expect(onSend).toHaveBeenCalledWith("Hello");
  });

  it("sends message on Enter", () => {
    const onSend = vi.fn();
    render(<MessageComposer disabled={false} onSend={onSend} />);
    const input = screen.getByLabelText("Message");
    fireEvent.change(input, { target: { value: "Hello" } });
    fireEvent.keyDown(input, { key: "Enter", code: "Enter", shiftKey: false });
    expect(onSend).toHaveBeenCalledWith("Hello");
  });

  it("does not send on Shift+Enter", () => {
    const onSend = vi.fn();
    render(<MessageComposer disabled={false} onSend={onSend} />);
    const input = screen.getByLabelText("Message");
    fireEvent.change(input, { target: { value: "Line 1" } });
    fireEvent.keyDown(input, { key: "Enter", code: "Enter", shiftKey: true });
    expect(onSend).not.toHaveBeenCalled();
  });

  it("does not send empty messages", () => {
    const onSend = vi.fn();
    render(<MessageComposer disabled={false} onSend={onSend} />);
    fireEvent.click(screen.getByRole("button", { name: "Send" }));
    expect(onSend).not.toHaveBeenCalled();
  });

  it("disables send while disabled", () => {
    render(<MessageComposer disabled onSend={() => {}} />);
    expect(screen.getByRole("button")).toBeDisabled();
    expect(screen.getByLabelText("Message")).toBeDisabled();
  });
});
```

- [ ] **Step 3: Run tests to verify they pass**

Run:
```bash
cd frontend && npm run test -- src/features/chat/MessageComposer.test.tsx
```

Expected: PASS.

- [ ] **Step 4: Commit**

```bash
git add frontend/src/features/chat/MessageComposer.tsx frontend/src/features/chat/MessageComposer.test.tsx
git commit -m "feat: auto-resize composer, keyboard send, and update tests"
```

---

## Task 7: Refactor ChatPage for States, Errors, and Mobile

**Files:**
- Modify: `frontend/src/features/chat/ChatPage.tsx`

- [ ] **Step 1: Add mobile view state and improve error/empty/loading UI**

Modify `frontend/src/features/chat/ChatPage.tsx`:

```tsx
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useEffect, useState } from "react";
import {
  listConversations,
  listMessages,
  sendMessage,
  type Conversation,
  type Message,
  type User,
} from "../../lib/api";
import { connectEvents } from "../../lib/ws";
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
  const [connected, setConnected] = useState(false);
  const [failedMessages, setFailedMessages] = useState<{ id: string; contentMarkdown: string }[]>([]);
  const [mobileShowChat, setMobileShowChat] = useState(false);

  const conversations = useQuery({
    queryKey: ["conversations"],
    queryFn: () => listConversations(token),
  });

  const activeId = selectedId ?? conversations.data?.conversations[0]?.id ?? null;

  const messages = useQuery({
    queryKey: ["messages", activeId],
    queryFn: () => listMessages(token, activeId!),
    enabled: activeId !== null,
  });

  const send = useMutation({
    mutationFn: (text: string) => sendMessage(token, activeId!, text),
    onSuccess: (message: Message) => {
      queryClient.setQueryData<{ messages: Message[] }>(["messages", activeId], (old) => ({
        messages: [message, ...(old?.messages ?? [])],
      }));
      void queryClient.invalidateQueries({ queryKey: ["conversations"] });
    },
    onError: (_error, text) => {
      setFailedMessages((items) => [...items, { id: crypto.randomUUID(), contentMarkdown: text }]);
    },
  });

  const activeConversation = (conversations.data?.conversations ?? []).find(
    (item: Conversation) => item.id === activeId,
  );

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
          if (activeId) {
            void queryClient.invalidateQueries({ queryKey: ["messages", activeId] });
          }
        }
      },
    );
  }, [token, queryClient, activeId]);

  function retryFailedMessage(id: string, contentMarkdown: string) {
    setFailedMessages((items) => items.filter((item) => item.id !== id));
    send.mutate(contentMarkdown);
  }

  function handleSelect(id: number) {
    setSelectedId(id);
    setMobileShowChat(true);
  }

  function handleBackToList() {
    setMobileShowChat(false);
  }

  const showListClass = mobileShowChat ? "mobile-hidden" : "";
  const showChatClass = mobileShowChat ? "mobile-visible" : "";

  return (
    <main className="chat-layout">
      <ConversationList
        className={showListClass}
        conversations={conversations.data?.conversations ?? []}
        selectedId={activeId}
        onSelect={handleSelect}
      />
      <section className={`chat-panel ${showChatClass}`.trim()}>
        {activeId ? (
          <>
            <header className="chat-header">
              <div className="chat-header-start">
                <button
                  className="back-button"
                  type="button"
                  onClick={handleBackToList}
                  aria-label="Back to conversations"
                >
                  ←
                </button>
                <span className="chat-header-title">
                  {activeConversation?.title ?? `Conversation ${activeId}`}
                </span>
              </div>
              <span className={connected ? "ws-status connected" : "ws-status"}>
                <span className="ws-status-dot" aria-hidden="true" />
                {connected ? "Live" : "Reconnecting"}
              </span>
            </header>
            {messages.isLoading ? (
              <div className="loading-state" aria-label="Loading messages">
                <div className="skeleton" style={{ height: 64 }} />
                <div className="skeleton" style={{ height: 64, width: "80%" }} />
                <div className="skeleton" style={{ height: 64, width: "60%" }} />
              </div>
            ) : messages.error ? (
              <div className="error-state">
                <p className="error-state-message">Could not load messages</p>
                <button
                  className="primary-button"
                  type="button"
                  onClick={() => messages.refetch()}
                >
                  Retry
                </button>
              </div>
            ) : (
              <MessageList
                currentUserId={user.id}
                failedMessages={failedMessages}
                messages={messages.data?.messages ?? []}
                onRetry={retryFailedMessage}
              />
            )}
            <MessageComposer disabled={send.isPending} onSend={(text) => send.mutate(text)} />
          </>
        ) : (
          <div className="empty-state">
            <div>
              <p className="empty-state-title">Select a conversation</p>
              <p className="empty-state-hint">Choose a chat from the list to start messaging</p>
            </div>
          </div>
        )}
      </section>
    </main>
  );
}
```

The `className` prop on `ConversationList` was already added in Task 4; ChatPage now passes `mobile-hidden` to hide the list when the chat view is active on mobile.

- [ ] **Step 2: Verify lint passes**

Run:
```bash
cd frontend && npm run lint
```

Expected: PASS.

- [ ] **Step 3: Commit**

```bash
git add frontend/src/features/chat/ChatPage.tsx
git commit -m "feat: chat page states, errors, and mobile navigation"
```

---

## Task 8: Refactor AdminPage

**Files:**
- Modify: `frontend/src/features/admin/AdminPage.tsx`

- [ ] **Step 1: Update AdminPage.tsx with card/table styling and loading state**

Modify `frontend/src/features/admin/AdminPage.tsx`:

```tsx
import { useQuery } from "@tanstack/react-query";
import { adminConversations, adminMessages, adminUsers } from "../../lib/api";

type Props = {
  token: string;
};

export function AdminPage({ token }: Props) {
  const users = useQuery({ queryKey: ["admin", "users"], queryFn: () => adminUsers(token) });
  const conversations = useQuery({
    queryKey: ["admin", "conversations"],
    queryFn: () => adminConversations(token),
  });
  const messages = useQuery({ queryKey: ["admin", "messages"], queryFn: () => adminMessages(token) });

  return (
    <main className="admin-page">
      <div className="admin-page-inner">
        <h1>Admin</h1>
        <AdminTable title="Users" rows={users.data?.users ?? []} isLoading={users.isLoading} />
        <AdminTable
          title="Conversations"
          rows={conversations.data?.conversations ?? []}
          isLoading={conversations.isLoading}
        />
        <AdminTable
          title="Recent Messages"
          rows={messages.data?.messages ?? []}
          isLoading={messages.isLoading}
        />
      </div>
    </main>
  );
}

function AdminTable({
  title,
  rows,
  isLoading,
}: {
  title: string;
  rows: Record<string, unknown>[];
  isLoading: boolean;
}) {
  const columns = Object.keys(rows[0] ?? {});

  return (
    <section className="admin-section">
      <h2>{title}</h2>
      <div className="admin-table-wrap">
        {isLoading ? (
          <div className="loading-state" aria-label={`Loading ${title.toLowerCase()}`}>
            <div className="skeleton" style={{ height: 40 }} />
            <div className="skeleton" style={{ height: 40 }} />
            <div className="skeleton" style={{ height: 40 }} />
          </div>
        ) : rows.length === 0 ? (
          <div className="empty-state">
            <p className="empty-state-hint">No {title.toLowerCase()} found</p>
          </div>
        ) : (
          <table>
            <thead>
              <tr>
                {columns.map((column) => (
                  <th key={column}>{column}</th>
                ))}
              </tr>
            </thead>
            <tbody>
              {rows.map((row, index) => (
                <tr key={index}>
                  {columns.map((column) => (
                    <td
                      key={column}
                      className={typeof row[column] === "number" ? "tabular-nums" : ""}
                    >
                      {String(row[column] ?? "")}
                    </td>
                  ))}
                </tr>
              ))}
            </tbody>
          </table>
        )}
      </div>
    </section>
  );
}
```

- [ ] **Step 2: Verify lint passes**

Run:
```bash
cd frontend && npm run lint
```

Expected: PASS.

- [ ] **Step 3: Commit**

```bash
git add frontend/src/features/admin/AdminPage.tsx
git commit -m "feat: restyle admin page with cards and loading states"
```

---

## Task 9: Final Verification

**Files:**
- All modified frontend files.

- [ ] **Step 1: Run all frontend tests**

Run:
```bash
cd frontend && npm run test
```

Expected: All tests pass.

- [ ] **Step 2: Run TypeScript lint**

Run:
```bash
cd frontend && npm run lint
```

Expected: PASS.

- [ ] **Step 3: Run build**

Run:
```bash
cd frontend && npm run build
```

Expected: Build succeeds with no errors.

- [ ] **Step 4: Manual verification checklist**

Start the backend and frontend (MySQL must be running):

```bash
docker compose up -d mysql
cd backend && go run ./cmd/migrate
cd backend && ENABLE_DEV_LOGIN=true go run ./cmd/api &
cd frontend && npm run dev
```

Verify in browser:

- Login page shows seed users with avatars and spinner on click.
- Login form validates empty username.
- Chat layout loads on desktop.
- Conversation list shows unread badge.
- Selecting a conversation loads messages.
- Sending a message works; Enter sends, Shift+Enter inserts newline.
- Failed message retry works (simulate by stopping backend briefly).
- Connection status pill shows "Reconnecting" when backend stops and "Live" when it returns.
- Mobile view (< 768px) shows list first, tapping conversation shows chat, back button returns to list.
- Admin page tables render with styling and hover states.

- [ ] **Step 5: Final commit**

```bash
git add frontend/
git commit -m "feat: polish frontend visual design and interaction states"
```

---

## Self-Review

### Spec Coverage

| Spec Requirement | Task |
|------------------|------|
| Geist font | Task 1 |
| Design tokens in `:root` | Task 2 |
| `100dvh` viewport | Task 2 |
| Focus rings | Task 2 |
| Reduced motion | Task 2 |
| Login loading/error states | Task 3 |
| ChatPage loading/error/empty states | Task 7 |
| ConversationList unread badge | Task 4 |
| MessageList metadata/failed state/animation | Task 5 |
| MessageComposer auto-resize/keyboard/spinner | Task 6 |
| AdminPage card/table styling | Task 8 |
| Mobile responsive behavior | Task 2 + Task 7 |

### Placeholder Scan

No TBD, TODO, or vague steps found. Each step includes exact file paths, code, commands, and expected output.

### Type Consistency

- `ConversationList` gains optional `className` prop in Task 4/7.
- `AdminTable` gains `isLoading` prop in Task 8.
- Message and User types come from existing `lib/api`.

No contradictions found.
