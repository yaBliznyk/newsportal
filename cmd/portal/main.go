package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yaBliznyk/newsportal/internal/repository"
	"github.com/yaBliznyk/newsportal/internal/service"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Настройка graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Загрузка конфигурации из переменных окружения
	dbURL := getEnv("DATABASE_URL", "postgres://test:test@localhost:5432/test?sslmode=disable")

	// Инициализация подключения к PostgreSQL
	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		log.Fatalf("failed to create connection pool: %v", err)
	}
	defer pool.Close()

	// Проверка подключения
	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("failed to ping database: %v", err)
	}
	fmt.Println("Successfully connected to PostgreSQL")

	// Инициализация слоёв приложения
	repo := repository.New(pool)
	svc := service.New(repo)

	// TODO: здесь можно запустить HTTP-сервер или использовать svc
	_ = svc

	// Ожидание сигнала завершения
	<-sigChan
	fmt.Println("\nShutting down...")
	cancel()
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
