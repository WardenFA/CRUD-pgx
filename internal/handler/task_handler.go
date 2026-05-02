package handler

import (
	"crud-pgx/internal/apperrors"
	"crud-pgx/internal/service"
	"crud-pgx/model"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type TaskHandler struct {
	service *service.TaskService
}

func NewTaskHandler(service *service.TaskService) *TaskHandler {
	return &TaskHandler{service: service}
}

func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	ct := r.Header.Get("Content-Type")
	if !strings.HasPrefix(ct, "application/json") {
		log.Println("invalid media type")
		writeError(w, http.StatusUnsupportedMediaType, errors.New("invalid content type"))
		return
	}
	var task model.Task
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&task); err != nil {
		log.Println("invalid JSON")
		writeError(w, http.StatusBadRequest, errors.New("invalid json"))
		return
	}
	madeTask, err := h.service.CreateTask(r.Context(), task.Title, task.User_id)
	if err != nil {
		switch err {
		case apperrors.ErrInvalidInput:
			log.Println("invalid input data")
			writeError(w, http.StatusBadRequest, apperrors.ErrInvalidInput)
			return
		default:
			log.Println("internal error")
			writeError(w, http.StatusInternalServerError, errors.New("internal server error"))
			return
		}
	}
	writeJSON(w, http.StatusCreated, madeTask)
}

func (h *TaskHandler) ListTasks(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.service.ListTasks(r.Context())
	if err != nil {
		log.Println("internal error")
		writeError(w, http.StatusInternalServerError, errors.New("internal server error"))
		return
	}
	writeJSON(w, http.StatusOK, tasks)
}

func (h *TaskHandler) GetTaskByID(w http.ResponseWriter, r *http.Request) {
	URLid := r.URL.Query().Get("id")
	id, err := strconv.Atoi(URLid)
	if err != nil {
		log.Println("invalid url query")
		writeError(w, http.StatusBadRequest, apperrors.ErrInvalidInput)
		return
	}
	gotTask, err := h.service.GetTaskByID(r.Context(), id)
	if err != nil {
		switch err {
		case apperrors.ErrInvalidInput:
			log.Println("invalid input")
			writeError(w, http.StatusBadRequest, apperrors.ErrInvalidInput)
			return
		case apperrors.ErrNotFound:
			log.Println("task not found")
			writeError(w, http.StatusNotFound, apperrors.ErrNotFound)
			return
		default:
			log.Println("internal error")
			writeError(w, http.StatusInternalServerError, errors.New("internal server error"))
			return
		}
	}
	writeJSON(w, http.StatusOK, gotTask)
}

func (h *TaskHandler) UpdateTaskStatus(w http.ResponseWriter, r *http.Request) {
	ct := r.Header.Get("Content-Type")
	if !strings.HasPrefix(ct, "application/json") {
		log.Println("invalid media type")
		writeError(w, http.StatusUnsupportedMediaType, errors.New("Unsupported media type"))
		return
	}

	var task model.Task
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&task); err != nil {
		log.Println("invalid json")
		writeError(w, http.StatusBadRequest, errors.New("invalid json"))
		return
	}
	updatedTask, err := h.service.UpdateTaskStatus(r.Context(), task.ID, task.Completed)
	if err != nil {
		switch err {
		case apperrors.ErrInvalidInput:
			log.Println("invalid input")
			writeError(w, http.StatusBadRequest, apperrors.ErrInvalidInput)
			return
		case apperrors.ErrNotFound:
			log.Println("task not found")
			writeError(w, http.StatusNotFound, apperrors.ErrNotFound)
			return
		default:
			log.Println("internal error")
			writeError(w, http.StatusInternalServerError, errors.New("internal server error"))
			return
		}
	}
	writeJSON(w, http.StatusOK, updatedTask)
}

func (h *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	URLid := r.URL.Query().Get("id")
	id, err := strconv.Atoi(URLid)
	if err != nil {
		log.Println("invalid url query")
		writeError(w, http.StatusBadRequest, apperrors.ErrInvalidInput)
		return
	}
	deletedTask, err := h.service.DeleteTask(r.Context(), id)
	if err != nil {
		switch err {
		case apperrors.ErrNotFound:
			log.Println("task not found")
			writeError(w, http.StatusNotFound, apperrors.ErrNotFound)
			return
		case apperrors.ErrInvalidInput:
			log.Println("invalid input")
			writeError(w, http.StatusBadRequest, apperrors.ErrInvalidInput)
			return
		default:
			log.Println("internal server error")
			writeError(w, http.StatusInternalServerError, errors.New("internal server error"))
			return
		}
	}
	writeJSON(w, http.StatusOK, deletedTask)
}
