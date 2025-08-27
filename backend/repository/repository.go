package repository

import (
	"database/sql"
	"fmt"
)

type Repository struct {
	db *sql.DB
}

func New(dbPath string) (*Repository, error) {
	db, err := sql.Open("sqlite3", fmt.Sprintf("%s?_foreign_keys=on", dbPath))
	if err != nil {
		return nil, fmt.Errorf("FATAL: could not determine database path: %v", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("FATAL: fail to open database: %v", err)
	}

	r := Repository{db: db}
	err = r.SetupSchemas()
	if err != nil {
		return nil, fmt.Errorf("FATAL: cannot create schemas: %v", err)
	}
	return &r, nil
}
