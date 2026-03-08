package portal

import (
	"time"

	"github.com/yaBliznyk/newsportal/internal/db"
	"github.com/yaBliznyk/newsportal/internal/domain"
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

// News полный формат новости (с content)
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

// NewNews преобразует DAO в доменную модель
func NewNews(news db.News) News {
	return News{
		ID:          int(news.ID),
		Title:       news.Title,
		Preamble:    news.Preamble,
		Content:     news.Content,
		Category:    domainews.Category{ID: int(news.CategoryID), Name: news.Category},
		Tags:        tags,
		Author:      news.Author,
		CreatedAt:   news.CreatedAt,
		PublishedAt: news.PublishedAt,
	}
}

func NewCategory(category db.Category) Category {
	return Category{
		ID:   category.ID,
		Name: category.Name,
	}
}
