package domain

import "time"

// ListNewsReq параметры запроса списка новостей
type ListNewsReq struct {
	CategoryID int       // Идентификатор категории
	TagID      int       // Идентификатор тега
	From       time.Time // Начало периода
	To         time.Time // Конец периода
	Page       int       // Номер страницы (по умолчанию 1)
	Limit      int       // Количество на страницу
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
