package apperrors

import "errors"

// кастомные ошибки
var (
	ErrUsedEmail     = errors.New("email is already used")
	ErrForeignKey    = errors.New("foreign key error")
	ErrNullViolation = errors.New("must be not null")
)
