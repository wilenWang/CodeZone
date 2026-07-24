package admin

import (
	"database/sql"
	"net/http"
	"strconv"

	"codezone/backend/internal/httpx"
)

type Handler struct {
	db *sql.DB
}

func NewHandler(db *sql.DB) *Handler {
	return &Handler{db: db}
}

func normalizeLimit(limit int) int {
	if limit <= 0 {
		return 50
	}
	if limit > 100 {
		return 100
	}
	return limit
}

func (h *Handler) Users(w http.ResponseWriter, r *http.Request) {
	rows, err := h.db.QueryContext(r.Context(), `
		SELECT id, username, display_name, user_type, created_at
		FROM users
		ORDER BY id
	`)
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, "admin_users_failed", "Could not load users")
		return
	}
	defer rows.Close()
	httpx.WriteJSON(w, http.StatusOK, map[string]any{"users": rowsToMaps(rows)})
}

func (h *Handler) Conversations(w http.ResponseWriter, r *http.Request) {
	rows, err := h.db.QueryContext(r.Context(), `
		SELECT id, type, title, last_message_at, created_at
		FROM conversations
		ORDER BY id DESC
		LIMIT 100
	`)
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, "admin_conversations_failed", "Could not load conversations")
		return
	}
	defer rows.Close()
	httpx.WriteJSON(w, http.StatusOK, map[string]any{"conversations": rowsToMaps(rows)})
}

func (h *Handler) Messages(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	rows, err := h.db.QueryContext(r.Context(), `
		SELECT id, conversation_id, sender_id, content_plain, created_at
		FROM messages
		ORDER BY id DESC
		LIMIT ?
	`, normalizeLimit(limit))
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, "admin_messages_failed", "Could not load messages")
		return
	}
	defer rows.Close()
	httpx.WriteJSON(w, http.StatusOK, map[string]any{"messages": rowsToMaps(rows)})
}

func rowsToMaps(rows *sql.Rows) []map[string]any {
	columns, err := rows.Columns()
	if err != nil {
		return []map[string]any{}
	}
	var result []map[string]any
	for rows.Next() {
		values := make([]any, len(columns))
		pointers := make([]any, len(columns))
		for i := range values {
			pointers[i] = &values[i]
		}
		if err := rows.Scan(pointers...); err != nil {
			return result
		}
		row := map[string]any{}
		for i, column := range columns {
			switch value := values[i].(type) {
			case []byte:
				row[column] = string(value)
			default:
				row[column] = value
			}
		}
		result = append(result, row)
	}
	return result
}
