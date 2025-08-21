package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/mestr-io/slack-export-viewer/backend/models"
)

func (r *Repository) CreateImport(ctx context.Context, path string) (*models.SlackImport, error) {
	sql := `INSERT INTO imports (path) VALUES (?);`
	result, err := r.db.ExecContext(ctx, sql, path)
	if err != nil {
		return nil, err
	}
	var i models.SlackImport
	i.ID, err = result.LastInsertId()
	if err != nil {
		return nil, err
	}
	return &i, nil
}

func (r *Repository) FinishedImport(ctx context.Context, ID int64) error {
	sql := `UPDATE imports SET finished = ?;`
	_, err := r.db.ExecContext(ctx, sql, time.Now().Format(time.RFC3339))
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) NoFinishedImports(ctx context.Context) (bool, error) {
	sql := `SELECT count(*) > 0 FROM imports WHERE finished IS NOT NULL;`
	rows := r.db.QueryRowContext(ctx, sql)
	var any bool
	err := rows.Scan(&any)
	if err != nil {
		return false, err
	}
	return any, nil
}

func (r *Repository) ListImports(ctx context.Context) (imports []models.SlackImport, err error) {
	sql := `SELECT id, created, finished, num_messages, path FROM imports;`
	rows, err := r.db.Query(sql)
	if err != nil {
		return nil, fmt.Errorf("FATAL: couldnt query imports: %w", err)
	}
	defer func() {
		err = errors.Join(rows.Close(), err)
	}()

	for rows.Next() {
		var imp models.SlackImport
		err = rows.Scan(&imp.ID, &imp.Created, &imp.Finished, &imp.NumMessages, &imp.Path)
		if err != nil {
			return nil, fmt.Errorf("FATAL: couldnt query imports: %w", err)
		}
		imports = append(imports, imp)
	}
	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("FATAL: couldnt query imports: %w", err)
	}
	return imports, nil
}
