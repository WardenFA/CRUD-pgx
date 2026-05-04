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

type UserHandler struct {
	service *service.UserService
}

func NewUserHandler(service *service.UserService) *UserHandler {
	return &UserHandler{service: service}
}

// остановился тут
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	ct := r.Header.Get("Content-Type")
	if !strings.HasPrefix(ct, "application/json") {
		log.Println("invalid media type")
		WriteError(w, http.StatusUnsupportedMediaType, errors.New("invalid content type"))
		return
	}
	var user model.User
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&user); err != nil {
		log.Println("invalid JSON")
		WriteError(w, http.StatusBadRequest, errors.New("invalid json"))
		return
	}
	newUser, err := h.service.CreateUser(r.Context(), user.Email)
	if err != nil {
		switch err {
		case apperrors.ErrInvalidInput:
			log.Println("invalid input data")
			WriteError(w, http.StatusBadRequest, apperrors.ErrInvalidInput)
			return
		case apperrors.ErrAlreadyExists:
			log.Println("email is already used")
			WriteError(w, http.StatusConflict, apperrors.ErrAlreadyExists)
			return
		default:
			log.Println("internal error")
			WriteError(w, http.StatusInternalServerError, errors.New("internal server error"))
			return
		}
	}
	WriteJSON(w, http.StatusCreated, newUser)
}

func (h *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.service.ListUsers(r.Context())
	if err != nil {
		log.Println("internal error")
		WriteError(w, http.StatusInternalServerError, errors.New("internal server error"))
		return
	}
	WriteJSON(w, http.StatusOK, users)
}

func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	URLid := r.URL.Query().Get("id")
	id, err := strconv.Atoi(URLid)
	if err != nil {
		log.Println("invalid url query")
		WriteError(w, http.StatusBadRequest, apperrors.ErrInvalidInput)
		return
	}
	gotUser, err := h.service.GetUserByID(r.Context(), id)
	if err != nil {
		switch err {
		case apperrors.ErrInvalidInput:
			log.Println("invalid input")
			WriteError(w, http.StatusBadRequest, apperrors.ErrInvalidInput)
			return
		case apperrors.ErrNotFound:
			log.Println("user not found")
			WriteError(w, http.StatusNotFound, apperrors.ErrNotFound)
			return
		default:
			log.Println("internal error")
			WriteError(w, http.StatusInternalServerError, errors.New("internal server error"))
			return
		}
	}
	WriteJSON(w, http.StatusOK, gotUser)
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	ct := r.Header.Get("Content-Type")
	if !strings.HasPrefix(ct, "application/json") {
		log.Println("invalid media type")
		WriteError(w, http.StatusUnsupportedMediaType, errors.New("Unsupported media type"))
		return
	}

	var user model.User
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&user); err != nil {
		log.Println("invalid json")
		WriteError(w, http.StatusBadRequest, errors.New("invalid json"))
		return
	}
	updatedUser, err := h.service.UpdateUser(r.Context(), user.Email, user.ID)
	if err != nil {
		switch err {
		case apperrors.ErrInvalidInput:
			log.Println("invalid input")
			WriteError(w, http.StatusBadRequest, apperrors.ErrInvalidInput)
			return
		case apperrors.ErrNotFound:
			log.Println("user not found")
			WriteError(w, http.StatusNotFound, apperrors.ErrNotFound)
			return
		case apperrors.ErrAlreadyExists:
			log.Println("email is already used")
			WriteError(w, http.StatusConflict, apperrors.ErrAlreadyExists)
			return
		default:
			log.Println("internal error")
			WriteError(w, http.StatusInternalServerError, errors.New("internal server error"))
			return
		}
	}
	WriteJSON(w, http.StatusOK, updatedUser)
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	URLid := r.URL.Query().Get("id")
	id, err := strconv.Atoi(URLid)
	if err != nil {
		log.Println("invalid url query")
		WriteError(w, http.StatusBadRequest, apperrors.ErrInvalidInput)
		return
	}
	deletedUser, err := h.service.DeleteUser(r.Context(), id)
	if err != nil {
		switch err {
		case apperrors.ErrNotFound:
			log.Println("user not found")
			WriteError(w, http.StatusNotFound, apperrors.ErrNotFound)
			return
		case apperrors.ErrInvalidInput:
			log.Println("invalid input")
			WriteError(w, http.StatusBadRequest, apperrors.ErrInvalidInput)
			return
		default:
			log.Println("internal server error")
			WriteError(w, http.StatusInternalServerError, errors.New("internal server error"))
			return
		}
	}
	WriteJSON(w, http.StatusOK, deletedUser)
}
