package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/mestr-io/slack-export-viewer/backend/models"
)

func (r *RepositorySql) SearchMessages(ctx context.Context, query string) ([]models.SearchResult, error) {

	matchQuery := fmt.Sprintf("text : %s", query)
	q := `select m.import_id, snippet(fts.messages_fts , 0, '<mark>', '</mark>', '...', ?) as hl, m.ts, m.channel_id, m.user_id, m.thread_ts
		from messages m JOIN messages_fts fts ON m.rowid = fts.rowid
		where fts.messages_fts match ?
		order by bm25(fts.messages_fts)
		limit ?;
	`
	args := []interface{}{100, matchQuery, 100}
	rows, err := r.db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, fmt.Errorf("could not search messages: %v", err)
	}
	defer func() { err = errors.Join(rows.Close(), err) }()

	var results []models.SearchResult
	for rows.Next() {
		var res models.SearchResult

		if err := rows.Scan(&res.ImportID, &res.MatchedText, &res.Ts, &res.ChannelID, &res.UserId, &res.ThreadTs); err != nil {
			return nil, fmt.Errorf("could not scan message row: %v", err)
		}
		results = append(results, res)
	}
	return results, nil
}
