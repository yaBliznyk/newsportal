package db

import (
	"context"
	"fmt"
	"sort"

	"github.com/jackc/pgx/v5"

	"github.com/yaBliznyk/newsportal/internal/domain"
)

func (r *NewsRepo) GetTags(ctx context.Context) ([]domain.Tag, error) {
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
		var dao TagDAO
		if err := rows.Scan(&dao.ID, &dao.Name); err != nil {
			return nil, fmt.Errorf("failed to scan tag: %w", err)
		}
		tags = append(tags, dao.ToDomain())
	}

	return tags, rows.Err()
}

func (r *NewsRepo) getTagsByIDs(ctx context.Context, tagIDs []int32) ([]domain.Tag, error) {
	tagsMap, err := r.getTagsMapByIDs(ctx, tagIDs)
	if err != nil {
		return nil, err
	}

	tags := make([]domain.Tag, 0, len(tagsMap))
	for _, tag := range tagsMap {
		tags = append(tags, tag)
	}

	// Сортировка по идентификатору для стабильного порядка
	sort.Slice(tags, func(i, j int) bool {
		return tags[i].ID < tags[j].ID
	})

	return tags, nil
}

// getTagsMapByIDs возвращает мапу тегов по их ID (для избежания N+1 запросов)
func (r *NewsRepo) getTagsMapByIDs(ctx context.Context, tagIDs []int32) (map[int32]domain.Tag, error) {
	query := `
		SELECT "tagId", "name"
		FROM "tags"
		WHERE "tagId" = ANY(@tagIDs)
		  AND "statusId" = @statusID
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

	tags := make(map[int32]domain.Tag)
	for rows.Next() {
		var dao TagDAO
		if err := rows.Scan(&dao.ID, &dao.Name); err != nil {
			return nil, err
		}
		tags[dao.ID] = dao.ToDomain()
	}

	return tags, rows.Err()
}
