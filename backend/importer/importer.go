// file: backend/importer/importer.go
package importer

import (
	"archive/zip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"strings"

	"github.com/mestr-io/slack-export-viewer/backend/models"
	"github.com/mestr-io/slack-export-viewer/backend/repository"
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
	ctx        context.Context
	path       string
	repository *repository.Repository
}

func NewImporter(ctx context.Context, path string, repository *repository.Repository) *Importer {
	return &Importer{ctx: ctx, path: path, repository: repository}
}

func (i *Importer) processUsers(f fs.File) error {
	defer f.Close()
	var data []byte
	data, err := io.ReadAll(f)
	if err != nil {
		return fmt.Errorf("fail to read 'users.json' file: %w", err)
	}

	var result []models.User
	err = json.Unmarshal(data, &result)
	if err != nil {
		return fmt.Errorf("fail to deserialize 'users.json': %v", err)
	}

	return i.repository.BulkInsertUsers(result)
}

func (i *Importer) processChannels(f fs.File) error {
	defer f.Close()
	var data []byte
	data, err := io.ReadAll(f)
	if err != nil {
		return fmt.Errorf("fail to read 'channels.json' file: %w", err)
	}

	var result []models.Channel
	err = json.Unmarshal(data, &result)
	if err != nil {
		return fmt.Errorf("fail to deserialize 'channels.json': %v", err)
	}

	return i.repository.BulkInsertChannels(result)
}

func getExportBaseFolder(r *zip.ReadCloser) (string, error) {
	if len(r.File) > 0 {
		firstFile := r.File[0]
		path := strings.Split(firstFile.Name, "/")
		if len(path) > 0 {
			return path[0], nil
		} else {
			return "", fmt.Errorf("no base folder found")
		}
	} else {
		return "", fmt.Errorf("no files in the zip")
	}
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

	// get the base folder
	baseFolder, err := getExportBaseFolder(r)
	if err != nil {
		return err
	}
	// users.json
	f, err := r.Open(fmt.Sprintf("%s/users.json", baseFolder))
	if err != nil {
		return fmt.Errorf("fail to open unzip users.json file: %w", err)
	}
	i.processUsers(f)

	// channels.json
	f, err = r.Open(fmt.Sprintf("%s/channels.json", baseFolder))
	if err != nil {
		return fmt.Errorf("fail to open unzip channels.json file: %w", err)
	}
	i.processChannels(f)

	// messages & threads
	for idx, f := range r.File {
		if len(strings.Split(f.Name, "/")) == 1 {
			// Skip files that are not in a subfolder (e.g., users.json, channels.json)
			// NOTE: progress is not superaccurate here, but it's fine for now
			continue
		}

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

	}
	return nil
}
