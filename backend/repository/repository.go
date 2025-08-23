package repository

import (
	"database/sql"
	"context"	
	"fmt"

	"github.com/mestr-io/slack-export-viewer/backend/models"
)

type Repository interface {
	//channels
	BulkInsertChannels(channels []models.Channel) (err error)
	GetChannels() (channels []models.Channel, err error)

	//imports
	StoreImport(ctx context.Context, path string) (*models.Import, error)
	FinishedImport(ctx context.Context, ID int64) error	
	NoFinishedImports(ctx context.Context) (bool, error)

	//messages
	BulkInserMessages(ctx context.Context, ChannelID string, messages []models.Message) (err error)
	GetMessages(ctx context.Context, channelID string, cursor string, pageSize int) ([]models.Message, error)
	ListChannelFiles(ctx context.Context, channelID string, pageSize int, cursor int) ([]models.File, error)
	GetThreadByTs(ctx context.Context, threadTs string, channelID string) ([]models.Message, error)
	GetFileById(ctx context.Context, fileId string) (*models.File, error)

	//setup
	SetupSchemas() error

	//users
	BulkInsertUsers(users []models.User) (err error)	
	GetUserAvatar(ctx context.Context, userId string) (string, error)	
}

type RepositorySql struct {
	db *sql.DB
}

func New(dbPath string) (*RepositorySql, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("FATAL: could not determine database path: %v", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("FATAL: fail to open database: %v", err)
	}

	r := RepositorySql{db: db}
	err = r.SetupSchemas()
	if err != nil {
		return nil, fmt.Errorf("FATAL: cannot create schemas: %v", err)
	}
	return &r, nil
}
