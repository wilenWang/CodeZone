# Apple glass frontend implementation plan

> Execute from `/Users/wilenwang/works/CodeZone`.

## 1. Establish silver glass design tokens

1. Replace the flat neutral surface variables in `frontend/src/styles.css` with a cold silver-gray palette, translucent surface variants, material edge colors, inset highlights, and cool diffusion shadows.
2. Add a fixed, pointer-events-none background layer through the page root pseudo-element. Use layered radial gradients and no filter on scrolling elements.
3. Add an `@supports (backdrop-filter: blur(1px))` material utility/fallback pattern so translucent panels remain readable without browser blur support.
4. Keep teal as the only accent and retain the reduced-motion override.

## 2. Restyle application shells

1. Apply the silver material system to the login page, form controls, seed-user choices, errors, and primary actions.
2. Apply a denser but still readable silver material treatment to the admin page, headings, sections, and table rows.
3. Preserve existing semantic markup, route behavior, and form control states.

## 3. Restyle the chat workspace

1. Convert the chat layout to edge-to-edge glass: a more opaque sidebar navigation layer, a moderate header/composer material layer, and a transparent message field over the fixed background.
2. Restyle the current-user block, group/direct section labels, group creation button, conversation rows, direct-user rows, selected states, unread badges, and inline sidebar feedback.
3. Use fine dividers, inset highlights, and cool shadows instead of white card backgrounds. Keep desktop widths and the existing mobile single-column collapse behavior unchanged.

## 4. Restyle messages and composer

1. Give received messages low-elevation white-gray glass bubbles and sent messages translucent teal glass bubbles, retaining sufficient text contrast and failure styling.
2. Update the chat header, live indicator, back button, text area, and send controls to match the material system.
3. Preserve loading skeleton geometry, shimmer, empty states, retries, message animation, focus states, and press feedback.

## 5. Restyle group dialog

1. Upgrade the backdrop to a subdued material dimmer and the dialog to the strongest elevated glass surface, with refractive edge, inset highlight, and diffusion shadow.
2. Update dialog header, close button, field controls, member options, validation errors, and action buttons without changing dialog state or submission behavior.

## 6. Validate

1. Run `cd frontend && npm run test`.
2. Run `cd frontend && npm run lint`.
3. Run `cd frontend && npm run build`.
4. Use the running app to inspect desktop and mobile chat, login, admin, dialog, selected/hover/focus states, loading, empty, and error states. Confirm no horizontal overflow and that backdrop-filter fallback remains legible.
