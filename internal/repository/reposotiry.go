package repository

import (
	"context"
	"crud-pgx/internal/apperrors"
	"crud-pgx/model"
	"errors"
	"log"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// реп на таски
type TaskRepository struct {
	pool *pgxpool.Pool
}

// реп за юзеров
type UserRepository struct {
	pool *pgxpool.Pool
}

// реализация CRUD

func (r *TaskRepository) CreateTask(ctx context.Context, title string, user_id int) (model.Task, error) {
	var task model.Task
	err := r.pool.QueryRow(ctx, `INSERT INTO tasks(title, user_id) VALUES ($1, $2) RETURNING id, title, user_id`, title, user_id).Scan(&task.ID, &task.Title, &task.User_id)
	if err != nil {
		var pgErr *pgconn.PgError // обработка ошибки по коду
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23503" {
				log.Println(apperrors.ErrForeignKey)
				return model.Task{}, apperrors.ErrForeignKey
			}
			if pgErr.Code == "23502" {
				log.Println(apperrors.ErrForeignKey)
				return model.Task{}, apperrors.ErrForeignKey
			}
		}
		log.Println(err)
		return model.Task{}, err
	}
	return task, nil
}

func (r *TaskRepository) ListTasks(ctx context.Context) ([]model.Task, error) {
	var TaskSlice []model.Task
	rows, err := r.pool.Query(ctx, `SELECT id, title, completed, user_id FROM tasks`)
	if err != nil {
		log.Println(err)
		return TaskSlice, err
	}
	defer rows.Close()
	for rows.Next() {
		var task model.Task
		if err := rows.Scan(&task.ID, &task.Title, &task.Completed, &task.User_id); err != nil {
			log.Println(err)
			return TaskSlice, nil
		}
		TaskSlice = append(TaskSlice, task)
	}

	if err := rows.Err(); err != nil {
		log.Println(err)
		return TaskSlice, err
	}
	return TaskSlice, nil
}

func (r *TaskRepository) GetTaskByID(ctx context.Context, id int) (model.Task, error) {
	var task model.Task
	err := r.pool.QueryRow(ctx, `SELECT * FROM tasks WHERE id = $1 RETURNING *`, id).Scan(&task.ID, &task.Title, &task.Completed, &task.Created_at, &task.User_id)
	if err != nil {
		return model.Task{}, err
	}
	return task, nil
}

func (r *TaskRepository) UpdateTaskStatus(ctx context.Context, id int, completed bool) (model.Task, error) {
	var task model.Task
	err := r.pool.QueryRow(ctx, `UPDATE tasks SET completed = $1 WHERE id=$2 RETURNING id, title, completed, user_id`, completed, id).Scan(&task.ID, &task.Title, &task.Completed, &task.User_id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23502" {
				log.Println(apperrors.ErrNullViolation)
				return model.Task{}, apperrors.ErrNullViolation
			}
		}
		log.Println(err)
		return model.Task{}, err
	}
	return task, nil
}

// DELETE
