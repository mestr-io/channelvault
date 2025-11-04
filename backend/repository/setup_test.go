package repository

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func TestSetupSchemas(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("failed to open in-memory database: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			panic(err)
		}
	}()

	repo := &RepositorySql{db: db}

	if err := repo.SetupSchemas(); err != nil {
		t.Fatalf("SetupSchemas() error = %v", err)
	}

	tables := []string{"imports", "files", "messages", "channels", "users"}
	for _, table := range tables {
		var name string
		err := db.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name=?", table).Scan(&name)
		if err != nil {
			if err == sql.ErrNoRows {
				t.Errorf("table %s was not created", table)
			} else {
				t.Errorf("error checking for table %s: %v", table, err)
			}
		}
	}
}
