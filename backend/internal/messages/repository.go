package messages

import (
	"context"
	"database/sql"
)

type SQLRepository struct {
	db *sql.DB
}

func NewSQLRepository(db *sql.DB) *SQLRepository {
	return &SQLRepository{db: db}
}

func (r *SQLRepository) Create(ctx context.Context, input SendInput, contentPlain string) (Message, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return Message{}, err
	}
	defer tx.Rollback()

	if err := requireMember(ctx, tx, input.ConversationID, input.SenderID); err != nil {
		return Message{}, err
	}

	result, err := tx.ExecContext(ctx, `
		INSERT INTO messages (conversation_id, sender_id, content_markdown, content_plain)
		VALUES (?, ?, ?, ?)
	`, input.ConversationID, input.SenderID, input.ContentMarkdown, contentPlain)
	if err != nil {
		return Message{}, err
	}
	messageID, err := result.LastInsertId()
	if err != nil {
		return Message{}, err
	}
	var createdAt string
	if err := tx.QueryRowContext(ctx, `
		SELECT CAST(created_at AS CHAR)
		FROM messages
		WHERE id = ?
	`, messageID).Scan(&createdAt); err != nil {
		return Message{}, err
	}
	if _, err := tx.ExecContext(ctx, `
		UPDATE conversations
		SET last_message_id = ?, last_message_at = (SELECT created_at FROM messages WHERE id = ?)
		WHERE id = ?
	`, messageID, messageID, input.ConversationID); err != nil {
		return Message{}, err
	}
	if _, err := tx.ExecContext(ctx, `
		UPDATE conversation_members
		SET unread_count = unread_count + 1
		WHERE conversation_id = ? AND user_id <> ?
	`, input.ConversationID, input.SenderID); err != nil {
		return Message{}, err
	}
	if err := tx.Commit(); err != nil {
		return Message{}, err
	}
	return Message{
		ID:              messageID,
		ConversationID:  input.ConversationID,
		SenderID:        input.SenderID,
		ContentMarkdown: input.ContentMarkdown,
		ContentPlain:    contentPlain,
		CreatedAt:       createdAt,
	}, nil
}

func (r *SQLRepository) ListBefore(ctx context.Context, conversationID int64, userID int64, beforeID int64, limit int) ([]Message, error) {
	if err := requireMember(ctx, r.db, conversationID, userID); err != nil {
		return nil, err
	}
	rows, err := r.db.QueryContext(ctx, `
		SELECT m.id, m.conversation_id, m.sender_id, m.content_markdown, m.content_plain, CAST(m.created_at AS CHAR)
		FROM messages m
		JOIN conversation_members cm ON cm.conversation_id = m.conversation_id AND cm.user_id = ?
		WHERE m.conversation_id = ? AND (? = 0 OR m.id < ?)
		ORDER BY m.id DESC
		LIMIT ?
	`, userID, conversationID, beforeID, beforeID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []Message
	for rows.Next() {
		var message Message
		if err := rows.Scan(
			&message.ID,
			&message.ConversationID,
			&message.SenderID,
			&message.ContentMarkdown,
			&message.ContentPlain,
			&message.CreatedAt,
		); err != nil {
			return nil, err
		}
		out = append(out, message)
	}
	return out, rows.Err()
}

func (r *SQLRepository) MarkRead(ctx context.Context, conversationID int64, userID int64) error {
	result, err := r.db.ExecContext(ctx, `
		UPDATE conversation_members
		SET unread_count = 0,
		    last_read_message_id = (SELECT last_message_id FROM conversations WHERE id = ?)
		WHERE conversation_id = ? AND user_id = ?
	`, conversationID, conversationID, userID)
	if err != nil {
		return err
	}
	return errIfNoRowsAffected(result)
}

type memberChecker interface {
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

func requireMember(ctx context.Context, db memberChecker, conversationID int64, userID int64) error {
	var one int
	err := db.QueryRowContext(ctx, `
		SELECT 1
		FROM conversation_members
		WHERE conversation_id = ? AND user_id = ?
		LIMIT 1
	`, conversationID, userID).Scan(&one)
	if err == sql.ErrNoRows {
		return ErrNotFound
	}
	return err
}

func errIfNoRowsAffected(result sql.Result) error {
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}
