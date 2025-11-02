package repository

import (
	"context"
	"database/sql"
	"testing"

	"github.com/mestr-io/slack-export-viewer/backend/models"
	_ "modernc.org/sqlite"
)

func TestBulkInsertUsers(t *testing.T) {
	repo := setup(t)
	defer func() {
		if err := repo.db.Close(); err != nil {
			panic(err)
		}
	}()
	ctx := context.Background()
	imp := models.SlackImport{}
	imp.ID = 66
	ctx = context.WithValue(ctx, "import", &imp)
	users := []models.User{
		{ID: "U1", Name: "user1", RealName: "User One", Profile: models.Profile{Image72: "http://example.com/avatar1.png"}},
		{ID: "U2", Name: "user2", RealName: "User Two", IsBot: true},
	}

	if err := repo.BulkInsertUsers(ctx, users); err != nil {
		t.Fatalf("BulkInsertUsers() error = %v", err)
	}

	var count int
	err := repo.db.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
	if err != nil {
		t.Fatalf("failed to query user count: %v", err)
	}
	if count != 2 {
		t.Errorf("expected 2 users, got %d", count)
	}

	var name, realName, profileImageUrl string
	var isBot bool
	err = repo.db.QueryRow("SELECT name, real_name, profile_image_url, is_bot FROM users WHERE id = 'U1'").Scan(&name, &realName, &profileImageUrl, &isBot)
	if err != nil {
		t.Fatalf("failed to query user U1: %v", err)
	}

	if name != "user1" || realName != "User One" || profileImageUrl != "http://example.com/avatar1.png" || isBot {
		t.Errorf("unexpected data for user U1")
	}
}

func TestGetUserAvatar(t *testing.T) {
	repo := setup(t)
	defer func() {
		if err := repo.db.Close(); err != nil {
			panic(err)
		}
	}()

	ctx := context.Background()
	imp := models.SlackImport{}
	imp.ID = 66
	ctx = context.WithValue(ctx, "import", &imp)

	users := []models.User{
		{ID: "U1", Name: "user1", Profile: models.Profile{Image72: "http://example.com/avatar1.png"}},
	}
	if err := repo.BulkInsertUsers(ctx, users); err != nil {
		t.Fatalf("failed to insert users for testing: %v", err)
	}

	avatar, err := repo.GetUserAvatar(ctx, "U1")
	if err != nil {
		t.Fatalf("GetUserAvatar() error = %v", err)
	}

	if avatar != "http://example.com/avatar1.png" {
		t.Errorf("expected avatar %s, got %s", "http://example.com/avatar1.png", avatar)
	}

	_, err = repo.GetUserAvatar(context.Background(), "U2")
	if err != sql.ErrNoRows {
		t.Errorf("expected sql.ErrNoRows for non-existent user, got %v", err)
	}
}
