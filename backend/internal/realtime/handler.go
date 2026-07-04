package realtime

import (
	"context"
	"net/http"
	"net/url"
	"time"

	"github.com/gorilla/websocket"

	"vibework-chat/backend/internal/config"
	"vibework-chat/backend/internal/httpx"
)

const websocketWriteTimeout = 10 * time.Second

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

	subscription := h.hub.Register(userID)
	defer h.hub.Unregister(userID, subscription)

	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()
	go func() {
		defer cancel()
		defer conn.Close()
		for {
			if _, _, err := conn.NextReader(); err != nil {
				return
			}
		}
	}()

	for {
		event, ok := subscription.Next(ctx)
		if !ok {
			return
		}
		if err := conn.SetWriteDeadline(time.Now().Add(websocketWriteTimeout)); err != nil {
			return
		}
		if err := conn.WriteJSON(event); err != nil {
			return
		}
	}
}
