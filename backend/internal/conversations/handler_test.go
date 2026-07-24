package conversations

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"codezone/backend/internal/httpx"
)

func TestCreateMapsValidationErrorsToBadRequest(t *testing.T) {
	handler := NewHandler(NewService(&fakeRepo{}))
	req := httptest.NewRequest(http.MethodPost, "/api/conversations", strings.NewReader(`{"type":"channel","memberIds":[2]}`))
	req = req.WithContext(httpx.WithUserID(req.Context(), 1))
	rec := httptest.NewRecorder()

	handler.Create(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", rec.Code)
	}
}
