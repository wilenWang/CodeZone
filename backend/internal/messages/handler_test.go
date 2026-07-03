package messages

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
)

type fakeMessageRepo struct {
	listCalled bool
}

func (f *fakeMessageRepo) Create(ctx context.Context, input SendInput, contentPlain string) (Message, error) {
	return Message{}, nil
}

func (f *fakeMessageRepo) ListBefore(ctx context.Context, conversationID int64, beforeID int64, limit int) ([]Message, error) {
	f.listCalled = true
	return []Message{}, nil
}

func (f *fakeMessageRepo) MarkRead(ctx context.Context, conversationID int64, userID int64) error {
	return nil
}

func TestListMessagesRequiresUserContext(t *testing.T) {
	repo := &fakeMessageRepo{}
	handler := NewHandler(NewService(repo))
	router := chi.NewRouter()
	router.Get("/api/conversations/{id}/messages", handler.ListMessages)

	req := httptest.NewRequest(http.MethodGet, "/api/conversations/42/messages", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", rec.Code)
	}
	if repo.listCalled {
		t.Fatal("expected list repository not to be called")
	}
}
