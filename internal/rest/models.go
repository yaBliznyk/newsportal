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

type ListNewsItem struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Category    Category  `json:"category"`
	Tags        []Tag     `json:"tags,omitempty"`
	Author      string    `json:"author"`
	CreatedAt   time.Time `json:"created_at"`
	PublishedAt time.Time `json:"published_at,omitempty"`
}

type News struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Preamble    string    `json:"preamble"`
	Content     string    `json:"content,omitempty"`
	Category    Category  `json:"category"`
	Tags        []Tag     `json:"tags,omitempty"`
	Author      string    `json:"author"`
	CreatedAt   time.Time `json:"created_at"`
	PublishedAt time.Time `json:"published_at,omitempty"`
}

type ListNewsResp struct {
	News []ListNewsItem `json:"news"`
}

type CountNewsResp struct {
	Count int `json:"count"`
}

type GetNewsResp struct {
	News News `json:"news"`
}

type GetCategoriesResp struct {
	Categories []Category `json:"categories"`
}

type GetTagsResp struct {
	Tags []Tag `json:"tags"`
}

func NewCategory(c portal.Category) Category {
	return Category{
		ID:   c.ID,
		Name: c.Name,
	}
}

func NewCategories(cats []portal.Category) []Category {
	res := make([]Category, 0, len(cats))
	for _, c := range cats {
		res = append(res, NewCategory(c))
	}
	return res
}

func NewTag(t portal.Tag) Tag {
	return Tag{
		ID:   t.ID,
		Name: t.Name,
	}
}

func NewTags(tags []portal.Tag) []Tag {
	res := make([]Tag, 0, len(tags))
	for _, t := range tags {
		res = append(res, NewTag(t))
	}
	return res
}

func NewListNewsItem(n portal.ShortNews) ListNewsItem {
	return ListNewsItem{
		ID:          n.ID,
		Title:       n.Title,
		Category:    NewCategory(n.Category),
		Tags:        NewTags(n.Tags),
		Author:      n.Author,
		CreatedAt:   n.CreatedAt,
		PublishedAt: n.PublishedAt,
	}
}

func NewListNewsItems(news []portal.ShortNews) []ListNewsItem {
	res := make([]ListNewsItem, 0, len(news))
	for _, n := range news {
		res = append(res, NewListNewsItem(n))
	}
	return res
}

func NewNews(n portal.News) News {
	return News{
		ID:          n.ID,
		Title:       n.Title,
		Preamble:    n.Preamble,
		Content:     n.Content,
		Category:    NewCategory(n.Category),
		Tags:        NewTags(n.Tags),
		Author:      n.Author,
		CreatedAt:   n.CreatedAt,
		PublishedAt: n.PublishedAt,
	}
}
