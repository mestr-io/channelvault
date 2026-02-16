# Plan 3: Enhance Search Experience

**Status:** 🟢 Done

## Description
Improve the search experience by highlighting matching terms in the message view and ensuring the view refreshes even when navigating to a result within the currently active channel.

## Subtasks

- [x] **Implement Highlighting**
  - Description: Highlight the search term in the message body.
  - Actions:
    - Add `highlightTerm` state to `App.svelte`.
    - Pass `highlightTerm` to `Message.svelte`.
    - Modify `Message.svelte` to wrap the matching term in `<mark>` tags within the Markdown content.

- [x] **Fix Same-Channel Navigation**
  - Description: Ensure the view refreshes when clicking a search result for the current channel.
  - Actions:
    - Update `selectChannel` in `App.svelte` to accept an optional `force` parameter.
    - Update `handleSearchNavigation` to call `selectChannel(channel, true)`.

## Notes/Updates
- {2026-02-16}: Plan created.
- {2026-02-16}: Implemented highlighting and forced navigation. Verified by running the app.
