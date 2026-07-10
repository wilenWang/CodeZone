# Frontend Visual & UX Optimization — Design Spec

Date: 2026-07-10

## Summary

Upgrade the existing React + Vite chat frontend from a generic default look to a polished, product-quality interface without migrating styling frameworks or adding heavy animation libraries.

The work is scoped to the `frontend/` directory. It preserves all existing functionality (login, chat, admin, WebSocket recovery) and focuses on:

- A unified design-token system.
- Geist typography.
- Zinc + single Teal accent palette.
- Full UI states: loading, empty, error, hover, active, focus.
- Mobile-safe viewport handling and responsive improvements.
- Subtle CSS-driven micro-interactions.

## Goals

- Make the app feel intentional and finished.
- Keep the existing vanilla CSS + React stack; no Tailwind or Framer Motion migration.
- Add low-risk font packages: `@fontsource-variable/geist` and `@fontsource-variable/geist-mono`.
- Improve accessibility with visible focus rings and semantic markup.
- Ensure all interactive elements have clear feedback cycles.

## Non-Goals

- No backend changes.
- No new features such as attachments, presence, typing indicators, or dark mode.
- No framework or build-tool migration.
- No heavy scroll-driven animations or parallax.

## Design Tokens

All tokens live in CSS custom properties on `:root` in `frontend/src/styles.css`.

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
```

Rules:

- Exactly one accent color (Teal). No purple/blue gradients.
- All grays come from the Zinc family; no warm/cool gray mixing.
- No pure black (`#000000`). Darkest text is Zinc-950.

## Typography

- Primary typeface: **Geist** for all UI text.
- Monospace: **Geist Mono** for timestamps, IDs, and admin numeric data.
- Display headings: `font-weight: 600`, `letter-spacing: -0.025em`, `line-height: 1.1`.
- Body: `font-weight: 400`, `line-height: 1.5`.
- Labels and buttons: `font-weight: 500`.
- Small caps/labels: `font-size: 12px`, `font-weight: 500`, `letter-spacing: 0.01em`, `text-transform: uppercase` where semantically appropriate.

## Global Foundations

### Viewport

Replace every `min-height: 100vh` with `min-height: 100dvh` to prevent mobile Safari viewport jumps.

### Focus

All focusable elements use `:focus-visible` with a 2px solid accent outline and a 2px offset.

### Transitions

Interactive elements use:

```css
transition: color 0.2s var(--ease-out),
            background-color 0.2s var(--ease-out),
            border-color 0.2s var(--ease-out),
            transform 0.2s var(--ease-out),
            box-shadow 0.2s var(--ease-out);
```

### Reduced Motion

```css
@media (prefers-reduced-motion: reduce) {
  *, *::before, *::after {
    animation-duration: 0.01ms !important;
    animation-iteration-count: 1 !important;
    transition-duration: 0.01ms !important;
  }
}
```

## Component Changes

### `index.html`

- Preload Geist font files to minimize flash of unstyled text.
- Add `<meta name="description">` and `<meta name="theme-color">`.

### `LoginPage.tsx`

- Page layout: full `dvh` centered card with generous padding.
- Card: white surface, `radius-lg`, subtle shadow, `max-width: 420px`.
- Title: `text-2xl` (1.5rem), `font-weight: 600`, `tracking-tight`.
- Seed-user buttons: stacked list, each with a neutral avatar circle (first initial), username, and hover lift.
- Loading state: clicked seed button shows an inline spinner and disables itself.
- Manual login form:
  - Visible `<label>` above input.
  - Input focus ring.
  - Error state: red border + inline error text below input.
- Error banner: `error-bg` + `error-border` block below the form.

### `ChatPage.tsx`

- Layout: `min-height: 100dvh`, CSS Grid `320px 1fr` on desktop, single column on mobile.
- Header:
  - Left: conversation title, `font-weight: 600`.
  - Right: connection status pill.
    - Connected: Teal dot + "Live".
    - Reconnecting: pulsing amber dot + "Reconnecting".
- Loading: skeleton placeholders for conversation list and message list.
- Empty state (no selected conversation): centered illustration-free message with muted text and a small helper line.
- Error state: inline banner with retry button when conversation/message queries fail.

### `ConversationList.tsx`

- List title "Chats": `font-weight: 600`, muted color, small size.
- Conversation row:
  - Padding `10px 12px`, `radius-md`.
  - Hover: `background: #f4f4f5`.
  - Active press: `transform: scale(0.99)`.
  - Selected: `background: var(--color-accent-subtle)`, `color: var(--color-accent-text)`, 3px accent left border.
- Unread count: compact circular badge (`min-width: 18px`, `height: 18px`, `radius-full`), accent background, white text.
- Show conversation title only in MVP; last-message preview is out of scope unless data is already available.

### `MessageList.tsx`

- Message bubble:
  - Mine: aligned right, `radius-lg` with bottom-right tighter (`radius-md`), `accent-subtle` background, `accent-text` color.
  - Others: aligned left, `radius-lg` with bottom-left tighter, white background, thin border.
- Metadata line below each bubble: sender display name + timestamp, `text-secondary`, `font-size: 12px`.
- Failed message: `error-bg`, `error-border`, and a red action bar with "Failed to send" + "Retry" button.
- New-message entrance animation: CSS keyframes `message-in` (`opacity 0 → 1`, `translateY(8px) → 0`) over 0.25s.
- Code blocks in Markdown: light gray background, `radius-sm`, monospace font.

### `MessageComposer.tsx`

- Container: white background, top border, padding `12px`.
- Textarea:
  - Auto-growing: starts at one row, expands up to ~6 rows.
  - `radius-md`, subtle border, focus ring.
  - `Enter` sends; `Shift+Enter` inserts newline.
- Send button:
  - `radius-md`, accent background.
  - Hover: `translateY(-1px)`.
  - Active: `scale(0.98)`.
  - Disabled: shows a small spinner instead of "Send".

### `AdminPage.tsx`

- Container padding `24px`, max-width `1200px`, centered.
- Section cards: white surface, `radius-lg`, subtle shadow, `overflow: hidden`.
- Tables:
  - Header: `Zinc-50` background, uppercase small labels.
  - Rows: white with `border-bottom`, hover `Zinc-50`.
  - Numeric cells use `font-mono` and `tabular-nums`.
- Loading: skeleton row placeholders.
- Empty table: centered muted text.

## States

| State | Visual treatment |
|-------|------------------|
| Loading | Skeleton blocks matching layout shape; no generic spinners. |
| Empty | Composed centered text with a muted helper line. |
| Error | Inline banner or field-level message with retry where possible. |
| Hover | Background/color shift + cursor pointer. |
| Active/Pressed | `scale(0.98)` or `translateY(1px)`. |
| Focus | Visible accent outline via `:focus-visible`. |
| Disabled | Reduced opacity + cursor not-allowed; buttons show spinner. |

## Responsive Behavior

- **Desktop (≥ 768px):** two-column grid, conversation list always visible.
- **Mobile (< 768px):**
  - Single-column view.
  - Initially show conversation list full screen.
  - Tapping a conversation switches to the chat view.
  - Chat header gains a back button to return to the list.
  - No awkward 42vh split; views are mutually exclusive.

## Accessibility

- Semantic HTML: `<main>`, `<aside>`, `<section>`, `<header>`, `<article>`, `<nav>` where appropriate.
- All form inputs have associated `<label>` or `aria-label`.
- Focus rings are always visible.
- Buttons have explicit `type`.
- Reduced-motion preference is respected.

## Testing

- `npm run lint` must pass after all changes.
- Update `MessageComposer.test.tsx` for new loading/disabled behavior and Enter/Shift+Enter handling.
- Add or update tests for:
  - `ConversationList` unread badge rendering.
  - `MessageList` failed message retry action.
  - `LoginPage` error state rendering.
- Manual checks:
  - Login with seed users.
  - Send/receive messages.
  - Retry a failed send.
  - Switch conversations.
  - Verify layout on mobile viewport in DevTools.

## Dependencies

Add:

```bash
cd frontend && npm install @fontsource-variable/geist @fontsource-variable/geist-mono
```

No other new runtime dependencies. These are font-only packages and do not change the build pipeline.

## Files Likely to Change

- `frontend/src/styles.css` — design tokens + component styles.
- `frontend/index.html` — font preload + meta tags.
- `frontend/src/main.tsx` — import Geist font CSS.
- `frontend/src/features/auth/LoginPage.tsx`.
- `frontend/src/features/chat/ChatPage.tsx`.
- `frontend/src/features/chat/ConversationList.tsx`.
- `frontend/src/features/chat/MessageList.tsx`.
- `frontend/src/features/chat/MessageComposer.tsx`.
- `frontend/src/features/chat/MarkdownMessage.tsx` — code block styling.
- `frontend/src/features/admin/AdminPage.tsx`.
- `frontend/src/features/chat/MessageComposer.test.tsx` and related tests.

## Out of Scope

- Dark mode.
- New chat creation UI.
- Conversation search.
- User avatars from real URLs (use initials fallback).
- File attachments or rich content beyond Markdown.
