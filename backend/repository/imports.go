package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/mestr-io/slack-export-viewer/backend/models"
)

func (r *RepositorySql) CreateCheckpointPlaceHolders(ctx context.Context, importID int64, paths []string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("starting transaction failed: %w", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	stmt, err := tx.PrepareContext(ctx, `INSERT OR IGNORE INTO checkpoints (import_id, path) VALUES (?, ?)`)
	if err != nil {
		return fmt.Errorf("creating prepare statement failed: %w", err)
	}
	defer stmt.Close()

	for _, path := range paths {
		if _, err = stmt.ExecContext(ctx, importID, path); err != nil {
			return fmt.Errorf("inserting record for path %s failed: %w", path, err)
		}
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("commiting transaction failed: %w", err)
	}

	return nil
}

func (r *RepositorySql) MarkCheckpointAsCompleted(ctx context.Context, importID int64, path string, errMsg string) error {
	sql := `UPDATE checkpoints SET processed_at = CURRENT_TIMESTAMP, error = ? WHERE import_id = ? AND path = ?;`
	if _, err := r.db.ExecContext(ctx, sql, errMsg, importID, path); err != nil {
		return fmt.Errorf("update checkpoint for import %d and path %s failed: %w", importID, path, err)
	}
	return nil
}

func (r *RepositorySql) GetImport(ctx context.Context, id int64) (*models.SlackImport, error) {
	sqlStatement := `SELECT id, name, path, created, finished, num_messages, error FROM imports WHERE id = ?`
	var imp models.SlackImport
	var finished sql.NullString
	err := r.db.QueryRowContext(ctx, sqlStatement, id).Scan(
		&imp.ID,
		&imp.Name,
		&imp.Path,
		&imp.Created,
		&finished,
		&imp.NumFiles,
		&imp.Error,
	)
	if err != nil {
		return nil, fmt.Errorf("retrieving import %d failed: %w", id, err)
	}
	if finished.Valid {
		imp.Finished = finished.String
	}
	return &imp, nil
}

func (r *RepositorySql) FinishedImportUpdate(ctx context.Context, imp *models.SlackImport) error {
	sql := `UPDATE imports SET name = ?, finished = CURRENT_TIMESTAMP, num_messages = ?, error = ? WHERE id = ?`
	if _, err := r.db.ExecContext(ctx, sql,
		imp.Name,
		imp.NumFiles,
		imp.Error,
		imp.ID,
	); err != nil {
		return fmt.Errorf("updating import %d failed: %w", imp.ID, err)
	}
	return nil
}

func (r *RepositorySql) FinishedImportUpdateWithError(ctx context.Context, imp *models.SlackImport) error {
	sql := `UPDATE imports SET name = ?, num_messages = ?, error = ? WHERE id = ?`
	if _, err := r.db.ExecContext(ctx, sql,
		imp.Name,
		imp.NumFiles,
		imp.Error,
		imp.ID,
	); err != nil {
		return fmt.Errorf("updating import %d failed: %w", imp.ID, err)
	}
	return nil
}

func (r *RepositorySql) DeleteImport(ctx context.Context, id int64) error {
	sql := `DELETE FROM imports WHERE id = ?`
	_, err := r.db.ExecContext(ctx, sql, id)
	return err
}

func (r *RepositorySql) RenameImport(ctx context.Context, importID int64, name string) error {
	sql := `UPDATE imports SET name = ? WHERE id = ?`
	if _, err := r.db.ExecContext(ctx, sql, name, importID); err != nil {
		return fmt.Errorf("updating import %d failed: %v", importID, err)
	}
	return nil
}

func (r *RepositorySql) CreateImport(ctx context.Context, path string, name string) (*models.SlackImport, error) {
	sql := `INSERT INTO imports (path, name, error) VALUES (?, ?, '');`
	result, err := r.db.ExecContext(ctx, sql, path, name)
	if err != nil {
		return nil, err
	}
	var i models.SlackImport
	i.ID, err = result.LastInsertId()
	if err != nil {
		return nil, err
	}
	i.Name = name
	i.Path = path
	return &i, nil
}

func (r *RepositorySql) NoFinishedImports(ctx context.Context) (bool, error) {
	sql := `SELECT count(*) > 0 FROM imports WHERE finished IS NOT NULL;`
	rows := r.db.QueryRowContext(ctx, sql)
	var any bool
	err := rows.Scan(&any)
	if err != nil {
		return false, err
	}
	return any, nil
}

func (r *RepositorySql) ListImports(ctx context.Context) (imports []models.SlackImport, err error) {
	sqlStatement := `SELECT id, created, finished, num_messages, path, name, error FROM imports ORDER BY id DESC;`
	rows, err := r.db.Query(sqlStatement)
	if err != nil {
		return nil, fmt.Errorf("FATAL: couldnt query imports: %w", err)
	}
	defer func() {
		err = errors.Join(rows.Close(), err)
	}()

	for rows.Next() {
		var imp models.SlackImport
		var finished sql.NullString
		var errorMsg sql.NullString
		err = rows.Scan(&imp.ID, &imp.Created, &finished, &imp.NumFiles, &imp.Path, &imp.Name, &errorMsg)
		if err != nil {
			return nil, fmt.Errorf("FATAL: couldnt query imports: %w", err)
		}
		if finished.Valid {
			imp.Finished = finished.String
		}
		if errorMsg.Valid {
			imp.Error = errorMsg.String
		}

		imports = append(imports, imp)
	}
	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("FATAL: couldnt query imports: %w", err)
	}

	// Ensure empty slice instead of nil
	if imports == nil {
		imports = []models.SlackImport{}
	}

	return imports, nil
}
