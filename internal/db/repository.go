package db

import "github.com/jackc/pgx/v5/pgxpool"

// NewsRepo реализует service.NewsRepository
type NewsRepo struct {
	db *pgxpool.Pool
}

// NewNewsRepo создаёт экземпляр репозитория новостей
func NewNewsRepo(db *pgxpool.Pool) *NewsRepo {
	return &NewsRepo{db: db}
}
