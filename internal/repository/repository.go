package repository

import "github.com/jackc/pgx/v5/pgxpool"

// NewsRepository реализует service.NewsRepository
type NewsRepository struct {
	db *pgxpool.Pool
}

// New создаёт экземпляр репозитория новостей
func New(db *pgxpool.Pool) *NewsRepository {
	return &NewsRepository{db: db}
}
