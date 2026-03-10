package db

import (
	"time"
)

// Status статус записи
type Status int

const (
	StatusUndefined   Status = 0 // Статус не определен
	StatusPublished   Status = 1 // Опубликовано
	StatusUnpublished Status = 2 // Не опубликовано
	StatusDeleted     Status = 3 // Удалено
)

// News модель новости для работы с БД
type News struct {
	ID          int       // Идентификатор новости
	Title       string    // Заголовок
	Preamble    string    // Преамбула (краткое описание)
	Content     string    // Содержимое новости
	CategoryID  int       // Идентификатор категории
	TagIDs      []int     // Идентификаторы тегов
	Author      string    // Автор
	CreatedAt   time.Time // Дата создания
	PublishedAt time.Time // Дата публикации
	StatusID    int       // Идентификатор статуса
}

// ListNews краткая модель новости для списка
type ListNews struct {
	ID          int       // Идентификатор новости
	Title       string    // Заголовок
	CategoryID  int       // Идентификатор категории
	TagIDs      []int     // Идентификаторы тегов
	Author      string    // Автор
	CreatedAt   time.Time // Дата создания
	PublishedAt time.Time // Дата публикации
	StatusID    int       // Идентификатор статуса
}

// Tag модель тега
type Tag struct {
	ID       int    // Идентификатор тега
	Name     string // Название тега
	StatusID Status // Идентификатор статуса
}

// Category модель категории
type Category struct {
	ID        int    // Идентификатор категории
	Name      string // Название категории
	SortOrder int    // Порядок сортировки
	StatusID  Status // Идентификатор статуса
}

// PagedListNewsFilter фильтр списка новостей с пагинацией
type PagedListNewsFilter struct {
	ListNewsFilter
	Page  int // Номер страницы (по умолчанию 1)
	Limit int // Количество на страницу
}

// ListNewsFilter фильтр новостей
type ListNewsFilter struct {
	StatusID   Status    // Идентификатор статуса
	CategoryID int       // Идентификатор категории
	TagID      int       // Идентификатор тега
	From       time.Time // Начало периода
	To         time.Time // Конец периода
}
