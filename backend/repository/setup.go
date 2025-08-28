package repository

import (
	"fmt"
)

func (r *RepositorySql) SetupSchemas() error {
	createTableSQL := `
		CREATE TABLE IF NOT EXISTS imports (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT,
			path TEXT,
			created DATETIME DEFAULT CURRENT_TIMESTAMP,
			finished DATETIME,
			num_messages INTEGER DEFAULT 0,
			error TEXT
		);
		CREATE TABLE IF NOT EXISTS checkpoints (
			import_id INTEGER REFERENCES imports(id) ON DELETE CASCADE,
			path TEXT,
			processed_at DATETIME,
			error TEXT,
			PRIMARY KEY(import_id, path)
		);
		CREATE TABLE IF NOT EXISTS files (
			import_id INTEGER REFERENCES imports(id) ON DELETE CASCADE,
			id TEXT,
			name TEXT,
			title TEXT,
			mimetype TEXT,
			user TEXT,
			url_private TEXT,
			ts TEXT,
			channel_id TEXT,
			downloaded INTEGER DEFAULT 0,
			PRIMARY KEY (import_id, id)
		);
		CREATE TABLE IF NOT EXISTS messages (
			import_id INTEGER REFERENCES imports(id) ON DELETE CASCADE,
			ts TEXT NOT NULL,
			channel_id TEXT NOT NULL,
			user_id TEXT,
			text TEXT,
			thread_ts TEXT,
			event_type TEXT,
			event_subtype TEXT,
			reactions_json TEXT,
			files_json TEXT,
			reply_users_json TEXT,
			latest_reply TEXT,
			attachments_json TEXT,
			replies_json TEXT,
			reply_users_count INTEGER NOT NULL DEFAULT 0,
			reply_count INTEGER NOT NULL DEFAULT 0,
			PRIMARY KEY (import_id, ts, channel_id)
		);
		CREATE TABLE IF NOT EXISTS channels (
			import_id INTEGER REFERENCES imports(id) ON DELETE CASCADE,
			id TEXT,
			members_json TEXT,
			pins_json TEXT,
			topic TEXT,
			purpose TEXT,
			name TEXT NOT NULL,
			PRIMARY KEY (import_id, id)	
		);
		CREATE TABLE IF NOT EXISTS users (
			import_id INTEGER REFERENCES imports(id) ON DELETE CASCADE,
			id TEXT,
			name TEXT NOT NULL,
			real_name TEXT NOT NULL,
			is_deleted INTEGER DEFAULT 0,
			is_bot INTEGER DEFAULT 0,
			profile_image_url TEXT,
			PRIMARY KEY (import_id, id)	
		);`
	if _, err := r.db.Exec(createTableSQL); err != nil {
		return fmt.Errorf("failed to create table: %w", err)
	}
	return nil
}
