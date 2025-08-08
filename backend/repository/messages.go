package repository

import (
	"encoding/json"
	"fmt"

	"github.com/mestr-io/slack-export-viewer/backend/models"
)

func (r *Repository) BulkInserMessages(ChannelID string, messages []models.Message) (err error) {
	// Start a new transaction
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	// Defer a rollback in case of an error. The rollback will be ignored if the transaction is committed.
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// The SQL statement for insertion. Using 'INSERT OR IGNORE' will prevent errors
	// if you try to insert a message with a (ts, channel_id) that already exists.
	// Use 'INSERT OR REPLACE' if you want to update existing messages instead.
	sqlStmt := `
        INSERT OR IGNORE INTO messages (
            ts, channel_id, user_id, text, thread_ts, event_type, event_subtype,
            reactions_json, reply_users_json, replies_json, files_json,
            attachments_json, reply_count
        ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);
    `

	// Prepare the statement. This is compiled once and reused for every message.
	stmt, err := tx.Prepare(sqlStmt)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

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
			msg.User,
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
		)
		if err != nil {
			return fmt.Errorf("failed to execute statement for ts %s: %w", msg.Ts, err)
		}
	}

	// If the loop completes without errors, commit the transaction.
	return tx.Commit()

}
