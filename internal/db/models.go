package db

import (
	"time"
)

// StatusEnum статус записи
type StatusEnum int

const (
	StatusUndefined   StatusEnum = 0 // Статус не определен
	StatusPublished   StatusEnum = 1 // Опубликовано
	StatusUnpublished StatusEnum = 2 // Не опубликовано
	StatusDeleted     StatusEnum = 3 // Удалено
)

// Pagination пагинация
type Pagination struct {
	Page  int // Номер страницы (по умолчанию 1)
	Limit int // Количество на страницу
}

// NewsFilter фильтр новостей
type NewsFilter struct {
	StatusID         StatusEnum // Идентификатор статуса
	CategoryStatusID StatusEnum // Статус категории
	CategoryID       int        // Идентификатор категории
	TagID            int        // Идентификатор тега
	From             time.Time  // Начало периода
	To               time.Time  // Конец периода
}
