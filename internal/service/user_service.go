package service

import (
	"context"
	"crud-pgx/internal/apperrors"
	"crud-pgx/internal/repository"
	"crud-pgx/model"
	"strings"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
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
