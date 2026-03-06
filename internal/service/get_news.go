package service

import (
	"context"
	"fmt"

	"github.com/yaBliznyk/newsportal/internal/domain"
)

func (s *NewsService) GetNews(ctx context.Context, req domain.GetNewsReq) (*domain.GetNewsResp, error) {
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("req.Validate: %w", err)
	}

	news, err := s.repo.GetNews(ctx, req.ID)
	if err != nil {
		return nil, fmt.Errorf("repo.GetNews(%d): %w", req.ID, err)
	}

	return &domain.GetNewsResp{
		News: *news,
	}, nil
}
