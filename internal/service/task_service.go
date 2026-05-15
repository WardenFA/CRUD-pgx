package service

import (
	"context"
	"crud-pgx/internal/apperrors"
	"crud-pgx/internal/repository"
	"crud-pgx/model"
	"strings"
)

type TaskService struct {
	repo *repository.TaskRepository
}

func NewTaskService(repo *repository.TaskRepository) *TaskService {
	return &TaskService{repo: repo}
}

func (s *TaskService) CreateTask(ctx context.Context, title string, user_id int) (model.Task, error) {
	title = strings.TrimSpace(title)
	if title == "" {
		return model.Task{}, apperrors.ErrInvalidInput
	}

	if user_id <= 0 {
		return model.Task{}, apperrors.ErrInvalidInput
	}
	return s.repo.CreateTask(ctx, title, user_id)
}

func (s *TaskService) ListTasks(ctx context.Context, limit, offset int) ([]model.Task, error) {
	if limit <= 0 || offset < 0 {
		return []model.Task{}, apperrors.ErrInvalidInput
	}
	return s.repo.ListTasks(ctx, limit, offset)
}

func (s *TaskService) GetTaskByID(ctx context.Context, id int) (model.Task, error) {
	if id <= 0 {
		return model.Task{}, apperrors.ErrInvalidInput
	}
	return s.repo.GetTaskByID(ctx, id)
}

func (s *TaskService) UpdateTaskStatus(ctx context.Context, id int, completed bool) (model.Task, error) {
	if id <= 0 {
		return model.Task{}, apperrors.ErrInvalidInput
	}
	return s.repo.UpdateTaskStatus(ctx, id, completed)
}

func (s *TaskService) DeleteTask(ctx context.Context, id int) (model.Task, error) {
	if id <= 0 {
		return model.Task{}, apperrors.ErrInvalidInput
	}
	return s.repo.DeleteTask(ctx, id)
}
