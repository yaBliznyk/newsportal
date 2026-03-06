package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/yaBliznyk/newsportal/internal/domain"
)

func (r *NewsRepository) GetNews(ctx context.Context, id int) (*domain.News, error) {
	const query = `
		SELECT n."newsId", n."title", n."preamble", n."content",
		       c."categoryId", c."name",
		       n."tagIds", n."author", n."createdAt", n."publishedAt"
		FROM "news" n
		JOIN "categories" c ON n."categoryId" = c."categoryId"
		WHERE n."newsId" = @newsID
		  AND n."statusId" = @statusID
	`

	args := pgx.NamedArgs{
		"newsID":   id,
		"statusID": domain.StatusPublished,
	}

	var news domain.News
	var tagIDs []int32
	err := r.db.QueryRow(ctx, query, args).Scan(
		&news.ID,
		&news.Title,
		&news.Preamble,
		&news.Content,
		&news.Category.ID,
		&news.Category.Name,
		&tagIDs,
		&news.Author,
		&news.CreatedAt,
		&news.PublishedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get news: %w", err)
	}

	if len(tagIDs) > 0 {
		news.Tags, err = r.getTagsByIDs(ctx, tagIDs)
		if err != nil {
			return nil, fmt.Errorf("failed to get tags: %w", err)
		}
	}

	return &news, nil
}

func (r *NewsRepository) getTagsByIDs(ctx context.Context, tagIDs []int32) ([]domain.Tag, error) {
	query := `
		SELECT "tagId", "name"
		FROM "tags"
		WHERE "tagId" = ANY(@tagIDs)
		  AND "statusId" = @statusID
		ORDER BY "name"
	`

	args := pgx.NamedArgs{
		"tagIDs":   tagIDs,
		"statusID": domain.StatusPublished,
	}

	rows, err := r.db.Query(ctx, query, args)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []domain.Tag
	for rows.Next() {
		var t domain.Tag
		if err := rows.Scan(&t.ID, &t.Name); err != nil {
			return nil, err
		}
		tags = append(tags, t)
	}

	return tags, rows.Err()
}
