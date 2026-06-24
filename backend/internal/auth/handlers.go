package auth

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/mail"

	"github.com/go-playground/validator/v10"
	"github.com/zipshell/dev-learning-tracker/internal/jsonutil"
)

type handler struct {
	service AuthService
}

func NewAuthHandler(svc AuthService) *handler {
	return &handler{
		service: svc,
	}
}

type UserCredentials struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

var validate = validator.New()

func (h *handler) Login(w http.ResponseWriter, r *http.Request) {
	var params UserCredentials
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		jsonutil.Write(w, http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
		return
	}

	if err := validate.Struct(params); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if _, err := mail.ParseAddress(params.Email); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response, err := h.service.Login(r.Context(), params)
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, ErrUserNotFound) || errors.Is(err, ErrInvalidCredential) {
			status = http.StatusUnauthorized
		}
		http.Error(w, err.Error(), status)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    response.SessionToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		Expires:  response.ExpiredAt.Time,
	})
	jsonutil.Write(w, http.StatusOK, map[string]bool{
		"success": true,
	})
}
