package messages

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"codezone/backend/internal/httpx"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) ListMessages(w http.ResponseWriter, r *http.Request) {
	userID, ok := httpx.UserID(r.Context())
	if !ok {
		httpx.WriteError(w, http.StatusUnauthorized, "unauthorized", "Login required")
		return
	}

	conversationID, err := parsePositiveID(chi.URLParam(r, "id"))
	if err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "bad_conversation_id", "Invalid conversation id")
		return
	}
	beforeID, err := parseOptionalPositiveInt64(r.URL.Query().Get("before"))
	if err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "bad_before", "Invalid before message id")
		return
	}
	limit, err := parseOptionalLimit(r.URL.Query().Get("limit"))
	if err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "bad_limit", "Invalid limit")
		return
	}
	messages, err := h.service.ListBefore(r.Context(), conversationID, userID, beforeID, limit)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			httpx.WriteError(w, http.StatusNotFound, "conversation_not_found", "Conversation not found")
			return
		}
		httpx.WriteError(w, http.StatusInternalServerError, "messages_failed", "Could not load messages")
		return
	}
	httpx.WriteJSON(w, http.StatusOK, map[string]any{"messages": messages})
}

func (h *Handler) SendMessage(w http.ResponseWriter, r *http.Request) {
	userID, ok := httpx.UserID(r.Context())
	if !ok {
		httpx.WriteError(w, http.StatusUnauthorized, "unauthorized", "Login required")
		return
	}
	conversationID, err := parsePositiveID(chi.URLParam(r, "id"))
	if err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "bad_conversation_id", "Invalid conversation id")
		return
	}
	var req struct {
		ContentMarkdown string `json:"contentMarkdown"`
	}
	if err := httpx.ReadJSON(r, &req); err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "bad_json", "Invalid JSON body")
		return
	}
	message, err := h.service.Send(r.Context(), SendInput{
		ConversationID:  conversationID,
		SenderID:        userID,
		ContentMarkdown: req.ContentMarkdown,
	})
	if err != nil {
		if errors.Is(err, ErrEmptyMessage) {
			httpx.WriteError(w, http.StatusBadRequest, "message_invalid", err.Error())
			return
		}
		if errors.Is(err, ErrNotFound) {
			httpx.WriteError(w, http.StatusNotFound, "conversation_not_found", "Conversation not found")
			return
		}
		httpx.WriteError(w, http.StatusInternalServerError, "message_failed", "Could not send message")
		return
	}
	httpx.WriteJSON(w, http.StatusCreated, message)
}

func (h *Handler) MarkRead(w http.ResponseWriter, r *http.Request) {
	userID, ok := httpx.UserID(r.Context())
	if !ok {
		httpx.WriteError(w, http.StatusUnauthorized, "unauthorized", "Login required")
		return
	}
	conversationID, err := parsePositiveID(chi.URLParam(r, "id"))
	if err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "bad_conversation_id", "Invalid conversation id")
		return
	}
	if err := h.service.MarkRead(r.Context(), conversationID, userID); err != nil {
		if errors.Is(err, ErrNotFound) {
			httpx.WriteError(w, http.StatusNotFound, "conversation_not_found", "Conversation not found")
			return
		}
		httpx.WriteError(w, http.StatusInternalServerError, "read_failed", "Could not mark conversation read")
		return
	}
	httpx.WriteJSON(w, http.StatusOK, map[string]any{"ok": true})
}

func parsePositiveID(raw string) (int64, error) {
	id, err := strconv.ParseInt(raw, 10, 64)
	if err != nil || id <= 0 {
		return 0, strconv.ErrSyntax
	}
	return id, nil
}

func parseOptionalPositiveInt64(raw string) (int64, error) {
	if raw == "" {
		return 0, nil
	}
	value, err := strconv.ParseInt(raw, 10, 64)
	if err != nil || value <= 0 {
		return 0, strconv.ErrSyntax
	}
	return value, nil
}

func parseOptionalLimit(raw string) (int, error) {
	if raw == "" {
		return 0, nil
	}
	limit, err := strconv.Atoi(raw)
	if err != nil || limit < 1 || limit > 100 {
		return 0, strconv.ErrSyntax
	}
	return limit, nil
}
