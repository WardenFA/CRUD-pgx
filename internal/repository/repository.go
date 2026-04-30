package repository

import (
	"context"
	"crud-pgx/internal/apperrors"
	"crud-pgx/model"
	"errors"
	"log"

	"github.com/jackc/pgx/v5"
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

// Добавляем конструкторы для сборки в main
func NewTaskRepository(pool *pgxpool.Pool) *TaskRepository {
	return &TaskRepository{pool: pool}
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{pool: pool}
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
	rows, err := r.pool.Query(ctx, `SELECT id, title, completed, user_id FROM tasks`)
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
	err := r.pool.QueryRow(ctx, `SELECT * FROM tasks WHERE id = $1`, id).Scan(&task.ID, &task.Title, &task.Completed, &task.Created_at, &task.User_id)
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

func (r *UserRepository) CreateUser(ctx context.Context, email string) (model.User, error) {
	var u model.User
	err := r.pool.QueryRow(ctx, `INSERT INTO users(email) VALUES($1) RETURNING *`, email).Scan(&u.ID, &u.Email, &u.Created_at)
	if err != nil {
		// нет ошибки с несуществующим пользователем так как мы его вставляем
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23502" {
				return model.User{}, apperrors.ErrInvalidInput
			}
			if pgErr.Code == "23505" {
				return model.User{}, apperrors.ErrAlreadyExists
			}
		}
		return model.User{}, err
	}
	return u, nil
}

func (r *UserRepository) ListUsers(ctx context.Context) ([]model.User, error) {
	var users []model.User
	rows, err := r.pool.Query(ctx, `SELECT id, email FROM users`)
	if err != nil {
		return users, err
	}
	defer rows.Close()

	for rows.Next() {
		var u model.User
		err := rows.Scan(&u.ID, &u.Email)
		if err != nil {
			return users, err
		}
		users = append(users, u)
	}

	if err := rows.Err(); err != nil {
		return users, err
	}

	return users, nil
}

func (r *UserRepository) GetUserByID(ctx context.Context, id int) (model.User, error) {
	var u model.User
	err := r.pool.QueryRow(ctx, `SELECT id, email FROM users WHERE id=$1`, id).Scan(&u.ID, &u.Email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.User{}, apperrors.ErrNotFound
		}
		return model.User{}, err
	}
	return u, nil
}

func (r *UserRepository) UpdateUser(ctx context.Context, email string, id int) (model.User, error) {
	var u model.User
	err := r.pool.QueryRow(ctx, `UPDATE users SET email = $1 WHERE id = $2 RETURNING *`, email, id).Scan(&u.ID, &u.Email, &u.Created_at)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.User{}, apperrors.ErrNotFound
		}
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23502" {
				return model.User{}, apperrors.ErrInvalidInput
			}
			if pgErr.Code == "23505" {
				return model.User{}, apperrors.ErrAlreadyExists
			}
		}
		return model.User{}, err
	}

	return u, nil
}

func (r *UserRepository) DeleteUser(ctx context.Context, id int) (model.User, error) {
	// здесь можно продемонстрировать работу с транзакцией, потому что для удаления
	// юзера, нужно удалить сначала все его задачи

	// открывает транзакцию
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return model.User{}, err
	}
	defer tx.Rollback(ctx)
	// удаляем задачи
	cmdTag, err := tx.Exec(ctx, `DELETE FROM tasks WHERE user_id = $1`, id) // принимает cmdTag (узнай подробнее что это) (UPD: возвращает что-то типа DELETE 0 1 как в консоли с SQL)
	if err != nil {
		return model.User{}, err
	}
	// ошибка отсутствия изменений в Exec обрабатывается через cmdTag
	if cmdTag.RowsAffected() == 0 {
		log.Println("No data")
		// не прерываем, потому что у пользователя действительно может не быть задач
	}
	var u model.User
	err = tx.QueryRow(ctx, `DELETE FROM users WHERE id = $1 RETURNING *`, id).Scan(&u.ID, &u.Email, &u.Created_at)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.User{}, apperrors.ErrNotFound
		}
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23503" {
				return model.User{}, apperrors.ErrInvalidInput
			}
		}
		return model.User{}, err
	}
	return u, tx.Commit(ctx)
}
