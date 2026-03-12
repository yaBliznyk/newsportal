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

// GetNews возвращает детали новости
func (s *NewsManager) GetNews(ctx context.Context, id int) (*News, error) {
	// Получаем новость по ID со статусом "опубликована"
	dbNews, err := s.repo.NewsByIDAndStatus(ctx, id, db.StatusPublished, db.StatusPublished)
	if err != nil {
		return nil, fmt.Errorf("get news by id: %w", err)
	} else if dbNews == nil {
		return nil, nil
	}

	news := NewNews(dbNews)

	// Получаем опубликованные теги новости
	if len(news.TagIDs) > 0 {
		dbTags, err := s.repo.GetTagsByIDsAndStatusID(ctx, news.TagIDs, db.StatusPublished)
		if err != nil {
			return nil, fmt.Errorf("get tags by ids: %w", err)
		}
		news.Tags = NewTags(dbTags)
	}
	return news, nil
}

// ListNews список кратких новостей без текста
func (s *NewsManager) ListNews(ctx context.Context, filter ListNewsFilter, pager Pagination) (NewsList, error) {
	if err := filter.Validate(); err != nil {
		return nil, fmt.Errorf("validate filter: %w", err)
	}

	// Получаем список новостей
	dbNewsList, err := s.repo.ListNewsByFilter(ctx, filter.ToDB(), pager.ToDB())
	if err != nil {
		return nil, fmt.Errorf("list news by filter: %w", err)
	}
	if len(dbNewsList) == 0 {
		return nil, nil
	}

	nn := NewNewsList(dbNewsList)

	// Заполняем новости тегами
	tagIDs := nn.UniqueTagIDs()
	if len(tagIDs) > 0 {
		dbTags, err := s.repo.GetTagsByIDsAndStatusID(ctx, tagIDs, db.StatusPublished)
		if err != nil {
			return nil, fmt.Errorf("get tags by ids: %w", err)
		}
		nn.FillTags(NewTags(dbTags))
	}

	return nn, nil
}

// CountNews количество новостей по фильтру
func (s *NewsManager) CountNews(ctx context.Context, filter ListNewsFilter) (int, error) {
	if err := filter.Validate(); err != nil {
		return 0, fmt.Errorf("validate filter: %w", err)
	}

	// Получаем число новостей по фильтру
	count, err := s.repo.CountNews(ctx, filter.ToDB())
	if err != nil {
		return 0, fmt.Errorf("count news: %w", err)
	}

	return count, nil
}

// ListCategories список опубликованных категорий с сортировкой
func (s *NewsManager) ListCategories(ctx context.Context) ([]Category, error) {
	cc, err := s.repo.GetCategoriesByStatusID(ctx, db.StatusPublished)
	if err != nil {
		return nil, fmt.Errorf("get published categories: %w", err)
	}

	return NewCategories(cc), nil
}

// ListTags список опубликованных тегов
func (s *NewsManager) ListTags(ctx context.Context) ([]Tag, error) {
	tt, err := s.repo.GetTagsByStatusID(ctx, db.StatusPublished)
	if err != nil {
		return nil, fmt.Errorf("get published tags: %w", err)
	}

	return NewTags(tt), nil
}
