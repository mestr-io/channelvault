# Backend SQL Queries

This document lists the SQL queries used in the `backend/repository` package, verified against the SQLite schema.

## Channels (`channels.go`)

### `GetChannels`
Retrieves a list of unique channels available in the database.
```sql
SELECT DISTINCT(id), name FROM channels;
```
*   **Test Status:** ✅ Verified
*   **Notes:** Returns empty list if no channels found.

---

## Imports (`imports.go`)

### `ListPendingFiles`
Finds checkpoints that haven't been processed yet for a specific import.
```sql
SELECT path FROM checkpoints WHERE processed = 0 AND import_id = ?;
```
*   **Test Status:** ✅ Verified

### `GetImport`
Retrieves details for a specific import by ID.
```sql
SELECT id, name, path, created, finished, num_messages, error, download_started, download_finished 
FROM imports WHERE id = ?;
```
*   **Test Status:** ✅ Verified

### `NoFinishedImports`
Checks if there are any completed imports in the system.
```sql
SELECT count(*) > 0 FROM imports WHERE finished IS NOT NULL;
```
*   **Test Status:** ✅ Verified

### `ListImports`
Retrieves all imports with a computed count of pending binary downloads.
```sql
SELECT
    i.id,
    i.created,
    i.finished,
    i.num_messages,
    i.path,
    i.name,
    i.error,
    i.download_started,
    i.download_finished,
    COALESCE((
        SELECT COUNT(*)
        FROM checkpoints c
        WHERE c.import_id = i.id AND c.processed = 0
    ), 0) AS pending_binaries
FROM imports i
ORDER BY i.id DESC;
```
*   **Test Status:** ✅ Verified

---

## Messages (`messages.go`)

### `GetMessages`
Retrieves messages for a specific channel, joining with User and Channel tables for metadata. Handles pagination via cursor (timestamp) and limit.
```sql
SELECT 
    ts, text, 
    user_id, 
    u.real_name as author, 
    u.profile_image_url as profile_url,
    thread_ts,
    reply_count,
    latest_reply,
    reply_users_json,
    reply_users_count,
    files_json,
    attachments_json,
    m.import_id
FROM messages AS m
JOIN users AS u ON u.id = m.user_id AND u.import_id = m.import_id
JOIN channels AS c ON c.name = channel_id AND c.import_id = m.import_id
WHERE c.id = ? AND (thread_ts = ts OR thread_ts = "") 
-- Optional: AND ts < ? (if cursor provided)
ORDER BY ts DESC LIMIT ?;
```
*   **Test Status:** ✅ Verified

### `GetThreadByTs`
Retrieves all messages in a specific thread.
```sql
SELECT 
    ts, text, 
    user_id, 
    u.real_name as author, 
    u.profile_image_url as profile_url,
    '' as thread_ts,
    reply_count,
    latest_reply,
    reply_users_json,
    reply_users_count,
    reactions_json,
    m.import_id,
    files_json,
    attachments_json
FROM messages AS m
JOIN users AS u ON u.id = m.user_id AND u.import_id = m.import_id
JOIN channels AS c ON c.name = m.channel_id AND c.import_id = m.import_id
WHERE c.id = ? AND thread_ts = ? AND m.import_id = ?
ORDER BY ts ASC;
```
*   **Test Status:** ✅ Verified (Note: Inner Join may filter out system messages like USLACKBOT if not present in users table).

### `GetFileById`
Retrieves file metadata by ID.
```sql
SELECT id, name, title, mimetype, url_private, user, import_id 
FROM files WHERE id = ?;
```
*   **Test Status:** ✅ Verified

### `ListPendingBinaries`
Lists files that haven't been downloaded locally yet.
```sql
SELECT id, name, title, mimetype, url_private, ts 
FROM files 
WHERE downloaded = 0 AND import_id = ? AND url_private != '' 
ORDER BY ts DESC LIMIT ?;
```
*   **Test Status:** ✅ Verified

### `CountPendingBinaries`
Counts total files waiting for download.
```sql
SELECT COUNT(*) 
FROM files 
WHERE downloaded = 0 AND import_id = ? AND url_private != '';
```
*   **Test Status:** ✅ Verified

---

## Users (`users.go`)

### `GetUserAvatar`
Gets the profile image URL for a user.
```sql
SELECT profile_image_url FROM users WHERE id = ? LIMIT 1;
```
*   **Test Status:** ✅ Verified

### `GetUserHandle`
Gets the real name for a user.
```sql
SELECT real_name FROM users WHERE id = ? LIMIT 1;
```
*   **Test Status:** ✅ Verified

---

## Schema Reference
For full table definitions, refer to `schema_dump.sql` in this directory.
