package agent

import (
	"context"
	"database/sql"
	"log"
	"strings"
	"time"
)

const replyTimeout = 5 * time.Second

type MockRunner struct {
	prefix string
}

func NewMockRunner(prefix string) *MockRunner {
	return &MockRunner{prefix: prefix}
}

func (r *MockRunner) BuildReply(ctx context.Context, input string) string {
	trimmed := strings.TrimSpace(input)
	if trimmed == "" {
		return r.prefix + " empty message"
	}
	return r.prefix + " " + trimmed
}

type MessageSender interface {
	SendFromAgent(ctx context.Context, conversationID int64, agentUserID int64, contentMarkdown string) error
}

type ConversationAgentFinder interface {
	EnabledMockAgentForConversation(ctx context.Context, conversationID int64) (agentUserID int64, enabled bool, err error)
}

type Orchestrator struct {
	finder ConversationAgentFinder
	sender MessageSender
	runner *MockRunner
}

func NewOrchestrator(finder ConversationAgentFinder, sender MessageSender, runner *MockRunner) *Orchestrator {
	return &Orchestrator{finder: finder, sender: sender, runner: runner}
}

func (o *Orchestrator) MaybeReply(ctx context.Context, conversationID int64, humanMessage string) {
	if o == nil || o.finder == nil || o.sender == nil || o.runner == nil {
		return
	}

	go func() {
		replyCtx, cancel := context.WithTimeout(context.Background(), replyTimeout)
		defer cancel()

		agentUserID, enabled, err := o.finder.EnabledMockAgentForConversation(replyCtx, conversationID)
		if err != nil {
			log.Printf("mock agent lookup failed for conversation %d: %v", conversationID, err)
			return
		}
		if !enabled {
			return
		}
		reply := o.runner.BuildReply(replyCtx, humanMessage)
		if err := o.sender.SendFromAgent(replyCtx, conversationID, agentUserID, reply); err != nil {
			log.Printf("mock agent reply failed for conversation %d: %v", conversationID, err)
		}
	}()
}

type SQLFinder struct {
	db *sql.DB
}

func NewSQLFinder(db *sql.DB) *SQLFinder {
	return &SQLFinder{db: db}
}

func (f *SQLFinder) EnabledMockAgentForConversation(ctx context.Context, conversationID int64) (int64, bool, error) {
	var agentUserID int64
	err := f.db.QueryRowContext(ctx, `
		SELECT u.id
		FROM conversation_members cm
		JOIN users u ON u.id = cm.user_id
		JOIN agent_profiles ap ON ap.user_id = u.id
		WHERE cm.conversation_id = ?
		  AND u.user_type = 'agent'
		  AND ap.kind = 'mock'
		  AND ap.enabled = TRUE
		LIMIT 1
	`, conversationID).Scan(&agentUserID)
	if err == sql.ErrNoRows {
		return 0, false, nil
	}
	if err != nil {
		return 0, false, err
	}
	return agentUserID, true, nil
}
