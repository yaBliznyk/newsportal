package domain

// Category категория новости
type Category struct {
	ID   int    // Идентификатор категории
	Name string // Название категории
}

// Tag тег новости
type Tag struct {
	ID   int    // Идентификатор тега
	Name string // Название тега
}
