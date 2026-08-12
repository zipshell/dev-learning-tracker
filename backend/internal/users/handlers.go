package users

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/mail"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	repo "github.com/zipshell/dev-learning-tracker/internal/adapters/postgresql/sqlc"
	"github.com/zipshell/dev-learning-tracker/internal/jsonutil"
)

var validate = validator.New()

type handler struct {
	service UserService
}

func NewUserHandler(service UserService) *handler {
	return &handler{
		service: service,
	}
}

type CreateUserRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

func (h *handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var params repo.CreateUserParams

	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		jsonutil.Write(w, http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
		return
	}

	if err := validate.Struct(params); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// optional: stricter email validation
	if _, err := mail.ParseAddress(params.Email); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err := h.service.CreateUser(r.Context(), params)
	if err != nil {
		if errors.Is(err, ErrEmailAlreadyExists) {
			http.Error(w, "email already exists", http.StatusConflict)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	jsonutil.Write(w, http.StatusOK, map[string]bool{
		"success": true,
	})
}

func (h *handler) GetUserInfo(w http.ResponseWriter, r *http.Request) {
	userIdString := chi.URLParam(r, "user_id")
	userId := int64(0)
	if userIdString == "me" {
		userId = 0
	} else {
		parsedUserID, err := strconv.ParseInt(userIdString, 10, 64)
		if err != nil {
			jsonutil.Write(w, http.StatusBadRequest, map[string]string{
				"error": "Invalid user id",
			})
			return
		}
		userId = parsedUserID
	}

	userInfo, err := h.service.GetUserInfo(r.Context(), userId)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			http.NotFound(w, r)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	jsonutil.Write(w, http.StatusOK, userInfo)
}

func (h *handler) UpdateUserInfo(w http.ResponseWriter, r *http.Request) {
	var params repo.UpdateUserByIdParams

	userIdString := chi.URLParam(r, "user_id")
	userId := int64(0)
	if userIdString == "me" {
		userId = 0
	} else {
		parsedUserID, err := strconv.ParseInt(userIdString, 10, 64)
		if err != nil {
			jsonutil.Write(w, http.StatusBadRequest, map[string]string{
				"error": "Invalid user id",
			})
			return
		}
		userId = parsedUserID
	}

	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		jsonutil.Write(w, http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
		return
	}

	if params.ID == 0 {
		params.ID = int64(userId)
	} else if params.ID != int64(userId) {
		jsonutil.Write(w, http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
		return
	}

	updateResult, err := h.service.UpdateUserInfo(r.Context(), params)
	if err != nil {
		if errors.Is(err, ErrNoPermission) || errors.Is(err, ErrUserNotFound) {
			http.NotFound(w, r)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	jsonutil.Write(w, http.StatusOK, updateResult)
}

func (h *handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	err := h.service.DeleteUser(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	jsonutil.Write(w, http.StatusOK, map[string]bool{
		"success": true,
	})
}
