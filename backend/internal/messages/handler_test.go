package messages

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"

	"codezone/backend/internal/httpx"
)

type fakeMessageRepo struct {
	createErr      error
	listErr        error
	markReadErr    error
	listCalled     bool
	listUserID     int64
	createSenderID int64
}

func (f *fakeMessageRepo) Create(ctx context.Context, input SendInput, contentPlain string) (Message, error) {
	f.createSenderID = input.SenderID
	if f.createErr != nil {
		return Message{}, f.createErr
	}
	return Message{ID: 1, ConversationID: input.ConversationID, SenderID: input.SenderID}, nil
}

func (f *fakeMessageRepo) ListBefore(ctx context.Context, conversationID int64, userID int64, beforeID int64, limit int) ([]Message, error) {
	f.listCalled = true
	f.listUserID = userID
	if f.listErr != nil {
		return nil, f.listErr
	}
	return []Message{}, nil
}

func (f *fakeMessageRepo) MarkRead(ctx context.Context, conversationID int64, userID int64) error {
	if f.markReadErr != nil {
		return f.markReadErr
	}
	return nil
}

func TestListMessagesRequiresUserContext(t *testing.T) {
	repo := &fakeMessageRepo{}
	router := messageRouter(NewHandler(NewService(repo)))

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

func TestListMessagesPassesUserID(t *testing.T) {
	repo := &fakeMessageRepo{}
	router := messageRouter(NewHandler(NewService(repo)))

	req := httptest.NewRequest(http.MethodGet, "/api/conversations/42/messages", nil)
	req = req.WithContext(httpx.WithUserID(req.Context(), 7))
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}
	if repo.listUserID != 7 {
		t.Fatalf("expected list user 7, got %d", repo.listUserID)
	}
}

func TestListMessagesMapsNotFound(t *testing.T) {
	repo := &fakeMessageRepo{listErr: ErrNotFound}
	router := messageRouter(NewHandler(NewService(repo)))

	req := httptest.NewRequest(http.MethodGet, "/api/conversations/42/messages", nil)
	req = req.WithContext(httpx.WithUserID(req.Context(), 7))
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", rec.Code)
	}
}

func TestSendMessageMapsNotFound(t *testing.T) {
	repo := &fakeMessageRepo{createErr: ErrNotFound}
	router := messageRouter(NewHandler(NewService(repo)))

	req := httptest.NewRequest(http.MethodPost, "/api/conversations/42/messages", strings.NewReader(`{"contentMarkdown":"hi"}`))
	req = req.WithContext(httpx.WithUserID(req.Context(), 7))
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", rec.Code)
	}
	if repo.createSenderID != 7 {
		t.Fatalf("expected sender 7, got %d", repo.createSenderID)
	}
}

func TestMarkReadMapsNotFound(t *testing.T) {
	repo := &fakeMessageRepo{markReadErr: ErrNotFound}
	router := messageRouter(NewHandler(NewService(repo)))

	req := httptest.NewRequest(http.MethodPost, "/api/conversations/42/read", nil)
	req = req.WithContext(httpx.WithUserID(req.Context(), 7))
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", rec.Code)
	}
}

func TestMessageHandlerRejectsInvalidIDsAndQuery(t *testing.T) {
	handler := NewHandler(NewService(&fakeMessageRepo{}))
	tests := []struct {
		name   string
		method string
		path   string
	}{
		{name: "list zero id", method: http.MethodGet, path: "/api/conversations/0/messages"},
		{name: "send negative id", method: http.MethodPost, path: "/api/conversations/-1/messages"},
		{name: "read bad id", method: http.MethodPost, path: "/api/conversations/nope/read"},
		{name: "bad before", method: http.MethodGet, path: "/api/conversations/42/messages?before=nope"},
		{name: "bad limit", method: http.MethodGet, path: "/api/conversations/42/messages?limit=nope"},
		{name: "zero limit", method: http.MethodGet, path: "/api/conversations/42/messages?limit=0"},
		{name: "too large limit", method: http.MethodGet, path: "/api/conversations/42/messages?limit=101"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, strings.NewReader(`{"contentMarkdown":"hi"}`))
			req = req.WithContext(httpx.WithUserID(req.Context(), 7))
			rec := httptest.NewRecorder()
			messageRouter(handler).ServeHTTP(rec, req)
			if rec.Code != http.StatusBadRequest {
				t.Fatalf("expected status 400, got %d", rec.Code)
			}
		})
	}
}

func TestMessageHandlerKeepsUnexpectedErrorsAsServerErrors(t *testing.T) {
	repo := &fakeMessageRepo{listErr: errors.New("database down")}
	router := messageRouter(NewHandler(NewService(repo)))

	req := httptest.NewRequest(http.MethodGet, "/api/conversations/42/messages", nil)
	req = req.WithContext(httpx.WithUserID(req.Context(), 7))
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", rec.Code)
	}
}

func messageRouter(handler *Handler) http.Handler {
	router := chi.NewRouter()
	router.Get("/api/conversations/{id}/messages", handler.ListMessages)
	router.Post("/api/conversations/{id}/messages", handler.SendMessage)
	router.Post("/api/conversations/{id}/read", handler.MarkRead)
	return router
}
