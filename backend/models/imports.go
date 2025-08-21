package models

type SlackImport struct {
	ID          int64
	Path        string
	Created     string
	Finished    string
	NumMessages int64
	Name        string
	Error       string
}
