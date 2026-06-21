package folders

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	repo "github.com/zipshell/dev-learning-tracker/internal/adapters/postgresql/sqlc"
	"github.com/zipshell/dev-learning-tracker/internal/jsonutil"
)

type handler struct {
	service Service
}

func NewHandler(service Service) *handler {
	return &handler{
		service: service,
	}
}

func (h *handler) ListFoldersByUserId(w http.ResponseWriter, r *http.Request) {
	userIdString := chi.URLParam(r, "user_id")
	userId, err := strconv.Atoi(userIdString)
	if err != nil {
		log.Fatalf("Invalid user id: %v", err)
	}
	folders, err := h.service.ListFoldersByUserId(r.Context(), int64(userId))
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonutil.Write(w, http.StatusOK, folders)
}

func (h *handler) CreateFolder(w http.ResponseWriter, r *http.Request) {
	var params repo.CreateFolderParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		jsonutil.Write(w, http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
		return
	}

	newFolder, err := h.service.CreateFolder(r.Context(), params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	jsonutil.Write(w, http.StatusOK, newFolder)
}
