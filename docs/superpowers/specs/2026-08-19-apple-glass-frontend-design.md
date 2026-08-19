# Apple glass frontend redesign

## Goal

Rework the existing frontend into an Apple-inspired frosted desktop experience while preserving every existing chat, direct-chat, group-chat, authentication, and admin behavior.

## Chosen visual direction

The interface uses a cold silver-gray, edge-to-edge material system. The chat workspace fills the viewport rather than floating in a centered application window. Internal glass layers and fine separators establish hierarchy between the sidebar, chat header, message area, composer, and dialogs.

A single low-saturation teal remains the functional accent for selected rows, primary actions, online status, focus treatment, and the current user's message bubbles. All remaining colors are neutral silver, graphite, translucent white, and accessible muted text.

## Material system

The application background is a fixed, layered silver-gray radial field that provides depth without being attached to scroll containers. Panels use translucent white fills, `backdrop-filter` blur where available, a semitransparent white edge, an inset top highlight, and a restrained cool-gray diffusion shadow. CSS fallback colors preserve legibility on browsers without backdrop-filter support.

Glass opacity varies by hierarchy: the sidebar is more opaque for navigation readability; the header and composer are moderately translucent; messages use independent, low-elevation glass bubbles; and the group dialog gets the strongest elevated material and dimmed backdrop.

## Layout and components

The desktop chat layout remains full-viewport and two-column. The mobile one-column conversation/chat switching behavior remains unchanged. White card-like regions are replaced by material layers and separators.

The current-user block, group/direct sections, action rows, selected states, unread badges, status indicator, message bubbles, composer, loading skeletons, empty states, errors, login form, group dialog, and admin page are all restyled against the same tokens. No API contracts or component business behavior change.

## Interaction and accessibility

Interactive elements retain visible focus styles, inline errors, loading states, and tactile press feedback. Hover and press transitions use only opacity and transforms. The existing reduced-motion media query remains the motion fallback. The online indicator keeps a restrained pulse, and skeletons retain their shimmer.

## Implementation and validation

No new frontend dependency is required; the existing React and CSS architecture is sufficient. The implementation centralizes material tokens and updates the existing stylesheet plus only the markup needed for semantic visual hooks.

Run `npm run test`, `npm run lint`, and `npm run build`. Visually inspect desktop and mobile layouts, dialogs, selected states, loading/error/empty states, and glass fallback readability.
