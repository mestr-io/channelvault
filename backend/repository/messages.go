package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/mestr-io/slack-export-viewer/backend/models"
)

// TODO: add context to stop it if needed
func (r *RepositorySql) BulkInsertMessages(ctx context.Context, ChannelID string, messages []models.Message) (err error) {
	// Start a new transaction
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}
	// Defer a rollback in case of an error. The rollback will be ignored if the transaction is committed.
	defer func() {
		if err != nil {
			err = errors.Join(tx.Rollback(), err)
		}
	}()

	// extract importId from the context
	imp, ok := ctx.Value("import").(*models.SlackImport)
	if !ok {
		return fmt.Errorf("failed to retrieve context: %v", err)
	}

	sqlStmtMsg := `
        INSERT OR IGNORE INTO messages (
            import_id, ts, channel_id, user_id, text, thread_ts, event_type, event_subtype,
            reactions_json, reply_users_json, replies_json, files_json,
            attachments_json, reply_count, latest_reply, reply_users_count
        ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);
    `

	// Prepare the statement. This is compiled once and reused for every message.
	stmt1, err := tx.Prepare(sqlStmtMsg)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %v", err)
	}
	defer func() { err = errors.Join(stmt1.Close(), err) }()

	sqlStmtFiles := `
        INSERT OR IGNORE INTO files (
            import_id, ts, channel_id, user, url_private, mimetype, id, name, title
        ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?);
    `

	// Prepare the statement. This is compiled once and reused for every message.
	stmt2, err := tx.Prepare(sqlStmtFiles)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %v", err)
	}
	defer func() { err = errors.Join(stmt2.Close(), err) }()

	// Loop through each message in the slice
	for _, msg := range messages {
		// Marshal complex fields into JSON strings for storage in TEXT columns.
		reactionsJSON, err := json.Marshal(msg.Reactions)
		if err != nil {
			return fmt.Errorf("failed to marshal reactions for ts %s: %v", msg.Ts, err)
		}

		// The 'reply_users_json' field in your table seems to correspond to the users in the 'replies' field.
		// Let's extract those user IDs.
		var replyUsers []string
		for _, r := range msg.Replies {
			replyUsers = append(replyUsers, r.User)
		}
		replyUsersJSON, err := json.Marshal(replyUsers)
		if err != nil {
			return fmt.Errorf("failed to marshal reply users for ts %s: %v", msg.Ts, err)
		}

		repliesJSON, err := json.Marshal(msg.Replies)
		if err != nil {
			return fmt.Errorf("failed to marshal replies for ts %s: %v", msg.Ts, err)
		}

		filesJSON, err := json.Marshal(msg.Files)
		if err != nil {
			return fmt.Errorf("failed to marshal files for ts %s: %v", msg.Ts, err)
		}
		filesNullable := toNullString(filesJSON)
		if filesNullable.Valid {
			// insert into files table
			for _, file := range msg.Files {
				_, e := stmt2.ExecContext(ctx, imp.ID, msg.Ts, ChannelID, file.User, file.FileURL, file.Mimetype, file.ID, file.Name, file.Title)
				err = errors.Join(err, e)
			}
		}

		attachmentsJSON, err := json.Marshal(msg.Attachments)
		if err != nil {
			return fmt.Errorf("failed to marshal attachments for ts %s: %v", msg.Ts, err)
		}

		// Execute the prepared statement with the data for the current message.
		_, err = stmt1.Exec(
			imp.ID,
			msg.Ts,
			ChannelID, // Make sure this value is set in your Message struct before calling this function.
			msg.UserId,
			msg.Text,
			msg.ThreadTs,
			msg.Type,
			msg.Subtype,
			toNullString(reactionsJSON),
			toNullString(replyUsersJSON),
			toNullString(repliesJSON),
			filesNullable,
			toNullString(attachmentsJSON),
			msg.ReplyCount,
			msg.LatestReply,
			msg.ReplyUsersCount,
		)

		if err != nil {
			return fmt.Errorf("failed to execute statement for ts %s: %v", msg.Ts, err)
		}
	}

	// If the loop completes without errors, commit the transaction.
	return tx.Commit()

}

func (r *RepositorySql) GetMessages(ctx context.Context, channelID string, cursor string, pageSize int) ([]models.Message, error) {
	query := `
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
        WHERE c.id = $1 AND (thread_ts = ts OR thread_ts = "") 
    `
	args := []interface{}{channelID}

	// If a cursor is provided, add a condition to fetch messages older than the cursor.
	// This is a common pattern for keyset/cursor-based pagination.
	if cursor != "" {
		query += " AND ts < $2"
		args = append(args, cursor)
	}

	// Always order by timestamp to ensure consistent pagination, and limit the result size.
	query += " ORDER BY ts DESC LIMIT $3"
	args = append(args, pageSize)

	// Execute the query using the context for cancellation safety.
	// Pass query parameters as arguments to prevent SQL injection.
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("could not query messages for channel %s: %v", channelID, err)
	}
	defer func() { err = errors.Join(rows.Close(), err) }()

	var messages []models.Message
	for rows.Next() {
		var msg models.Message
		var usersJson sql.NullString
		var filesJson sql.NullString
		var attachmentsJson sql.NullString

		if err := rows.Scan(&msg.Ts, &msg.Text, &msg.UserId, &msg.UserName, &msg.UserProfileUrl, &msg.ThreadTs, &msg.ReplyCount, &msg.LatestReply, &usersJson, &msg.ReplyUsersCount, &filesJson, &attachmentsJson, &msg.ImportID); err != nil {
			// Error during row scanning is a server-side issue.
			return nil, fmt.Errorf("could not scan message row: %v", err)
		}

		if usersJson.Valid {
			err = json.Unmarshal([]byte(usersJson.String), &msg.ReplyUsers)
			if err != nil {
				return nil, fmt.Errorf("could not deserialize reply users: %v", err)
			}
		}
		if filesJson.Valid {
			err = json.Unmarshal([]byte(filesJson.String), &msg.Files)
			if err != nil {
				return nil, fmt.Errorf("could not deserialize files: %v", err)
			}
		}
		if attachmentsJson.Valid {
			err = json.Unmarshal([]byte(attachmentsJson.String), &msg.Attachments)
			if err != nil {
				return nil, fmt.Errorf("could not deserialize attachments: %w", err)
			}
		}

		messages = append(messages, msg)
	}

	// Check for any errors that occurred during iteration.
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating message rows: %v", err)
	}

	return messages, nil
}

func (r *RepositorySql) ListChannelFiles(ctx context.Context, channelID string, pageSize int, cursor int) ([]models.File, error) {
	panic("not implemented")
}

func (r *RepositorySql) GetThreadByTs(ctx context.Context, threadTs string, channelID string) (messages []models.Message, err error) {
	query := `
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
			m.import_id
        FROM messages AS m
		JOIN users AS u ON u.id = m.user_id AND u.import_id = m.import_id
		JOIN channels AS c ON c.name = m.channel_id AND c.import_id = m.import_id
        WHERE c.id = $1 AND thread_ts = $2 
    `
	args := []interface{}{channelID, threadTs}

	// Always order by timestamp to ensure consistent pagination, and limit the result size.
	query += " ORDER BY ts ASC"

	// Execute the query using the context for cancellation safety.
	// Pass query parameters as arguments to prevent SQL injection.
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("could not query messages for channel %s: %v", channelID, err)
	}
	defer func() {
		err = errors.Join(err, rows.Close())
	}()

	for rows.Next() {
		var msg models.Message
		var usersJson sql.NullString
		var reactionsJson sql.NullString

		if err := rows.Scan(&msg.Ts, &msg.Text, &msg.UserId, &msg.UserName, &msg.UserProfileUrl, &msg.ThreadTs, &msg.ReplyCount, &msg.LatestReply, &usersJson, &msg.ReplyUsersCount, &reactionsJson, &msg.ImportID); err != nil {
			// Error during row scanning is a server-side issue.
			return nil, fmt.Errorf("could not scan message row: %v", err)
		}
		if usersJson.Valid {
			err = json.Unmarshal([]byte(usersJson.String), &msg.ReplyUsers)
			if err != nil {
				return nil, fmt.Errorf("could not deserialize reply users: %v", err)
			}
		}
		if reactionsJson.Valid {
			err = json.Unmarshal([]byte(reactionsJson.String), &msg.Reactions)
			if err != nil {
				return nil, fmt.Errorf("could not deserialize reactions: %v", err)
			}
		}
		messages = append(messages, msg)
	}

	// Check for any errors that occurred during iteration.
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating message rows: %v", err)
	}

	return messages, nil
}

func (r *RepositorySql) GetFileById(ctx context.Context, fileId string) (*models.File, error) {
	query := `
        SELECT 
			id,
			name,
			title,
			mimetype,
			url_private,
			user,
			import_id
        FROM files 
        WHERE id = ?; 
    `
	args := []interface{}{fileId}

	file := models.File{}
	err := r.db.QueryRowContext(ctx, query, args...).Scan(&file.ID, &file.Name, &file.Title, &file.Mimetype, &file.FileURL, &file.User, &file.ImportID)
	if err != nil {
		return nil, fmt.Errorf("could not query file with id %s: %v", fileId, err)
	}
	return &file, nil
}

func (r *RepositorySql) ListPendingBinaries(ctx context.Context, importID int64, pageSize int) ([]models.File, error) {
	query := `SELECT id, name, title, mimetype, url_private, ts FROM files WHERE  downloaded = 0 AND import_id = ? AND url_private != '' ORDER BY ts DESC LIMIT ?;`
	rows, err := r.db.QueryContext(ctx, query, importID, pageSize)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	files := []models.File{}
	for rows.Next() {
		var file models.File
		if err := rows.Scan(&file.ID, &file.Name, &file.Title, &file.Mimetype, &file.FileURL, &file.Ts); err != nil {
			return nil, err
		}
		files = append(files, file)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return files, nil
}

func (r *RepositorySql) CountPendingBinaries(ctx context.Context, importID int64) (int64, error) {
	query := `SELECT COUNT(*) FROM files WHERE downloaded = 0 AND import_id = ? AND url_private != '';`
	row := r.db.QueryRowContext(ctx, query, importID)

	var count int64
	err := row.Scan(&count)

	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *RepositorySql) MarkFileAsDownloaded(ctx context.Context, importID int64, fileID string) error {
	query := `UPDATE files SET downloaded = 1 WHERE import_id = ? AND id = ?;`
	if _, err := r.db.ExecContext(ctx, query, importID, fileID); err != nil {
		return fmt.Errorf("failed update file for import %d and id %s failed: %v", importID, fileID, err)
	}
	return nil
}
