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

func (s *TaskService) ListTasks(ctx context.Context) ([]model.Task, error) {
	return s.repo.ListTasks(ctx)
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

type UserService struct {
	repo *repository.UserRepository
}

func (s *UserService) CreateUser(ctx context.Context, email string) (model.User, error) {
	email = strings.TrimSpace(email)
	if !strings.Contains(email, "@") || email == "" {
		return model.User{}, apperrors.ErrInvalidInput
	}
	return s.repo.CreateUser(ctx, email)
}

func (s *UserService) ListUsers(ctx context.Context) ([]model.User, error) {
	return s.repo.ListUsers(ctx)
}

func (s *UserService) GetUserByID(ctx context.Context, id int) (model.User, error) {
	if id <= 0 {
		return model.User{}, apperrors.ErrInvalidInput
	}
	return s.repo.GetUserByID(ctx, id)
}

func (s *UserService) UpdateUser(ctx context.Context, email string, id int) (model.User, error) {
	email = strings.TrimSpace(email)
	if email == "" || !strings.Contains(email, "@") || id <= 0 {
		return model.User{}, apperrors.ErrInvalidInput
	}
	return s.repo.UpdateUser(ctx, email, id)
}

func (s *UserService) DeleteUser(ctx context.Context, id int) (model.User, error) {
	if id <= 0 {
		return model.User{}, apperrors.ErrInvalidInput
	}
	return s.repo.DeleteUser(ctx, id)
}
