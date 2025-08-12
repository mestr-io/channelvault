package repository

import (
	"fmt"
	"log"
	"github.com/mestr-io/slack-export-viewer/backend/app_errors"
	"github.com/mestr-io/slack-export-viewer/backend/models"
)

func (r *Repository) BulkInsertUsers(users []models.User) (err error) {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		// in case of error then rollback the txn
		if err != nil {
			err = app_errors.Finally(tx.Rollback, err)
		}
	}()

	stmt, err := tx.Prepare("INSERT OR REPLACE INTO users (id, name, profile_image_url) VALUES (?, ?, ?)")
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer func(){err = app_errors.Finally(stmt.Close, err)}()

	// 5. Iterate over users, filter, and execute the prepared statement.
	insertCount := 0
	for _, user := range users {
		// Filter out deleted users and bots, as requested.
		if user.Deleted || user.IsBot {
			continue
		}

		_, err := stmt.Exec(user.ID, user.Name, user.Profile.Image32)
		if err != nil {
			return fmt.Errorf("failed to insert user %s: %w", user.ID, err)
		}
		insertCount++
	}

	log.Printf("Successfully prepared %d users for insertion.", insertCount)

	return tx.Commit()
}
