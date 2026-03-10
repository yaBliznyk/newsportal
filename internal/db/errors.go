package db

import (
	"errors"
)

var (
	ErrNewsNotFound     = errors.New("news not found")
	ErrCategoryNotFound = errors.New("category not found")
)
