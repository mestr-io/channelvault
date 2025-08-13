package repository

import (
	"context"
	"time"

	"github.com/mestr-io/slack-export-viewer/backend/models"
)

func (r *Repository) StoreImport(ctx context.Context, path string) (*models.Import, error) {
	sql := `INSERT INTO imports (path) VALUES (?);`
	result, err := r.db.ExecContext(ctx, sql, path)
	if err != nil {
		return nil, err
	}
	var i models.Import
	i.ID, err = result.LastInsertId()
	if err != nil {
		return nil, err
	}
	return &i, nil
}

func (r *Repository) FinishedImport(ctx context.Context, ID int64) error {
	sql := `UPDATE imports SET finished = ?;`
	_, err := r.db.ExecContext(ctx, sql, time.Now())
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) NoFinishedImports(ctx context.Context) (bool, error) {
	sql := `SELECT count(*) FROM imports WHERE finished IS NOT NULL;`
	rows := r.db.QueryRowContext(ctx, sql)
	var count int
	err := rows.Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
