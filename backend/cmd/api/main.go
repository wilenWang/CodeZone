package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"

	"vibework-chat/backend/internal/auth"
	"vibework-chat/backend/internal/config"
	"vibework-chat/backend/internal/conversations"
	"vibework-chat/backend/internal/db"
	"vibework-chat/backend/internal/httpx"
	"vibework-chat/backend/internal/messages"
	"vibework-chat/backend/internal/users"
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
	conversationHandler := conversations.NewHandler(conversationService, conversationRepo)
	messageRepo := messages.NewSQLRepository(conn)
	messageService := messages.NewService(messageRepo)
	messageHandler := messages.NewHandler(messageService)

	router := buildRouter(cfg, authHandler.DevLogin, sessionRepo, func(r chi.Router) {
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
			const prefix = "Bearer "
			header := r.Header.Get("Authorization")
			if !strings.HasPrefix(header, prefix) {
				httpx.WriteError(w, http.StatusUnauthorized, "unauthorized", "Login required")
				return
			}
			token := strings.TrimSpace(strings.TrimPrefix(header, prefix))
			if token == "" || sessions == nil {
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
