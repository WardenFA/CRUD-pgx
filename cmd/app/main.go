package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"crud-pgx/internal/db"
	"crud-pgx/internal/handler"
	"crud-pgx/internal/middleware"
	"crud-pgx/internal/repository"
	"crud-pgx/internal/service"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second) // контекст на время выполнения запроса
	defer cancel()
	// выполняем подключение
	pool, err := db.NewPool(ctx)
	if err != nil {
		log.Fatal("Failed to connect DB: ", err)
	}
	defer pool.Close()

	// пингуем базу для проверки
	if err := pool.Ping(ctx); err != nil {
		log.Fatal("ping error: ", err)
	}

	// сборка
	TaskRepository := repository.NewTaskRepository(pool)
	UserRepository := repository.NewUserRepository(pool)

	TaskService := service.NewTaskService(TaskRepository)
	UserService := service.NewUserService(UserRepository)

	TaskHandler := handler.NewTaskHandler(TaskService)
	UserHandler := handler.NewUserHandler(UserService)

	//роуты
	http.Handle("/task/create", middleware.LoggingMiddleWare(middleware.MethodMiddleware(http.MethodPost)(TaskHandler.CreateTask)))
	http.Handle("/task/list", middleware.LoggingMiddleWare(middleware.MethodMiddleware(http.MethodGet)(TaskHandler.ListTasks)))
	http.Handle("/task/get", middleware.LoggingMiddleWare(middleware.MethodMiddleware(http.MethodGet)(TaskHandler.GetTaskByID)))
	http.Handle("/task/update", middleware.LoggingMiddleWare(middleware.MethodMiddleware(http.MethodPatch)(TaskHandler.UpdateTaskStatus)))
	http.Handle("/task/delete", middleware.LoggingMiddleWare(middleware.MethodMiddleware(http.MethodDelete)(TaskHandler.DeleteTask)))

	http.Handle("/user/create", middleware.LoggingMiddleWare(middleware.MethodMiddleware(http.MethodPost)(UserHandler.CreateUser)))
	http.Handle("/user/list", middleware.LoggingMiddleWare(middleware.MethodMiddleware(http.MethodGet)(UserHandler.ListUsers)))
	http.Handle("/user/get", middleware.LoggingMiddleWare(middleware.MethodMiddleware(http.MethodGet)(UserHandler.GetUserByID)))
	http.Handle("/user/update", middleware.LoggingMiddleWare(middleware.MethodMiddleware(http.MethodPatch)(UserHandler.UpdateUser)))
	http.Handle("/user/delete", middleware.LoggingMiddleWare(middleware.MethodMiddleware(http.MethodDelete)(UserHandler.DeleteUser)))
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
