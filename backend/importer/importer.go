// file: backend/importer/importer.go
package importer

import (
	"archive/zip"
	"context"
	"encoding/json"
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
	ctx        context.Context
	path       string
	repository *repository.Repository
}

func NewImporter(ctx context.Context, path string, repository *repository.Repository) *Importer {
	return &Importer{ctx: ctx, path: path, repository: repository}
}

func (i *Importer) processUsers(zf *zip.File) error {
	f, err := zf.Open()
	if err != nil {
		return fmt.Errorf("FATAL: cannot open 'users.json': %v", err)
	}
	defer f.Close()
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

func (i *Importer) processChannels(zf *zip.File) error {
	f, err := zf.Open()
	if err != nil {
		return fmt.Errorf("FATAL: cannot open 'users.json': %v", err)
	}
	defer f.Close()
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

func (i *Importer) processMessages(ChannelID string, f *zip.File) error {
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

	return i.repository.BulkInserMessages(ChannelID, result)
}

func findFile(filename string, files []*zip.File) (*zip.File, error) {
	for _, f := range files {
		if filename == filepath.Base(f.Name) {
			return f, nil
		}
	}

	return nil, fmt.Errorf("FATAL: couldnt find file %s", filename)
}

func (i *Importer) Import(c chan Progress) error {
	defer close(c)

	r, err := zip.OpenReader(i.path)
	if err != nil {
		return fmt.Errorf("fail to open the zip export file: %w", err)
	}
	defer r.Close()

	totalFiles := len(r.File)

	// users.json
	f, err := findFile("users.json", r.File)
	if err != nil {
		return fmt.Errorf("fail to open unzip users.json file: %w", err)
	}
	i.processUsers(f)

	// channels.json
	f, err = findFile("channels.json", r.File)
	if err != nil {
		return fmt.Errorf("fail to open unzip channels.json file: %w", err)
	}
	i.processChannels(f)

	// messages & threads
	for idx, f := range r.File {
		select {
		case <-i.ctx.Done():
			// The context was cancelled (e.g., app is shutting down).
			return i.ctx.Err()
		default:
			// only process history files
			messagesFileRegex := regexp.MustCompile(`([^/]+)/(\d{4}-\d{2}-\d{2}\.json)$`)
			if matches := messagesFileRegex.FindStringSubmatch(f.Name); len(matches) > 1 {
				channelID := matches[1]
				c <- Progress{
					FileName:    fmt.Sprintf("%s / %s", channelID, matches[2]),
					Progress:    float32(idx+1) / float32(totalFiles),
					TotalFiles:  totalFiles,
					CurrentFile: idx + 1,
				}
				i.processMessages(channelID, f)
			}

		}
	}
	return nil
}
