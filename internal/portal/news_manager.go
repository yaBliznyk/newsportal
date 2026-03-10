package portal

import (
	"context"
	"errors"
	"fmt"
	"slices"

	"github.com/yaBliznyk/newsportal/internal/db"
)

// NewsManager менеджер новостей
type NewsManager struct {
	repo *db.NewsRepo
}

// NewNewsManager создаёт экземпляр сервиса новостей
func NewNewsManager(repo *db.NewsRepo) *NewsManager {
	return &NewsManager{
		repo: repo,
	}
}

// GetNews возвращает детали новости
func (s *NewsManager) GetNews(ctx context.Context, id int) (*News, error) {
	// Получаем новость по ID со статусом "опубликована"
	news, err := s.repo.NewsByIDAndStatus(ctx, id, db.StatusPublished)
	if err != nil {
		if errors.Is(err, db.ErrNewsNotFound) {
			return nil, ErrNewsNotFound
		}
		return nil, fmt.Errorf("get news by id: %w", err)
	}

	// Получаем категорию и проверяем что она активна
	category, err := s.repo.GetCategoryByIDAndStatusID(ctx, news.CategoryID, db.StatusPublished)
	if err != nil {
		if errors.Is(err, db.ErrCategoryNotFound) {
			return nil, ErrCategoryNotFound
		}
		return nil, fmt.Errorf("get category by id: %w", err)
	}

	// Получаем опубликованные теги новости
	var tags []Tag
	if len(news.TagIDs) > 0 {
		dbTags, err := s.repo.GetTagsByIDsAndStatusID(ctx, news.TagIDs, db.StatusPublished)
		if err != nil {
			return nil, fmt.Errorf("get tags by ids: %w", err)
		}
		tags = NewTags(dbTags)
	}

	return &News{
		ID:          news.ID,
		Title:       news.Title,
		Preamble:    news.Preamble,
		Content:     news.Content,
		Category:    Category{ID: category.ID, Name: category.Name},
		Tags:        tags,
		Author:      news.Author,
		CreatedAt:   news.CreatedAt,
		PublishedAt: news.PublishedAt,
	}, nil
}

// ListNews список кратких новостей без текста
func (s *NewsManager) ListNews(ctx context.Context, filter PagedListNewsFilter) ([]ShortNews, error) {
	if err := filter.Validate(); err != nil {
		return nil, fmt.Errorf("validate filter: %w", err)
	}

	// Проверяем категорию фильтра, если указана
	if filter.CategoryID != 0 {
		_, err := s.repo.GetCategoryByIDAndStatusID(ctx, filter.CategoryID, db.StatusPublished)
		if err != nil {
			if errors.Is(err, db.ErrCategoryNotFound) {
				return nil, ErrCategoryNotFound
			}
			return nil, fmt.Errorf("get category by id and statusID: %w", err)
		}
	}

	// Получаем список новостей
	news, err := s.repo.ListNewsByFilter(ctx, db.PagedListNewsFilter{
		ListNewsFilter: db.ListNewsFilter{
			StatusID:   db.StatusPublished,
			CategoryID: filter.CategoryID,
			TagID:      filter.TagID,
			From:       filter.From,
			To:         filter.To,
		},
		Page:  filter.Page,
		Limit: filter.Limit,
	})
	if err != nil {
		return nil, fmt.Errorf("list news by filter: %w", err)
	}

	if len(news) == 0 {
		return nil, nil
	}

	// Получаем уникальные идентификаторы категорий из новостей
	categoryIDs := uniqNewsCategoryIDs(news)

	// Получаем опубликованные категории
	categories, err := s.repo.GetCategoriesByIDsAndStatusID(ctx, categoryIDs, db.StatusPublished)
	if err != nil {
		return nil, fmt.Errorf("get categories by ids: %w", err)
	}

	// Строим маппу категорий по ID
	categoryMap := make(map[int]db.Category, len(categories))
	for _, cat := range categories {
		categoryMap[cat.ID] = cat
	}

	// Получаем уникальные идентификаторы тегов из новостей
	tagIDs := uniqNewsTagIDs(news)

	// Строим маппу тегов по ID
	tagMap := make(map[int]db.Tag, len(tagIDs))
	if len(tagIDs) > 0 {
		tags, err := s.repo.GetTagsByIDsAndStatusID(ctx, tagIDs, db.StatusPublished)
		if err != nil {
			return nil, fmt.Errorf("get tags by ids: %w", err)
		}
		for _, tag := range tags {
			tagMap[tag.ID] = tag
		}
	}

	// Собираем список коротких новостей
	result := make([]ShortNews, 0, len(news))
	for _, n := range news {
		cat, ok := categoryMap[n.CategoryID]
		if !ok {
			continue
		}

		var nTags []db.Tag
		for _, tagID := range n.TagIDs {
			if tag, ok := tagMap[tagID]; ok {
				nTags = append(nTags, tag)
			}
		}
		result = append(result, NewShortNews(n, cat, nTags))
	}

	return result, nil
}

func uniqNewsCategoryIDs(news []db.ListNews) []int {
	seen := make(map[int]struct{})
	var result []int
	for _, n := range news {
		if _, ok := seen[n.CategoryID]; !ok {
			seen[n.CategoryID] = struct{}{}
			result = append(result, n.CategoryID)
		}
	}
	return result
}

func uniqNewsTagIDs(news []db.ListNews) []int {
	seen := make(map[int]struct{})
	var result []int
	for _, n := range news {
		for _, tagID := range n.TagIDs {
			if _, ok := seen[tagID]; !ok {
				seen[tagID] = struct{}{}
				result = append(result, tagID)
			}
		}
	}
	return result
}

// CountNews количество новостей по фильтру
func (s *NewsManager) CountNews(ctx context.Context, filter ListNewsFilter) (int, error) {
	if err := filter.Validate(); err != nil {
		return 0, fmt.Errorf("validate filter: %w", err)
	}

	// Проверяем категорию фильтра, если указана
	if filter.CategoryID != 0 {
		_, err := s.repo.GetCategoryByIDAndStatusID(ctx, filter.CategoryID, db.StatusPublished)
		if err != nil {
			if errors.Is(err, db.ErrCategoryNotFound) {
				return 0, ErrCategoryNotFound
			}
			return 0, fmt.Errorf("get category by id and statusID: %w", err)
		}
	}

	// Получаем число новостей по фильтру
	count, err := s.repo.CountNews(ctx, db.ListNewsFilter{
		StatusID:   db.StatusPublished,
		CategoryID: filter.CategoryID,
		TagID:      filter.TagID,
		From:       filter.From,
		To:         filter.To,
	})
	if err != nil {
		return 0, fmt.Errorf("repo.CountNews: %w", err)
	}

	return count, nil
}

// ListCategories список опубликованных категорий с сортировкой
func (s *NewsManager) ListCategories(ctx context.Context) ([]Category, error) {
	categories, err := s.repo.GetCategoriesByStatusID(ctx, db.StatusPublished)
	if err != nil {
		return nil, fmt.Errorf("get published categories: %w", err)
	}

	// Сортируем по SortOrder и имени
	slices.SortStableFunc(categories, func(a, b db.Category) int {
		if a.SortOrder != b.SortOrder {
			return a.SortOrder - b.SortOrder
		}
		if a.Name < b.Name {
			return -1
		}
		if a.Name > b.Name {
			return 1
		}
		return 0
	})

	// Формируем результат
	resp := make([]Category, 0, len(categories))
	for _, cat := range categories {
		resp = append(resp, Category{
			ID:   cat.ID,
			Name: cat.Name,
		})
	}

	return resp, nil
}

// ListTags список опубликованных тегов
func (s *NewsManager) ListTags(ctx context.Context) ([]Tag, error) {
	tags, err := s.repo.GetTagsByStatusID(ctx, db.StatusPublished)
	if err != nil {
		return nil, fmt.Errorf("get published tags: %w", err)
	}

	resp := make([]Tag, 0, len(tags))
	for _, t := range tags {
		resp = append(resp, Tag{
			ID:   t.ID,
			Name: t.Name,
		})
	}

	return resp, nil
}
