package apperrors

import "errors"

// Кастомные общие ошибки, понятные всему приложению
var (
	ErrNotFound      = errors.New("not found")
	ErrInvalidInput  = errors.New("invalid input")
	ErrAlreadyExists = errors.New("already exists")
	ErrConflict      = errors.New("conflict")
	ErrInternal      = errors.New("internal error")
)
