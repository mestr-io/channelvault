// file: backend/importer/importer.go
package importer

import (
	"archive/zip"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"regexp"

	"github.com/mestr-io/slack-export-viewer/backend/models"
	"github.com/mestr-io/slack-export-viewer/backend/repository"
)

type Progress struct {
	FileName    string
	Progress    float32
	TotalFiles  int
	CurrentFile int
}

type Importer struct {
	ctx           context.Context
	repository    repository.Repository
	currentImport *models.SlackImport
}

func NewImporter(ctx context.Context, repository repository.Repository, currentImport *models.SlackImport) *Importer {
	ctx = context.WithValue(ctx, "import", currentImport)
	return &Importer{ctx: ctx, repository: repository, currentImport: currentImport}
}

func (i *Importer) processUsers(zf *zip.File) (err error) {
	if err = i.ctx.Err(); err != nil {
		return fmt.Errorf("context expired: %w", err)
	}
	f, err := zf.Open()
	if err != nil {
		return fmt.Errorf("FATAL: cannot open 'users.json': %v", err)
	}
	defer func() {
		err = errors.Join(f.Close(), err)
	}()

	data, err := io.ReadAll(f)
	if err != nil {
		return fmt.Errorf("fail to read 'users.json' file: %w", err)
	}

	var result []models.User
	err = json.Unmarshal(data, &result)
	if err != nil {
		return fmt.Errorf("fail to deserialize 'users.json': %v", err)
	}

	return i.repository.BulkInsertUsers(i.ctx, result)
}

func (i *Importer) processChannels(zf *zip.File) (err error) {
	if err = i.ctx.Err(); err != nil {
		return fmt.Errorf("context expired: %w", err)
	}
	f, err := zf.Open()
	if err != nil {
		return fmt.Errorf("FATAL: cannot open 'users.json': %v", err)
	}

	defer func() { err = errors.Join(f.Close(), err) }()

	data, err := io.ReadAll(f)
	if err != nil {
		return fmt.Errorf("fail to read 'channels.json' file: %w", err)
	}

	var result []models.Channel
	err = json.Unmarshal(data, &result)
	if err != nil {
		return fmt.Errorf("fail to deserialize 'channels.json': %v", err)
	}

	return i.repository.BulkInsertChannels(i.ctx, result)
}

func (i *Importer) processMessages(ChannelID string, f *zip.File) (err error) {
	if err = i.ctx.Err(); err != nil {
		return fmt.Errorf("context expired: %w", err)
	}
	file, err := f.Open()
	if err != nil {
		return fmt.Errorf("FATAL: fail to open messages file %s: %v", f.Name, err)
	}
	data, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("fail to read 'users.json' file: %w", err)
	}

	var result []models.Message
	err = json.Unmarshal(data, &result)
	if err != nil {
		return fmt.Errorf("fail to deserialize 'users.json': %v", err)
	}

	return i.repository.BulkInsertMessages(i.ctx, ChannelID, result)
}

func findFile(filename string, files []*zip.File) (*zip.File, error) {
	for _, f := range files {
		if filename == filepath.Base(f.Name) {
			return f, nil
		}
	}

	return nil, fmt.Errorf("FATAL: couldnt find file %s", filename)
}

func (i *Importer) Import(c chan models.Progress) error {
	defer close(c)

	r, err := zip.OpenReader(i.currentImport.Path)
	if err != nil {
		return fmt.Errorf("fail to open the zip export file: %w", err)
	}

	defer func() {
		err = errors.Join(r.Close(), err)
	}()

	totalFiles := len(r.File)

	files, err := i.repository.ListPendingFiles(i.ctx, i.currentImport.ID)
	// checkpoints
	if err != nil {
		return fmt.Errorf("fail to list checkpoints: %v", err)
	}

	if files == nil {
		// first sync: create the checkpoints
		for _, f := range r.File {
			files = append(files, f.Name)
		}
		err = i.repository.CreateCheckpointPlaceHolders(i.ctx, i.currentImport.ID, files)
		if err != nil {
			return fmt.Errorf("fail to create import placeholder: %w", err)
		}
	}

	// users.json
	f, err := findFile("users.json", r.File)
	if err != nil {
		return fmt.Errorf("fail to open unzip users.json file: %w", err)
	}

	found := false
	for _, p := range files {
		if f.Name == p {
			found = true
			break
		}
	}
	if found {
		err = i.processUsers(f)
		if err != nil {
			return fmt.Errorf("fail to process users: %w", err)
		}
		err = i.repository.MarkCheckpointAsCompleted(i.ctx, i.currentImport.ID, f.Name, "")
		if err != nil {
			return fmt.Errorf("fail to mark checkpoint as completed: %w", err)
		}
	}

	// channels.json
	f, err = findFile("channels.json", r.File)
	if err != nil {
		return fmt.Errorf("fail to open unzip channels.json file: %w", err)
	}
	for _, p := range files {
		if f.Name == p {
			found = true
			break
		}
	}
	if found {
		err = i.processChannels(f)
		if err != nil {
			return fmt.Errorf("fail to process channels: %w", err)
		}
		err = i.repository.MarkCheckpointAsCompleted(i.ctx, i.currentImport.ID, f.Name, "")
		if err != nil {
			return fmt.Errorf("fail to mark checkpoint as completed: %w", err)
		}
	}
	// messages & threads
	for idx, f := range r.File {
		select {
		case <-i.ctx.Done():
			// The context was cancelled (e.g., app is shutting down).
			return i.ctx.Err()
		default:
			for _, p := range files {
				if f.Name == p {
					found = true
					break
				}
			}
			if !found {
				continue
			}

			// only process history files
			messagesFileRegex := regexp.MustCompile(`([^/]+)/(\d{4}-\d{2}-\d{2}\.json)$`)
			if matches := messagesFileRegex.FindStringSubmatch(f.Name); len(matches) > 1 {
				channelID := matches[1]
				c <- models.Progress{
					ImportID:    i.currentImport.ID,
					FileName:    fmt.Sprintf("%s / %s", channelID, matches[2]),
					Progress:    100 * (idx + 1) / totalFiles,
					TotalFiles:  totalFiles,
					CurrentFile: idx + 1,
				}
				err = i.processMessages(channelID, f)
				if err != nil {
					err = errors.Join(err, i.repository.MarkCheckpointAsCompleted(i.ctx, i.currentImport.ID, f.Name, fmt.Sprintf("%v", err)))
					return fmt.Errorf("fail to process messages for channel %s: %w", channelID, err)
				}
			}
			err = i.repository.MarkCheckpointAsCompleted(i.ctx, i.currentImport.ID, f.Name, "")
			if err != nil {
				return fmt.Errorf("fail to mark checkpoint as completed: %w", err)
			}
		}
	}
	return nil
}
