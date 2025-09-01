package repository

import (
	"context"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/mestr-io/slack-export-viewer/backend/models"
)

func TestBulkInsertMessages(t *testing.T) {
	repo := setup(t)
	defer repo.db.Close()
	ctx := context.Background()
	imp := models.SlackImport{}
	imp.ID = 66
	ctx = context.WithValue(ctx, "import", &imp)

	// Pre-insert a user and a channel for foreign key constraints
	repo.BulkInsertUsers(ctx, []models.User{{ID: "U1", Name: "user1"}})
	repo.BulkInsertChannels(ctx, []models.Channel{{ID: "C1", Name: "channel1"}})

	messages := []models.Message{
		{Ts: "123.456", UserId: "U1", Text: "Hello", ThreadTs: "123.456"},
		{Ts: "123.457", UserId: "U1", Text: "World", ThreadTs: "123.456"},
	}

	if err := repo.BulkInsertMessages(ctx, "C1", messages); err != nil {
		t.Fatalf("BulkInsertMessages() error = %v", err)
	}

	var count int
	err := repo.db.QueryRow("SELECT COUNT(*) FROM messages").Scan(&count)
	if err != nil {
		t.Fatalf("failed to query message count: %v", err)
	}
	if count != 2 {
		t.Errorf("expected 2 messages, got %d", count)
	}
}

func TestGetMessages(t *testing.T) {
	repo := setup(t)
	defer repo.db.Close()
	ctx := context.Background()
	imp := models.SlackImport{}
	imp.ID = 66
	ctx = context.WithValue(ctx, "import", &imp)

	repo.BulkInsertUsers(ctx, []models.User{{ID: "U1", Name: "user1", RealName: "User One"}})
	repo.BulkInsertChannels(ctx, []models.Channel{{ID: "C1", Name: "channel1"}})
	messages := []models.Message{
		{Ts: "123.457", UserId: "U1", Text: "Not in thread"},
		{Ts: "123.456", UserId: "U1", Text: "In thread", ThreadTs: "123.456"},
	}

	//note channel name is used on the insert and not channel ID
	repo.BulkInsertMessages(ctx, "channel1", messages)

	retrieved, err := repo.GetMessages(context.Background(), "C1", "", 10)
	if err != nil {
		t.Fatalf("GetMessages() error = %v", err)
	}

	if len(retrieved) != 2 {
		t.Fatalf("expected 2 message, got %d", len(retrieved))
	}

	//GetMessages newest first
	if retrieved[0].Text != "Not in thread" {
		t.Errorf("unexpected message text: %s", retrieved[0].Text)
	}

	if retrieved[1].Text != "In thread" {
		t.Errorf("unexpected message text: %s", retrieved[1].Text)
	}
}

func TestGetThreadByTs(t *testing.T) {
	repo := setup(t)
	defer repo.db.Close()
	ctx := context.Background()
	imp := models.SlackImport{}
	imp.ID = 66
	ctx = context.WithValue(ctx, "import", &imp)

	repo.BulkInsertUsers(ctx, []models.User{{ID: "U1", Name: "user1", RealName: "User One"}})
	repo.BulkInsertChannels(ctx, []models.Channel{{ID: "C1", Name: "channel1"}})
	messages := []models.Message{
		{Ts: "123.455", UserId: "U1", Text: "Outside thread"},
		{Ts: "123.457", UserId: "U1", Text: "Reply", ThreadTs: "123.456"},
		{Ts: "123.456", UserId: "U1", Text: "Main message", ThreadTs: "123.456"},
		{Ts: "123.458", UserId: "U1", Text: "Outside thread 2"},
	}

	//note channel name is used on the insert and not channel ID
	repo.BulkInsertMessages(ctx, "channel1", messages)

	retrieved, err := repo.GetThreadByTs(context.Background(), "123.456", "C1")
	if err != nil {
		t.Fatalf("GetThreadByTs() error = %v", err)
	}
	if len(retrieved) != 2 {
		t.Fatalf("expected 2 message in thread, got %d", len(retrieved))
	}

	//GetThreadByTs oldest first
	if retrieved[0].Text != "Main message" {
		t.Errorf("unexpected message text: %s", retrieved[0].Text)
	}
}

func TestGetFileById(t *testing.T) {
	repo := setup(t)
	defer repo.db.Close()
	ctx := context.Background()
	imp := models.SlackImport{}
	imp.ID = 66
	ctx = context.WithValue(ctx, "import", &imp)

	repo.BulkInsertChannels(ctx, []models.Channel{{ID: "C1", Name: "channel1"}})
	messages := []models.Message{
		{Ts: "123.456", Files: []models.File{{ID: "F1", Name: "file1.txt"}}},
	}
	repo.BulkInsertMessages(ctx, "C1", messages)

	file, err := repo.GetFileById(ctx, "F1")
	if err != nil {
		t.Fatalf("GetFileById() error = %v", err)
	}

	if file.Name != "file1.txt" {
		t.Errorf("unexpected file name: %s", file.Name)
	}
}
