package service

import (
	"context"

	"github.com/yaBliznyk/newsportal/internal/domain"
)

// NewsRepository интерфейс репозитория для работы с данными новостей
type NewsRepository interface {
	ListNews(ctx context.Context, req domain.ListNewsReq) ([]domain.ListNewsItem, error)
	CountNews(ctx context.Context, req domain.CountNewsReq) (int, error)
	GetNews(ctx context.Context, id int) (*domain.News, error)
	GetCategories(ctx context.Context) ([]domain.Category, error)
	GetTags(ctx context.Context) ([]domain.Tag, error)
}

// newsService реализует domain.Service
type newsService struct {
	repo NewsRepository
}

// New создаёт экземпляр сервиса новостей
func New(repo NewsRepository) *newsService {
	return &newsService{repo: repo}
}
