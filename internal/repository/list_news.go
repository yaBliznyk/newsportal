package repository

import (
	"context"

	"github.com/yaBliznyk/newsportal/internal/domain"
)

func (r *NewsRepository) ListNews(ctx context.Context, req domain.ListNewsReq) ([]domain.ListNewsItem, error) {
	// TODO: реализовать запрос к БД
	return nil, nil
}
