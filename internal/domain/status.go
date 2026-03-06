package domain

// Status статус записи
type Status int

const (
	StatusPublished   Status = 1 // Опубликовано
	StatusUnpublished Status = 2 // Не опубликовано
	StatusDeleted     Status = 3 // Удалено
)
