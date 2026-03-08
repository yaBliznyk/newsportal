package db

import (
	"context"

	"github.com/jackc/pgx/v5"

	"github.com/yaBliznyk/newsportal/internal/domain"
)

func (r *NewsRepo) GetCategories(ctx context.Context) ([]domain.Category, error) {
	const query = `
		SELECT "categoryId", "name"
		FROM "categories"
		WHERE "statusId" = @statusID
		ORDER BY "sortOrder", "name"
	`

	args := pgx.NamedArgs{"statusID": domain.StatusPublished}

	rows, err := r.db.Query(ctx, query, args)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []domain.Category
	for rows.Next() {
		var c domain.Category
		if err := rows.Scan(&c.ID, &c.Name); err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}

	return categories, rows.Err()
}
