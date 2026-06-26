package folders

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
	repo "github.com/zipshell/dev-learning-tracker/internal/adapters/postgresql/sqlc"
)

type Service interface {
	ListFoldersByUserId(ctx context.Context, userId int64) ([]repo.ListFoldersByUserIdRow, error)
	CreateFolder(ctx context.Context, newFolder repo.CreateFolderParams) (repo.UserFolder, error)
}

type svc struct {
	db   *pgx.Conn
	repo repo.Querier
}

func NewService(db *pgx.Conn) Service {
	return &svc{
		db:   db,
		repo: repo.New(db),
	}
}

var (
	ErrFolderListNotFound = errors.New("folder list not found")
)

func (s *svc) ListFoldersByUserId(ctx context.Context, userId int64) ([]repo.ListFoldersByUserIdRow, error) {
	value := ctx.Value("user")
	if value == nil {
		log.Println("No user info found that matches the session")
		return []repo.ListFoldersByUserIdRow{}, fmt.Errorf("No user info found that matches the session")
	}

	userInfo, ok := value.(repo.User)
	if !ok {
		log.Println("invalid type for user info for current session")
		return []repo.ListFoldersByUserIdRow{}, fmt.Errorf("invalid type for user info for current session")
	}

	if userInfo.ID != userId {
		log.Println("Token doesn't match user id param")
		return []repo.ListFoldersByUserIdRow{}, fmt.Errorf("%w", ErrFolderListNotFound)
	}

	return s.repo.ListFoldersByUserId(ctx, userId)
}

func (s *svc) CreateFolder(ctx context.Context, newFolder repo.CreateFolderParams) (repo.UserFolder, error) {
	if newFolder.Name == "" {
		return repo.UserFolder{}, fmt.Errorf("UserFolder name is required")
	}

	value := ctx.Value("user")
	if value == nil {
		return repo.UserFolder{}, fmt.Errorf("session invalid")
	}

	userInfo, ok := value.(repo.User)
	if !ok {
		return repo.UserFolder{}, fmt.Errorf("session invalid")
	}

	tx, err := s.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return repo.UserFolder{}, fmt.Errorf("begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()

	q := repo.New(tx)

	createdFolder, err := q.CreateFolder(ctx, newFolder)
	if err != nil {
		return repo.UserFolder{}, fmt.Errorf("create folder: %w", err)
	}

	userFolder, err := q.AssignFolderToUser(ctx, repo.AssignFolderToUserParams{
		UserID:   userInfo.ID,
		FolderID: createdFolder.ID,
	})
	if err != nil {
		return repo.UserFolder{}, fmt.Errorf("assign folder to user: %w", err)
	}

	if err = tx.Commit(ctx); err != nil {
		return repo.UserFolder{}, fmt.Errorf("commit transaction: %w", err)
	}

	return userFolder, nil
}
