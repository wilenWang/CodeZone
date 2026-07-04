package realtime

import (
	"net/http/httptest"
	"testing"
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
