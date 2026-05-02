package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"crud-pgx/internal/db"
	"crud-pgx/internal/handler"
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

	go func() {
		if err := http.ListenAndServe(":8080", nil); err != nil {
			log.Println("internal server error")
			return
		}
	}()

	// сборка
	TaskRepository := repository.NewTaskRepository(pool)
	UserRepository := repository.NewUserRepository(pool)

	TaskService := service.NewTaskService(TaskRepository)
	UserService := service.NewUserService(UserRepository)

	TaskHandler := handler.NewTaskHandler(TaskService)
	UserHandler := handler.NewUserHandler(UserService)

	//роуты
	http.HandleFunc("task/create", TaskHandler.CreateTask)
	http.HandleFunc("task/list", TaskHandler.ListTasks)
	http.HandleFunc("task/get", TaskHandler.GetTaskByID)
	http.HandleFunc("task/update", TaskHandler.UpdateTaskStatus)
	http.HandleFunc("task/delete", TaskHandler.DeleteTask)

	http.HandleFunc("user/create", UserHandler.CreateUser)
	http.HandleFunc("user/list", UserHandler.ListUsers)
	http.HandleFunc("user/get", UserHandler.GetUserByID)
	http.HandleFunc("user/update", UserHandler.UpdateUser)
	http.HandleFunc("user/delete", UserHandler.DeleteUser)
}
