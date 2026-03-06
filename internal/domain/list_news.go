package domain

import (
	"time"

	"github.com/yaBliznyk/newsportal/internal/svcerrs"
)

// ListNewsReq параметры запроса списка новостей
type ListNewsReq struct {
	CategoryID int       // Идентификатор категории
	TagID      int       // Идентификатор тега
	From       time.Time // Начало периода
	To         time.Time // Конец периода
	Page       int       // Номер страницы (по умолчанию 1)
	Limit      int       // Количество на страницу
}

// Validate проверяет корректность параметров запроса.
func (r ListNewsReq) Validate() error {
	if r.CategoryID < 0 {
		return svcerrs.NewInvalidFieldError("category_id", "must be non-negative")
	}
	if r.TagID < 0 {
		return svcerrs.NewInvalidFieldError("tag_id", "must be non-negative")
	}
	if !r.From.IsZero() && !r.To.IsZero() && r.From.After(r.To) {
		return svcerrs.NewInvalidFieldError("from", "must be before to")
	}
	if r.Page < 0 {
		return svcerrs.NewInvalidFieldError("page", "must be non-negative")
	}
	if r.Limit < 0 {
		return svcerrs.NewInvalidFieldError("limit", "must be non-negative")
	}
	return nil
}

// ListNewsItem краткий формат новости (без content)
type ListNewsItem struct {
	ID          int       // Идентификатор новости
	Title       string    // Заголовок новости
	Category    Category  // Категория новости
	Tags        []Tag     // Теги новости
	Author      string    // Автор новости
	CreatedAt   time.Time // Дата создания
	PublishedAt time.Time // Дата публикации
}

// ListNewsResp ответ со списком новостей
type ListNewsResp struct {
	News []ListNewsItem // Список новостей
}
