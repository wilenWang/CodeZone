package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"time"
)

var ErrInvalidCredentials = errors.New("invalid credentials")

type User struct {
	ID          int64   `json:"id"`
	WorkspaceID int64   `json:"workspaceId"`
	Username    string  `json:"username"`
	DisplayName string  `json:"displayName"`
	AvatarURL   *string `json:"avatarUrl"`
	UserType    string  `json:"userType"`
}

type UserFinder interface {
	FindByWorkspaceUsername(ctx context.Context, workspaceID int64, username string) (User, error)
}

type SessionStore interface {
	Create(ctx context.Context, userID int64, tokenHash string, expiresAt time.Time) error
}

type Service struct {
	users       UserFinder
	sessions    SessionStore
	secret      string
	workspaceID int64
}

type LoginResult struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

func NewService(users UserFinder, sessions SessionStore, secret string, workspaceID int64) *Service {
	return &Service{
		users:       users,
		sessions:    sessions,
		secret:      secret,
		workspaceID: workspaceID,
	}
}

func (s *Service) DevLogin(ctx context.Context, username string) (LoginResult, error) {
	if username == "" {
		return LoginResult{}, ErrInvalidCredentials
	}

	user, err := s.users.FindByWorkspaceUsername(ctx, s.workspaceID, username)
	if errors.Is(err, sql.ErrNoRows) {
		return LoginResult{}, ErrInvalidCredentials
	}
	if err != nil {
		return LoginResult{}, err
	}
	if user.ID == 0 {
		return LoginResult{}, ErrInvalidCredentials
	}

	token, err := newToken()
	if err != nil {
		return LoginResult{}, err
	}
	expiresAt := time.Now().Add(24 * time.Hour)
	if err := s.sessions.Create(ctx, user.ID, HashToken(token, s.secret), expiresAt); err != nil {
		return LoginResult{}, err
	}

	return LoginResult{Token: token, User: user}, nil
}

func newToken() (string, error) {
	var raw [32]byte
	if _, err := rand.Read(raw[:]); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(raw[:]), nil
}

func HashToken(token string, secret string) string {
	sum := sha256.Sum256([]byte(token + secret))
	return hex.EncodeToString(sum[:])
}
