# Chat current-user display design

## Goal

Make the logged-in identity visible in the upper-left of the chat page so users can immediately tell which account is active.

## Scope

The chat page sidebar will render a small user-identity block above the conversation list heading. It contains a `当前用户` label, the existing `User.displayName` as the primary line, and `@` plus `User.username` as supporting text.

The component uses the `user` prop already passed to `ChatPage`; it adds no API calls, state, persistence, or backend changes.

The block does not render on the admin route, login page, or elsewhere.

## Styling

Styles will use the existing surface, border, secondary text, spacing, and type tokens. The identity block is visually separated from the conversation list with a bottom border while preserving the existing sidebar scrolling and mobile behavior.

## Testing

Extend `ChatPage.test.tsx` to assert that the supplied display name and `@username` are rendered.
