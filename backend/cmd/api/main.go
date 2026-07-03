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
	authService := auth.NewService(authUserFinder{repo: userRepo}, sessionRepo, cfg.SessionSecret)

	authHandler := auth.NewHandler(authService)
	userHandler := users.NewHandler(userService)

	router := chi.NewRouter()
	router.Post("/api/auth/dev-login", authHandler.DevLogin)
	router.Get("/api/users", userHandler.List)

	addr := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("api listening on %s", addr)
	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatal(err)
	}
}

type authUserFinder struct {
	repo *users.Repository
}

func (f authUserFinder) FindByUsername(ctx context.Context, username string) (auth.User, error) {
	user, err := f.repo.FindByUsername(ctx, username)
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
