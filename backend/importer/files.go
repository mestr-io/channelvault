package importer

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/mestr-io/slack-export-viewer/backend/models"
	"github.com/mestr-io/slack-export-viewer/backend/repository"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

func DownloadBinaries(ctx context.Context, filesFolder string, repository repository.Repository, currentImport *models.SlackImport) error {
	runtime.EventsEmit(ctx, "downloadBinariesStarted", currentImport.ID)

	pageSize := 20
	for {
		page, err := repository.ListPendingBinaries(ctx, currentImport.ID, pageSize)
		if len(page) == 0 {
			break
		}
		if err != nil {
			return err
		}
		for _, file := range page {
			// download file
			destPath := filepath.Join(filesFolder, fmt.Sprintf("%s_%s", file.ID, file.Name))

			resp, err := http.Get(file.FileURL)
			if err != nil {
				return fmt.Errorf("fail downloading file %s : %v", file.FileURL, err)
			}
			defer resp.Body.Close()

			dest, err := os.Create(destPath)

			if err != nil {
				return fmt.Errorf("failed to create destination file %s : %v", destPath, err)
			}
			defer dest.Close()

			_, err = io.Copy(dest, resp.Body)
			if err != nil {
				return fmt.Errorf("failed to download file %s : %v", destPath, err)
			}

			// mark as downloaded
			if err = repository.MarkFileAsDownloaded(ctx, currentImport.ID, file.ID); err != nil {
				return fmt.Errorf("failed to modify file with id %s: %v", file.ID, err)
			}
			// notify progress
			runtime.EventsEmit(ctx, "downloadedBinary", currentImport.ID, file)
		}
	}

	runtime.EventsEmit(ctx, "downloadBinariesDone", currentImport.ID)
	return nil
}
