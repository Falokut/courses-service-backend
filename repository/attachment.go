package repository

import (
	"context"

	"github.com/Falokut/go-kit/client/db"
	"github.com/pkg/errors"
)

type Attachment struct {
	db db.DB
}

func NewAttachment(db db.DB) Attachment {
	return Attachment{
		db: db,
	}
}

func (r Attachment) RemoveOrphanedAttachments(ctx context.Context, maxRecordsDelete int) ([]string, error) {
	const query = `
	DELETE FROM lesson_attachments
    WHERE id IN (
        SELECT id FROM lesson_attachments
        WHERE lesson_id = 0
        ORDER BY id
        LIMIT $1
        FOR UPDATE SKIP LOCKED
    )
    RETURNING url;`
	urls := make([]string, 0)
	err := r.db.Select(ctx, &urls, query, maxRecordsDelete)
	if err != nil {
		return nil, errors.WithMessagef(err, "exec query: %s", query)
	}
	return urls, nil
}
