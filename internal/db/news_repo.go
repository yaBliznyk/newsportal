package db

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
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
		Where(`?.? = ?`, pg.Ident(Tables.News.Alias), pg.Ident(Columns.News.StatusID), filter.StatusID).
		Relation(Columns.News.Category)

	applyNewsFilter(filter, query)
	query.OrderExpr(`?.? DESC`, pg.Ident(Tables.News.Alias), pg.Ident(Columns.News.PublishedAt))

	// Пагинация
	limit := pager.Limit
	if limit <= 0 {
		limit = defaultLimit
	}
	query.Limit(limit)

	page := pager.Page
	if page <= 0 {
		page = 1
	}
	query.Offset((page - 1) * limit)

	err := query.Select()
	if err != nil {
		return nil, fmt.Errorf("failed to list news: %w", err)
	}
	return news, nil
}

// CountNews количество новостей по фильтру
func (r *NewsRepo) CountNews(ctx context.Context, filter NewsFilter) (int, error) {
	query := r.db.ModelContext(ctx, (*News)(nil)).
		Where(`?.? = ?`, pg.Ident(Tables.News.Alias), pg.Ident(Columns.News.StatusID), filter.StatusID).
		Relation(Columns.News.Category)

	applyNewsFilter(filter, query)

	count, err := query.Count()
	if err != nil {
		return 0, fmt.Errorf("failed to count news: %w", err)
	}

	return count, nil
}

// NewsByIDAndStatus получение полной новости по идентификатору и статусу
func (r *NewsRepo) NewsByIDAndStatus(ctx context.Context, id int, statusID, categoryStatusID StatusEnum) (*News, error) {
	news := &News{}

	query := r.db.ModelContext(ctx, news).
		Relation("Category").
		Where(`?.? = ?`, pg.Ident(Tables.News.Alias), pg.Ident(Columns.News.ID), id)

	if categoryStatusID != StatusUndefined {
		query.Where(`Category.? = ?`, pg.Ident(Columns.Category.StatusID), categoryStatusID)
	}

	if statusID != StatusUndefined {
		query.Where(`?.? = ?`, pg.Ident(Tables.News.Alias), pg.Ident(Columns.News.StatusID), statusID)
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
func (r *NewsRepo) GetCategoryByIDAndStatusID(ctx context.Context, id int, statusID StatusEnum) (*Category, error) {
	category := &Category{}

	query := r.db.ModelContext(ctx, category).Where(`? = ?`, pg.Ident(Columns.Category.ID), id)

	if statusID != StatusUndefined {
		query.Where(`? = ?`, pg.Ident(Columns.Category.StatusID), statusID)
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
func (r *NewsRepo) GetCategoriesByStatusID(ctx context.Context, statusID StatusEnum) ([]Category, error) {
	var categories []Category

	query := r.db.ModelContext(ctx, &categories)

	if statusID != StatusUndefined {
		query.Where(`? = ?`, pg.Ident(Columns.Category.StatusID), statusID)
	}

	err := query.Select()
	if err != nil {
		return nil, fmt.Errorf("failed to get categories by status id")
	}

	return categories, nil
}

// GetCategoriesByIDsAndStatusID получение категорий по идентификаторам и статусу
func (r *NewsRepo) GetCategoriesByIDsAndStatusID(ctx context.Context, ids []int, statusID StatusEnum) ([]Category, error) {
	if len(ids) == 0 {
		return nil, nil
	}

	var categories []Category

	query := r.db.ModelContext(ctx, &categories).Where(`? IN (?)`, pg.Ident(Columns.Category.ID), pg.In(ids))

	if statusID != StatusUndefined {
		query.Where(`? = ?`, pg.Ident(Columns.Category.StatusID), statusID)
	}

	err := query.Select()
	if err != nil {
		return nil, fmt.Errorf("failed to get categories by ids: %w", err)
	}

	return categories, nil
}

// GetTagsByIDsAndStatusID получение тегов по идентификаторам и статусу
func (r *NewsRepo) GetTagsByIDsAndStatusID(ctx context.Context, ids []int, statusID StatusEnum) ([]Tag, error) {
	if len(ids) == 0 {
		return nil, nil
	}

	var tags []Tag

	query := r.db.ModelContext(ctx, &tags).Where(`? IN (?)`, pg.Ident(Columns.Tag.ID), pg.In(ids))

	if statusID != StatusUndefined {
		query.Where(`? = ?`, pg.Ident(Columns.Tag.StatusID), statusID)
	}

	query.OrderExpr(`? ASC`, pg.Ident(Columns.Tag.Name))

	err := query.Select()
	if err != nil {
		return nil, fmt.Errorf("failed to get tags by ids: %w", err)
	}

	return tags, nil
}

// GetTagsByStatusID получение списка тегов по фильтру
func (r *NewsRepo) GetTagsByStatusID(ctx context.Context, ids []int, statusID StatusEnum) ([]Tag, error) {
	var tags []Tag

	query := r.db.ModelContext(ctx, &tags)

	if len(ids) > 0 {
		query.Where(`? IN (?)`, pg.Ident(Columns.Tag.ID), pg.In(ids))
	}

	if statusID != StatusUndefined {
		query.Where(`? = ?`, pg.Ident(Columns.Tag.StatusID), statusID)
	}

	query.OrderExpr(`? ASC`, pg.Ident(Columns.Tag.Name))

	err := query.Select()
	if err != nil {
		return nil, fmt.Errorf("failed to get tags: %w", err)
	}

	return tags, nil
}

func applyNewsFilter(filter NewsFilter, query *orm.Query) {
	if filter.CategoryStatusID != StatusUndefined {
		query.Where(`category.? = ?`, pg.Ident(Columns.Category.StatusID), filter.CategoryStatusID)
	}

	if filter.CategoryID != 0 {
		query.Where(`?.? = ?`, pg.Ident(Tables.News.Alias), pg.Ident(Columns.News.CategoryID), filter.CategoryID)
	}
	if filter.TagID != 0 {
		query.Where(`? = ANY(?.?)`, filter.TagID, pg.Ident(Tables.News.Alias), pg.Ident(Columns.News.TagIDs))
	}
	if !filter.From.IsZero() {
		query.Where(`?.? >= ?`, pg.Ident(Tables.News.Alias), pg.Ident(Columns.News.PublishedAt), filter.From)
	}
	if !filter.To.IsZero() {
		query.Where(`?.? <= ?`, pg.Ident(Tables.News.Alias), pg.Ident(Columns.News.PublishedAt), filter.To)
	}
}
