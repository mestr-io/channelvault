# Plan 5: Implement Full Text Search (FTS5)

**Status:** 🟡 In Progress

## Description
Replace the simple SQL `LIKE` search with SQLite's **FTS5** engine to support advanced features: Boolean operators (`AND`, `OR`, `NOT`), prefix wildcards (`term*`), and ranked results. This involves database schema changes, backend query updates, and smarter frontend highlighting.

## Subtasks

- [ ] **Database Migration & Schema Setup**
  - Description: Create `messages_fts` virtual table and sync triggers.
  - Actions:
    - Update `backend/repository/setup.go` to create the virtual table (External Content mode).
    - Add Triggers (INSERT, DELETE, UPDATE) to sync `messages` and `messages_fts`.
    - Add logic to rebuild/populate the index for existing data.

- [ ] **Update Repository Logic**
  - Description: Switch `SearchMessages` to use FTS MATCH.
  - Actions:
    - Update `backend/repository/messages.go`.
    - Use `highlight()` function from SQLite to return preview snippets.
    - Handle query syntax errors gracefully (fallback or error message).

- [ ] **Frontend: Display DB Highlights**
  - Description: Use the snippet returned by DB for the search results list.
  - Actions:
    - Update `SearchManager.svelte` to display the highlighted preview text.

- [ ] **Frontend: Advanced Message Highlighting**
  - Description: Parse boolean/wildcard queries to highlight terms correctly in the message view.
  - Actions:
    - Implement a query parser in `Message.svelte` (or a utility file) to convert FTS syntax to Regex for the in-message view.

## Notes/Updates
- {2026-02-16}: Plan created. Verified FTS5 support.
