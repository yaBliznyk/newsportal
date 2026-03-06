package service

import (
	"context"
	"fmt"

	"github.com/yaBliznyk/newsportal/internal/domain"
)

func (s *NewsService) GetCategories(ctx context.Context) (*domain.GetCategoriesResp, error) {
	categories, err := s.repo.GetCategories(ctx)
	if err != nil {
		return nil, fmt.Errorf("repo.GetCategories: %w", err)
	}

	return &domain.GetCategoriesResp{
		Categories: categories,
	}, nil
}
