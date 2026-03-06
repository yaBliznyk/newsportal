package domain

import (
	"time"

	"github.com/yaBliznyk/newsportal/internal/svcerrs"
)

// CountNewsReq параметры запроса количества новостей
type CountNewsReq struct {
	CategoryID int       // Идентификатор категории
	TagID      int       // Идентификатор тега
	From       time.Time // Начало периода
	To         time.Time // Конец периода
}

// Validate проверяет корректность параметров запроса.
func (r CountNewsReq) Validate() error {
	if r.CategoryID < 0 {
		return svcerrs.NewInvalidFieldError("category_id", "must be non-negative")
	}
	if r.TagID < 0 {
		return svcerrs.NewInvalidFieldError("tag_id", "must be non-negative")
	}
	if !r.From.IsZero() && !r.To.IsZero() && r.From.After(r.To) {
		return svcerrs.NewInvalidFieldError("from", "must be before to")
	}
	return nil
}

// CountNewsResp ответ с количеством новостей
type CountNewsResp struct {
	Count int // Количество новостей
}
