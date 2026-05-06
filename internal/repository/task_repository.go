package repository

import (
	"context"
	"crud-pgx/internal/apperrors"
	"crud-pgx/model"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// реп на таски
type TaskRepository struct {
	pool *pgxpool.Pool
}

// Добавляем конструкторы для сборки в main
func NewTaskRepository(pool *pgxpool.Pool) *TaskRepository {
	return &TaskRepository{pool: pool}
}

// реализация CRUD

func (r *TaskRepository) CreateTask(ctx context.Context, title string, user_id int) (model.Task, error) {
	var task model.Task
	err := r.pool.QueryRow(ctx, `INSERT INTO tasks(title, user_id) VALUES ($1, $2) RETURNING id, title, user_id`, title, user_id).Scan(&task.ID, &task.Title, &task.User_id)
	if err != nil {
		var pgErr *pgconn.PgError // обработка ошибки по коду
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23503" {
				return model.Task{}, apperrors.ErrInvalidInput
			}
			if pgErr.Code == "23502" {
				return model.Task{}, apperrors.ErrInvalidInput
			}
		}
		return model.Task{}, err
	}
	return task, nil
}

func (r *TaskRepository) ListTasks(ctx context.Context) ([]model.Task, error) {
	var TaskSlice []model.Task
	rows, err := r.pool.Query(ctx, `SELECT id, title, completed, created_at, user_id FROM tasks ORDER BY id`)
	if err != nil {
		return TaskSlice, err
	}
	defer rows.Close()
	for rows.Next() {
		var task model.Task
		if err := rows.Scan(&task.ID, &task.Title, &task.Completed, &task.Created_at, &task.User_id); err != nil {
			return TaskSlice, err
		}
		TaskSlice = append(TaskSlice, task)
	}

	if err := rows.Err(); err != nil {
		return TaskSlice, err
	}
	return TaskSlice, nil
}

func (r *TaskRepository) GetTaskByID(ctx context.Context, id int) (model.Task, error) {
	var task model.Task
	err := r.pool.QueryRow(ctx, `SELECT id, title, completed, created_at, user_id FROM tasks WHERE id = $1`, id).Scan(&task.ID, &task.Title, &task.Completed, &task.Created_at, &task.User_id)
	// Правило обработки ошибок с использованием QueryRow. Конкретно ошибка NotFound.
	// Если мы используем QueryRow и запрос с RETURNING (Не распространяется на SELECT,
	// так как он уже и так возвращает данные) - нам стоит обработать ошибку через
	// pgx.ErrNoRows. В случае если мы используем Exec (Returning тут уже не работает) -
	// обрабатываем ошибку через cmdTag.RowsAffected
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Task{}, apperrors.ErrNotFound
		}
		return model.Task{}, err
	}
	return task, nil
}

func (r *TaskRepository) UpdateTaskStatus(ctx context.Context, id int, completed bool) (model.Task, error) {
	var task model.Task
	err := r.pool.QueryRow(ctx, `UPDATE tasks SET completed = $1 WHERE id=$2 RETURNING id, title, completed, user_id`, completed, id).Scan(&task.ID, &task.Title, &task.Completed, &task.User_id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Task{}, apperrors.ErrNotFound
		}
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23502" {
				return model.Task{}, apperrors.ErrInvalidInput
			}
		}
		return model.Task{}, err
	}
	return task, nil
}

func (r *TaskRepository) DeleteTask(ctx context.Context, id int) (model.Task, error) {
	var task model.Task
	err := r.pool.QueryRow(ctx, `DELETE FROM tasks WHERE id =$1 RETURNING *`, id).Scan(&task.ID, &task.Title, &task.Completed, &task.Created_at, &task.User_id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Task{}, apperrors.ErrNotFound
		}
		return model.Task{}, err
	}
	return task, err
}
