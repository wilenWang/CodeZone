package auth

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

type fakeUserFinder struct {
	users map[string]User
	err   error
}

func (f fakeUserFinder) FindByUsername(_ context.Context, username string) (User, error) {
	if f.err != nil {
		return User{}, f.err
	}
	user, ok := f.users[username]
	if !ok {
		return User{}, sql.ErrNoRows
	}
	return user, nil
}

type fakeSessionStore struct {
	userID    int64
	tokenHash string
	expiresAt time.Time
}

func (f *fakeSessionStore) Create(_ context.Context, userID int64, tokenHash string, expiresAt time.Time) error {
	f.userID = userID
	f.tokenHash = tokenHash
	f.expiresAt = expiresAt
	return nil
}

func TestDevLoginCreatesSessionForSeedUser(t *testing.T) {
	ctx := context.Background()
	sessions := &fakeSessionStore{}
	users := fakeUserFinder{users: map[string]User{
		"alice": {
			ID:          1,
			WorkspaceID: 1,
			Username:    "alice",
			DisplayName: "Alice Chen",
			UserType:    "human",
		},
	}}
	svc := NewService(users, sessions, "test-secret")

	result, err := svc.DevLogin(ctx, "alice")
	if err != nil {
		t.Fatalf("DevLogin returned error: %v", err)
	}

	if result.Token == "" {
		t.Fatal("expected token")
	}
	if result.User.ID != 1 || result.User.Username != "alice" {
		t.Fatalf("unexpected user: %+v", result.User)
	}
	if sessions.userID != 1 {
		t.Fatalf("expected session for user 1, got %d", sessions.userID)
	}
	if sessions.tokenHash == "" {
		t.Fatal("expected token hash")
	}
	if sessions.tokenHash == result.Token {
		t.Fatal("expected stored token hash to differ from raw token")
	}
	if !sessions.expiresAt.After(time.Now().Add(23 * time.Hour)) {
		t.Fatalf("expected expiry about 24h in future, got %s", sessions.expiresAt)
	}
	if got := HashToken(result.Token, "test-secret"); got != sessions.tokenHash {
		t.Fatalf("expected stored hash %q, got %q", got, sessions.tokenHash)
	}
}

func TestDevLoginReturnsInvalidCredentialsForMissingUser(t *testing.T) {
	svc := NewService(fakeUserFinder{users: map[string]User{}}, &fakeSessionStore{}, "test-secret")

	_, err := svc.DevLogin(context.Background(), "missing")

	if !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("expected ErrInvalidCredentials, got %v", err)
	}
}

func TestDevLoginPropagatesLookupFailure(t *testing.T) {
	lookupErr := errors.New("lookup failed")
	svc := NewService(fakeUserFinder{err: lookupErr}, &fakeSessionStore{}, "test-secret")

	_, err := svc.DevLogin(context.Background(), "alice")

	if !errors.Is(err, lookupErr) {
		t.Fatalf("expected lookup error, got %v", err)
	}
}

func TestDevLoginHandlerWritesInvalidCredentialsCode(t *testing.T) {
	handler := NewHandler(NewService(fakeUserFinder{users: map[string]User{}}, &fakeSessionStore{}, "test-secret"))
	req := httptest.NewRequest(http.MethodPost, "/api/auth/dev-login", strings.NewReader(`{"username":"missing"}`))
	rec := httptest.NewRecorder()

	handler.DevLogin(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", rec.Code)
	}
	var body map[string]any
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if body["code"] != "invalid_credentials" {
		t.Fatalf("expected invalid_credentials code, got %v", body["code"])
	}
}

func TestDevLoginHandlerWritesAuthFailedCodeForLookupFailure(t *testing.T) {
	handler := NewHandler(NewService(fakeUserFinder{err: errors.New("lookup failed")}, &fakeSessionStore{}, "test-secret"))
	req := httptest.NewRequest(http.MethodPost, "/api/auth/dev-login", strings.NewReader(`{"username":"alice"}`))
	rec := httptest.NewRecorder()

	handler.DevLogin(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", rec.Code)
	}
	var body map[string]any
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if body["code"] != "auth_failed" {
		t.Fatalf("expected auth_failed code, got %v", body["code"])
	}
}
