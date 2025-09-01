package models

type SlackImport struct {
	ID                       int64
	Path                     string
	Created                  string
	Finished                 string
	NumFiles                 int64
	Name                     string
	Error                    string
	PendingBinaries          int64
	BinariesDownloadStarted  string
	BinariesDownloadFinished string
}
