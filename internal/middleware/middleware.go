package middleware

import (
	"crud-pgx/internal/handler"
	"errors"
	"log"
	"net/http"
	"time"
)

// Закончил тут
func LoggingMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Println("request", "method", r.Method, "path", r.URL.Path, "duration", time.Since(start).String())
	})
}
func MethodMiddleware(method string) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {

		return func(w http.ResponseWriter, r *http.Request) {
			if r.Method != method {
				handler.WriteError(w, http.StatusMethodNotAllowed, errors.New("Method not allowed"))
				log.Println("invalid request: метод не поддерживается")
				return
			}

			next(w, r)
		}
	}
}
