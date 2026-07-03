package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"

	"vibework-chat/backend/internal/auth"
	"vibework-chat/backend/internal/config"
	"vibework-chat/backend/internal/db"
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

	router := buildRouter(cfg, authHandler.DevLogin, userHandler.List)

	addr := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("api listening on %s", addr)
	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatal(err)
	}
}

type authUserFinder struct {
	repo *users.Repository
}

func buildRouter(cfg config.Config, devLogin http.HandlerFunc, listUsers http.HandlerFunc) http.Handler {
	router := chi.NewRouter()
	if cfg.EnableDevLogin {
		router.Post("/api/auth/dev-login", devLogin)
	}
	router.Get("/api/users", listUsers)
	return router
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
