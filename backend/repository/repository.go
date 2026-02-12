package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/mestr-io/slack-export-viewer/backend/models"
)

type Repository interface {
	//channels
	GetChannels() (channels []models.Channel, err error)
	BulkInsertChannels(ctx context.Context, channels []models.Channel) (err error)

	//imports
	MarkCheckpointAsCompleted(ctx context.Context, importID int64, path string, errMsg string) error
	CreateCheckpointPlaceHolders(ctx context.Context, importID int64, paths []string) error
	GetImport(ctx context.Context, id int64) (*models.SlackImport, error)
	FinishedImportUpdate(ctx context.Context, imp *models.SlackImport) error
	FinishedImportUpdateWithError(ctx context.Context, imp *models.SlackImport) error
	DeleteImport(ctx context.Context, id int64) error
	RenameImport(ctx context.Context, importID int64, name string) error
	CreateImport(ctx context.Context, path string, name string) (*models.SlackImport, error)
	NoFinishedImports(ctx context.Context) (bool, error)
	ListImports(ctx context.Context) (imports []models.SlackImport, err error)
	ListPendingFiles(ctx context.Context, importID int64) ([]string, error)
	StartDownload(ctx context.Context, importID int64) error
	FinishDownload(ctx context.Context, importID int64) error

	//messages
	BulkInsertMessages(ctx context.Context, ChannelID string, messages []models.Message) (err error)
	GetMessages(ctx context.Context, channelID string, cursor string, pageSize int) ([]models.Message, error)
	ListChannelFiles(ctx context.Context, channelID string, pageSize int, cursor int) ([]models.File, error)
	GetThreadByTs(ctx context.Context, threadTs string, channelID string, importID int64) ([]models.Message, error)
	GetFileById(ctx context.Context, fileId string) (*models.File, error)
	ListPendingBinaries(ctx context.Context, importID int64, pageSize int) ([]models.File, error)
	MarkFileAsDownloaded(ctx context.Context, importID int64, fileID string) error
	CountPendingBinaries(ctx context.Context, importID int64) (int64, error)

	// search
	SearchMessages(ctx context.Context, ChannelID string, userId string, query string) ([]models.SearchResult, error)

	//setup
	SetupSchemas() error

	//users
	BulkInsertUsers(ctx context.Context, users []models.User) (err error)
	GetUserAvatar(ctx context.Context, userId string) (string, error)
	GetUserHandle(ctx context.Context, userId string) (string, error)
}

type RepositorySql struct {
	db *sql.DB
}

func New(dbPath string) (*RepositorySql, error) {
	db, err := sql.Open("sqlite3", fmt.Sprintf("%s?_foreign_keys=on&_journal_mode=WAL&busy_timeout=10000&Synchronous=NORMAL", dbPath))
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
