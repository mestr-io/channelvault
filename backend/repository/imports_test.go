package repository

import (
	"context"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func TestStoreImport(t *testing.T) {
	repo := setup(t)
	defer repo.db.Close()

	imp, err := repo.CreateImport(context.Background(), "/path/to/import", "an import")
	if err != nil {
		t.Fatalf("StoreImport() error = %v", err)
	}

	if imp.ID != 1 {
		t.Errorf("expected import ID 1, got %d", imp.ID)
	}

	var path string
	err = repo.db.QueryRow("SELECT path FROM imports WHERE id = 1").Scan(&path)
	if err != nil {
		t.Fatalf("failed to query import: %v", err)
	}
	if path != "/path/to/import" {
		t.Errorf("expected path /path/to/import, got %s", path)
	}
}

func TestFinishedImport(t *testing.T) {
	repo := setup(t)
	defer repo.db.Close()

	imp, err := repo.CreateImport(context.Background(), "/path/to/import", "an import")
	if err != nil {
		t.Fatalf("StoreImport() error = %v", err)
	}

	if err := repo.FinishedImportUpdate(context.Background(), imp); err != nil {
		t.Fatalf("FinishedImport() error = %v", err)
	}

	var finished time.Time
	err = repo.db.QueryRow("SELECT finished FROM imports WHERE id = ?", imp.ID).Scan(&finished)
	if err != nil {
		t.Fatalf("failed to query import: %v", err)
	}
	if finished.IsZero() {
		t.Error("expected finished time to be set")
	}
}

func TestNoFinishedImports(t *testing.T) {
	repo := setup(t)
	defer repo.db.Close()

	t.Run("No imports", func(t *testing.T) {
		any, err := repo.NoFinishedImports(context.Background())
		if err != nil {
			t.Fatalf("NoFinishedImports() error = %v", err)
		}
		if any {
			t.Error("expected false when no imports exist")
		}
	})

	imp, err := repo.CreateImport(context.Background(), "/path/to/import", "an import")
	if err != nil {
		t.Fatalf("StoreImport() error = %v", err)
	}

	t.Run("Unfinished import", func(t *testing.T) {
		any, err := repo.NoFinishedImports(context.Background())
		if err != nil {
			t.Fatalf("NoFinishedImports() error = %v", err)
		}
		if any {
			t.Error("expected false when only unfinished imports exist")
		}
	})

	if err := repo.FinishedImportUpdate(context.Background(), imp); err != nil {
		t.Fatalf("FinishedImport() error = %v", err)
	}

	t.Run("Finished import", func(t *testing.T) {
		any, err := repo.NoFinishedImports(context.Background())
		if err != nil {
			t.Fatalf("NoFinishedImports() error = %v", err)
		}
		if !any {
			t.Error("expected true when a finished import exists")
		}
	})
}
