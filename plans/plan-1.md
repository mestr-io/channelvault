# Plan 1: Implement Simple Message Search

**Status:** 🟢 Done

## Description
Implement a simple text search functionality for messages within the backend. This will allow the frontend to query messages containing specific text. We will use a standard SQL `LIKE` query for this initial implementation.

## Subtasks

- [x] **Define SearchResult Model**
  - Description: Create a struct in `backend/models` to hold search results (Message content, Channel Name, Author, Timestamp).

- [x] **Update Repository Interface**
  - Description: Add `SearchMessages(ctx context.Context, query string) ([]models.SearchResult, error)` to `backend/repository/repository.go`.

- [x] **Implement Repository Logic**
  - Description: Implement the SQL query in `backend/repository/messages.go` using `LIKE %...%`.
  - Details: Join `messages` with `channels` and `users`. Limit results to 100.

- [x] **Expose API Method**
  - Description: Add `SearchMessages(query string)` to `backend/app/api.go` to expose it to Wails.

- [x] **Add Unit Tests**
  - Description: Create `TestSearchMessages` in `backend/repository/messages_test.go` to verify the query.

- [x] **Update Documentation**
  - Description: Add the new query to `backend/repository/queries.md`.

## Notes/Updates
- This is a "simple" search implementation. Full Text Search (FTS5) is deferred for future optimizations if needed.
- {2026-02-16}: Implemented SQL LIKE query, verified with unit tests and app startup.
