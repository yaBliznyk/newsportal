# Makefile для news_portal

# Переменные
MIGRATION_NAME=migration
GO=go
GENNA=genna

# Цвета для вывода
COLOR_RESET=\033[0m
COLOR_GREEN=\033[32m
COLOR_YELLOW=\033[33m
COLOR_BLUE=\033[34m

.PHONY: generate models migrate-up migrate-down run install-tools
# Генерация кода
generate:
	@echo "$(COLOR_YELLOW)Запуск go generate...$(COLOR_RESET)"
	@$(GO) generate ./...
	@echo "$(COLOR_GREEN)Готово!$(COLOR_RESET)"

# Применение миграций
migrate-up:
	@echo "$(COLOR_YELLOW)Применение миграций...$(COLOR_RESET)"
	@./bin/$(MIGRATION_NAME) up
	@echo "$(COLOR_GREEN)Миграции применены!$(COLOR_RESET)"

# Откат миграций
migrate-down:
	@echo "$(COLOR_YELLOW)Откат миграций...$(COLOR_RESET)"
	@./bin/$(MIGRATION_NAME) down
	@echo "$(COLOR_GREEN)Миграции откачены!$(COLOR_RESET)"

# Запуск приложения
run:
	@echo "$(COLOR_YELLOW)Запуск приложения...$(COLOR_RESET)"
	@$(GO) run ./cmd/portal/main.go

# Тестирование API
test-api:
	@echo "$(COLOR_YELLOW)Запуск тестов API...$(COLOR_RESET)"
	@./test_api.sh
