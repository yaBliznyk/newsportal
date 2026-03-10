package portal

import (
	"errors"
)

var (
	ErrInvalidData       = errors.New("invalid data")
	ErrInvalidCategoryID = errors.New("invalid category id")
	ErrInvalidPage       = errors.New("invalid page number")
	ErrInvalidLimit      = errors.New("invalid limit")
	ErrInvalidDateRange  = errors.New("invalid date range: from is after to")
	ErrCategoryNotFound  = errors.New("category not found")
	ErrNewsNotFound      = errors.New("news not found")
)
