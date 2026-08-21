# Profile settings and OSS avatar implementation plan

> Execute from `/Users/wilenwang/works/CodeZone`.

## 1. Add OSS configuration and dependency

1. Add the `codeup.aliyun.com/61e54b0e0bb300d827e1ae27/backend/golib` module dependency used by development-line and verify it resolves through the existing global Go configuration.
2. Extend `backend/internal/config.Config` with nested OSS configuration: endpoint, access key ID, access key secret, bucket name, public URL, and path prefix.
3. Add credential-free OSS placeholders to all versioned CodeZone YAML templates. Use `dev/codezone`, `test/codezone`, `gray/codezone`, and `prod/codezone` as their environment-specific prefixes; document that actual credential values are injected locally or by deployment secrets.
4. Add an OSS initialization step at API startup. Missing OSS configuration does not stop API startup; profile avatar upload returns a clear server error instead.

## 2. Add storage and avatar backend support

1. Add `backend/internal/storage` with a narrow public-object storage interface and an Alibaba OSS implementation wrapping the shared private module.
2. Normalize the configured public URL and construct avatar keys as `{path_prefix}/avatars/{user_id}/{uuid}.{ext}`.
3. Add a migration that changes `users.avatar_url` to `TEXT NULL`, using information-schema conditional SQL so repeated migrations remain safe.
4. Extend the users repository/service with authenticated-self profile update operations for display name and avatar URL.
5. Add `internal/profile` handler/service or focused users handler methods for `PATCH /api/me` and `POST /api/me/avatar`.
6. For avatar upload, cap multipart input at 5 MB, require JPEG/PNG/WebP declared MIME types, verify magic bytes, derive a trusted extension, generate a UUID filename, upload to OSS, then update the user's stored URL only after upload succeeds.
7. Register both endpoints inside the existing protected router group.
8. Add unit and handler tests for authentication, display-name validation, image validation, key construction, storage failure, and successful response payloads.

## 3. Add frontend profile state and APIs

1. Extend the frontend `User` type with nullable `avatarUrl`, and add API methods for updating display name and uploading an avatar with `FormData` (without setting an explicit JSON content type).
2. Build a reusable `Avatar` component that renders a real image when available and a deterministic username-hash color plus display-name initial otherwise; fall back to initials after image failure.
3. Build a `ProfileDrawer` component with local display-name draft, read-only username, file chooser, object-URL preview, size/type client checks, independent upload/save loading states, inline errors, and object-URL cleanup.

## 4. Connect profile UI

1. Change the chat sidebar current-user block into an accessible button that opens the right-side drawer.
2. Render `Avatar` in the current-user block, direct user list, group-member selector, and login seed buttons.
3. On a successful profile response, replace the user in `App` state through a callback, update/invalidate the users query, and close the relevant completed drawer action. Preserve drafts and preview when requests fail.
4. Add silver-glass drawer and avatar styles consistent with the existing Apple material design; retain responsive behavior and reduced-motion support.

## 5. Verify

1. Run `cd backend && gofmt -w ... && go mod tidy && go test ./...`.
2. Run `cd frontend && npm run test && npm run lint && npm run build`.
3. With OSS credentials present, run `make run`, update a display name, upload JPEG/PNG/WebP files under 5 MB, confirm `avatar_url` stores a public OSS URL under the correct environment prefix, and confirm invalid files are rejected.
