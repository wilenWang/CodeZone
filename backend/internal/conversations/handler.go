package conversations

import (
	"net/http"

	"codezone/backend/internal/httpx"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	userID, ok := httpx.UserID(r.Context())
	if !ok {
		httpx.WriteError(w, http.StatusUnauthorized, "unauthorized", "Login required")
		return
	}
	conversations, err := h.service.ListForUser(r.Context(), 1, userID)
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, "conversations_failed", "Could not load conversations")
		return
	}
	httpx.WriteJSON(w, http.StatusOK, map[string]any{"conversations": conversations})
}

func (h *Handler) EnsureDirect(w http.ResponseWriter, r *http.Request) {
	userID, ok := httpx.UserID(r.Context())
	if !ok {
		httpx.WriteError(w, http.StatusUnauthorized, "unauthorized", "Login required")
		return
	}
	var req struct {
		UserID int64 `json:"userId"`
	}
	if err := httpx.ReadJSON(r, &req); err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "bad_json", "Invalid JSON body")
		return
	}
	conversation, err := h.service.EnsureDirect(r.Context(), 1, userID, req.UserID)
	if err != nil {
		if isValidationError(err) {
			httpx.WriteError(w, http.StatusBadRequest, "direct_conversation_invalid", err.Error())
			return
		}
		httpx.WriteError(w, http.StatusInternalServerError, "direct_conversation_failed", "Could not open direct conversation")
		return
	}
	httpx.WriteJSON(w, http.StatusOK, conversation)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := httpx.UserID(r.Context())
	if !ok {
		httpx.WriteError(w, http.StatusUnauthorized, "unauthorized", "Login required")
		return
	}
	var req struct {
		Type      string  `json:"type"`
		Title     *string `json:"title"`
		MemberIDs []int64 `json:"memberIds"`
	}
	if err := httpx.ReadJSON(r, &req); err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "bad_json", "Invalid JSON body")
		return
	}
	input := CreateConversationInput{
		WorkspaceID: 1,
		CreatedBy:   userID,
		Type:        req.Type,
		Title:       req.Title,
		MemberIDs:   append(req.MemberIDs, userID),
	}
	conversation, err := h.service.Create(r.Context(), input)
	if err != nil {
		if isValidationError(err) {
			httpx.WriteError(w, http.StatusBadRequest, "conversation_invalid", err.Error())
			return
		}
		httpx.WriteError(w, http.StatusInternalServerError, "conversation_failed", "Could not create conversation")
		return
	}
	httpx.WriteJSON(w, http.StatusCreated, conversation)
}
