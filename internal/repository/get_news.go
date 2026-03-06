package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/yaBliznyk/newsportal/internal/domain"
	"github.com/yaBliznyk/newsportal/internal/svcerrs"
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

	var dao NewsDAO
	err := r.db.QueryRow(ctx, query, args).Scan(
		&dao.ID,
		&dao.Title,
		&dao.Preamble,
		&dao.Content,
		&dao.CategoryID,
		&dao.Category,
		&dao.TagIDs,
		&dao.Author,
		&dao.CreatedAt,
		&dao.PublishedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, svcerrs.ErrDataNotFound
		}
		return nil, fmt.Errorf("failed to get news: %w", err)
	}

	tags, err := r.getTagsByIDs(ctx, dao.TagIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to get tags: %w", err)
	}

	news := dao.ToDomain(tags)
	return &news, nil
}
