package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/yaBliznyk/newsportal/internal/domain"
)

func (r *NewsRepository) GetTags(ctx context.Context) ([]domain.Tag, error) {
	const query = `
		SELECT "tagId", "name"
		FROM "tags"
		WHERE "statusId" = @statusID
		ORDER BY "name"
	`

	args := pgx.NamedArgs{"statusID": domain.StatusPublished}

	rows, err := r.db.Query(ctx, query, args)
	if err != nil {
		return nil, fmt.Errorf("failed to get tags: %w", err)
	}
	defer rows.Close()

	var tags []domain.Tag
	for rows.Next() {
		var t domain.Tag
		if err := rows.Scan(&t.ID, &t.Name); err != nil {
			return nil, fmt.Errorf("failed to scan tag: %w", err)
		}
		tags = append(tags, t)
	}

	return tags, rows.Err()
}
