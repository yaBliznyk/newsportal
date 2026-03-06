package domain

import "context"

//go:generate go tool mockery --name=Service --output=mocks --outpkg=mocks --with-expecter

// Service основной сервис предметной области новостного портала
type Service interface {
	ListNews(ctx context.Context, req ListNewsReq) (*ListNewsResp, error)
	CountNews(ctx context.Context, req CountNewsReq) (*CountNewsResp, error)
	GetNews(ctx context.Context, req GetNewsReq) (*GetNewsResp, error)
	GetCategories(ctx context.Context) (*GetCategoriesResp, error)
	GetTags(ctx context.Context) (*GetTagsResp, error)
}
