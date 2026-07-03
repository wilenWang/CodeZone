package users

import (
	"net/http"

	"vibework-chat/backend/internal/httpx"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	users, err := h.service.List(r.Context(), 1)
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, "users_list_failed", "failed to list users")
		return
	}
	httpx.WriteJSON(w, http.StatusOK, map[string]any{"users": users})
}
