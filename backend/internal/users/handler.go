package users

import (
	"net/http"

	"codezone/backend/internal/httpx"
)

type Handler struct {
	service    *Service
	pathPrefix string
}

func NewHandler(service *Service, pathPrefix ...string) *Handler {
	prefix := "dev/codezone"
	if len(pathPrefix) > 0 && pathPrefix[0] != "" {
		prefix = pathPrefix[0]
	}
	return &Handler{service: service, pathPrefix: prefix}
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	users, err := h.service.List(r.Context(), 1)
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, "users_list_failed", "failed to list users")
		return
	}
	httpx.WriteJSON(w, http.StatusOK, map[string]any{"users": users})
}
