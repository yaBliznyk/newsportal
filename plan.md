# Анализ архитектуры news_portal: Dave Cheney & Rob Pike Perspective

## Текущая структура проекта

```
news_portal/
├── cmd/
│   ├── migration/main.go
│   └── portal/main.go
├── internal/
│   ├── config/config.go
│   ├── domain/
│   │   ├── common.go
│   │   ├── count_news.go
│   │   ├── get_categories.go
│   │   ├── get_news.go
│   │   ├── get_tags.go
│   │   ├── list_news.go
¤   │   ├── service.go         # интерфейс Service
│   │   └── status.go
│   ├── endpoints/public/
│   │   ├── controller.go
│   │   ├── count_news.go
│   │   ├── get_categories.go
│   │   ├── get_news.go
│   │   ├── get_tags.go
│   │   └── list_news.go
│   ├── repository/
│   │   ├── dao.go
│   │   ├── repository.go
│   │   └── *.go               # методы репозитория
│   ├── service/
│   │   ├── service.go         # интерфейс NewsRepository
│   │   └── *.go               # методы сервиса
│   └── svcerrs/
│       └── errors.go
├── migrations/
└── docs/
```

---

## 1. Анализ пакетов

### Проблема: Смешение ответственности в `service`

**Текущее состояние:**
```go
// internal/portal/portal.go
package service

type NewsRepository interface {
    ListNews(ctx context.Context, req domain.ListNewsReq) ([]domain.ListNewsItem, error)
    // ...
}

type NewsService struct {
    repo NewsRepository
}
```

**Проблема (Cheney):** Интерфейс `NewsRepository` определён в пакете `service`, но реализован в `repository`. Это нарушает **Accept interfaces, return structs**.

**Решение:**
```go
// internal/repository/repository.go
package repository

// Интерфейс должен быть рядом с потребителем или в domain
// Но лучше: потребитель определяет интерфейс

// internal/domain/repository.go
package domain

type NewsReader interface {
    ListNews(ctx context.Context, req ListNewsReq) ([]ListNewsItem, error)
    GetNews(ctx context.Context, id int) (*News, error)
    // Только чтение - segregation
}

type NewsWriter interface {
    CreateNews(ctx context.Context, news *News) error
    UpdateNews(ctx context.Context, news *News) error
}
```

### Проблема: Название `endpoints/public`

**Проблема (Pike):** `endpoints` - существительное, описывающее техническую деталь. `public` - признак, а не сущность.

**Решение:**
```
internal/
├── transport/
│   └── http/
│       └── public/            # или просто http/handler
```

Или ещё проще (Cheney style):
```
internal/
├── http/                      # пакет, а не директория
│   ├── handler.go
│   └── *.go
```

### Проблема: `svcerrs` сокращение

**Проблема (Pike):** "SVO errors" непонятно без контекста. `svcerrs` = service errors?

**Решение:**
```go
package errors  // или apierrors, apperrors
```

---

## 2. Анализ нейминга

### Проблема: Verb-Noun vs Noun-Verb

**Текущее:**
```go
// domain/portal.go
type Service interface {
    ListNews(ctx context.Context, ...) (*ListNewsResp, error)
    CountNews(ctx context.Context, ...) (*CountNewsResp, error)
    GetNews(ctx context.Context, ...) (*GetNewsResp, error)
}
```

**Анализ:**
- `ListNews`, `CountNews`, `GetNews` - Verb-Noun паттерн
- Это OK для HTTP API, но не идиоматично для Go

**Рекомендация (Pike):**
```go
type Service interface {
    News(ctx context.Context, id int) (*News, error)           // GetNews -> News
    NewsList(ctx context.Context, req ListReq) (*NewsList, error)  // ListNews -> NewsList
    NewsCount(ctx context.Context, req CountReq) (int, error)  // CountNews -> NewsCount
}
```

### Проблема: `ListNewsReq`, `ListNewsResp`, `ListNewsItem`

**Проблема (Cheney):** Избыточные суффиксы `Req`, `Resp`, `Item`.

**Текущее:**
```go
type ListNewsReq struct { ... }
type ListNewsResp struct { ... }
type ListNewsItem struct { ... }
```

**Решение:**
```go
package news

// Внутри пакета news контекст ясен
type ListRequest struct { ... }
type ListResponse struct { ... }
type Summary struct { ... }  // вместо ListNewsItem

// Или используя functional options:
type ListOption func(*listConfig)

func List(ctx context.Context, opts ...ListOption) ([]Summary, error)
```

### Проблема: `Controller`

**Проблема (Cheney):** `Controller` - это термин из MVC, не из Go. В Go мы пишем handlers.

**Текущее:**
```go
type Controller struct {
    log *slog.Logger
    svc domain.Service
}
```

**Решение:**
```go
package http

type Handler struct {
    log  *slog.Logger
    news domain.NewsService
}

func NewHandler(log *slog.Logger, news domain.NewsService) *Handler
```

---

## 3. Анализ интерфейсов

### Проблема: Бог-интерфейс `Service`

**Текущее:**
```go
type Service interface {
    ListNews(ctx context.Context, req ListNewsReq) (*ListNewsResp, error)
    CountNews(ctx context.Context, req CountNewsReq) (*CountNewsResp, error)
    GetNews(ctx context.Context, req GetNewsReq) (*GetNewsResp, error)
    GetCategories(ctx context.Context) (*GetCategoriesResp, error)
    GetTags(ctx context.Context) (*GetTagsResp, error)
}
```

**Проблема (Cheney):** Один интерфейс делает всё. Нарушает ISP (Interface Segregation Principle).

**Решение:**
```go
// domain/news.go
type NewsService interface {
    ByID(ctx context.Context, id int) (*News, error)
    List(ctx context.Context, req ListRequest) (*ListResponse, error)
    Count(ctx context.Context, req CountRequest) (int, error)
}

// domain/category.go
type CategoryService interface {
    All(ctx context.Context) ([]Category, error)
}

// domain/tag.go
type TagService interface {
    All(ctx context.Context) ([]Tag, error)
}
```

### Проблема: Интерфейс в неправильном пакете

**Текущее:**
- `domain.Service` - интерфейс бизнес-логики
- `service.NewsRepository` - интерфейс репозитория в слое сервиса

**Проблема (Cheney):** Интерфейсы должны определяться потребителем.

**Решение:**
```go
// portal/news.go - потребитель определяет интерфейс
package service

type newsRepository interface {  // приватный!
    ListNews(ctx context.Context, req domain.ListRequest) ([]domain.Summary, error)
    GetNews(ctx context.Context, id int) (*domain.News, error)
}

type NewsService struct {
    repo newsRepository
}

// repository/repository.go - реализация
package repository

type NewsRepository struct { db *pgxpool.Pool }

func (r *NewsRepository) ListNews(...) ([]domain.Summary, error) { ... }
```

---

## 4. Рекомендуемая структура

```
news_portal/
├── cmd/
│   └── server/main.go         # один entrypoint
├── internal/
│   ├── news/                  # доменный пакет
│   │   ├── news.go            # type News, Summary
│   │   ├── service.go         # NewsService interface
│   │   └── repository.go      # NewsRepository interface (или в service)
│   ├── category/
│   │   ├── category.go
│   │   └── service.go
│   ├── tag/
│   │   ├── tag.go
│   │   └── service.go
│   ├── postgres/              # реализация репозиториев
│   │   ├── news.go
│   │   ├── category.go
│   │   └── tag.go
│   ├── http/                  # HTTP layer
│   │   ├── handler.go
│   │   ├── news.go
│   │   ├── category.go
│   │   └── tag.go
│   ├── config/
│   │   └── config.go
│   └── errors/
│       └── errors.go
├── migrations/
└── docs/
```

---

## 5. План рефакторинга

### Фаза 1: Интерфейсы (низкий риск)

1. **Переместить `NewsRepository` интерфейс в `service` пакет (приватный)**
   - Файл: `internal/service/service.go`
   - Изменить: `type NewsRepository interface` → `type newsRepository interface`

2. **Разделить `domain.Service` на специализированные интерфейсы**
   - Файл: `internal/domain/service.go`
   - Создать: `NewsService`, `CategoryService`, `TagService`

### Фаза 2: Нейминг (средний риск)

3. **Переименовать `endpoints/public` → `http`**
   - Переместить весь пакет
   - Обновить импорты

4. **Переименовать `Controller` → `Handler`**
   - Файл: `internal/http/handler.go`

5. **Переименовать `svcerrs` → `errors` или `apierrors`**

### Фаза 3: Структура пакетов (высокий риск)

6. **Разделить `domain` на доменные пакеты**
   - `internal/domain/` → `internal/news/`, `internal/category/`, `internal/tag/`

7. **Разделить `repository` на `postgres`**
   - `internal/repository/` → `internal/postgres/`

8. **Упростить типы: убрать `Req`/`Resp` суффиксы**
   - Использовать контекст пакета

---

## 6. Принципы Cheney & Pike для этого проекта

| Принцип | Текущее | Рекомендация |
|---------|---------|--------------|
| Пакеты = существительные | `endpoints`, `service` | `http`, `news`, `postgres` |
| Интерфейсы у потребителя | `service.NewsRepository` (public) | `service.newsRepository` (private) |
| Маленькие интерфейсы | `domain.Service` (5 методов) | `NewsService`, `CategoryService` |
| Без суффиксов | `ListNewsReq`, `ListNewsResp` | `ListRequest`, `ListResponse` |
| Без сокращений | `svcerrs` | `errors`, `apierrors` |
| Без технических терминов | `Controller` | `Handler` |

---

## 7. Приоритеты

1. **Высокий приоритет:** Разделение `domain.Service` на маленькие интерфейсы
2. **Средний приоритет:** Перемещение `NewsRepository` интерфейса
3. **Низкий приоритет:** Переименование пакетов (breaking change для внешних потребителей)

---

## 8. Что уже хорошо

- Чистая слоистая архитектура (domain → service → repository → postgres)
- DAO паттерн для изоляции БД
- Валидация в доменных типах
- Консистентное использование контекста
- Правильная обработка ошибок с `errors.Is()`
- Graceful shutdown в main.go
