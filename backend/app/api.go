package app

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/mestr-io/slack-export-viewer/backend/models"
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

func (a *App) GetBase64Image(fileId string) (string, error) {
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
		defer resp.Body.Close()

		dest, err := os.Create(destPath)

		if err != nil {
			return "", fmt.Errorf("failed to create destination file %s : %v", destPath, err)
		}
		defer dest.Close()

		_, err = io.Copy(dest, resp.Body)
		if err != nil {
			return "", fmt.Errorf("failed to download file %s : %v", destPath, err)
		}

		// generate thumbnails 400x?
	}

	data, err := os.ReadFile(destPath) // os.ReadFile is simpler than Open + ReadAll
	if err != nil {
		return "", fmt.Errorf("failed to read file %s: %w", destPath, err)
	}

	// 4. Encode the data to Base64 and format the data URI string
	encodedData := base64.StdEncoding.EncodeToString(data)
	return fmt.Sprintf("data:%s;base64,%s", file.Mimetype, encodedData), nil
}

func (a *App) ListImports() ([]models.Import, error) {
	return a.repository.ListImports(a.ctx)
}
