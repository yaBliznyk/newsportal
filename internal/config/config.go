package config

import (
	"errors"
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	HTTP     HTTP
	Database Database
}

type HTTP struct {
	Addr string
}

type Database struct {
	URL string
}

func Load() (*Config, error) {
	return LoadWithPath(".")
}

func LoadWithPath(configPath string) (*Config, error) {
	v := viper.New()

	// Настройка viper для чтения .env файла
	v.SetConfigName(".env")
	v.SetConfigType("env")
	v.AddConfigPath(configPath)

	// Разрешаем читать из переменных окружения
	v.AutomaticEnv()

	// Устанавливаем значения по умолчанию
	v.SetDefault("HTTP_ADDR", ":8080")
	v.SetDefault("DATABASE_URL", "postgres://test:test@localhost:5432/test?sslmode=disable")

	// Bind environment variables
	_ = v.BindEnv("HTTP_ADDR")
	_ = v.BindEnv("DATABASE_URL")

	// Пытаемся прочитать .env файл (игнорируем ошибку, если файл не найден)
	if err := v.ReadInConfig(); err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if !errors.As(err, &configFileNotFoundError) {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
		// Файл не найден - это нормально, используем env vars и defaults
	}

	// Читаем значения напрямую
	cfg := &Config{
		HTTP: HTTP{
			Addr: v.GetString("HTTP_ADDR"),
		},
		Database: Database{
			URL: v.GetString("DATABASE_URL"),
		},
	}

	return cfg, nil
}
