package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"vibework-chat/backend/internal/config"
)

func TestRouterDevLoginRequiresExplicitEnablement(t *testing.T) {
	devLogin := func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}
	listUsers := func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}

	disabled := buildRouter(config.Config{EnableDevLogin: false}, devLogin, listUsers)
	req := httptest.NewRequest(http.MethodPost, "/api/auth/dev-login", nil)
	rec := httptest.NewRecorder()
	disabled.ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound && rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected dev-login to be unavailable when disabled, got %d", rec.Code)
	}

	enabled := buildRouter(config.Config{EnableDevLogin: true}, devLogin, listUsers)
	rec = httptest.NewRecorder()
	enabled.ServeHTTP(rec, req)
	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected dev-login to be available when enabled, got %d", rec.Code)
	}
}
