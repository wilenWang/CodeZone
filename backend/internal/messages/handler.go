package messages

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"vibework-chat/backend/internal/httpx"
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
	_ = userID

	conversationID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "bad_conversation_id", "Invalid conversation id")
		return
	}
	beforeID, _ := strconv.ParseInt(r.URL.Query().Get("before"), 10, 64)
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	messages, err := h.service.ListBefore(r.Context(), conversationID, beforeID, limit)
	if err != nil {
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
	conversationID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
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
	conversationID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "bad_conversation_id", "Invalid conversation id")
		return
	}
	if err := h.service.MarkRead(r.Context(), conversationID, userID); err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, "read_failed", "Could not mark conversation read")
		return
	}
	httpx.WriteJSON(w, http.StatusOK, map[string]any{"ok": true})
}
