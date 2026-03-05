package domain

import "time"

// CountNewsReq параметры запроса количества новостей
type CountNewsReq struct {
	CategoryID int       // Идентификатор категории
	TagID      int       // Идентификатор тега
	From       time.Time // Начало периода
	To         time.Time // Конец периода
}

// CountNewsResp ответ с количеством новостей
type CountNewsResp struct {
	Count int // Количество новостей
}
