package handler

import (
	"encoding/json"
	"net/http"
)

type response struct {
	Success bool   `json:"success"`
	Data    any    `json:"data,omitempty"`
	Error   string `json:"error,omitempty"`
}

func WriteJSON(w http.ResponseWriter, status int, data any) {
	// единый формат успешного JSON-ответа во всём API.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(response{Success: true, Data: data})
}

func WriteError(w http.ResponseWriter, status int, err error) {
	// единый формат ошибки
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(response{Success: false, Error: err.Error()})
}
