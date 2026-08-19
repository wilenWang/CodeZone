# IME Enter composition design

## Goal

Prevent Enter from sending a chat message while an input method editor (IME) is composing text or presenting candidates.

## Behavior

`MessageComposer` will track composition lifecycle events from its textarea. While composition is active, Enter remains owned by the IME so it can confirm or dismiss candidates; the component does not prevent the default event or submit the form.

The key handler also checks the native composing state and `keyCode === 229` as compatibility guards for browser event ordering differences. After composition ends, unmodified Enter sends the trimmed message as it does today. Shift+Enter continues to insert a newline.

## Scope and testing

The change is isolated to the message composer and does not change API calls, message persistence, or button submission. Add tests that Enter does not call `onSend` during composition and does call it after composition ends. Existing Enter, Shift+Enter, disabled, and button-submit behavior remains covered.
