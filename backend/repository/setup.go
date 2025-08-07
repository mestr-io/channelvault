package repository

import (
	"fmt"
)

func (r *Repository) SetupSchemas() error {
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS users (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		profile_image_url TEXT
	);`
	if _, err := r.db.Exec(createTableSQL); err != nil {
		return fmt.Errorf("failed to create users table: %w", err)
	}

	createTableSQL = `
	CREATE TABLE IF NOT EXISTS channels (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL		
	);`
	if _, err := r.db.Exec(createTableSQL); err != nil {
		return fmt.Errorf("failed to create channels table: %w", err)
	}	
	
	return nil
}
