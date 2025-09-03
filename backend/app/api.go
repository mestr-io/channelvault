package app

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/mestr-io/slack-export-viewer/backend/importer"
	"github.com/mestr-io/slack-export-viewer/backend/models"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

func (a *App) GetChannels() ([]models.Channel, error) {
	return a.repository.GetChannels()
}

func (a *App) GetMessages(channelID string, cursor string, pageSize int) ([]models.Message, error) {
	return a.repository.GetMessages(a.ctx, channelID, cursor, pageSize)
}

func (a *App) GetUserAvatar(userId string) (string, error) {
	return a.repository.GetUserAvatar(a.ctx, userId)
}

func (a *App) GetThreadByTs(threadTs string, channelID string) ([]models.Message, error) {
	return a.repository.GetThreadByTs(a.ctx, threadTs, channelID)
}

func (a *App) GetBase64Image(fileId string) (img string, err error) {
	//TODO: move this?
	file, err := a.repository.GetFileById(a.ctx, fileId)
	if err != nil {
		return "", fmt.Errorf("file not found: %v", err)
	}
	destPath := filepath.Join(a.filesFolder, fmt.Sprintf("%s_%s", fileId, file.Name))
	// check local FS for the image
	if _, err := os.Stat(destPath); err != nil {
		// get the file
		resp, err := http.Get(file.FileURL)
		if err != nil {
			return "", fmt.Errorf("fail downloading file %s : %v", file.FileURL, err)
		}
		if resp.StatusCode != http.StatusOK {
			return "", fmt.Errorf("failed to download file %s : got http status code %v", file.FileURL, resp.StatusCode)
		}
		defer func() {
			err = errors.Join(resp.Body.Close(), err)
		}()

		dest, err := os.Create(destPath)
		if err != nil {
			return "", fmt.Errorf("failed to create destination file %s : %v", destPath, err)
		}
		defer func() {
			err = errors.Join(dest.Close(), err)
		}()

		_, err = io.Copy(dest, resp.Body)
		if err != nil {
			return "", fmt.Errorf("failed to download file %s : %v", destPath, err)
		}
	}

	data, err := os.ReadFile(destPath) // os.ReadFile is simpler than Open + ReadAll
	if err != nil {
		return "", fmt.Errorf("failed to read file %s: %w", destPath, err)
	}

	// 4. Encode the data to Base64 and format the data URI string
	encodedData := base64.StdEncoding.EncodeToString(data)
	return fmt.Sprintf("data:%s;base64,%s", file.Mimetype, encodedData), nil
}

func (a *App) ListImports() ([]models.SlackImport, error) {
	return a.repository.ListImports(a.ctx)
}

func (a *App) DeleteImport(importID int64) error {
	selection, err := runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
		Type:          runtime.QuestionDialog,
		Title:         "Confirm Delete",
		Message:       "Are you sure you want to delete this import? This action cannot be undone.",
		Buttons:       []string{"Cancel", "Delete"},
		DefaultButton: "Cancel",
	})
	if err != nil {
		return err
	}
	if selection != "Delete" {
		return fmt.Errorf("user cancelled deletion")
	}
	err = a.repository.DeleteImport(a.ctx, importID)

	if err == nil {
		runtime.EventsEmit(a.ctx, "importDeleted", importID)
	}
	return err
}

func (a *App) ResumeImport(importID int64) error {
	imp, err := a.repository.GetImport(a.ctx, importID)
	if err != nil {
		return fmt.Errorf("cannot find import %d: %v", importID, err)
	}

	a.processImport(imp)
	return nil
}

func (a *App) RenameImport(importID int64, name string) error {
	return a.repository.RenameImport(a.ctx, importID, name)
}

func (a *App) DownloadBinaries(importID int64, confirm bool) error {
	if confirm {
		selection, err := runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Type:          runtime.QuestionDialog,
			Title:         "Download binaries",
			Message:       "Are you sure you want to download all workspace binaries? This action can take a lot of disk space.",
			Buttons:       []string{"Cancel", "OK"},
			DefaultButton: "OK",
		})
		if err != nil {
			return err
		}
		if selection != "OK" {
			return fmt.Errorf("user cancelled deletion")
		}
	}

	currentImport, err := a.repository.GetImport(a.ctx, importID)
	if err != nil {
		return fmt.Errorf("cannot find import %d: %v", importID, err)
	}

	// TODO: check the task correctly, keep a task list on the app struct
	//if currentImport.DownloadStarted.Valid {
	//	return fmt.Errorf("download already in progress")
	//}

	// run this in background
	go func() {
		importer.DownloadBinaries(a.ctx, a.filesFolder, a.repository, currentImport)
	}()
	return nil
}
