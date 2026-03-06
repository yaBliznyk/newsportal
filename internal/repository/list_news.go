package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/yaBliznyk/newsportal/internal/domain"
)

func (r *NewsRepository) ListNews(ctx context.Context, req domain.ListNewsReq) ([]domain.ListNewsItem, error) {
	query := `
		SELECT n."newsId", n."title",
		       c."categoryId", c."name",
		       n."tagIds", n."author", n."createdAt", n."publishedAt"
		FROM "news" n
		JOIN "categories" c ON n."categoryId" = c."categoryId"
		WHERE n."statusId" = @statusID
	`
	args := pgx.NamedArgs{"statusID": domain.StatusPublished}

	if req.CategoryID != 0 {
		query += ` AND n."categoryId" = @categoryID`
		args["categoryID"] = req.CategoryID
	}

	if req.TagID != 0 {
		query += ` AND @tagID = ANY(n."tagIds")`
		args["tagID"] = req.TagID
	}

	if !req.From.IsZero() {
		query += ` AND n."publishedAt" >= @from`
		args["from"] = req.From
	}

	if !req.To.IsZero() {
		query += ` AND n."publishedAt" <= @to`
		args["to"] = req.To
	}

	query += ` ORDER BY n."publishedAt" DESC`

	if req.Limit > 0 {
		query += ` LIMIT @limit`
		args["limit"] = req.Limit
	}

	if req.Page > 0 && req.Limit > 0 {
		query += ` OFFSET @offset`
		args["offset"] = (req.Page - 1) * req.Limit
	}

	rows, err := r.db.Query(ctx, query, args)
	if err != nil {
		return nil, fmt.Errorf("failed to list news: %w", err)
	}
	defer rows.Close()

	var newsList []domain.ListNewsItem
	for rows.Next() {
		var item domain.ListNewsItem
		var tagIDs []int32
		err := rows.Scan(
			&item.ID,
			&item.Title,
			&item.Category.ID,
			&item.Category.Name,
			&tagIDs,
			&item.Author,
			&item.CreatedAt,
			&item.PublishedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan news item: %w", err)
		}

		if len(tagIDs) > 0 {
			item.Tags, err = r.getTagsByIDs(ctx, tagIDs)
			if err != nil {
				return nil, fmt.Errorf("failed to get tags: %w", err)
			}
		}

		newsList = append(newsList, item)
	}

	return newsList, rows.Err()
}
