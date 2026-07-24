package realtime

import (
	"context"
	"database/sql"

	"codezone/backend/internal/messages"
)

type SQLNotifier struct {
	db  *sql.DB
	hub *Hub
}

func NewSQLNotifier(db *sql.DB, hub *Hub) *SQLNotifier {
	return &SQLNotifier{db: db, hub: hub}
}

func (n *SQLNotifier) MessageCreated(ctx context.Context, message messages.Message) error {
	memberIDs, err := n.memberIDs(ctx, message.ConversationID)
	if err != nil {
		return err
	}
	for _, userID := range memberIDs {
		n.hub.SendToUser(userID, Event{
			Type:    "message.created",
			Payload: message,
		})
	}
	return nil
}

func (n *SQLNotifier) ConversationUpdated(ctx context.Context, conversationID int64) error {
	memberIDs, err := n.memberIDs(ctx, conversationID)
	if err != nil {
		return err
	}
	for _, userID := range memberIDs {
		n.hub.SendToUser(userID, Event{
			Type:    "conversation.updated",
			Payload: map[string]int64{"conversationId": conversationID},
		})
	}
	return nil
}

func (n *SQLNotifier) memberIDs(ctx context.Context, conversationID int64) ([]int64, error) {
	rows, err := n.db.QueryContext(ctx, `
		SELECT user_id
		FROM conversation_members
		WHERE conversation_id = ?
	`, conversationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var memberIDs []int64
	for rows.Next() {
		var userID int64
		if err := rows.Scan(&userID); err != nil {
			return nil, err
		}
		memberIDs = append(memberIDs, userID)
	}
	return memberIDs, rows.Err()
}
