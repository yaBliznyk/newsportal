package db

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const defaultLimit = 20

// NewsRepo репозиторий новостей
type NewsRepo struct {
	db *pgxpool.Pool
}

// NewNewsRepo создаёт экземпляр репозитория новостей
func NewNewsRepo(db *pgxpool.Pool) *NewsRepo {
	return &NewsRepo{db: db}
}

// ListNewsByFilter список сокращенных новостей по фильтру
func (r *NewsRepo) ListNewsByFilter(ctx context.Context, filter PagedListNewsFilter) ([]ListNews, error) {
	conditions, args := buildFilterConditions(filter.ListNewsFilter, "n")

	query := `
		SELECT n."newsId", n."title", n."categoryId", n."tagIds", 
		       n."author", n."createdAt", n."publishedAt", n."statusId"
		FROM "news" n
	`

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	query += ` ORDER BY n."publishedAt" DESC`

	limit := filter.Limit
	if limit <= 0 {
		limit = defaultLimit
	}

	query += ` LIMIT @limit`
	args["limit"] = limit

	page := filter.Page
	if page <= 0 {
		page = 1
	}

	query += ` OFFSET @offset`
	args["offset"] = (page - 1) * limit

	rows, err := r.db.Query(ctx, query, args)
	if err != nil {
		return nil, fmt.Errorf("failed to list news: %w", err)
	}
	defer rows.Close()

	return pgx.CollectRows(rows, pgx.RowToStructByPos[ListNews])
}

// CountNews количество новостей по фильтру
func (r *NewsRepo) CountNews(ctx context.Context, filter ListNewsFilter) (int, error) {
	conditions, args := buildFilterConditions(filter, "")

	query := `
		SELECT COUNT(*)
		FROM "news"
	`

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	var count int
	err := r.db.QueryRow(ctx, query, args).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count news: %w", err)
	}

	return count, nil
}

// NewsByIDAndStatus получение полной новости по идентификатору и статусу
func (r *NewsRepo) NewsByIDAndStatus(ctx context.Context, id int, statusID Status) (*News, error) {
	query := `
		SELECT n."newsId", n."title", n."preamble", n."content", n."categoryId",
		       n."tagIds", n."author", n."createdAt", n."publishedAt", n."statusId"
		FROM "news" n
		WHERE n."newsId" = @newsID
	`

	args := pgx.NamedArgs{"newsID": id}

	if statusID != StatusUndefined {
		query += ` AND n."statusId" = @statusID`
		args["statusID"] = statusID
	}

	var news News
	err := r.db.QueryRow(ctx, query, args).Scan(
		&news.ID,
		&news.Title,
		&news.Preamble,
		&news.Content,
		&news.CategoryID,
		&news.TagIDs,
		&news.Author,
		&news.CreatedAt,
		&news.PublishedAt,
		&news.StatusID,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNewsNotFound
		}
		return nil, fmt.Errorf("failed to get news: %w", err)
	}

	return &news, nil
}

// GetCategoryByIDAndStatusID получение одной категории по идентификатору и статусу
func (r *NewsRepo) GetCategoryByIDAndStatusID(ctx context.Context, id int, statusID Status) (*Category, error) {
	query := `
		SELECT "categoryId", "name", "sortOrder", "statusId"
		FROM "categories"
		WHERE "categoryId" = @categoryID
	`

	args := pgx.NamedArgs{"categoryID": id}

	if statusID != StatusUndefined {
		query += ` AND "statusId" = @statusID`
		args["statusID"] = statusID
	}

	var category Category
	err := r.db.QueryRow(ctx, query, args).Scan(
		&category.ID,
		&category.Name,
		&category.SortOrder,
		&category.StatusID,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrCategoryNotFound
		}
		return nil, fmt.Errorf("failed to get category: %w", err)
	}

	return &category, nil
}

// GetCategoriesByStatusID получение списка категорий по статусу
func (r *NewsRepo) GetCategoriesByStatusID(ctx context.Context, statusID Status) ([]Category, error) {
	query := `
		SELECT "categoryId", "name", "sortOrder", "statusId"
		FROM "categories"
	`

	args := pgx.NamedArgs{}

	if statusID != StatusUndefined {
		query += ` WHERE "statusId" = @statusID`
		args["statusID"] = statusID
	}

	rows, err := r.db.Query(ctx, query, args)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return pgx.CollectRows(rows, pgx.RowToStructByPos[Category])
}

// GetTagsByStatusID получение списка тегов по фильтру
func (r *NewsRepo) GetTagsByStatusID(ctx context.Context, statusID Status) ([]Tag, error) {
	query := `
		SELECT "tagId", "name", "statusId"
		FROM "tags"
	`

	args := pgx.NamedArgs{}

	if statusID != StatusUndefined {
		query += ` WHERE "statusId" = @statusID`
		args["statusID"] = statusID
	}

	query += ` ORDER BY "name"`

	rows, err := r.db.Query(ctx, query, args)
	if err != nil {
		return nil, fmt.Errorf("failed to get tags: %w", err)
	}
	defer rows.Close()

	return pgx.CollectRows(rows, pgx.RowToStructByPos[Tag])
}

// buildFilterConditions формирует условия фильтрации и именованные аргументы.
// prefix — алиас таблицы (например "n"), если пустой — имена колонок без префикса.
func buildFilterConditions(filter ListNewsFilter, prefix string) ([]string, pgx.NamedArgs) {
	col := func(name string) string {
		if prefix != "" {
			return prefix + `."` + name + `"`
		}
		return `"` + name + `"`
	}

	var conditions []string
	args := pgx.NamedArgs{}

	if filter.StatusID != StatusUndefined {
		conditions = append(conditions, col("statusId")+" = @statusID")
		args["statusID"] = filter.StatusID
	}

	if filter.CategoryID != 0 {
		conditions = append(conditions, col("categoryId")+" = @categoryID")
		args["categoryID"] = filter.CategoryID
	}

	if filter.TagID != 0 {
		conditions = append(conditions, "@tagID = ANY("+col("tagIds")+")")
		args["tagID"] = filter.TagID
	}

	if !filter.From.IsZero() {
		conditions = append(conditions, col("publishedAt")+" >= @from")
		args["from"] = filter.From
	}

	if !filter.To.IsZero() {
		conditions = append(conditions, col("publishedAt")+" <= @to")
		args["to"] = filter.To
	}

	return conditions, args
}
