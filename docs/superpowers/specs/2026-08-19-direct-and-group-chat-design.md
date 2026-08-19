# Direct and group chat design

## Goal

Let a user find every other workspace member and start a one-to-one conversation by clicking that member. Let a user explicitly create named group chats by selecting members.

## Sidebar

The chat sidebar contains, in order: the existing current-user identity block, a `新建群聊` button, a `群聊` section, and a `单聊` section.

`群聊` contains only existing conversations with `type = group`. `单聊` is a directory of all workspace users other than the logged-in user, including agents. Each row shows display name and `@username`. Existing direct conversations are reached from their matching user row rather than being independently listed.

## Direct conversations

Clicking a user calls a new authenticated `POST /api/conversations/direct` endpoint with that user's ID. The server validates that the target exists in the same workspace and is not the caller. It returns the unique existing direct conversation between exactly those two users or creates it transactionally when missing.

The frontend invalidates the conversations query, selects the returned conversation, and opens the chat panel on mobile. While the request is in progress, the selected user row is disabled. A failure leaves the current conversation selected and shows a retryable error in the sidebar.

A database-level pair identity or locking strategy prevents duplicate direct conversations for concurrent requests by the same pair.

## Group conversations

`新建群聊` opens a modal. The creator supplies a non-empty name and chooses at least two other users. The creator is automatically included. Submitting uses the existing conversation creation endpoint with type `group`. On success, the modal closes, conversations refresh, and the returned group is selected. On failure, the form's name and selected members remain intact and an error is shown.

## API and components

The frontend adds list-users, ensure-direct-conversation, and create-conversation client methods. The sidebar component receives the user directory and separates group conversations from direct-user rows. A focused group-dialog component owns its form state.

The backend adds a direct-conversation service and repository operation, route, request validation, and repository integration coverage. The existing group creation route remains the group-chat backend contract.

## Error handling and testing

Backend tests cover reusing a direct chat, creating one, and rejecting self, nonexistent, or out-of-workspace targets. Tests also cover duplicate prevention under concurrent requests where the integration setup supports it.

Frontend tests cover filtering the current user from the direct directory, starting/reusing a direct chat from a user click, sidebar section filtering, group member minimum validation, successful group creation and failure preservation. Existing chat, message, and mobile navigation tests remain passing.
