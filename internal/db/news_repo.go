package db

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-pg/pg/v10"
)

const defaultLimit = 20

// NewsRepo репозиторий новостей
type NewsRepo struct {
	db *pg.DB
}

// NewNewsRepo создаёт экземпляр репозитория новостей
func NewNewsRepo(db *pg.DB) *NewsRepo {
	return &NewsRepo{db: db}
}

// ListNewsByFilter список сокращенных новостей по фильтру
func (r *NewsRepo) ListNewsByFilter(ctx context.Context, filter NewsFilter, pager Pagination) ([]News, error) {
	var news []News

	query := r.db.ModelContext(ctx, &news).
		ColumnExpr(`n."newsId", n."title", n."categoryId", n."tagIds", n."author", n."createdAt", n."publishedAt", n."statusId"`).
		TableExpr("news AS n").
		Relation("Category", func(q *pg.Query) (*pg.Query, error) {
			if filter.CategoryStatusID != StatusUndefined {
				return q.Where("c.statusId = ?", filter.CategoryStatusID), nil
			}
			return q, nil
		})

	// Применяем фильтры
	if filter.StatusID != StatusUndefined {
		query = query.Where(`n."statusId" = ?`, filter.StatusID)
	}

	if filter.CategoryID != 0 {
		query = query.Where(`n."categoryId" = ?`, filter.CategoryID)
	}

	if filter.TagID != 0 {
		query = query.Where(`? = ANY(n."tagIds")`, filter.TagID)
	}

	if !filter.From.IsZero() {
		query = query.Where(`n."publishedAt" >= ?`, filter.From)
	}

	if !filter.To.IsZero() {
		query = query.Where(`n."publishedAt" <= ?`, filter.To)
	}

	// Сортировка
	query = query.Order(`n."publishedAt" DESC`)

	// Пагинация
	limit := pager.Limit
	if limit <= 0 {
		limit = defaultLimit
	}
	query = query.Limit(limit)

	page := pager.Page
	if page <= 0 {
		page = 1
	}
	query = query.Offset((page - 1) * limit)

	err := query.Select()
	if err != nil {
		return nil, fmt.Errorf("failed to list news: %w", err)
	}

	return news, nil
}

// CountNews количество новостей по фильтру
func (r *NewsRepo) CountNews(ctx context.Context, filter NewsFilter) (int, error) {
	query := r.db.ModelContext(ctx, (*News)(nil)).
		TableExpr("news").
		Relation("Category", func(q *pg.Query) (*pg.Query, error) {
			if filter.CategoryStatusID != StatusUndefined {
				return q.Where("c.statusId = ?", filter.CategoryStatusID), nil
			}
			return q, nil
		})

	// Применяем фильтры
	if filter.StatusID != StatusUndefined {
		query = query.Where(`"statusId" = ?`, filter.StatusID)
	}

	if filter.CategoryID != 0 {
		query = query.Where(`"categoryId" = ?`, filter.CategoryID)
	}

	if filter.TagID != 0 {
		query = query.Where(`? = ANY("tagIds")`, filter.TagID)
	}

	if !filter.From.IsZero() {
		query = query.Where(`"publishedAt" >= ?`, filter.From)
	}

	if !filter.To.IsZero() {
		query = query.Where(`"publishedAt" <= ?`, filter.To)
	}

	count, err := query.Count()
	if err != nil {
		return 0, fmt.Errorf("failed to count news: %w", err)
	}

	return count, nil
}

// NewsByIDAndStatus получение полной новости по идентификатору и статусу
func (r *NewsRepo) NewsByIDAndStatus(ctx context.Context, id int, statusID, categoryStatusID Status) (*News, error) {
	news := &News{}

	query := r.db.ModelContext(ctx, news).
		Where(`"newsId" = ?`, id).
		Relation("Category", func(q *pg.Query) (*pg.Query, error) {
			if categoryStatusID != StatusUndefined {
				return q.Where("c.statusId = ?", categoryStatusID), nil
			}
			return q, nil
		})

	if statusID != StatusUndefined {
		query = query.Where(`"statusId" = ?`, statusID)
	}

	err := query.Select()
	if errors.Is(err, pg.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("failed to get news: %w", err)
	}

	return news, nil
}

// GetCategoryByIDAndStatusID получение одной категории по идентификатору и статусу
func (r *NewsRepo) GetCategoryByIDAndStatusID(ctx context.Context, id int, statusID Status) (*Category, error) {
	category := &Category{}

	query := r.db.ModelContext(ctx, category).Where(`"categoryId" = ?`, id)

	if statusID != StatusUndefined {
		query = query.Where(`"statusId" = ?`, statusID)
	}

	err := query.Select()
	if errors.Is(err, pg.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("failed to get category: %w", err)
	}

	return category, nil
}

// GetCategoriesByStatusID получение списка категорий по статусу
func (r *NewsRepo) GetCategoriesByStatusID(ctx context.Context, statusID Status) ([]Category, error) {
	var categories []Category

	query := r.db.ModelContext(ctx, &categories)

	if statusID != StatusUndefined {
		query = query.Where(`"statusId" = ?`, statusID)
	}

	err := query.Select()
	if err != nil {
		return nil, fmt.Errorf("failed to get categories by status id")
	}

	return categories, nil
}

// GetCategoriesByIDsAndStatusID получение категорий по идентификаторам и статусу
func (r *NewsRepo) GetCategoriesByIDsAndStatusID(ctx context.Context, ids []int, statusID Status) ([]Category, error) {
	if len(ids) == 0 {
		return nil, nil
	}

	var categories []Category

	query := r.db.ModelContext(ctx, &categories).Where(`"categoryId" IN (?)`, pg.In(ids))

	if statusID != StatusUndefined {
		query = query.Where(`"statusId" = ?`, statusID)
	}

	err := query.Select()
	if err != nil {
		return nil, fmt.Errorf("failed to get categories by ids: %w", err)
	}

	return categories, nil
}

// GetTagsByIDsAndStatusID получение тегов по идентификаторам и статусу
func (r *NewsRepo) GetTagsByIDsAndStatusID(ctx context.Context, ids []int, statusID Status) ([]Tag, error) {
	if len(ids) == 0 {
		return nil, nil
	}

	var tags []Tag

	query := r.db.ModelContext(ctx, &tags).Where(`"tagId" IN (?)`, pg.In(ids))

	if statusID != StatusUndefined {
		query = query.Where(`"statusId" = ?`, statusID)
	}

	query = query.Order(`"name" ASC`)

	err := query.Select()
	if err != nil {
		return nil, fmt.Errorf("failed to get tags by ids: %w", err)
	}

	return tags, nil
}

// GetTagsByStatusID получение списка тегов по фильтру
func (r *NewsRepo) GetTagsByStatusID(ctx context.Context, statusID Status) ([]Tag, error) {
	var tags []Tag

	query := r.db.ModelContext(ctx, &tags)

	if statusID != StatusUndefined {
		query = query.Where(`"statusId" = ?`, statusID)
	}

	query = query.Order(`"name" ASC`)

	err := query.Select()
	if err != nil {
		return nil, fmt.Errorf("failed to get tags: %w", err)
	}

	return tags, nil
}
