package auth

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	repo "github.com/zipshell/dev-learning-tracker/internal/adapters/postgresql/sqlc"
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

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrInvalidCredential = errors.New("invalid credentials")
)

type Tokens struct {
	SessionToken string             `json:"session_token"`
	ExpiredAt    pgtype.Timestamptz `json:"expired_at"`
}

var maxSessionRecordNumber = 5

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

	existingSessions, err := s.repo.FindSessionsByUserId(ctx, existingUserInfo.ID)
	if len(existingSessions) >= maxSessionRecordNumber {
		redundantSessionCount := len(existingSessions) - maxSessionRecordNumber + 1
		sessionsToDelete := existingSessions[:redundantSessionCount]

		tokenIDs := make([]int64, len(sessionsToDelete))
		for i, token := range sessionsToDelete {
			tokenIDs[i] = token.ID
		}

		err := s.repo.DeleteSessionsByIds(ctx, tokenIDs)
		if err != nil {
			log.Printf("Failed cleaning up old sessions %v", err)
		}
	}

	newSessionToken, err := tokens.GenerateOpaqueToken()
	if err != nil {
		log.Printf("Generating Session Token failed: %v", err)
		return Tokens{}, fmt.Errorf("create session token: %w", err)
	}

	newSessionCreation, err := s.repo.CreateSession(ctx, repo.CreateSessionParams{
		Token:  newSessionToken,
		UserID: existingUserInfo.ID,
	})
	if err != nil {
		log.Printf("Inserting refresh token to database failed: %v", err)
		return Tokens{}, fmt.Errorf("create refresh token: %w", err)
	}

	return Tokens{
		SessionToken: newSessionCreation.Token,
		ExpiredAt:    newSessionCreation.ExpiredAt,
	}, nil
}
