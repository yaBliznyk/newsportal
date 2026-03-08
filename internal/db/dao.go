package db

import (
	"time"

	"github.com/yaBliznyk/newsportal/internal/domain"
)

// News модель новости для работы с БД
type News struct {
	ID          int32
	Title       string
	Preamble    string
	Content     string
	CategoryID  int32
	Category    string
	TagIDs      []int32
	Author      string
	CreatedAt   time.Time
	PublishedAt time.Time
}

// ToDomain преобразует DAO в доменную модель
func (n News) ToDomain(tags []domain.Tag) domain.News {
	return domain.News{
		ID:          int(n.ID),
		Title:       n.Title,
		Preamble:    n.Preamble,
		Content:     n.Content,
		Category:    domain.Category{ID: int(n.CategoryID), Name: n.Category},
		Tags:        tags,
		Author:      n.Author,
		CreatedAt:   n.CreatedAt,
		PublishedAt: n.PublishedAt,
	}
}

// NewsListItemDAO краткая модель новости для списка
type NewsListItemDAO struct {
	ID          int32
	Title       string
	CategoryID  int32
	Category    string
	TagIDs      []int32
	Author      string
	CreatedAt   time.Time
	PublishedAt time.Time
}

// ToDomain преобразует DAO в доменную модель
func (n NewsListItemDAO) ToDomain(tags []domain.Tag) domain.ListNewsItem {
	return domain.ListNewsItem{
		ID:          int(n.ID),
		Title:       n.Title,
		Category:    domain.Category{ID: int(n.CategoryID), Name: n.Category},
		Tags:        tags,
		Author:      n.Author,
		CreatedAt:   n.CreatedAt,
		PublishedAt: n.PublishedAt,
	}
}

// TagDAO модель тега для работы с БД
type TagDAO struct {
	ID   int32
	Name string
}

// ToDomain преобразует DAO в доменную модель
func (t TagDAO) ToDomain() domain.Tag {
	return domain.Tag{
		ID:   int(t.ID),
		Name: t.Name,
	}
}
