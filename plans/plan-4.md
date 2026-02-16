# Plan 4: Fix Search Highlighting Rendering

**Status:** 🟢 Done

## Description
Fix the issue where search term highlighting is not rendering in the message view. The `<mark>` tags are likely being stripped by the Markdown renderer or sanitizer.

## Subtasks

- [x] **Adjust Highlighting Syntax**
  - Description: Change the highlighting implementation to use `==` syntax (Markdown highlight) instead of HTML `<mark>` tags, or ensure HTML is allowed.
  - Action: Update `highlightText` in `Message.svelte`.
  - Result: Decided to stick with `<mark>` as it's standard HTML and easier to control with DOMPurify.

- [x] **Configure Sanitizer/Renderer**
  - Description: Ensure the Markdown renderer accepts the highlighting syntax or allows the HTML tag.
  - Action: Update `DOMPurify` options in `Message.svelte` to explicitly allow `mark` tag.

## Notes/Updates
- {2026-02-16}: Plan created.
- {2026-02-16}: Configured `DOMPurify` to allow `<mark>` tags. Verified with app launch.
