package service

import (
	"context"
	"fmt"

	"github.com/yaBliznyk/newsportal/internal/domain"
)

func (s *NewsService) CountNews(ctx context.Context, req domain.CountNewsReq) (*domain.CountNewsResp, error) {
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("req.Validate: %w", err)
	}

	count, err := s.repo.CountNews(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("repo.CountNews: %w", err)
	}

	return &domain.CountNewsResp{
		Count: count,
	}, nil
}
