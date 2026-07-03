package conversations

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"vibework-chat/backend/internal/httpx"
)

type fakeLister struct{}

func (fakeLister) ListForUser(ctx context.Context, workspaceID int64, userID int64) ([]Conversation, error) {
	return []Conversation{}, nil
}

func TestCreateMapsValidationErrorsToBadRequest(t *testing.T) {
	handler := NewHandler(NewService(&fakeRepo{}), fakeLister{})
	req := httptest.NewRequest(http.MethodPost, "/api/conversations", strings.NewReader(`{"type":"channel","memberIds":[2]}`))
	req = req.WithContext(httpx.WithUserID(req.Context(), 1))
	rec := httptest.NewRecorder()

	handler.Create(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", rec.Code)
	}
}
