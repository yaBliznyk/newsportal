package rpc

import (
	"time"

	"github.com/yaBliznyk/newsportal/internal/portal"
)

type Category struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Tag struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type News struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Preamble    string    `json:"preamble"`
	Content     *string   `json:"content,omitempty"`
	Author      string    `json:"author"`
	CreatedAt   time.Time `json:"created_at"`
	PublishedAt time.Time `json:"published_at"`
	Category    *Category `json:"category"`
	Tags        []Tag     `json:"tags"`
}

func NewCategory(c *portal.Category) *Category {
	if c == nil {
		return nil
	}

	return &Category{
		ID:   c.ID,
		Name: c.Name,
	}
}

func NewTag(t *portal.Tag) *Tag {
	if t == nil {
		return nil
	}

	return &Tag{
		ID:   t.ID,
		Name: t.Name,
	}
}

func NewNews(n *portal.News) *News {
	if n == nil {
		return nil
	}

	return &News{
		ID:          n.ID,
		Title:       n.Title,
		Preamble:    n.Preamble,
		Content:     n.Content,
		Author:      n.Author,
		CreatedAt:   n.CreatedAt,
		PublishedAt: n.PublishedAt,
		Category:    NewCategory(n.Category),
		Tags:        NewTags(n.Tags),
	}
}

type NewsFilter struct {
	CategoryID int       `json:"category_id,omitempty"` // Идентификатор категории
	TagID      int       `json:"tag_id,omitempty"`      // Идентификатор тега
	From       time.Time `json:"from,omitempty"`        // Начало периода
	To         time.Time `json:"to,omitempty"`          // Конец периода
}

func (f NewsFilter) ToPortal() portal.ListNewsFilter {
	return portal.ListNewsFilter{
		CategoryID: f.CategoryID,
		TagID:      f.TagID,
		From:       f.From,
		To:         f.To,
	}
}

type Pager struct {
	Page  int `json:"page,omitempty"`  // Номер страницы
	Limit int `json:"limit,omitempty"` // Количество на страницу
}

func (p Pager) ToPortal() portal.Pagination {
	return portal.Pagination{
		Page:  p.Page,
		Limit: p.Limit,
	}
}
