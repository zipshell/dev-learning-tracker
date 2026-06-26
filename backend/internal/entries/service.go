package entries

import (
	"context"
	"fmt"
	"slices"

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

	value := ctx.Value("user")
	if value == nil {
		return repo.Entry{}, fmt.Errorf("No user session found")
	}

	userInfo, ok := value.(repo.User)
	if !ok {
		return repo.Entry{}, fmt.Errorf("invalid user session")
	}

	folderOwners, err := s.repo.FindUsersIdsByFolderId(ctx, newEntry.FolderID)
	if err != nil {
		return repo.Entry{}, fmt.Errorf("failed finding user ids")
	}

	userIds := make([]int64, len(folderOwners))

	for i, folderUser := range folderOwners {
		userIds[i] = folderUser.UserID
	}

	if !slices.Contains(userIds, userInfo.ID) {
		return repo.Entry{}, fmt.Errorf("not allowed to create entry")
	}

	return s.repo.CreateEntry(ctx, newEntry)
}
