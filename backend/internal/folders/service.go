package folders

import (
	"context"
	"fmt"

	repo "github.com/zipshell/dev-learning-tracker/internal/adapters/postgresql/sqlc"
)

type Service interface {
	ListFoldersByUserId(ctx context.Context, userId int64) ([]repo.ListFoldersByUserIdRow, error)
	CreateFolder(ctx context.Context, newFolder repo.CreateFolderParams) (repo.Folder, error)
}

type svc struct {
	repo repo.Querier
}

func NewService(repo repo.Querier) Service {
	return &svc{
		repo: repo,
	}
}

func (s *svc) ListFoldersByUserId(ctx context.Context, userId int64) ([]repo.ListFoldersByUserIdRow, error) {
	return s.repo.ListFoldersByUserId(ctx, userId)
}

func (s *svc) CreateFolder(ctx context.Context, newFolder repo.CreateFolderParams) (repo.Folder, error) {
	if newFolder.Name == "" {
		return repo.Folder{}, fmt.Errorf("Folder name is required")
	}

	return s.repo.CreateFolder(ctx, newFolder)
}
