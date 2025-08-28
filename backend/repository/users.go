package repository

import (
	"context"
	"errors"
	"fmt"

	ah "github.com/mestr-io/slack-export-viewer/backend/app_helpers"
	"github.com/mestr-io/slack-export-viewer/backend/models"
)

func (r *RepositorySql) BulkInsertUsers(ctx context.Context, users []models.User) (err error) {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		// in case of error then rollback the txn
		if err != nil {
			err = errors.Join(tx.Rollback(), err)
		}
	}()

	// extract importId from the context
	imp, ok := ctx.Value("import").(*models.SlackImport)
	if !ok {
		return fmt.Errorf("failed to retrieve context: %w", err)
	}
	stmt, err := tx.Prepare("INSERT OR REPLACE INTO users (import_id, id, name, profile_image_url, real_name, is_deleted, is_bot) VALUES (?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer func() { err = errors.Join(stmt.Close(), err) }()

	// 5. Iterate over users, filter, and execute the prepared statement.
	insertCount := 0
	for _, user := range users {
		_, err := stmt.Exec(imp.ID, user.ID, user.Name, user.Profile.Image32, user.RealName, user.Deleted, user.IsBot)
		if err != nil {
			return fmt.Errorf("failed to insert user %s: %w", user.ID, err)
		}
		insertCount++
	}

	ah.Infof("Successfully prepared %d users for insertion.", insertCount)

	return tx.Commit()
}

func (r *RepositorySql) GetUserAvatar(ctx context.Context, userId string) (string, error) {
	sql := `SELECT profile_image_url FROM users WHERE id = ? LIMIT 1;`
	rows := r.db.QueryRowContext(ctx, sql, userId)
	var url string
	err := rows.Scan(&url)
	if err != nil {
		return "", err
	}
	return url, nil
}
