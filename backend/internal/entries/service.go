package entries

import (
	"context"
	"fmt"

	repo "github.com/zipshell/dev-learning-tracker/internal/adapters/postgresql/sqlc"
)

type EntryService interface {
	CreateEntry(ctx context.Context, newEntry repo.CreateEntryParams) (repo.Entry, error)
}

type svc struct {
	repo repo.Querier
}

func NewEntryService(repo repo.Querier) EntryService {
	return &svc{
		repo: repo,
	}
}

func (s *svc) CreateEntry(ctx context.Context, newEntry repo.CreateEntryParams) (repo.Entry, error) {
	if newEntry.Name == "" {
		return repo.Entry{}, fmt.Errorf("Entry's name is required")
	}

	if newEntry.FolderID == 0 {
		return repo.Entry{}, fmt.Errorf("Entry's folder_id is required")
	}

	return s.repo.CreateEntry(ctx, newEntry)
}
