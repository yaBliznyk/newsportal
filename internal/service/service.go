package service

import (
	"context"

	"github.com/yaBliznyk/newsportal/internal/domain"
)

//go:generate go tool mockery --name=NewsRepository --output=mocks --outpkg=mocks --with-expecter

// NewsRepository интерфейс репозитория для работы с данными новостей
type NewsRepository interface {
	ListNews(ctx context.Context, req domain.ListNewsReq) ([]domain.ListNewsItem, error)
	CountNews(ctx context.Context, req domain.CountNewsReq) (int, error)
	GetNews(ctx context.Context, id int) (*domain.News, error)
	GetCategories(ctx context.Context) ([]domain.Category, error)
	GetTags(ctx context.Context) ([]domain.Tag, error)
}

// NewsService реализует domain.Service
type NewsService struct {
	repo NewsRepository
}

// New создаёт экземпляр сервиса новостей
func New(repo NewsRepository) *NewsService {
	return &NewsService{
		repo: repo,
	}
}
