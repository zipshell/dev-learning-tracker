package entries

import (
	"encoding/json"
	"net/http"

	repo "github.com/zipshell/dev-learning-tracker/internal/adapters/postgresql/sqlc"
	"github.com/zipshell/dev-learning-tracker/internal/jsonutil"
)

type handler struct {
	service EntryService
}

func NewEntryHandler(svc EntryService) *handler {
	return &handler{
		service: svc,
	}
}

func (h *handler) CreateEntry(w http.ResponseWriter, r *http.Request) {
	var params repo.CreateEntryParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		jsonutil.Write(w, http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
		return
	}

	response, err := h.service.CreateEntry(r.Context(), params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	jsonutil.Write(w, http.StatusOK, response)
}
