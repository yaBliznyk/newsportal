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
	tableName   struct{}  `pg:"news"`
	ID          int       `pg:"newsId,pk"`                      // Идентификатор новости
	Title       string    `pg:"title"`                          // Заголовок
	Preamble    string    `pg:"preamble"`                       // Преамбула (краткое описание)
	Content     string    `pg:"content"`                        // Содержимое новости
	CategoryID  int       `pg:"categoryId"`                     // Идентификатор категории
	TagIDs      []int     `pg:"tagIds,array"`                   // Идентификаторы тегов
	Author      string    `pg:"author"`                         // Автор
	CreatedAt   time.Time `pg:"createdAt"`                      // Дата создания
	PublishedAt time.Time `pg:"publishedAt"`                    // Дата публикации
	StatusID    Status    `pg:"statusId"`                       // Идентификатор статуса
	Category    *Category `pg:"rel:has-one,join_fk:categoryId"` // Категория
}

// Tag модель тега
type Tag struct {
	tableName struct{} `pg:"tags,alias:t"`
	ID        int      `pg:"tagId,pk"` // Идентификатор тега
	Name      string   `pg:"name"`     // Название тега
	StatusID  Status   `pg:"statusId"` // Идентификатор статуса
}

// Category модель категории
type Category struct {
	tableName struct{} `pg:"categories"`
	ID        int      `pg:"categoryId,pk"` // Идентификатор категории
	Name      string   `pg:"name"`          // Название категории
	SortOrder int      `pg:"sortOrder"`     // Порядок сортировки
	StatusID  Status   `pg:"statusId"`      // Идентификатор статуса
}

// Pagination пагинация
type Pagination struct {
	Page  int // Номер страницы (по умолчанию 1)
	Limit int // Количество на страницу
}

// NewsFilter фильтр новостей
type NewsFilter struct {
	StatusID         Status    // Идентификатор статуса
	CategoryStatusID Status    // Статус категории
	CategoryID       int       // Идентификатор категории
	TagID            int       // Идентификатор тега
	From             time.Time // Начало периода
	To               time.Time // Конец периода
}
