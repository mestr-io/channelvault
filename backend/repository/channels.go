package repository

import (
	"errors"
	"fmt"
	"log"
	"github.com/mestr-io/slack-export-viewer/backend/app_errors"
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
	
	defer func(){ err = app_errors.Finally(stmt.Close, err) }()

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

func (r *Repository) GetChannels() (channels []models.Channel, err error) {
	rows, err := r.db.Query("select id, name from channels")
	if err != nil {
		return nil, fmt.Errorf("FATAL: couldnt query channels: %v", err)
	}
	defer func(){ err = app_errors.Finally(rows.Close, err) }()

	for rows.Next() {
		var channel models.Channel
		err = rows.Scan(&channel.ID, &channel.Name)
		if err != nil {
			return nil, fmt.Errorf("FATAL: couldnt query channels: %v", err)
		}
		channels = append(channels, channel)
	}
	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("FATAL: couldnt query channels: %v", err)
	}
	return channels, nil
}
