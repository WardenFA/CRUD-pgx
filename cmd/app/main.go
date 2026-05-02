package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"crud-pgx/internal/db"
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
}
