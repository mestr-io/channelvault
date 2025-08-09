package app

import "github.com/mestr-io/slack-export-viewer/backend/models"

func (a *App) GetChannels() ([]models.Channel, error) {
	return a.repository.GetChannels()
}

func (a *App) GetMessages(channelID string, cursor string) ([]models.Message, error) {
	return a.repository.GetMessages(a.ctx, channelID, cursor, 100)
}
