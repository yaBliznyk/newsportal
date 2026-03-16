package rpc

import (
	"context"
	"log/slog"

	"github.com/vmkteam/zenrpc/v2"

	"github.com/yaBliznyk/newsportal/internal/portal"
)

//go:generate zenrpc

type NewsService struct {
	log         *slog.Logger
	newsManager *portal.NewsManager
	zenrpc.Service
}

func NewNewsService(log *slog.Logger, newsManager *portal.NewsManager) *NewsService {
	return &NewsService{log: log, newsManager: newsManager}
}

func (h *NewsService) List(ctx context.Context, filter NewsFilter, pager Pager) ([]News, error) {
	news, err := h.newsManager.ListNews(ctx, filter.ToPortal(), pager.ToPortal())
	if err != nil {
		return nil, err
	}

	return NewNewsList(news), nil
}

func (h *NewsService) Count(ctx context.Context, filter NewsFilter) (int, error) {
	count, err := h.newsManager.CountNews(ctx, filter.ToPortal())
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (h *NewsService) Get(ctx context.Context, id int) (*News, error) {
	news, err := h.newsManager.GetNews(ctx, id)
	if err != nil {
		return nil, err
	}

	return NewNews(news), nil
}

func (h *NewsService) Categories(ctx context.Context) ([]Category, error) {
	categories, err := h.newsManager.ListCategories(ctx)
	if err != nil {
		return nil, err
	}

	return NewCategories(categories), nil
}

func (h *NewsService) Tags(ctx context.Context) ([]Tag, error) {
	tags, err := h.newsManager.ListTags(ctx)
	if err != nil {
		return nil, err
	}

	return NewTags(tags), nil
}
