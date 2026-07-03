package main

import (
	"context"
	"database/sql"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"

	"vibework-chat/backend/internal/auth"
	"vibework-chat/backend/internal/config"
	"vibework-chat/backend/internal/httpx"
)

type fakeSessionLookup struct {
	users map[string]int64
}

func (f fakeSessionLookup) UserIDByToken(_ context.Context, tokenHash string) (int64, error) {
	userID, ok := f.users[tokenHash]
	if !ok {
		return 0, sql.ErrNoRows
	}
	return userID, nil
}

func TestRouterDevLoginRequiresExplicitEnablement(t *testing.T) {
	devLogin := func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}
	registerProtected := func(r chi.Router) {
		r.Get("/api/users", func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusNoContent)
		})
	}
	sessions := fakeSessionLookup{users: map[string]int64{}}

	disabled := buildRouter(config.Config{EnableDevLogin: false}, devLogin, sessions, registerProtected)
	req := httptest.NewRequest(http.MethodPost, "/api/auth/dev-login", nil)
	rec := httptest.NewRecorder()
	disabled.ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound && rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected dev-login to be unavailable when disabled, got %d", rec.Code)
	}

	enabled := buildRouter(config.Config{EnableDevLogin: true}, devLogin, sessions, registerProtected)
	rec = httptest.NewRecorder()
	enabled.ServeHTTP(rec, req)
	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected dev-login to be available when enabled, got %d", rec.Code)
	}
}

func TestProtectedRoutesRequireBearerToken(t *testing.T) {
	const secret = "test-secret"
	tokenHash := auth.HashToken("valid-token", secret)
	sessions := fakeSessionLookup{users: map[string]int64{tokenHash: 7}}
	router := buildRouter(config.Config{SessionSecret: secret}, nil, sessions, func(r chi.Router) {
		r.Get("/api/users", func(w http.ResponseWriter, r *http.Request) {
			userID, ok := httpx.UserID(r.Context())
			if !ok || userID != 7 {
				t.Fatalf("expected user 7 in context, got %d ok=%v", userID, ok)
			}
			w.WriteHeader(http.StatusNoContent)
		})
	})

	req := httptest.NewRequest(http.MethodGet, "/api/users", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected missing token to return 401, got %d", rec.Code)
	}

	req = httptest.NewRequest(http.MethodGet, "/api/users", nil)
	req.Header.Set("Authorization", "Bearer")
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected malformed bearer to return 401, got %d", rec.Code)
	}

	req = httptest.NewRequest(http.MethodGet, "/api/users", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected valid token to pass, got %d", rec.Code)
	}
}

func TestProtectedRoutesAcceptQueryToken(t *testing.T) {
	const secret = "test-secret"
	tokenHash := auth.HashToken("query-token", secret)
	sessions := fakeSessionLookup{users: map[string]int64{tokenHash: 7}}
	router := buildRouter(config.Config{SessionSecret: secret}, nil, sessions, func(r chi.Router) {
		r.Get("/api/users", func(w http.ResponseWriter, r *http.Request) {
			userID, ok := httpx.UserID(r.Context())
			if !ok || userID != 7 {
				t.Fatalf("expected user 7 in context, got %d ok=%v", userID, ok)
			}
			w.WriteHeader(http.StatusNoContent)
		})
	})

	req := httptest.NewRequest(http.MethodGet, "/api/users?access_token=query-token", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected query token to pass, got %d", rec.Code)
	}
}

func TestExtractAccessToken(t *testing.T) {
	tests := []struct {
		name string
		auth string
		path string
		want string
		ok   bool
	}{
		{name: "missing", path: "/api/users", ok: false},
		{name: "malformed bearer", auth: "Bearer", path: "/api/users", ok: false},
		{name: "wrong scheme", auth: "Token abc", path: "/api/users", ok: false},
		{name: "bearer", auth: "Bearer abc", path: "/api/users", want: "abc", ok: true},
		{name: "query", path: "/api/users?access_token=abc", want: "abc", ok: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			if tt.auth != "" {
				req.Header.Set("Authorization", tt.auth)
			}
			got, ok := extractAccessToken(req)
			if ok != tt.ok || got != tt.want {
				t.Fatalf("got token=%q ok=%v want token=%q ok=%v", got, ok, tt.want, tt.ok)
			}
		})
	}
}
