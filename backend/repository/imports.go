package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/mestr-io/slack-export-viewer/backend/models"
)

func (r *RepositorySql) CreateCheckpointPlaceHolders(ctx context.Context, importID int64, paths []string) (err error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("starting transaction failed: %w", err)
	}
	defer func() {
		if err != nil {
			err = errors.Join(err, tx.Rollback())
		}
	}()

	stmt, err := tx.PrepareContext(ctx, `INSERT OR IGNORE INTO checkpoints (import_id, path) VALUES (?, ?)`)
	if err != nil {
		return fmt.Errorf("creating prepare statement failed: %w", err)
	}
	defer func() { err = errors.Join(err, stmt.Close()) }()

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
	sql := `UPDATE checkpoints SET processed = 1, error = ? WHERE import_id = ? AND path = ?;`
	if _, err := r.db.ExecContext(ctx, sql, errMsg, importID, path); err != nil {
		return fmt.Errorf("update checkpoint for import %d and path %s failed: %w", importID, path, err)
	}
	return nil
}

func (r *RepositorySql) ListPendingFiles(ctx context.Context, importID int64) ([]string, error) {
	query := `SELECT path FROM checkpoints WHERE processed = 0 AND import_id = ?;`
	rows, err := r.db.Query(query, importID)
	if err != nil {
		return nil, err
	}
	defer func() {
		err = errors.Join(err, rows.Close())
	}()

	var checkpoints []string
	for rows.Next() {
		var path string
		if err := rows.Scan(&path); err != nil {
			return nil, err
		}
		checkpoints = append(checkpoints, path)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return checkpoints, nil
}

func (r *RepositorySql) GetImport(ctx context.Context, id int64) (*models.SlackImport, error) {
	sqlStatement := `SELECT id, name, path, created, finished, num_messages, error, download_started, download_finished FROM imports WHERE id = ?`
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
		&imp.DownloadStarted,
		&imp.DownloadFinished,
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

func (r *RepositorySql) StartDownload(ctx context.Context, importID int64) error {
	sql := `UPDATE imports SET download_started = CURRENT_TIMESTAMP WHERE id = ?`
	if _, err := r.db.ExecContext(ctx, sql, importID); err != nil {
		return fmt.Errorf("updating import %d failed: %v", importID, err)
	}
	return nil
}

func (r *RepositorySql) FinishDownload(ctx context.Context, importID int64) error {
	sql := `UPDATE imports SET download_finished = CURRENT_TIMESTAMP WHERE id = ?`
	if _, err := r.db.ExecContext(ctx, sql, importID); err != nil {
		return fmt.Errorf("updating import %d failed: %v", importID, err)
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
	const sqlStatement = `
SELECT
    i.id,
    i.created,
    i.finished,
    i.num_messages,
    i.path,
    i.name,
    i.error,
    i.download_started,
    i.download_finished,
    COALESCE((
        SELECT COUNT(*)
        FROM checkpoints c
        WHERE c.import_id = i.id AND c.processed = 0
    ), 0) AS pending_binaries
FROM imports i
ORDER BY i.id DESC;
`
	rows, err := r.db.QueryContext(ctx, sqlStatement)
	if err != nil {
		return nil, fmt.Errorf("FATAL: couldn't query imports: %w", err)
	}
	defer func() {
		if cerr := rows.Close(); cerr != nil {
			err = errors.Join(err, cerr)
		}
	}()

	for rows.Next() {
		var imp models.SlackImport
		var finished, errorMsg sql.NullString

		err = rows.Scan(
			&imp.ID,
			&imp.Created,
			&finished,
			&imp.NumFiles,
			&imp.Path,
			&imp.Name,
			&errorMsg,
			&imp.DownloadStarted,
			&imp.DownloadFinished,
			&imp.PendingBinaries,
		)
		if err != nil {
			return nil, fmt.Errorf("FATAL: couldn't scan import row: %w", err)
		}

		if finished.Valid {
			imp.Finished = finished.String
		}
		if errorMsg.Valid {
			imp.Error = errorMsg.String
		}

		imports = append(imports, imp)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("FATAL: couldn't iterate imports: %w", err)
	}

	if imports == nil {
		imports = []models.SlackImport{}
	}
	return imports, nil
}
