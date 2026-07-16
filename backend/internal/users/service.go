package users

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	repo "github.com/zipshell/dev-learning-tracker/internal/adapters/postgresql/sqlc"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	CreateUser(ctx context.Context, newUser repo.CreateUserParams) error
	GetUserInfo(ctx context.Context) (UserInfo, error)
	UpdateUserInfo(ctx context.Context, updateInfo repo.UpdateUserByIdParams) (UserInfo, error)
	DeleteUser(ctx context.Context) error
}

var (
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrUserNotFound       = errors.New("user not found")
	ErrNoPermission       = errors.New("no permission")
	ErrUpdateFailure      = errors.New("user update failed")
	ErrDeleteFailure      = errors.New("user delete failed")
)

type UserInfo struct {
	ID        int64              `json:"id"`
	Email     string             `json:"email"`
	CreatedAt pgtype.Timestamptz `json:"created_at"`
}

func NewUserService(db *pgx.Conn) UserService {
	return &svc{
		db:   db,
		repo: repo.New(db),
	}
}

type Claims struct {
	UserID int64 `json:"user_id"`
	jwt.RegisteredClaims
}

type svc struct {
	db   *pgx.Conn
	repo repo.Querier
}

type Tokens struct {
	SessionToken string `json:"session_token"`
}

func (s *svc) CreateUser(ctx context.Context, newUser repo.CreateUserParams) error {
	tx, err := s.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()

	hashedPassword, err := hashPassword(newUser.Password)
	if err != nil {
		return err
	}

	q := repo.New(tx)

	newUser.Password = hashedPassword

	_, err = q.CreateUser(ctx, newUser)
	var pgErr *pgconn.PgError
	if err != nil {
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return fmt.Errorf("%w", ErrEmailAlreadyExists)
		}
		return fmt.Errorf("New user creation failed in write to db step")
	}

	if err = tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}

func (s *svc) GetUserInfo(ctx context.Context) (UserInfo, error) {
	value := ctx.Value("user")
	if value == nil {
		log.Println("No user info found that matches the session")
		return UserInfo{}, fmt.Errorf("%w", ErrUserNotFound)
	}

	userInfo, ok := value.(repo.User)
	if !ok {
		log.Println("invalid type for user info for current session")
		return UserInfo{}, fmt.Errorf("%w", ErrUserNotFound)
	}

	return UserInfo{
		ID:        userInfo.ID,
		Email:     userInfo.Email,
		CreatedAt: userInfo.CreatedAt,
	}, nil
}

func (s *svc) UpdateUserInfo(ctx context.Context, updateInfo repo.UpdateUserByIdParams) (UserInfo, error) {
	value := ctx.Value("user")
	if value == nil {
		log.Println("No user info found that matches the session")
		return UserInfo{}, fmt.Errorf("%w", ErrUserNotFound)
	}

	userInfo, ok := value.(repo.User)
	if !ok {
		log.Println("invalid type for user info for current session")
		return UserInfo{}, fmt.Errorf("%w", ErrUserNotFound)
	}

	if userInfo.ID != updateInfo.ID {
		log.Println("no permission to delete")
		return UserInfo{}, fmt.Errorf("%w", ErrNoPermission)
	}

	updateResult, err := s.repo.UpdateUserById(ctx, updateInfo)
	if err != nil {
		log.Println("update failed")
		return UserInfo{}, fmt.Errorf("%w", ErrUpdateFailure)
	}

	return UserInfo{
		ID:        updateResult.ID,
		Email:     updateResult.Email,
		CreatedAt: updateResult.CreatedAt,
	}, nil
}

func (s *svc) DeleteUser(ctx context.Context) error {
	value := ctx.Value("user")
	if value == nil {
		log.Println("No user info found that matches the session")
		return fmt.Errorf("%w", ErrUserNotFound)
	}

	userInfo, ok := value.(repo.User)
	if !ok {
		log.Println("invalid type for user info for current session")
		return fmt.Errorf("%w", ErrUserNotFound)
	}

	err := s.repo.DeleteUserById(ctx, userInfo.ID)
	if err != nil {
		log.Println("user delete failed")
		return fmt.Errorf("%w", ErrDeleteFailure)
	}
	return nil
}

func hashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}
