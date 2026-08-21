# Direct and group chat implementation plan

> Execute from `/Users/wilenwang/works/CodeZone`.

## 1. Add direct-chat persistence support

1. Create a migration that adds a nullable canonical direct-pair key to `conversations` and a unique key over workspace plus that pair key. The value is only populated for direct conversations and is the two user IDs in ascending order, joined by a delimiter.
2. Extend `CreateConversationInput` with the internal direct-pair key and make `SQLRepository.Create` persist it.
3. Add an SQL repository method that starts a transaction, locks the two target user rows in ascending ID order, finds a direct conversation by workspace and canonical pair key, and returns it if present. Otherwise it creates the conversation plus exactly two membership rows, commits, and returns it. If a duplicate-key race occurs, re-query the existing row and return it.
4. Add a service-level `EnsureDirect` method. Validate positive and distinct IDs, verify both members are in the workspace, calculate the canonical pair key, then delegate to the repository.
5. Add focused unit tests for validation, workspace membership checks, reuse, and create delegation. Add repository integration coverage for lookup/create and duplicate prevention if the existing integration test harness supports it.

## 2. Expose the direct-chat endpoint

1. Add `Handler.EnsureDirect`, decoding `{ "userId": number }`, using the authenticated caller as the source user, and calling `EnsureDirect` for workspace 1.
2. Return `400` for malformed input, self-targeting, and validation/membership failures; return the conversation with `200` for both reused and newly created direct chats.
3. Register `POST /api/conversations/direct` inside the protected router group.
4. Add handler tests for unauthorized context, malformed input, successful delegation, and invalid target responses.

## 3. Add frontend client operations

1. Extend `frontend/src/lib/api.ts` with `listUsers`, `ensureDirectConversation`, and `createConversation` methods and their request/response types.
2. Add API-client tests verifying authorization headers, request paths, and JSON bodies.

## 4. Build sidebar user directory and group dialog

1. Refactor the sidebar component to accept workspace users, an active direct-user ID, direct-start loading/error state, callbacks, and a group-dialog trigger.
2. Render `群聊` using only group conversations. Render `单聊` from workspace users excluding the signed-in user; each row displays display name and `@username` and selects/ensures its direct conversation.
3. Add a `CreateGroupDialog` component with a group-name input, checkbox member list excluding the creator, at-least-two-member client validation, loading state, close action, and retained form/error state after failed submissions.
4. Add matching styles for the sidebar actions, user rows, direct selection/loading state, and accessible modal overlay; preserve the existing mobile sidebar behavior.

## 5. Connect ChatPage data and actions

1. Query `listUsers(token)` beside conversations and pass the directory to the sidebar.
2. On direct-user click, call the ensure-direct mutation; on success invalidate conversations, select the returned ID, and set mobile chat visibility. On failure show a sidebar retry message without changing the active conversation.
3. On group submission, call `createConversation` with `group`, name, and selected IDs. On success close/reset dialog, invalidate conversations, select the returned conversation, and open the mobile chat panel.
4. Keep message loading keyed by the selected conversation ID, so entering a newly created or reused chat behaves like selecting an existing group.

## 6. Verify

1. Run `cd backend && gofmt -w ... && go test ./...`.
2. Run `cd frontend && npm run test && npm run lint && npm run build`.
3. Start local services with `make run`; as Carol, verify the group section, user directory, direct chat with Alice/agent, group creation with at least two users, and validation/error behavior.
