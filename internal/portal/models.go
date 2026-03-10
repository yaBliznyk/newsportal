package portal

import (
	"time"

	"github.com/yaBliznyk/newsportal/internal/db"
)

// Category категория новости
type Category struct {
	ID   int    // Идентификатор категории
	Name string // Название категории
}

// Tag тег новости
type Tag struct {
	ID   int    // Идентификатор тега
	Name string // Название тега
}

// News полный формат новости
type News struct {
	ID          int       // Идентификатор новости
	Title       string    // Заголовок новости
	Preamble    string    // Преамбула (краткое описание)
	Content     string    // Полный контент новости
	Category    Category  // Категория новости
	Tags        []Tag     // Теги новости
	Author      string    // Автор новости
	CreatedAt   time.Time // Дата создания
	PublishedAt time.Time // Дата публикации
}

// ShortNews краткая новость для списка новостей
type ShortNews struct {
	ID          int       // Идентификатор новости
	Title       string    // Заголовок новости
	Category    Category  // Категория новости
	Tags        []Tag     // Теги новости
	Author      string    // Автор новости
	CreatedAt   time.Time // Дата создания
	PublishedAt time.Time // Дата публикации
}

// PagedListNewsFilter фильтр списка новостей с пагинацией
type PagedListNewsFilter struct {
	ListNewsFilter
	Page  int // Номер страницы (по умолчанию 1)
	Limit int // Количество на страницу
}

// Validate проверяет корректность фильтра с пагинацией
func (f PagedListNewsFilter) Validate() error {
	if err := f.ListNewsFilter.Validate(); err != nil {
		return err
	}
	if f.Page < 0 {
		return ErrInvalidPage
	}
	if f.Limit < 0 {
		return ErrInvalidLimit
	}
	return nil
}

// ListNewsFilter фильтр новостей
type ListNewsFilter struct {
	CategoryID int       // Идентификатор категории
	TagID      int       // Идентификатор тега
	From       time.Time // Начало периода
	To         time.Time // Конец периода
}

// Validate проверяет корректность фильтра
func (f ListNewsFilter) Validate() error {
	if f.CategoryID <= 0 {
		return ErrInvalidCategoryID
	}
	if !f.From.IsZero() && !f.To.IsZero() && f.From.After(f.To) {
		return ErrInvalidDateRange
	}
	return nil
}

func NewCategory(category db.Category) Category {
	return Category{
		ID:   category.ID,
		Name: category.Name,
	}
}

func NewTags(tags []db.Tag) []Tag {
	res := make([]Tag, 0, len(tags))
	for _, tag := range tags {
		res = append(res, Tag{
			ID:   tag.ID,
			Name: tag.Name,
		})
	}
	return res
}

func NewShortNews(news db.ListNews, category db.Category, tags []db.Tag) ShortNews {
	return ShortNews{
		ID:          news.ID,
		Title:       news.Title,
		Category:    NewCategory(category),
		Tags:        NewTags(tags),
		Author:      news.Author,
		CreatedAt:   news.CreatedAt,
		PublishedAt: news.PublishedAt,
	}
}

func NewNews(news db.News, category db.Category, tags []db.Tag) News {
	return News{
		ID:          news.ID,
		Title:       news.Title,
		Preamble:    news.Preamble,
		Content:     news.Content,
		Category:    NewCategory(category),
		Tags:        NewTags(tags),
		Author:      news.Author,
		CreatedAt:   news.CreatedAt,
		PublishedAt: news.PublishedAt,
	}
}
