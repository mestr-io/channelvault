package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"reflect"
	"testing"

	"github.com/mestr-io/slack-export-viewer/backend/models"
	_ "modernc.org/sqlite"
)

func setup(t *testing.T) *RepositorySql {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("failed to open in-memory database: %v", err)
	}

	repo := &RepositorySql{db: db}
	if err := repo.SetupSchemas(); err != nil {
		t.Fatalf("failed to setup schemas: %v", err)
	}

	return repo
}

func TestBulkInsertChannels(t *testing.T) {
	repo := setup(t)
	defer func() {
		err := repo.db.Close()
		if err != nil {
			panic(err)
		}
	}()
	ctx := context.Background()
	imp := models.SlackImport{}
	imp.ID = 66
	ctx = context.WithValue(ctx, "import", &imp)

	channels := []models.Channel{
		{ID: "C1", Name: "channel1", Members: []string{"U1"}, Pins: []models.ChannelPin{{ID: "p1"}}},
		{ID: "C2", Name: "channel2", Topic: models.ChannelTopic{Value: "topic2"}, Purpose: models.ChannelPurpose{Value: "purpose2"}},
	}

	if err := repo.BulkInsertChannels(ctx, channels); err != nil {
		t.Fatalf("BulkInsertChannels() error = %v", err)
	}

	var count int
	err := repo.db.QueryRow("SELECT COUNT(*) FROM channels").Scan(&count)
	if err != nil {
		t.Fatalf("failed to query channel count: %v", err)
	}
	if count != 2 {
		t.Errorf("expected 2 channels, got %d", count)
	}

	var id, name, topic, purpose string
	var members, pins []byte
	err = repo.db.QueryRow("SELECT id, name, topic, purpose, members_json, pins_json FROM channels WHERE id = 'C1'").Scan(&id, &name, &topic, &purpose, &members, &pins)
	if err != nil {
		t.Fatalf("failed to query channel C1: %v", err)
	}

	if id != "C1" || name != "channel1" {
		t.Errorf("unexpected data for channel C1")
	}

	var membersSlice []string
	if err := json.Unmarshal(members, &membersSlice); err != nil {
		t.Fatalf("failed to unmarshal members: %v", err)
	}
	if !reflect.DeepEqual(membersSlice, []string{"U1"}) {
		t.Errorf("expected members %v, got %v", []string{"U1"}, membersSlice)
	}
}

func TestGetChannels(t *testing.T) {
	repo := setup(t)
	defer func() {
		err := repo.db.Close()
		if err != nil {
			panic(err)
		}
	}()
	ctx := context.Background()
	imp := models.SlackImport{}
	imp.ID = 66
	ctx = context.WithValue(ctx, "import", &imp)

	channels := []models.Channel{
		{ID: "C1", Name: "channel1"},
		{ID: "C2", Name: "channel2"},
	}
	if err := repo.BulkInsertChannels(ctx, channels); err != nil {
		t.Fatalf("failed to insert channels for testing: %v", err)
	}

	retrievedChannels, err := repo.GetChannels()
	if err != nil {
		t.Fatalf("GetChannels() error = %v", err)
	}

	if len(retrievedChannels) != 2 {
		t.Fatalf("expected 2 channels, got %d", len(retrievedChannels))
	}

	expected := []models.Channel{
		{ID: "C1", Name: "channel1"},
		{ID: "C2", Name: "channel2"},
	}

	// The order is not guaranteed, so we need to check for existence
	for _, exp := range expected {
		found := false
		for _, got := range retrievedChannels {
			if exp.ID == got.ID && exp.Name == got.Name {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected channel %v not found", exp)
		}
	}
}
