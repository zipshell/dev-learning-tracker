package auth

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
	repo "github.com/zipshell/dev-learning-tracker/internal/adapters/postgresql/sqlc"
	"github.com/zipshell/dev-learning-tracker/internal/env"
	"github.com/zipshell/dev-learning-tracker/internal/tokens"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Login(ctx context.Context, userInfo UserCredentials) (Tokens, error)
}

type svc struct {
	repo repo.Querier
}

func NewAuthService(repo repo.Querier) AuthService {
	return &svc{
		repo: repo,
	}
}

const maxRefreshTokens = 3

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrInvalidCredential = errors.New("invalid credentials")
)

type Tokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (s *svc) Login(ctx context.Context, userCredentials UserCredentials) (Tokens, error) {
	existingUserInfo, err := s.repo.FindUserByEmail(ctx, userCredentials.Email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Tokens{}, fmt.Errorf("%w", ErrUserNotFound)
		}
		return Tokens{}, fmt.Errorf("user lookup failed: %w", err)
	}
	err = bcrypt.CompareHashAndPassword([]byte(existingUserInfo.Password), []byte(userCredentials.Password))
	if err != nil {
		return Tokens{}, fmt.Errorf("%w", ErrInvalidCredential)
	}

	existingRefreshToken, err := s.repo.FindRefreshTokensByUserId(ctx, existingUserInfo.ID)
	if err != nil && err != pgx.ErrNoRows {
		return Tokens{}, fmt.Errorf("Failed querying tokens table")
	}
	if len(existingRefreshToken) >= maxRefreshTokens {
		redundantCount := len(existingRefreshToken) - maxRefreshTokens + 1
		tokensToDelete := existingRefreshToken[:redundantCount]

		// Extract IDs
		tokenIDs := make([]int64, len(tokensToDelete))
		for i, token := range tokensToDelete {
			tokenIDs[i] = token.ID
		}

		err = s.repo.DeleteRefreshTokensByIds(ctx, tokenIDs)
		if err != nil {
			log.Println(err)
		}
	}
	newRefreshToken, err := tokens.GenerateOpaqueToken()
	if err != nil {
		log.Printf("Generate Token failed: %v", err)
		return Tokens{}, fmt.Errorf("create refresh token: %w", err)
	}

	newRefreshTokenCreation, err := s.repo.CreateRefreshToken(ctx, repo.CreateRefreshTokenParams{
		Token:  newRefreshToken,
		UserID: existingUserInfo.ID,
	})
	if err != nil {
		log.Printf("Inserting refresh token to database failed: %v", err)
		return Tokens{}, fmt.Errorf("create refresh token: %w", err)
	}

	secret := env.GetString("JWT_SECRET", "random string")

	newJwt, err := tokens.CreateJwt([]byte(secret), existingUserInfo.ID)
	if err != nil {
		log.Printf("Access token creation failed: %v", err)
		return Tokens{}, fmt.Errorf("create access token: %w", err)
	}

	return Tokens{
		AccessToken:  newJwt,
		RefreshToken: newRefreshTokenCreation.Token,
	}, nil
}
