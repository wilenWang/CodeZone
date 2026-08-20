# Profile settings and OSS avatar design

## Goal

Add a right-side personal settings drawer where authenticated users can edit their display name and upload an avatar. Store avatars in the same Alibaba OSS infrastructure and configuration convention used by `development-line`.

## User model and defaults

`users.avatar_url` changes from `VARCHAR(500)` to `TEXT NULL`. An absent URL remains the normal state for newly registered or existing users. The frontend renders a deterministic avatar fallback from the username hash and display-name initial whenever the URL is absent or an image fails to load.

The frontend `User` contract adds `avatarUrl: string | null`, matching the existing backend user structs and login/list responses.

## Personal settings

The current-user block in the chat sidebar becomes an accessible button. It opens a right-side drawer without leaving the active chat. The drawer has an avatar picker and preview, editable non-empty display name, and read-only `@username`.

Display-name changes use `PATCH /api/me`. Avatar changes use `POST /api/me/avatar` with multipart field `file`. Each response returns the updated authenticated user. On success the client updates current user state and invalidates/updates user-directory data so the sidebar, direct-chat list, group-member selector, and login seed list reflect the change. A failure leaves draft input and preview intact and appears inline.

## OSS storage

CodeZone adds an `internal/storage` interface with upload and public-URL methods, implemented through the private `codeup.aliyun.com/61e54b0e0bb300d827e1ae27/backend/golib/oss` module used by `development-line`.

YAML config gains the matching OSS fields: `end_point`, `access_key_ID`, `access_key_secret`, `bucket_name`, `oss_url`, and `path_prefix`. Credentials remain environment-injected placeholders in versioned templates. Environment prefixes are `dev/codezone`, `test/codezone`, `gray/codezone`, and `prod/codezone`.

Avatar object keys are `{path_prefix}/avatars/{user_id}/{uuid}.{ext}`. The users table saves only the public OSS URL. Database work and OSS upload remain outside the same transaction; the profile URL is updated only after an upload succeeds. Replaced objects are not deleted in this MVP.

## Validation and security

Only the authenticated user can update their own profile. Avatar upload accepts JPEG, PNG, and WebP only, with a 5 MB maximum. The backend validates declared MIME type and image magic bytes, generates the filename with a random UUID, and never accepts client-controlled paths.

## Testing

Backend tests cover authorization, display-name validation, image MIME/size/signature checks, object-key generation, storage failure, and user update responses. Frontend tests cover default and failed-image avatar fallbacks, image selection preview, API request construction, drawer validation, and successful/failed profile updates.
