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
	Logout(ctx context.Context) error
	ValidateSession(ctx context.Context, sessionToken string) (UserWithToken, error)
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

type UserWithToken struct {
	User         repo.User
	SessionToken string
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

func (s *svc) Logout(ctx context.Context) error {
	value := ctx.Value("session_token")
	if value == nil {
		log.Println("No session token found")
		return nil
	}

	sessionToken, ok := value.(string)
	if !ok {
		log.Println("Session token is of invalid type")
		return nil
	}

	err := s.repo.DeleteSessionByToken(ctx, sessionToken)
	if err != nil {
		log.Printf("failed deleting session cookie: %v", err)
	}
	return nil
}

func (s *svc) ValidateSession(ctx context.Context, sessionToken string) (UserWithToken, error) {
	existingSession, err := s.repo.FindActiveSessionByToken(ctx, sessionToken)
	if err != nil {
		log.Printf("Invalid session token: %v", err)
		return UserWithToken{}, fmt.Errorf("invalid session token: %w", err)
	}
	sessionUser, err := s.repo.FindUserById(ctx, existingSession.UserID)
	if err != nil {
		log.Printf("Cannot find user for session token: %v", err)
		return UserWithToken{}, fmt.Errorf("invalid session token: %w", err)
	}
	return UserWithToken{
		User:         sessionUser,
		SessionToken: sessionToken,
	}, nil
}
