package portal

import (
	"context"
	"fmt"

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

// ListNews список кратких новостей без текста
func (s *NewsManager) ListNews(ctx context.Context, filter PagedListNewsFilter) ([]ShortNews, error) {
	if err := filter.Validate(); err != nil {
		return nil, fmt.Errorf("filter.Validate: %w", err)
	}

	// Получаем категорию, по которой будет идти поиск
	category, err := s.repo.GetCategoryByIDAndStatusID(ctx, filter.CategoryID, db.StatusPublished)
	if err != nil {
		return nil, fmt.Errorf("get category by id and statusID: %w", err)
	}

	// Получаем список новостей
	news, err := s.repo.ListNewsByFilter(ctx, db.PagedListNewsFilter{
		ListNewsFilter: db.ListNewsFilter{
			StatusID:   db.StatusPublished,
			CategoryID: category.ID,
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

	// Получаем уникальные идентификаторы тегов из новостей
	tagIDs := uniqNewsTagIDs(news)

	if len(tagIDs) == 0 {
		result := make([]ShortNews, 0, len(news))
		for _, n := range news {
			result = append(result, NewShortNews(n, *category, nil))
		}
		return result, nil
	}

	// Получаем все опубликованные теги
	allTags, err := s.repo.GetTagsByStatusID(ctx, db.StatusPublished)
	if err != nil {
		return nil, fmt.Errorf("get tags: %w", err)
	}

	// Строим маппу тегов по ID
	tagMap := make(map[int]db.Tag, len(allTags))
	for _, tag := range allTags {
		tagMap[tag.ID] = tag
	}

	// Собираем список коротких новостей
	result := make([]ShortNews, 0, len(news))
	for _, n := range news {
		var tags []db.Tag
		for _, tagID := range n.TagIDs {
			if tag, ok := tagMap[tagID]; ok {
				tags = append(tags, tag)
			}
		}
		result = append(result, NewShortNews(n, *category, tags))
	}

	return result, nil
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
		return 0, fmt.Errorf("filter.Validate: %w", err)
	}

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
