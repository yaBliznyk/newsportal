package svcerrs

import (
	"errors"
	"fmt"
)

var (
	// ErrDataNotFound запрашиваемые данные не удалось найти.
	ErrDataNotFound = errors.New("data not found")
	// ErrAlreadyExist ошибка существования создаваемых данных.
	ErrAlreadyExist = errors.New("already exists")
	// ErrInvalidData универсальная ошибка про невалидные данные.
	ErrInvalidData = errors.New("invalid data")
	// ErrBusinessLogic универсальная ошибка бизнес-логики.
	ErrBusinessLogic = errors.New("wrong business logic")
	// ErrAccessDenied универсальная ошибка доступа к данным.
	ErrAccessDenied = errors.New("access denied")
	// ErrUnauthorized универсальная ошибка авторизации
	ErrUnauthorized = errors.New("unauthorized")
)

// ---------------------------------------
// BusinessLogicError
// ---------------------------------------

// BusinessLogicError универсальная ошибка бизнес-логики.
type BusinessLogicError struct {
	Alias string // Псевдоним ошибки, понятный для фронта.
}

// Error реализация error интерфейс.
func (e *BusinessLogicError) Error() string {
	return fmt.Sprintf("wrong business logic: %s", e.Alias)
}

// Is ошибка бизнес-логики ведет себя как ErrBusinessLogic.
func (e *BusinessLogicError) Is(target error) bool {
	return errors.Is(ErrBusinessLogic, target)
}

// NewBusinessLogicError конструктор ошибки бизнес-логики.
func NewBusinessLogicError(alias string) error {
	return &BusinessLogicError{
		Alias: alias,
	}
}

// ---------------------------------------
// InvalidFieldError
// ---------------------------------------

// InvalidFieldError ошибка валидации поля.
type InvalidFieldError struct {
	Field  string // Код проверяемого поля, понятный для фронта.
	Reason string // Причина ошибки валидации.
}

// Error реализация error интерфейс.
func (e *InvalidFieldError) Error() string {
	return fmt.Sprintf("invalid field: %q, reason: %s", e.Field, e.Reason)
}

// Is ошибка валидации поля ведет себя как ErrInvalidData.
func (e *InvalidFieldError) Is(target error) bool {
	return errors.Is(ErrInvalidData, target)
}

// NewInvalidFieldError конструктор ошибки валидации поля.
func NewInvalidFieldError(field, reason string) error {
	return &InvalidFieldError{
		Field:  field,
		Reason: reason,
	}
}
