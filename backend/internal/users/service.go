package users

import (
	"context"
	"fmt"
	"log"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
	repo "github.com/zipshell/dev-learning-tracker/internal/adapters/postgresql/sqlc"
	"github.com/zipshell/dev-learning-tracker/internal/env"
	"github.com/zipshell/dev-learning-tracker/internal/tokens"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	CreateUser(ctx context.Context, newUser repo.CreateUserParams) (UserWithToken, error)
}

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

type UserWithToken struct {
	*repo.User   `json:"user"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (s *svc) CreateUser(ctx context.Context, newUser repo.CreateUserParams) (UserWithToken, error) {
	tx, err := s.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return UserWithToken{}, err
	}

	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()

	hashedPassword, err := hashPassword(newUser.Password)
	if err != nil {
		return UserWithToken{}, err
	}

	q := repo.New(tx)

	newUser.Password = hashedPassword
	newUserCreation, err := q.CreateUser(ctx, newUser)
	if err != nil {
		return UserWithToken{}, fmt.Errorf("New user creation failed in write to db step")
	}

	newRefreshToken, err := tokens.GenerateOpaqueToken()
	if err != nil {
		log.Printf("Generate Token failed: %v", err)
		return UserWithToken{}, fmt.Errorf("create refresh token: %w", err)
	}

	newTokenCreation, err := q.CreateRefreshToken(ctx, repo.CreateRefreshTokenParams{
		Token:  newRefreshToken,
		UserID: newUserCreation.ID,
	})
	if err != nil {
		log.Printf("Inserting refresh token to database failed: %v", err)
		return UserWithToken{}, fmt.Errorf("create refresh token: %w", err)
	}

	if err = tx.Commit(ctx); err != nil {
		return UserWithToken{}, err
	}

	secret := env.GetString("JWT_SECRET", "random string")

	newJwt, err := tokens.CreateJwt([]byte(secret), newUserCreation.ID)
	if err != nil {
		log.Printf("CreateToken failed: %v", err)
		return UserWithToken{}, fmt.Errorf("create access token: %w", err)
	}

	return UserWithToken{
		User:         &newUserCreation,
		AccessToken:  newJwt,
		RefreshToken: newTokenCreation.Token,
	}, nil
}

func hashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}
