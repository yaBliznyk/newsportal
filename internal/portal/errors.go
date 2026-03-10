package portal

import (
	"errors"
)

var (
	ErrInvalidCategoryID = errors.New("invalid category id")
	ErrInvalidPage       = errors.New("invalid page number")
	ErrInvalidLimit      = errors.New("invalid limit")
	ErrInvalidDateRange  = errors.New("invalid date range: from is after to")
)
