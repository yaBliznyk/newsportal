package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/yaBliznyk/newsportal/internal/domain"
)

func (r *NewsRepo) ListNews(ctx context.Context, req domain.ListNewsReq) ([]domain.ListNewsItem, error) {
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

	// Временная структура для хранения DAO с tagIDs
	type newsRow struct {
		dao    NewsListItemDAO
		tagIDs []int32
	}

	var newsRows []newsRow
	allTagIDs := make(map[int32]struct{})

	// Первый проход: собираем DAO и все уникальные tagIDs
	for rows.Next() {
		var dao NewsListItemDAO
		err := rows.Scan(
			&dao.ID,
			&dao.Title,
			&dao.CategoryID,
			&dao.Category,
			&dao.TagIDs,
			&dao.Author,
			&dao.CreatedAt,
			&dao.PublishedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan news item: %w", err)
		}

		newsRows = append(newsRows, newsRow{dao: dao, tagIDs: dao.TagIDs})
		for _, id := range dao.TagIDs {
			allTagIDs[id] = struct{}{}
		}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	// Один запрос для получения всех тегов
	var tagsMap map[int32]domain.Tag
	if len(allTagIDs) > 0 {
		tagIDList := make([]int32, 0, len(allTagIDs))
		for id := range allTagIDs {
			tagIDList = append(tagIDList, id)
		}
		tagsMap, err = r.getTagsMapByIDs(ctx, tagIDList)
		if err != nil {
			return nil, fmt.Errorf("failed to get tags: %w", err)
		}
	}

	// Второй проход: конвертируем DAO в domain модели
	newsList := make([]domain.ListNewsItem, 0, len(newsRows))
	for _, row := range newsRows {
		tags := make([]domain.Tag, 0, len(row.tagIDs))
		for _, tagID := range row.tagIDs {
			if tag, ok := tagsMap[tagID]; ok {
				tags = append(tags, tag)
			}
		}
		newsList = append(newsList, row.dao.ToDomain(tags))
	}

	return newsList, nil
}
