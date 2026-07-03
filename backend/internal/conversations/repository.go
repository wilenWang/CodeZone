package conversations

import (
	"context"
	"database/sql"
	"strings"
)

type SQLRepository struct {
	db *sql.DB
}

func NewSQLRepository(db *sql.DB) *SQLRepository {
	return &SQLRepository{db: db}
}

func (r *SQLRepository) Create(ctx context.Context, input CreateConversationInput) (Conversation, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return Conversation{}, err
	}
	defer tx.Rollback()

	result, err := tx.ExecContext(ctx, `
		INSERT INTO conversations (workspace_id, type, title, created_by)
		VALUES (?, ?, ?, ?)
	`, input.WorkspaceID, input.Type, input.Title, input.CreatedBy)
	if err != nil {
		return Conversation{}, err
	}
	conversationID, err := result.LastInsertId()
	if err != nil {
		return Conversation{}, err
	}
	for _, memberID := range input.MemberIDs {
		role := "member"
		if memberID == input.CreatedBy {
			role = "owner"
		}
		if _, err := tx.ExecContext(ctx, `
			INSERT INTO conversation_members (conversation_id, user_id, role)
			VALUES (?, ?, ?)
		`, conversationID, memberID, role); err != nil {
			return Conversation{}, err
		}
	}
	if err := tx.Commit(); err != nil {
		return Conversation{}, err
	}
	return Conversation{
		ID:          conversationID,
		WorkspaceID: input.WorkspaceID,
		Type:        input.Type,
		Title:       input.Title,
	}, nil
}

func (r *SQLRepository) ListForUser(ctx context.Context, workspaceID int64, userID int64) ([]Conversation, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT
		  c.id,
		  c.workspace_id,
		  c.type,
		  COALESCE(c.title, GROUP_CONCAT(CASE WHEN u.id <> ? THEN u.display_name END ORDER BY u.display_name SEPARATOR ', ')) AS title,
		  c.last_message_id,
		  CAST(c.last_message_at AS CHAR),
		  cm.unread_count
		FROM conversations c
		JOIN conversation_members cm ON cm.conversation_id = c.id AND cm.user_id = ?
		JOIN conversation_members all_cm ON all_cm.conversation_id = c.id
		JOIN users u ON u.id = all_cm.user_id
		WHERE c.workspace_id = ?
		GROUP BY c.id, c.workspace_id, c.type, c.title, c.last_message_id, c.last_message_at, cm.unread_count
		ORDER BY COALESCE(c.last_message_at, c.created_at) DESC
	`, userID, userID, workspaceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []Conversation
	for rows.Next() {
		var conversation Conversation
		var title sql.NullString
		var lastMessageID sql.NullInt64
		var lastMessageAt sql.NullString
		if err := rows.Scan(
			&conversation.ID,
			&conversation.WorkspaceID,
			&conversation.Type,
			&title,
			&lastMessageID,
			&lastMessageAt,
			&conversation.UnreadCount,
		); err != nil {
			return nil, err
		}
		if title.Valid {
			conversation.Title = &title.String
		}
		if lastMessageID.Valid {
			conversation.LastMessageID = &lastMessageID.Int64
		}
		if lastMessageAt.Valid {
			conversation.LastMessageAt = &lastMessageAt.String
		}
		out = append(out, conversation)
	}
	return out, rows.Err()
}

func (r *SQLRepository) CountUsersInWorkspace(ctx context.Context, workspaceID int64, memberIDs []int64) (int, error) {
	if len(memberIDs) == 0 {
		return 0, nil
	}
	placeholders := strings.TrimRight(strings.Repeat("?,", len(memberIDs)), ",")
	args := make([]any, 0, len(memberIDs)+1)
	args = append(args, workspaceID)
	for _, memberID := range memberIDs {
		args = append(args, memberID)
	}

	var count int
	err := r.db.QueryRowContext(ctx, `
		SELECT COUNT(DISTINCT id)
		FROM users
		WHERE workspace_id = ? AND id IN (`+placeholders+`)
	`, args...).Scan(&count)
	return count, err
}
