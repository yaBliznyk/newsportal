package service

import (
	"context"
	"fmt"

	"github.com/yaBliznyk/newsportal/internal/domain"
)

func (s *NewsService) GetTags(ctx context.Context) (*domain.GetTagsResp, error) {
	tags, err := s.repo.GetTags(ctx)
	if err != nil {
		return nil, fmt.Errorf("repo.GetTags: %w", err)
	}

	return &domain.GetTagsResp{
		Tags: tags,
	}, nil
}
