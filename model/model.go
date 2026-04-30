package model

import (
	"time"
)

type Task struct {
	ID         int
	Title      string
	Completed  bool
	Created_at time.Time
	User_id    int
}

type User struct {
	ID         int
	Email      string
	Created_at time.Time
}
