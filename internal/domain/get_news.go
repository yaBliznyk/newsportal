package domain

import (
	"time"

	"github.com/yaBliznyk/newsportal/internal/svcerrs"
)

// GetNewsReq параметры запроса конкретной новости
type GetNewsReq struct {
	ID int // Идентификатор новости
}

// Validate проверяет корректность параметров запроса.
func (r GetNewsReq) Validate() error {
	if r.ID <= 0 {
		return svcerrs.NewInvalidFieldError("id", "must be positive")
	}
	return nil
}

// News полный формат новости (с content)
type News struct {
	ID          int       // Идентификатор новости
	Title       string    // Заголовок новости
	Preamble    string    // Преамбула (краткое описание)
	Content     string    // Полный контент новости
	Category    Category  // Категория новости
	Tags        []Tag     // Теги новости
	Author      string    // Автор новости
	CreatedAt   time.Time // Дата создания
	PublishedAt time.Time // Дата публикации
}

// GetNewsResp ответ с конкретной новостью
type GetNewsResp struct {
	News News // Новость
}
