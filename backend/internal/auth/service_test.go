package auth

import (
	"context"
	"testing"
	"time"
)

type fakeUserFinder struct {
	users map[string]User
}

func (f fakeUserFinder) FindByUsername(_ context.Context, username string) (User, error) {
	return f.users[username], nil
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
