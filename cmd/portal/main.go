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

	"github.com/go-pg/pg/v10"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"

	"github.com/yaBliznyk/newsportal/internal/config"
	"github.com/yaBliznyk/newsportal/internal/db"
	"github.com/yaBliznyk/newsportal/internal/portal"
	"github.com/yaBliznyk/newsportal/internal/rest"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Инициализация логгера
	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))

	// Загрузка конфигурации
	cfg, err := config.Load()
	if err != nil {
		log.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	// Инициализация подключения к PostgreSQL
	opt, err := pg.ParseURL(cfg.Database.URL)
	if err != nil {
		log.Error("failed to parse database url", "error", err)
		os.Exit(1)
	}

	pool := pg.Connect(opt)
	defer pool.Close()

	// Проверка подключения
	if err := pool.Ping(ctx); err != nil {
		log.Error("failed to ping database", "error", err)
		os.Exit(1)
	}
	log.Info("Successfully connected to PostgreSQL")

	// Инициализация слоёв приложения
	repo := db.NewNewsRepo(pool)
	newsManager := portal.NewNewsManager(repo)

	// Инициализация HTTP-сервера
	newsHandler := rest.NewNewsHandler(log, newsManager)

	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Recover())

	// Routes

	// Start server
	if err := e.Start(":8080"); err != nil {
		slog.Error("failed to start server", "error", err)
	}

	server := &http.Server{Addr: cfg.HTTP.Addr, Handler: newsHandler.Handle()}

	go func() {
		log.Info("HTTP server starting", "addr", cfg.HTTP.Addr)
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
