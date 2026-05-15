package model

import (
	"time"
)

type Task struct {
	ID         int       `json:"id"`
	Title      string    `json:"title"`
	Completed  bool      `json:"completed"`
	Created_at time.Time `json:"created_at"`
	User_id    int       `json:"user_id"`
}

type User struct {
	ID         int       `json:"id"`
	Email      string    `json:"email"`
	Created_at time.Time `json:"created_at"`
}

// DTO модели

// tasks

// "/task/create"
type CreateTaskRequest struct {
	Title   string `json:"title"`
	User_id int    `json:"user_id"`
}

// "/task/update"
type UpdateTaskRequest struct {
	ID        int  `json:"id"`
	Completed bool `json:"completed"`
}

// users

// "/user/create"
type CreateUserRequest struct {
	Email string `json:"email"`
}

// "/user/update"
type UpdateUserRequest struct {
	Email string `json:"email"`
	ID    int    `json:"id"`
}

// Структура для фильтра litsTasks
type TaskFilter struct {
	Limit     int
	Offset    int
	Completed *bool // используем адрес, потому что пользователь может не использовать фильтр
	User_id   *int  // тоже самое
}
