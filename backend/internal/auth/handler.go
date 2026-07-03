package auth

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"time"

	"vibework-chat/backend/internal/httpx"
)

type SessionRepository struct {
	db     *sql.DB
	secret string
}

func NewSessionRepository(db *sql.DB, secret string) *SessionRepository {
	return &SessionRepository{db: db, secret: secret}
}

func (r *SessionRepository) Create(ctx context.Context, userID int64, tokenHash string, expiresAt time.Time) error {
	const query = `
		INSERT INTO sessions (user_id, token_hash, expires_at)
		VALUES (?, ?, ?)`
	_, err := r.db.ExecContext(ctx, query, userID, tokenHash, expiresAt)
	return err
}

func (r *SessionRepository) UserIDByToken(ctx context.Context, token string) (int64, error) {
	const query = `
		SELECT user_id
		FROM sessions
		WHERE token_hash = ? AND expires_at > ?
		LIMIT 1`
	var userID int64
	if err := r.db.QueryRowContext(ctx, query, HashToken(token, r.secret), time.Now()).Scan(&userID); err != nil {
		return 0, err
	}
	return userID, nil
}

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

type devLoginRequest struct {
	Username string `json:"username"`
}

func (h *Handler) DevLogin(w http.ResponseWriter, r *http.Request) {
	var req devLoginRequest
	if err := httpx.ReadJSON(r, &req); err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "invalid json")
		return
	}

	result, err := h.service.DevLogin(r.Context(), req.Username)
	if err != nil {
		if errors.Is(err, ErrInvalidCredentials) {
			httpx.WriteError(w, http.StatusUnauthorized, "invalid credentials")
			return
		}
		httpx.WriteError(w, http.StatusInternalServerError, "failed to login")
		return
	}

	httpx.WriteJSON(w, http.StatusOK, result)
}
