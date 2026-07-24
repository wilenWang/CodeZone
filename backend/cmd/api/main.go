package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"

	"codezone/backend/internal/admin"
	"codezone/backend/internal/agent"
	"codezone/backend/internal/auth"
	"codezone/backend/internal/config"
	"codezone/backend/internal/conversations"
	"codezone/backend/internal/db"
	"codezone/backend/internal/httpx"
	"codezone/backend/internal/messages"
	"codezone/backend/internal/realtime"
	"codezone/backend/internal/users"
)

func main() {
	cfg := config.Load()

	conn, err := db.Open(cfg.MySQLDSN)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	userRepo := users.NewRepository(conn)
	userService := users.NewService(userRepo)
	sessionRepo := auth.NewSessionRepository(conn, cfg.SessionSecret)
	authService := auth.NewService(authUserFinder{repo: userRepo}, sessionRepo, cfg.SessionSecret, 1)

	authHandler := auth.NewHandler(authService)
	userHandler := users.NewHandler(userService)
	conversationRepo := conversations.NewSQLRepository(conn)
	conversationService := conversations.NewService(conversationRepo)
	conversationHandler := conversations.NewHandler(conversationService)
	messageRepo := messages.NewSQLRepository(conn)
	realtimeHub := realtime.NewHub()
	realtimeHandler := realtime.NewHandler(realtimeHub, cfg)
	realtimeNotifier := realtime.NewSQLNotifier(conn, realtimeHub)
	agentFinder := agent.NewSQLFinder(conn)
	agentRunner := agent.NewMockRunner("Mock Agent received:")
	messageService := messages.NewServiceWithNotifier(messageRepo, realtimeNotifier)
	agentOrchestrator := agent.NewOrchestrator(agentFinder, messageService, agentRunner)
	messageService = messages.NewServiceWithRealtime(messageRepo, realtimeNotifier, agentOrchestrator)
	messageHandler := messages.NewHandler(messageService)
	adminHandler := admin.NewHandler(conn)

	router := buildRouter(cfg, authHandler.DevLogin, sessionRepo, func(r chi.Router) {
		r.Get("/api/ws", realtimeHandler.ServeWS)
		r.Get("/api/admin/users", adminHandler.Users)
		r.Get("/api/admin/conversations", adminHandler.Conversations)
		r.Get("/api/admin/messages", adminHandler.Messages)
		r.Get("/api/users", userHandler.List)
		r.Get("/api/conversations", conversationHandler.List)
		r.Post("/api/conversations", conversationHandler.Create)
		r.Get("/api/conversations/{id}/messages", messageHandler.ListMessages)
		r.Post("/api/conversations/{id}/messages", messageHandler.SendMessage)
		r.Post("/api/conversations/{id}/read", messageHandler.MarkRead)
	})

	addr := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("api listening on %s", addr)
	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatal(err)
	}
}

type authUserFinder struct {
	repo *users.Repository
}

type tokenUserLookup interface {
	UserIDByToken(ctx context.Context, tokenHash string) (int64, error)
}

func buildRouter(cfg config.Config, devLogin http.HandlerFunc, sessions tokenUserLookup, registerProtected func(chi.Router)) http.Handler {
	router := chi.NewRouter()
	if cfg.EnableDevLogin {
		router.Post("/api/auth/dev-login", devLogin)
	}
	if registerProtected != nil {
		router.Group(func(r chi.Router) {
			r.Use(authMiddleware(cfg, sessions))
			registerProtected(r)
		})
	}
	return router
}

func authMiddleware(cfg config.Config, sessions tokenUserLookup) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token, ok := extractAccessToken(r)
			if !ok || sessions == nil {
				httpx.WriteError(w, http.StatusUnauthorized, "unauthorized", "Login required")
				return
			}
			userID, err := sessions.UserIDByToken(r.Context(), auth.HashToken(token, cfg.SessionSecret))
			if err != nil {
				httpx.WriteError(w, http.StatusUnauthorized, "unauthorized", "Login required")
				return
			}
			next.ServeHTTP(w, r.WithContext(httpx.WithUserID(r.Context(), userID)))
		})
	}
}

func extractAccessToken(r *http.Request) (string, bool) {
	const prefix = "Bearer "
	header := r.Header.Get("Authorization")
	if header != "" {
		if !strings.HasPrefix(header, prefix) {
			return "", false
		}
		token := strings.TrimSpace(strings.TrimPrefix(header, prefix))
		return token, token != ""
	}
	token := strings.TrimSpace(r.URL.Query().Get("access_token"))
	return token, token != ""
}

func (f authUserFinder) FindByWorkspaceUsername(ctx context.Context, workspaceID int64, username string) (auth.User, error) {
	user, err := f.repo.FindByWorkspaceUsername(ctx, workspaceID, username)
	if err != nil {
		return auth.User{}, err
	}
	return auth.User{
		ID:          user.ID,
		WorkspaceID: user.WorkspaceID,
		Username:    user.Username,
		DisplayName: user.DisplayName,
		AvatarURL:   user.AvatarURL,
		UserType:    user.UserType,
	}, nil
}
