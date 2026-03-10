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
	ID          int       `db:"newsId"`      // Идентификатор новости
	Title       string    `db:"title"`       // Заголовок
	Preamble    string    `db:"preamble"`    // Преамбула (краткое описание)
	Content     string    `db:"content"`     // Содержимое новости
	CategoryID  int       `db:"categoryId"`  // Идентификатор категории
	TagIDs      []int     `db:"tagIds"`      // Идентификаторы тегов
	Author      string    `db:"author"`      // Автор
	CreatedAt   time.Time `db:"createdAt"`   // Дата создания
	PublishedAt time.Time `db:"publishedAt"` // Дата публикации
	StatusID    Status    `db:"statusId"`    // Идентификатор статуса
}

// ListNews краткая модель новости для списка
type ListNews struct {
	ID          int       `db:"newsId"`      // Идентификатор новости
	Title       string    `db:"title"`       // Заголовок
	CategoryID  int       `db:"categoryId"`  // Идентификатор категории
	TagIDs      []int     `db:"tagIds"`      // Идентификаторы тегов
	Author      string    `db:"author"`      // Автор
	CreatedAt   time.Time `db:"createdAt"`   // Дата создания
	PublishedAt time.Time `db:"publishedAt"` // Дата публикации
	StatusID    Status    `db:"statusId"`    // Идентификатор статуса
}

// Tag модель тега
type Tag struct {
	ID       int    `db:"tagId"`    // Идентификатор тега
	Name     string `db:"name"`     // Название тега
	StatusID Status `db:"statusId"` // Идентификатор статуса
}

// Category модель категории
type Category struct {
	ID        int    `db:"categoryId"` // Идентификатор категории
	Name      string `db:"name"`       // Название категории
	SortOrder int    `db:"sortOrder"`  // Порядок сортировки
	StatusID  Status `db:"statusId"`   // Идентификатор статуса
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
