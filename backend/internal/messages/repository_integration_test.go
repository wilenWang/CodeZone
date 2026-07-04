package messages

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"testing"
	"time"

	"vibework-chat/backend/internal/db"
)

func TestSQLRepositoryIntegrationAuthorization(t *testing.T) {
	conn := openIntegrationDB(t)
	ctx := context.Background()
	repo := NewSQLRepository(conn)

	workspaceID, userAID, userBID, userCID, conversationID := createMessageRepositoryFixture(t, ctx, conn)
	t.Cleanup(func() {
		cleanupMessageRepositoryFixture(t, conn, workspaceID, []int64{userAID, userBID, userCID}, conversationID)
	})

	message, err := repo.Create(ctx, SendInput{
		ConversationID:  conversationID,
		SenderID:        userAID,
		ContentMarkdown: "hello from **member A**",
	}, "hello from member A")
	if err != nil {
		t.Fatalf("member Create returned error: %v", err)
	}
	if message.ConversationID != conversationID || message.SenderID != userAID {
		t.Fatalf("created message got conversation=%d sender=%d, want conversation=%d sender=%d", message.ConversationID, message.SenderID, conversationID, userAID)
	}

	assertUnreadCount(t, ctx, conn, conversationID, userAID, 0)
	assertUnreadCount(t, ctx, conn, conversationID, userBID, 1)

	messageCount := countMessages(t, ctx, conn, conversationID)
	_, err = repo.Create(ctx, SendInput{
		ConversationID:  conversationID,
		SenderID:        userCID,
		ContentMarkdown: "non-member attempt",
	}, "non-member attempt")
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("non-member Create got error %v, want ErrNotFound", err)
	}
	if got := countMessages(t, ctx, conn, conversationID); got != messageCount {
		t.Fatalf("non-member Create inserted message: got %d messages, want %d", got, messageCount)
	}

	messages, err := repo.ListBefore(ctx, conversationID, userBID, 0, 50)
	if err != nil {
		t.Fatalf("member ListBefore returned error: %v", err)
	}
	if len(messages) != 1 {
		t.Fatalf("member ListBefore returned %d messages, want 1", len(messages))
	}
	if messages[0].ID != message.ID || messages[0].ContentPlain != "hello from member A" {
		t.Fatalf("member ListBefore got message id=%d content=%q, want id=%d content=%q", messages[0].ID, messages[0].ContentPlain, message.ID, "hello from member A")
	}

	if _, err := repo.ListBefore(ctx, conversationID, userCID, 0, 50); !errors.Is(err, ErrNotFound) {
		t.Fatalf("non-member ListBefore got error %v, want ErrNotFound", err)
	}

	if err := repo.MarkRead(ctx, conversationID, userBID); err != nil {
		t.Fatalf("member MarkRead returned error: %v", err)
	}
	assertUnreadCount(t, ctx, conn, conversationID, userBID, 0)

	if err := repo.MarkRead(ctx, conversationID, userCID); !errors.Is(err, ErrNotFound) {
		t.Fatalf("non-member MarkRead got error %v, want ErrNotFound", err)
	}
}

func openIntegrationDB(t *testing.T) *sql.DB {
	t.Helper()

	dsn := os.Getenv("CHAT_TEST_MYSQL_DSN")
	if dsn == "" {
		dsn = os.Getenv("MYSQL_DSN")
	}
	if dsn == "" {
		t.Skip("skipping MySQL integration test; set CHAT_TEST_MYSQL_DSN or MYSQL_DSN")
	}

	conn, err := db.Open(dsn)
	if err != nil {
		t.Fatalf("open MySQL integration database: %v", err)
	}
	t.Cleanup(func() {
		if err := conn.Close(); err != nil {
			t.Errorf("close MySQL integration database: %v", err)
		}
	})
	return conn
}

func createMessageRepositoryFixture(t *testing.T, ctx context.Context, conn *sql.DB) (int64, int64, int64, int64, int64) {
	t.Helper()

	suffix := time.Now().UTC().Format("20060102150405.000000000")
	workspaceName := fmt.Sprintf("messages-it-%s", suffix)

	result, err := conn.ExecContext(ctx, `INSERT INTO workspaces (name) VALUES (?)`, workspaceName)
	if err != nil {
		t.Fatalf("insert workspace: %v", err)
	}
	workspaceID := lastInsertID(t, result, "workspace")

	userAID := insertMessageRepositoryUser(t, ctx, conn, workspaceID, "a", suffix)
	userBID := insertMessageRepositoryUser(t, ctx, conn, workspaceID, "b", suffix)
	userCID := insertMessageRepositoryUser(t, ctx, conn, workspaceID, "c", suffix)

	result, err = conn.ExecContext(ctx, `
		INSERT INTO conversations (workspace_id, type, title, created_by)
		VALUES (?, 'group', ?, ?)
	`, workspaceID, fmt.Sprintf("messages-it-%s", suffix), userAID)
	if err != nil {
		t.Fatalf("insert conversation: %v", err)
	}
	conversationID := lastInsertID(t, result, "conversation")

	for _, member := range []struct {
		userID int64
		role   string
	}{
		{userID: userAID, role: "owner"},
		{userID: userBID, role: "member"},
	} {
		if _, err := conn.ExecContext(ctx, `
			INSERT INTO conversation_members (conversation_id, user_id, role)
			VALUES (?, ?, ?)
		`, conversationID, member.userID, member.role); err != nil {
			t.Fatalf("insert conversation member %d: %v", member.userID, err)
		}
	}

	return workspaceID, userAID, userBID, userCID, conversationID
}

func insertMessageRepositoryUser(t *testing.T, ctx context.Context, conn *sql.DB, workspaceID int64, label string, suffix string) int64 {
	t.Helper()

	username := fmt.Sprintf("messages_it_%s_%s", label, suffix)
	result, err := conn.ExecContext(ctx, `
		INSERT INTO users (workspace_id, username, display_name, user_type)
		VALUES (?, ?, ?, 'human')
	`, workspaceID, username, fmt.Sprintf("Messages IT %s", label))
	if err != nil {
		t.Fatalf("insert user %s: %v", label, err)
	}
	return lastInsertID(t, result, "user")
}

func cleanupMessageRepositoryFixture(t *testing.T, conn *sql.DB, workspaceID int64, userIDs []int64, conversationID int64) {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err := conn.ExecContext(ctx, `
		UPDATE conversations
		SET last_message_id = NULL, last_message_at = NULL
		WHERE id = ?
	`, conversationID); err != nil {
		t.Errorf("cleanup conversation last message: %v", err)
	}
	if _, err := conn.ExecContext(ctx, `DELETE FROM conversation_members WHERE conversation_id = ?`, conversationID); err != nil {
		t.Errorf("cleanup conversation members: %v", err)
	}
	if _, err := conn.ExecContext(ctx, `DELETE FROM messages WHERE conversation_id = ?`, conversationID); err != nil {
		t.Errorf("cleanup messages: %v", err)
	}
	if _, err := conn.ExecContext(ctx, `DELETE FROM conversations WHERE id = ?`, conversationID); err != nil {
		t.Errorf("cleanup conversation: %v", err)
	}
	for _, userID := range userIDs {
		if _, err := conn.ExecContext(ctx, `DELETE FROM users WHERE id = ?`, userID); err != nil {
			t.Errorf("cleanup user %d: %v", userID, err)
		}
	}
	if _, err := conn.ExecContext(ctx, `DELETE FROM workspaces WHERE id = ?`, workspaceID); err != nil {
		t.Errorf("cleanup workspace: %v", err)
	}
}

func assertUnreadCount(t *testing.T, ctx context.Context, conn *sql.DB, conversationID int64, userID int64, want int) {
	t.Helper()

	var got int
	if err := conn.QueryRowContext(ctx, `
		SELECT unread_count
		FROM conversation_members
		WHERE conversation_id = ? AND user_id = ?
	`, conversationID, userID).Scan(&got); err != nil {
		t.Fatalf("query unread count for user %d: %v", userID, err)
	}
	if got != want {
		t.Fatalf("unread count for user %d got %d, want %d", userID, got, want)
	}
}

func countMessages(t *testing.T, ctx context.Context, conn *sql.DB, conversationID int64) int {
	t.Helper()

	var count int
	if err := conn.QueryRowContext(ctx, `
		SELECT COUNT(*)
		FROM messages
		WHERE conversation_id = ?
	`, conversationID).Scan(&count); err != nil {
		t.Fatalf("count messages: %v", err)
	}
	return count
}

func lastInsertID(t *testing.T, result sql.Result, label string) int64 {
	t.Helper()

	id, err := result.LastInsertId()
	if err != nil {
		t.Fatalf("read %s insert id: %v", label, err)
	}
	return id
}
