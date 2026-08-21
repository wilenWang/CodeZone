# IME Enter composition implementation plan

1. In `frontend/src/features/chat/MessageComposer.tsx`, add a ref-backed IME composition flag so composition lifecycle changes do not require an intermediate render.
2. Set that flag on `onCompositionStart` and clear it on `onCompositionEnd`.
3. Update the textarea Enter handler to return without preventing default or sending when the ref is active, `event.nativeEvent.isComposing` is true, or `event.keyCode === 229`. Retain Shift+Enter newline behavior and ordinary Enter submission after composition ends.
4. Extend `MessageComposer.test.tsx` with coverage for Enter during composition not sending and Enter after composition ends sending.
5. Run `cd frontend && npm run test && npm run lint && npm run build`.
