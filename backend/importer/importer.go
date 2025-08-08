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

func getExportBaseFolder(files []*zip.File) (string, error) {
	// Dani: i found that in an older export all the files are contained
	// in a root folder with an arbitrary name. Newer ones are deflatted in root
	// So. thanks slack. ill look i.e. users.json file, and from there i can calculate the
	// base folder.

	for _, f := range files {
		path := strings.Split(f.Name, "/")
		lastIndex := len(path) - 1
		fileName := path[lastIndex]
		if fileName == "users.json" {
			// found it! return the base
			if len(path) == 1 {
				return "", nil
			} else {
				return fmt.Sprintf("%s/", path[0]), nil
			}
		}
	}

	return "", fmt.Errorf("no base folder found")
}

func isNotValidMessageFile(path []string) bool {
	// This is pretty weird considering that the zip file has been produced
	// at slack side
	lastIndex := len(path) - 1
	return path[0] == "__MACOSX" || path[lastIndex] == "" || len(path) < 3
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
	baseFolder, err := getExportBaseFolder(r.File)
	if err != nil {
		return err
	}
	// users.json
	f, err := r.Open(fmt.Sprintf("%susers.json", baseFolder))
	if err != nil {
		return fmt.Errorf("fail to open unzip users.json file: %w", err)
	}
	i.processUsers(f)

	// channels.json
	f, err = r.Open(fmt.Sprintf("%schannels.json", baseFolder))
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

			// only process history files
			path := strings.Split(f.Name, "/")
			if isNotValidMessageFile(path) {
				continue
			}

			i.processMessages(path[1], f)
		}

	}
	return nil
}
