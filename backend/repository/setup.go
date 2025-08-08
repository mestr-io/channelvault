package repository

import (
	"fmt"
)

func (r *Repository) SetupSchemas() error {
	createTableSQL := `
		PRAGMA foreign_keys = ON;
		CREATE TABLE IF NOT EXISTS messages (
			ts TEXT NOT NULL,
			channel_id TEXT NOT NULL,
			user_id TEXT,
			text TEXT,
			thread_ts TEXT,
			event_type TEXT,
			event_subtype TEXT,
			reactions_json TEXT,
			reply_users_json TEXT,
			replies_json TEXT,
			files_json TEXT,
			attachments_json TEXT,
			reply_count INTEGER NOT NULL DEFAULT 0,
			PRIMARY KEY (ts, channel_id)
		);
		CREATE TABLE IF NOT EXISTS channels (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL		
		);
		CREATE TABLE IF NOT EXISTS users (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			profile_image_url TEXT
		);`
	if _, err := r.db.Exec(createTableSQL); err != nil {
		return fmt.Errorf("failed to create table: %w", err)
	}
	return nil
}
