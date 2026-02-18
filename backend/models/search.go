package models

type SearchResult struct {
	ChannelID   string
	Ts          string
	MatchedText string
	UserId      string
	ThreadTs    string
	ImportID    int64
	ProfileUrl  string
	Username    string
}
