package repository_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/yaBliznyk/newsportal/internal/domain"
	"github.com/yaBliznyk/newsportal/internal/repository"
)

// Константы для справочных данных из сидов
const (
	// Категории
	categoryTech     = 1 // Технологии
	categoryBusiness = 2 // Бизнес
	categorySport    = 3 // Спорт
	categoryScience  = 5 // Наука

	// Теги
	tagAI       = 1
	tagStartups = 2
	tagFootball = 3
	tagSpace    = 5
	tagFinance  = 6
)

// testDBConfig конфигурация подключения к тестовой БД
const testDBConfig = "postgres://test:test@localhost:5432/test?sslmode=disable"

// setupTestDB создает подключение к БД и возвращает репозиторий
func setupTestDB(t *testing.T) *repository.NewsRepository {
	ctx := t.Context()

	pool, err := pgxpool.New(ctx, testDBConfig)
	require.NoError(t, err, "failed to connect to database")

	t.Cleanup(func() {
		pool.Close()
	})

	return repository.New(pool)
}

// clearNews очищает таблицу новостей перед тестом
func clearNews(t *testing.T) {
	ctx := t.Context()

	pool, err := pgxpool.New(ctx, testDBConfig)
	require.NoError(t, err)
	defer pool.Close()

	_, err = pool.Exec(ctx, `DELETE FROM "news"`)
	require.NoError(t, err, "failed to clear news table")
}

// insertNews вставляет тестовую новость и возвращает её ID
func insertNews(t *testing.T, title string, categoryID int, tagIDs []int32, author, publishedAt string, statusID int) int {
	ctx := t.Context()

	pool, err := pgxpool.New(ctx, testDBConfig)
	require.NoError(t, err)
	defer pool.Close()

	var newsID int
	err = pool.QueryRow(ctx, `
		INSERT INTO "news" ("title", "categoryId", "tagIds", "author", "preamble", "content", "publishedAt", "statusId")
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING "newsId"
	`, title, categoryID, tagIDs, author, "Test preamble", "Test content", publishedAt, statusID).Scan(&newsID)
	require.NoError(t, err)

	return newsID
}

// ========================================
// GetCategories
// ========================================

func TestGetCategories_Success(t *testing.T) {
	repo := setupTestDB(t)

	categories, err := repo.GetCategories(t.Context())

	require.NoError(t, err)
	require.NotEmpty(t, categories, "categories should not be empty")

	// Проверяем, что категории приходят отсортированными по sortOrder, name
	assert.Equal(t, "Наука", categories[0].Name, "first category should be Наука (sortOrder=0)")
}

// ========================================
// GetTags
// ========================================

func TestGetTags_Success(t *testing.T) {
	repo := setupTestDB(t)

	tags, err := repo.GetTags(t.Context())

	require.NoError(t, err)
	require.NotEmpty(t, tags, "tags should not be empty")

	// Проверяем, что теги приходят отсортированными по name
	assert.Equal(t, "AI", tags[0].Name, "first tag should be AI (alphabetically first)")
}

// ========================================
// GetNews
// ========================================

func TestGetNews_Success(t *testing.T) {
	repo := setupTestDB(t)
	clearNews(t)

	// Создаём опубликованную новость
	newsID := insertNews(t, "Test News", categoryTech, []int32{tagAI}, "Author", "2024-06-15 12:00:00", int(domain.StatusPublished))

	news, err := repo.GetNews(t.Context(), newsID)

	require.NoError(t, err)
	require.NotNil(t, news)
	assert.Equal(t, newsID, news.ID)
	assert.Equal(t, "Test News", news.Title)
	assert.Equal(t, "Author", news.Author)
	assert.Equal(t, categoryTech, news.Category.ID)
	assert.Equal(t, "Технологии", news.Category.Name)
	require.Len(t, news.Tags, 1)
	assert.Equal(t, tagAI, news.Tags[0].ID)
	assert.Equal(t, "AI", news.Tags[0].Name)
}

func TestGetNews_NotFound(t *testing.T) {
	repo := setupTestDB(t)
	clearNews(t)

	news, err := repo.GetNews(t.Context(), 99999)

	require.NoError(t, err)
	assert.Nil(t, news, "should return nil for non-existent news")
}

func TestGetNews_DraftNotReturned(t *testing.T) {
	repo := setupTestDB(t)
	clearNews(t)

	// Создаём черновик (не опубликованную новость)
	newsID := insertNews(t, "Draft News", categoryTech, []int32{tagAI}, "Author", "2024-06-15 12:00:00", int(domain.StatusUnpublished))

	news, err := repo.GetNews(t.Context(), newsID)

	require.NoError(t, err)
	assert.Nil(t, news, "draft news should not be returned")
}

func TestGetNews_WithMultipleTags(t *testing.T) {
	repo := setupTestDB(t)
	clearNews(t)

	// Создаём новость с несколькими тегами
	newsID := insertNews(t, "Multi-tag News", categoryBusiness, []int32{tagStartups, tagFinance}, "Author", "2024-06-15 12:00:00", int(domain.StatusPublished))

	news, err := repo.GetNews(t.Context(), newsID)

	require.NoError(t, err)
	require.NotNil(t, news)
	require.Len(t, news.Tags, 2)
	// Теги должны быть отсортированы по имени (латиница идёт раньше кириллицы)
	assert.Equal(t, "Startups", news.Tags[0].Name)
	assert.Equal(t, "Финансы", news.Tags[1].Name)
}

// ========================================
// ListNews
// ========================================

func TestListNews_EmptyResult(t *testing.T) {
	repo := setupTestDB(t)
	clearNews(t)

	news, err := repo.ListNews(t.Context(), domain.ListNewsReq{Page: 1, Limit: 10})

	require.NoError(t, err)
	assert.Empty(t, news)
}

func TestListNews_Success(t *testing.T) {
	repo := setupTestDB(t)
	clearNews(t)

	// Создаём несколько новостей с разными датами
	insertNews(t, "News 1", categoryTech, []int32{tagAI}, "Author 1", "2024-06-15 12:00:00", int(domain.StatusPublished))
	insertNews(t, "News 2", categoryBusiness, []int32{tagStartups}, "Author 2", "2024-06-14 12:00:00", int(domain.StatusPublished))
	insertNews(t, "News 3", categoryScience, []int32{tagSpace}, "Author 3", "2024-06-16 12:00:00", int(domain.StatusPublished))

	news, err := repo.ListNews(t.Context(), domain.ListNewsReq{Page: 1, Limit: 10})

	require.NoError(t, err)
	require.Len(t, news, 3)

	// Должны быть отсортированы по publishedAt DESC
	assert.Equal(t, "News 3", news[0].Title) // 2024-06-16
	assert.Equal(t, "News 1", news[1].Title) // 2024-06-15
	assert.Equal(t, "News 2", news[2].Title) // 2024-06-14
}

func TestListNews_FilterByCategory(t *testing.T) {
	repo := setupTestDB(t)
	clearNews(t)

	// Создаём новости в разных категориях
	insertNews(t, "Tech News 1", categoryTech, []int32{tagAI}, "Author", "2024-06-15 12:00:00", int(domain.StatusPublished))
	insertNews(t, "Business News", categoryBusiness, []int32{tagFinance}, "Author", "2024-06-15 11:00:00", int(domain.StatusPublished))
	insertNews(t, "Tech News 2", categoryTech, []int32{tagStartups}, "Author", "2024-06-15 10:00:00", int(domain.StatusPublished))

	// Фильтруем по категории "Технологии"
	news, err := repo.ListNews(t.Context(), domain.ListNewsReq{
		CategoryID: categoryTech,
		Page:       1,
		Limit:      10,
	})

	require.NoError(t, err)
	require.Len(t, news, 2)
	assert.Equal(t, "Tech News 1", news[0].Title)
	assert.Equal(t, "Tech News 2", news[1].Title)
	assert.Equal(t, categoryTech, news[0].Category.ID)
	assert.Equal(t, categoryTech, news[1].Category.ID)
}

func TestListNews_FilterByTag(t *testing.T) {
	repo := setupTestDB(t)
	clearNews(t)

	// Создаём новости с разными тегами
	insertNews(t, "AI News", categoryTech, []int32{tagAI}, "Author", "2024-06-15 12:00:00", int(domain.StatusPublished))
	insertNews(t, "Sport News", categorySport, []int32{tagFootball}, "Author", "2024-06-15 11:00:00", int(domain.StatusPublished))
	insertNews(t, "AI & Startups", categoryTech, []int32{tagAI, tagStartups}, "Author", "2024-06-15 10:00:00", int(domain.StatusPublished))

	// Фильтруем по тегу AI
	news, err := repo.ListNews(t.Context(), domain.ListNewsReq{
		TagID: tagAI,
		Page:  1,
		Limit: 10,
	})

	require.NoError(t, err)
	require.Len(t, news, 2)
	assert.Equal(t, "AI News", news[0].Title)
	assert.Equal(t, "AI & Startups", news[1].Title)
}

func TestListNews_FilterByDateRange(t *testing.T) {
	repo := setupTestDB(t)
	clearNews(t)

	// Создаём новости с разными датами публикации
	insertNews(t, "Old News", categoryTech, []int32{tagAI}, "Author", "2024-06-10 12:00:00", int(domain.StatusPublished))
	insertNews(t, "Middle News", categoryTech, []int32{tagAI}, "Author", "2024-06-15 12:00:00", int(domain.StatusPublished))
	insertNews(t, "New News", categoryTech, []int32{tagAI}, "Author", "2024-06-20 12:00:00", int(domain.StatusPublished))

	from := time.Date(2024, 6, 12, 0, 0, 0, 0, time.UTC)
	to := time.Date(2024, 6, 18, 23, 59, 59, 0, time.UTC)

	news, err := repo.ListNews(t.Context(), domain.ListNewsReq{
		From:  from,
		To:    to,
		Page:  1,
		Limit: 10,
	})

	require.NoError(t, err)
	require.Len(t, news, 1)
	assert.Equal(t, "Middle News", news[0].Title)
}

func TestListNews_FilterByFromDate(t *testing.T) {
	repo := setupTestDB(t)
	clearNews(t)

	insertNews(t, "Old News", categoryTech, []int32{tagAI}, "Author", "2024-06-10 12:00:00", int(domain.StatusPublished))
	insertNews(t, "New News", categoryTech, []int32{tagAI}, "Author", "2024-06-20 12:00:00", int(domain.StatusPublished))

	from := time.Date(2024, 6, 15, 0, 0, 0, 0, time.UTC)

	news, err := repo.ListNews(t.Context(), domain.ListNewsReq{
		From:  from,
		Page:  1,
		Limit: 10,
	})

	require.NoError(t, err)
	require.Len(t, news, 1)
	assert.Equal(t, "New News", news[0].Title)
}

func TestListNews_FilterByToDate(t *testing.T) {
	repo := setupTestDB(t)
	clearNews(t)

	insertNews(t, "Old News", categoryTech, []int32{tagAI}, "Author", "2024-06-10 12:00:00", int(domain.StatusPublished))
	insertNews(t, "New News", categoryTech, []int32{tagAI}, "Author", "2024-06-20 12:00:00", int(domain.StatusPublished))

	to := time.Date(2024, 6, 15, 23, 59, 59, 0, time.UTC)

	news, err := repo.ListNews(t.Context(), domain.ListNewsReq{
		To:    to,
		Page:  1,
		Limit: 10,
	})

	require.NoError(t, err)
	require.Len(t, news, 1)
	assert.Equal(t, "Old News", news[0].Title)
}

func TestListNews_Pagination(t *testing.T) {
	repo := setupTestDB(t)
	clearNews(t)

	// Создаём 5 новостей
	for i := 1; i <= 5; i++ {
		date := fmt.Sprintf("2024-06-%02d 12:00:00", 10+i)
		insertNews(t, fmt.Sprintf("News %d", i), categoryTech, []int32{tagAI}, "Author", date, int(domain.StatusPublished))
	}

	// Первая страница (2 элемента)
	news, err := repo.ListNews(t.Context(), domain.ListNewsReq{
		Page:  1,
		Limit: 2,
	})
	require.NoError(t, err)
	require.Len(t, news, 2)
	assert.Equal(t, "News 5", news[0].Title) // самая новая
	assert.Equal(t, "News 4", news[1].Title)

	// Вторая страница
	news, err = repo.ListNews(t.Context(), domain.ListNewsReq{
		Page:  2,
		Limit: 2,
	})
	require.NoError(t, err)
	require.Len(t, news, 2)
	assert.Equal(t, "News 3", news[0].Title)
	assert.Equal(t, "News 2", news[1].Title)

	// Третья страница (1 элемент)
	news, err = repo.ListNews(t.Context(), domain.ListNewsReq{
		Page:  3,
		Limit: 2,
	})
	require.NoError(t, err)
	require.Len(t, news, 1)
	assert.Equal(t, "News 1", news[0].Title)

	// Четвёртая страница (пустая)
	news, err = repo.ListNews(t.Context(), domain.ListNewsReq{
		Page:  4,
		Limit: 2,
	})
	require.NoError(t, err)
	assert.Empty(t, news)
}

func TestListNews_CombinedFilters(t *testing.T) {
	repo := setupTestDB(t)
	clearNews(t)

	// Создаём новости с разными параметрами
	insertNews(t, "Tech AI Old", categoryTech, []int32{tagAI}, "Author", "2024-06-10 12:00:00", int(domain.StatusPublished))
	insertNews(t, "Tech AI New", categoryTech, []int32{tagAI}, "Author", "2024-06-20 12:00:00", int(domain.StatusPublished))
	insertNews(t, "Business AI", categoryBusiness, []int32{tagAI}, "Author", "2024-06-15 12:00:00", int(domain.StatusPublished))
	insertNews(t, "Tech Startups", categoryTech, []int32{tagStartups}, "Author", "2024-06-15 12:00:00", int(domain.StatusPublished))

	from := time.Date(2024, 6, 12, 0, 0, 0, 0, time.UTC)
	to := time.Date(2024, 6, 18, 23, 59, 59, 0, time.UTC)

	// Фильтр: категория Tech И тег AI И дата в диапазоне
	news, err := repo.ListNews(t.Context(), domain.ListNewsReq{
		CategoryID: categoryTech,
		TagID:      tagAI,
		From:       from,
		To:         to,
		Page:       1,
		Limit:      10,
	})

	require.NoError(t, err)
	// Только одна новость соответствует всем фильтрам
	assert.Empty(t, news, "no news should match all filters")
}

func TestListNews_ExcludeDrafts(t *testing.T) {
	repo := setupTestDB(t)
	clearNews(t)

	// Создаём опубликованную и черновую новости
	insertNews(t, "Published", categoryTech, []int32{tagAI}, "Author", "2024-06-15 12:00:00", int(domain.StatusPublished))
	insertNews(t, "Draft", categoryTech, []int32{tagAI}, "Author", "2024-06-15 12:00:00", int(domain.StatusUnpublished))

	news, err := repo.ListNews(t.Context(), domain.ListNewsReq{Page: 1, Limit: 10})

	require.NoError(t, err)
	require.Len(t, news, 1)
	assert.Equal(t, "Published", news[0].Title)
}

// ========================================
// CountNews
// ========================================

func TestCountNews_Empty(t *testing.T) {
	repo := setupTestDB(t)
	clearNews(t)

	count, err := repo.CountNews(t.Context(), domain.CountNewsReq{})

	require.NoError(t, err)
	assert.Equal(t, 0, count)
}

func TestCountNews_Success(t *testing.T) {
	repo := setupTestDB(t)
	clearNews(t)

	// Создаём 3 опубликованные новости
	insertNews(t, "News 1", categoryTech, []int32{tagAI}, "Author", "2024-06-15 12:00:00", int(domain.StatusPublished))
	insertNews(t, "News 2", categoryBusiness, []int32{tagFinance}, "Author", "2024-06-14 12:00:00", int(domain.StatusPublished))
	insertNews(t, "News 3", categoryScience, []int32{tagSpace}, "Author", "2024-06-13 12:00:00", int(domain.StatusPublished))

	count, err := repo.CountNews(t.Context(), domain.CountNewsReq{})

	require.NoError(t, err)
	assert.Equal(t, 3, count)
}

func TestCountNews_FilterByCategory(t *testing.T) {
	repo := setupTestDB(t)
	clearNews(t)

	insertNews(t, "Tech 1", categoryTech, []int32{tagAI}, "Author", "2024-06-15 12:00:00", int(domain.StatusPublished))
	insertNews(t, "Business", categoryBusiness, []int32{tagFinance}, "Author", "2024-06-14 12:00:00", int(domain.StatusPublished))
	insertNews(t, "Tech 2", categoryTech, []int32{tagStartups}, "Author", "2024-06-13 12:00:00", int(domain.StatusPublished))

	count, err := repo.CountNews(t.Context(), domain.CountNewsReq{
		CategoryID: categoryTech,
	})

	require.NoError(t, err)
	assert.Equal(t, 2, count)
}

func TestCountNews_FilterByTag(t *testing.T) {
	repo := setupTestDB(t)
	clearNews(t)

	insertNews(t, "AI 1", categoryTech, []int32{tagAI}, "Author", "2024-06-15 12:00:00", int(domain.StatusPublished))
	insertNews(t, "Football", categorySport, []int32{tagFootball}, "Author", "2024-06-14 12:00:00", int(domain.StatusPublished))
	insertNews(t, "AI 2", categoryBusiness, []int32{tagAI, tagStartups}, "Author", "2024-06-13 12:00:00", int(domain.StatusPublished))

	count, err := repo.CountNews(t.Context(), domain.CountNewsReq{
		TagID: tagAI,
	})

	require.NoError(t, err)
	assert.Equal(t, 2, count)
}

func TestCountNews_FilterByDateRange(t *testing.T) {
	repo := setupTestDB(t)
	clearNews(t)

	insertNews(t, "Old", categoryTech, []int32{tagAI}, "Author", "2024-06-10 12:00:00", int(domain.StatusPublished))
	insertNews(t, "Middle", categoryTech, []int32{tagAI}, "Author", "2024-06-15 12:00:00", int(domain.StatusPublished))
	insertNews(t, "New", categoryTech, []int32{tagAI}, "Author", "2024-06-20 12:00:00", int(domain.StatusPublished))

	from := time.Date(2024, 6, 12, 0, 0, 0, 0, time.UTC)
	to := time.Date(2024, 6, 18, 23, 59, 59, 0, time.UTC)

	count, err := repo.CountNews(t.Context(), domain.CountNewsReq{
		From: from,
		To:   to,
	})

	require.NoError(t, err)
	assert.Equal(t, 1, count)
}

func TestCountNews_FilterByFromDate(t *testing.T) {
	repo := setupTestDB(t)
	clearNews(t)

	insertNews(t, "Old", categoryTech, []int32{tagAI}, "Author", "2024-06-10 12:00:00", int(domain.StatusPublished))
	insertNews(t, "New", categoryTech, []int32{tagAI}, "Author", "2024-06-20 12:00:00", int(domain.StatusPublished))

	from := time.Date(2024, 6, 15, 0, 0, 0, 0, time.UTC)

	count, err := repo.CountNews(t.Context(), domain.CountNewsReq{
		From: from,
	})

	require.NoError(t, err)
	assert.Equal(t, 1, count)
}

func TestCountNews_FilterByToDate(t *testing.T) {
	repo := setupTestDB(t)
	clearNews(t)

	insertNews(t, "Old", categoryTech, []int32{tagAI}, "Author", "2024-06-10 12:00:00", int(domain.StatusPublished))
	insertNews(t, "New", categoryTech, []int32{tagAI}, "Author", "2024-06-20 12:00:00", int(domain.StatusPublished))

	to := time.Date(2024, 6, 15, 23, 59, 59, 0, time.UTC)

	count, err := repo.CountNews(t.Context(), domain.CountNewsReq{
		To: to,
	})

	require.NoError(t, err)
	assert.Equal(t, 1, count)
}

func TestCountNews_CombinedFilters(t *testing.T) {
	repo := setupTestDB(t)
	clearNews(t)

	insertNews(t, "Tech AI Old", categoryTech, []int32{tagAI}, "Author", "2024-06-10 12:00:00", int(domain.StatusPublished))
	insertNews(t, "Tech AI New", categoryTech, []int32{tagAI}, "Author", "2024-06-20 12:00:00", int(domain.StatusPublished))
	insertNews(t, "Business AI", categoryBusiness, []int32{tagAI}, "Author", "2024-06-15 12:00:00", int(domain.StatusPublished))
	insertNews(t, "Tech Startups", categoryTech, []int32{tagStartups}, "Author", "2024-06-15 12:00:00", int(domain.StatusPublished))

	from := time.Date(2024, 6, 12, 0, 0, 0, 0, time.UTC)
	to := time.Date(2024, 6, 18, 23, 59, 59, 0, time.UTC)

	count, err := repo.CountNews(t.Context(), domain.CountNewsReq{
		CategoryID: categoryTech,
		TagID:      tagAI,
		From:       from,
		To:         to,
	})

	require.NoError(t, err)
	assert.Equal(t, 0, count, "no news should match all filters")
}

func TestCountNews_ExcludeDrafts(t *testing.T) {
	repo := setupTestDB(t)
	clearNews(t)

	insertNews(t, "Published", categoryTech, []int32{tagAI}, "Author", "2024-06-15 12:00:00", int(domain.StatusPublished))
	insertNews(t, "Draft", categoryTech, []int32{tagAI}, "Author", "2024-06-15 12:00:00", int(domain.StatusUnpublished))
	insertNews(t, "Published 2", categoryTech, []int32{tagAI}, "Author", "2024-06-14 12:00:00", int(domain.StatusPublished))

	count, err := repo.CountNews(t.Context(), domain.CountNewsReq{})

	require.NoError(t, err)
	assert.Equal(t, 2, count, "drafts should not be counted")
}
