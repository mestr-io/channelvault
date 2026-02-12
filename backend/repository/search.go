package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/mestr-io/slack-export-viewer/backend/models"
)

func (r *RepositorySql) SearchMessages(ctx context.Context, ChannelID string, userId string, query string) ([]models.SearchResult, error) {
	rows, err := r.db.Query("select distinct(id), name from channels")
	if err != nil {
		return nil, fmt.Errorf("FATAL: couldnt query channels: %w", err)
	}
	defer func() {
		err = errors.Join(rows.Close(), err)
	}()
	return nil, nil
}
