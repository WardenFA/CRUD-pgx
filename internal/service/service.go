package service

import (
	"context"
	"crud-pgx/internal/repository"
)

type TaskService struct {
	repo *repository.TaskRepository
}

func (s *TaskService) CreateTask(ctx context.Context)

type UserService struct {
	repo *repository.UserRepository
}
