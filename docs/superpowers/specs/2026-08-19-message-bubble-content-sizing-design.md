# Message bubble content sizing design

## Goal

Make short chat messages use only the width required by their content instead of stretching toward the existing maximum width.

## Behavior

The message list will align items to the start, preventing flexbox cross-axis stretching. Incoming bubbles will size to their intrinsic content width while retaining the existing maximum width and wrapping behavior for longer messages. Outgoing bubbles will explicitly align themselves to the end, preserving their right-hand position with the same intrinsic sizing.

Message metadata, failed-message controls, long content wrapping, mobile layout, and the existing 78% maximum bubble width remain unchanged.

## Validation

Run the frontend tests, type check, and production build. Visually verify one-character, short, and long messages from both sender directions at desktop and mobile widths.
