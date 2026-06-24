package users

import (
	"context"
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	repo "github.com/zipshell/dev-learning-tracker/internal/adapters/postgresql/sqlc"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	CreateUser(ctx context.Context, newUser repo.CreateUserParams) error
}

var ErrEmailAlreadyExists = errors.New("email already exists")

func NewUserService(db *pgx.Conn) UserService {
	return &svc{db: db}
}

type Claims struct {
	UserID int64 `json:"user_id"`
	jwt.RegisteredClaims
}

type svc struct {
	db *pgx.Conn
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

func hashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}
