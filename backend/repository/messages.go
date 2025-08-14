package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"errors"

	"github.com/mestr-io/slack-export-viewer/backend/models"

)

// TODO: add context to stop it if needed
func (r *Repository) BulkInserMessages(ChannelID string, messages []models.Message) (err error) {
	// Start a new transaction
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	// Defer a rollback in case of an error. The rollback will be ignored if the transaction is committed.
	defer func() {
		if err != nil {	
			err = errors.Join(tx.Rollback(), err)	
		}
	}()

	// The SQL statement for insertion. Using 'INSERT OR IGNORE' will prevent errors
	// if you try to insert a message with a (ts, channel_id) that already exists.
	// Use 'INSERT OR REPLACE' if you want to update existing messages instead.
	sqlStmt := `
        INSERT OR IGNORE INTO messages (
            ts, channel_id, user_id, text, thread_ts, event_type, event_subtype,
            reactions_json, reply_users_json, replies_json, files_json,
            attachments_json, reply_count, latest_reply, reply_users_count
        ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);
    `

	// Prepare the statement. This is compiled once and reused for every message.
	stmt, err := tx.Prepare(sqlStmt)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer func(){err = errors.Join(stmt.Close(), err)}()

	// Loop through each message in the slice
	for _, msg := range messages {
		// Marshal complex fields into JSON strings for storage in TEXT columns.
		reactionsJSON, err := json.Marshal(msg.Reactions)
		if err != nil {
			return fmt.Errorf("failed to marshal reactions for ts %s: %w", msg.Ts, err)
		}

		// The 'reply_users_json' field in your table seems to correspond to the users in the 'replies' field.
		// Let's extract those user IDs.
		var replyUsers []string
		for _, r := range msg.Replies {
			replyUsers = append(replyUsers, r.User)
		}
		replyUsersJSON, err := json.Marshal(replyUsers)
		if err != nil {
			return fmt.Errorf("failed to marshal reply users for ts %s: %w", msg.Ts, err)
		}

		repliesJSON, err := json.Marshal(msg.Replies)
		if err != nil {
			return fmt.Errorf("failed to marshal replies for ts %s: %w", msg.Ts, err)
		}

		filesJSON, err := json.Marshal(msg.Files)
		if err != nil {
			return fmt.Errorf("failed to marshal files for ts %s: %w", msg.Ts, err)
		}

		attachmentsJSON, err := json.Marshal(msg.Attachments)
		if err != nil {
			return fmt.Errorf("failed to marshal attachments for ts %s: %w", msg.Ts, err)
		}

		// Execute the prepared statement with the data for the current message.
		_, err = stmt.Exec(
			msg.Ts,
			ChannelID, // Make sure this value is set in your Message struct before calling this function.
			msg.UserId,
			msg.Text,
			msg.ThreadTs,
			msg.Type,
			msg.Subtype,
			string(reactionsJSON),
			string(replyUsersJSON),
			string(repliesJSON),
			string(filesJSON),
			string(attachmentsJSON),
			msg.ReplyCount,
			msg.LatestReply,
			msg.ReplyUsersCount,
		)

		if err != nil {
			return fmt.Errorf("failed to execute statement for ts %s: %w", msg.Ts, err)
		}
	}

	// If the loop completes without errors, commit the transaction.
	return tx.Commit()

}

func (r *Repository) GetMessages(ctx context.Context, channelID string, cursor string, pageSize int) ([]models.Message, error) {
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
			reply_users_count
        FROM messages 
		JOIN users AS u ON u.id = user_id
		JOIN channels AS c ON c.name = channel_id
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
		return nil, fmt.Errorf("could not query messages for channel %s: %w", channelID, err)
	}
	defer func(){err = errors.Join(rows.Close(), err)}()

	var messages []models.Message
	for rows.Next() {
		var msg models.Message
		var usersJson []byte
		if err := rows.Scan(&msg.Ts, &msg.Text, &msg.UserId, &msg.UserName, &msg.UserProfileUrl, &msg.ThreadTs, &msg.ReplyCount, &msg.LatestReply, &usersJson, &msg.ReplyUsersCount); err != nil {
			// Error during row scanning is a server-side issue.
			return nil, fmt.Errorf("could not scan message row: %w", err)
		}
		err = json.Unmarshal(usersJson, &msg.ReplyUsers)
		if err != nil {
			return nil, fmt.Errorf("could not deserialize reply users: %w", err)

		}
		messages = append(messages, msg)
	}

	// Check for any errors that occurred during iteration.
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating message rows: %w", err)
	}

	return messages, nil
}

func (r *Repository) GetThreadByTs(ctx context.Context, threadTs string, channelID string) ([]models.Message, error) {
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
			reply_users_count
        FROM messages 
		JOIN users AS u ON u.id = user_id
		JOIN channels AS c ON c.name = channel_id
        WHERE c.id = $1 AND thread_ts = $2 
    `
	args := []interface{}{channelID, threadTs}

	// Always order by timestamp to ensure consistent pagination, and limit the result size.
	query += " ORDER BY ts DESC"

	// Execute the query using the context for cancellation safety.
	// Pass query parameters as arguments to prevent SQL injection.
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("could not query messages for channel %s: %w", channelID, err)
	}
	defer rows.Close()

	var messages []models.Message
	for rows.Next() {
		var msg models.Message
		var usersJson []byte
		if err := rows.Scan(&msg.Ts, &msg.Text, &msg.UserId, &msg.UserName, &msg.UserProfileUrl, &msg.ThreadTs, &msg.ReplyCount, &msg.LatestReply, &usersJson, &msg.ReplyUsersCount); err != nil {
			// Error during row scanning is a server-side issue.
			return nil, fmt.Errorf("could not scan message row: %w", err)
		}
		err = json.Unmarshal(usersJson, &msg.ReplyUsers)
		if err != nil {
			return nil, fmt.Errorf("could not deserialize reply users: %w", err)

		}
		messages = append(messages, msg)
	}

	// Check for any errors that occurred during iteration.
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating message rows: %w", err)
	}

	return messages, nil
}
