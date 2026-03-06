package service

import (
	"context"
	"fmt"

	"github.com/yaBliznyk/newsportal/internal/domain"
)

func (s *NewsService) ListNews(ctx context.Context, req domain.ListNewsReq) (*domain.ListNewsResp, error) {
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("req.Validate: %w", err)
	}

	news, err := s.repo.ListNews(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("repo.ListNews: %w", err)
	}

	return &domain.ListNewsResp{
		News: news,
	}, nil
}
