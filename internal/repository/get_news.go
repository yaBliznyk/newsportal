package repository

import (
	"context"

	"github.com/yaBliznyk/newsportal/internal/domain"
)

func (r *NewsRepository) GetNews(ctx context.Context, id int) (*domain.News, error) {
	// TODO: реализовать запрос к БД
	return nil, nil
}
