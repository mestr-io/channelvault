package repository

import (
	"encoding/json"
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

	stmt, err := tx.Prepare("INSERT OR REPLACE INTO channels (id, name, topic, purpose, members_json, pins_json) VALUES (?, ?, ?, ?, ?, ?)")
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	insertCount := 0
	for _, channel := range channels {

		membersJson, err := json.Marshal(channel.Members)
		if err != nil {
			return fmt.Errorf("failed to marshal members for channel %s: %v", channel.ID, err)
		}
		pinsJson, err := json.Marshal(channel.Pins)
		if err != nil {
			return fmt.Errorf("failed to marshal pins for channel %s: %v", channel.ID, err)
		}
		_, err = stmt.Exec(channel.ID, channel.Name, channel.Topic.Value, channel.Purpose.Value, membersJson, pinsJson)
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
	defer rows.Close()

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
