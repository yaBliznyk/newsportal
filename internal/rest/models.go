package rest

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
	Content     string    `json:"content,omitempty"`
	Author      string    `json:"author"`
	CreatedAt   time.Time `json:"created_at"`
	PublishedAt time.Time `json:"published_at"`
	Category    *Category `json:"category"`
	Tags        []Tag     `json:"tags"`
}

type ListNewsReq struct {
	CategoryID int       `query:"category"` // Идентификатор категории
	TagID      int       `query:"tag"`      // Идентификатор тега
	From       time.Time `query:"from"`     // Начало периода
	To         time.Time `query:"to"`       // Конец периода
	Page       int       `query:"page"`     // Номер страницы
	Limit      int       `query:"limit"`    // Количество на страницу
}

func (r ListNewsReq) ToPortalFilter() portal.ListNewsFilter {
	return portal.ListNewsFilter{
		CategoryID: r.CategoryID,
		TagID:      r.TagID,
		From:       r.From,
		To:         r.To,
	}
}

func (r ListNewsReq) ToPortalPagination() portal.Pagination {
	return portal.Pagination{
		Page:  r.Page,
		Limit: r.Limit,
	}
}

type ListNewsResp struct {
	News []News `json:"news"`
}

type CountNewsResp struct {
	Count int `json:"count"`
}

type GetNewsResp struct {
	News *News `json:"news"`
}

type GetCategoriesResp struct {
	Categories []Category `json:"categories"`
}

type GetTagsResp struct {
	Tags []Tag `json:"tags"`
}

type ErrorResp struct {
	Error string
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
