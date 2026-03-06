package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/yaBliznyk/newsportal/internal/domain"
)

func (r *NewsRepository) CountNews(ctx context.Context, req domain.CountNewsReq) (int, error) {
	query := `
		SELECT COUNT(*)
		FROM "news"
		WHERE "statusId" = @statusID
	`
	args := pgx.NamedArgs{"statusID": domain.StatusPublished}

	if req.CategoryID != 0 {
		query += ` AND "categoryId" = @categoryID`
		args["categoryID"] = req.CategoryID
	}

	if req.TagID != 0 {
		query += ` AND @tagID = ANY("tagIds")`
		args["tagID"] = req.TagID
	}

	if !req.From.IsZero() {
		query += ` AND "publishedAt" >= @from`
		args["from"] = req.From
	}

	if !req.To.IsZero() {
		query += ` AND "publishedAt" <= @to`
		args["to"] = req.To
	}

	var count int
	err := r.db.QueryRow(ctx, query, args).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count news: %w", err)
	}

	return count, nil
}
