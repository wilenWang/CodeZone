package realtime

import (
	"net/http"
	"net/url"

	"github.com/gorilla/websocket"

	"vibework-chat/backend/internal/config"
	"vibework-chat/backend/internal/httpx"
)

type Handler struct {
	hub      *Hub
	upgrader websocket.Upgrader
}

func NewHandler(hub *Hub, cfg config.Config) *Handler {
	return &Handler{
		hub: hub,
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return isAllowedOrigin(r, cfg.CORSOrigin)
			},
		},
	}
}

func isAllowedOrigin(r *http.Request, allowedOrigin string) bool {
	origin := r.Header.Get("Origin")
	if origin == "" || origin == allowedOrigin {
		return true
	}

	originURL, err := url.Parse(origin)
	if err != nil {
		return false
	}
	return originURL.Host == r.Host
}

func (h *Handler) ServeWS(w http.ResponseWriter, r *http.Request) {
	userID, ok := httpx.UserID(r.Context())
	if !ok {
		httpx.WriteError(w, http.StatusUnauthorized, "unauthorized", "Login required")
		return
	}

	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	events := make(chan Event, 16)
	h.hub.Register(userID, events)
	defer h.hub.Unregister(userID, events)

	done := make(chan struct{})
	go func() {
		defer close(done)
		for {
			if _, _, err := conn.NextReader(); err != nil {
				return
			}
		}
	}()

	for {
		select {
		case event := <-events:
			if err := conn.WriteJSON(event); err != nil {
				return
			}
		case <-done:
			return
		case <-r.Context().Done():
			return
		}
	}
}
