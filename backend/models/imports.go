package models

import "database/sql"

type SlackImport struct {
	ID               int64
	Path             string
	Created          string
	Finished         string
	NumFiles         int64
	Name             string
	Error            string
	PendingBinaries  int64
	DownloadStarted  sql.NullString
	DownloadFinished sql.NullString
}
