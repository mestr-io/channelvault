package repository

import (
	"errors"
	"fmt"
	"log"

	"github.com/mestr-io/slack-export-viewer/backend/models"
)

func (r *Repository) BulkInsertChannels(channels []models.Channel) (err error) {
	tx, err := r.db.Begin()

	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if err != nil {
			terr := tx.Rollback()
			if terr != nil {
				err = errors.Join(err, terr)
			}
		}
	}()

	stmt, err := tx.Prepare("INSERT OR REPLACE INTO channels (id, name) VALUES (?, ?)")
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	insertCount := 0
	for _, channel := range channels {
		_, err := stmt.Exec(channel.ID, channel.Name)
		if err != nil {
			return fmt.Errorf("failed to insert channel %s: %w", channel.ID, err)
		}
		insertCount++
	}

	log.Printf("Successfully prepared %d channels for insertion.", insertCount)

	return tx.Commit()
}
