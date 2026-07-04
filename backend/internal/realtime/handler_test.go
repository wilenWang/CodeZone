package realtime

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/websocket"

	"vibework-chat/backend/internal/config"
	"vibework-chat/backend/internal/httpx"
)

func TestIsAllowedOrigin(t *testing.T) {
	tests := []struct {
		name          string
		requestHost   string
		origin        string
		allowedOrigin string
		want          bool
	}{
		{
			name:        "allows missing origin",
			requestHost: "localhost:8080",
			want:        true,
		},
		{
			name:          "allows configured origin",
			requestHost:   "localhost:8080",
			origin:        "http://localhost:5173",
			allowedOrigin: "http://localhost:5173",
			want:          true,
		},
		{
			name:        "allows same origin host",
			requestHost: "localhost:8080",
			origin:      "http://localhost:8080",
			want:        true,
		},
		{
			name:        "rejects different origin",
			requestHost: "localhost:8080",
			origin:      "http://localhost:5173",
			want:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest("GET", "http://"+tt.requestHost+"/api/ws", nil)
			if tt.origin != "" {
				request.Header.Set("Origin", tt.origin)
			}

			if got := isAllowedOrigin(request, tt.allowedOrigin); got != tt.want {
				t.Fatalf("got %t want %t", got, tt.want)
			}
		})
	}
}

func TestServeWSDeliversEventsForAuthenticatedUser(t *testing.T) {
	hub := NewHub()
	handler := NewHandler(hub, config.Config{})

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler.ServeWS(w, r.WithContext(httpx.WithUserID(r.Context(), 42)))
	}))
	defer server.Close()

	wsURL := "ws" + server.URL[len("http"):]
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("Dial returned error: %v", err)
	}
	defer conn.Close()

	hub.SendToUser(42, Event{Type: "message.created", Payload: "ok"})

	if err := conn.SetReadDeadline(time.Now().Add(time.Second)); err != nil {
		t.Fatalf("SetReadDeadline returned error: %v", err)
	}
	var got Event
	if err := conn.ReadJSON(&got); err != nil {
		t.Fatalf("ReadJSON returned error: %v", err)
	}
	if got.Type != "message.created" || got.Payload != "ok" {
		t.Fatalf("got event %#v want message.created payload ok", got)
	}
}
