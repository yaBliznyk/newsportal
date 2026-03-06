package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/yaBliznyk/newsportal/internal/endpoints/public"
	"github.com/yaBliznyk/newsportal/internal/repository"
	"github.com/yaBliznyk/newsportal/internal/service"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Инициализация логгера
	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))

	// Загрузка конфигурации из переменных окружения
	dbURL := getEnv("DATABASE_URL", "postgres://test:test@localhost:5432/test?sslmode=disable")

	// Инициализация подключения к PostgreSQL
	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		log.Error("failed to create connection pool", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	// Проверка подключения
	if err := pool.Ping(ctx); err != nil {
		log.Error("failed to ping database", "error", err)
		os.Exit(1)
	}
	log.Info("Successfully connected to PostgreSQL")

	// Инициализация слоёв приложения
	repo := repository.New(pool)
	svc := service.New(repo)

	// Инициализация HTTP-сервера
	mux := http.NewServeMux()
	ctrl := public.NewController(log, svc)
	ctrl.Init(mux)

	addr := getEnv("HTTP_ADDR", ":8080")
	server := &http.Server{Addr: addr, Handler: mux}

	go func() {
		log.Info("HTTP server starting", "addr", addr)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("HTTP server error", "error", err)
			os.Exit(1)
		}
	}()

	// Настройка graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Ожидание сигнала завершения
	<-sigChan
	log.Info("Shutting down...")

	// Graceful shutdown HTTP сервера
	timeoutCtx, timeoutCancel := context.WithTimeout(context.Background(), time.Second*10)
	defer timeoutCancel()
	if err := server.Shutdown(timeoutCtx); err != nil {
		log.Error("HTTP server shutdown error", "error", err)
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
