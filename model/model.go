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
