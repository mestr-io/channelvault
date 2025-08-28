package models

type Progress struct {
	ImportID    int64
	FileName    string
	Progress    int
	TotalFiles  int
	CurrentFile int
}
