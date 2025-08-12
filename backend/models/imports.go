package models

import "time"

type Import struct {
	ID          int64
	Path        string
	Created     time.Time
	Finished    time.Time
	NumMessages int64
}
