package repository

import (
	"context"

	"github.com/yaBliznyk/newsportal/internal/domain"
)

func (r *NewsRepository) CountNews(ctx context.Context, req domain.CountNewsReq) (int, error) {
	// TODO: реализовать запрос к БД
	return 0, nil
}
