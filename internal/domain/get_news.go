package domain

import "time"

// GetNewsReq параметры запроса конкретной новости
type GetNewsReq struct {
	ID int // Идентификатор новости
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
