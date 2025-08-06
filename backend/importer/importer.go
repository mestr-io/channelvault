// file: backend/importer/importer.go
package importer

import (
	"archive/zip"
	"context"
	"database/sql"
	"fmt"
	"time"
)

type Progress struct {
	// It's good practice to make struct fields exported (uppercase)
	// if they are part of a public API, which this is.
	FileName    string
	Progress    float32
	TotalFiles  int
	CurrentFile int
}

type Importer struct {
	ctx  context.Context
	path string
	db   *sql.DB
}

func NewImporter(ctx context.Context, path string, db *sql.DB) *Importer {
	return &Importer{ctx: ctx, path: path, db: db}
}

// Import now returns an error instead of panicking
func (i *Importer) Import(c chan Progress) error {
	defer close(c)

	r, err := zip.OpenReader(i.path)
	if err != nil {
		return fmt.Errorf("fail to open the zip export file: %w", err)
	}
	defer r.Close()

	totalFiles := len(r.File)
	for idx, f := range r.File {
		select {
		case <-i.ctx.Done():
			// The context was cancelled (e.g., app is shutting down).
			return i.ctx.Err()
		default:
			c <- Progress{
				FileName:    f.Name,
				Progress:    float32(idx+1) / float32(totalFiles),
				TotalFiles:  totalFiles,
				CurrentFile: idx + 1,
			}
		}
		time.Sleep(1 * time.Millisecond)
	}
	return nil
}
