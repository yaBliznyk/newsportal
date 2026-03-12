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
	Author      string    // Автор новости
	CategoryID  int       // Идентификатор категории
	TagIDs      []int     // Идентификаторы тегов
	CreatedAt   time.Time // Дата создания
	PublishedAt time.Time // Дата публикации
	Category    *Category // Категория новости
	Tags        []Tag     // Теги новости
}

func NewCategory(c *db.Category) *Category {
	if c == nil {
		return nil
	}

	return &Category{
		ID:   c.ID,
		Name: c.Name,
	}
}

func NewTag(t *db.Tag) *Tag {
	if t == nil {
		return nil
	}

	return &Tag{
		ID:   t.ID,
		Name: t.Name,
	}
}

func NewNews(n *db.News) *News {
	if n == nil {
		return nil
	}

	return &News{
		ID:          n.ID,
		Title:       n.Title,
		Preamble:    n.Preamble,
		Content:     n.Content,
		Author:      n.Author,
		CategoryID:  n.CategoryID,
		TagIDs:      n.TagIDs,
		CreatedAt:   n.CreatedAt,
		PublishedAt: n.PublishedAt,
	}
}

// Pagination пагинация
type Pagination struct {
	Page  int // Номер страницы (по умолчанию 1)
	Limit int // Количество на страницу
}

func (p Pagination) ToDB() db.Pagination {
	return db.Pagination{
		Page:  p.Page,
		Limit: p.Limit,
	}
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
	if f.CategoryID < 0 {
		return ErrInvalidCategoryID
	}
	if !f.From.IsZero() && !f.To.IsZero() && f.From.After(f.To) {
		return ErrInvalidDateRange
	}
	return nil
}

func (f ListNewsFilter) ToDB() db.NewsFilter {
	return db.NewsFilter{
		StatusID:         db.StatusPublished,
		CategoryStatusID: db.StatusPublished,
		CategoryID:       f.CategoryID,
		TagID:            f.TagID,
		From:             f.From,
		To:               f.To,
	}
}
