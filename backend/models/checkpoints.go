package models

type Checkpoint struct {
	ImportID    int64
	Path        string
	ProcessedAt string
	Error       string
}
