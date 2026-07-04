package agent

import (
	"context"
	"testing"
	"time"
)

func TestMockRunnerBuildsReply(t *testing.T) {
	runner := NewMockRunner("Mock Agent received:")
	reply := runner.BuildReply(context.Background(), "hello")
	want := "Mock Agent received: hello"
	if reply != want {
		t.Fatalf("got %q want %q", reply, want)
	}
}

type fakeFinder struct {
	agentUserID int64
	enabled     bool
	err         error
}

func (f fakeFinder) EnabledMockAgentForConversation(ctx context.Context, conversationID int64) (int64, bool, error) {
	return f.agentUserID, f.enabled, f.err
}

type fakeSender struct {
	called       chan struct{}
	conversation int64
	agentUserID  int64
	content      string
}

func (s *fakeSender) SendFromAgent(ctx context.Context, conversationID int64, agentUserID int64, contentMarkdown string) error {
	s.conversation = conversationID
	s.agentUserID = agentUserID
	s.content = contentMarkdown
	close(s.called)
	return nil
}

func TestOrchestratorRepliesWhenMockAgentEnabled(t *testing.T) {
	sender := &fakeSender{called: make(chan struct{})}
	orchestrator := NewOrchestrator(
		fakeFinder{agentUserID: 10, enabled: true},
		sender,
		NewMockRunner("Mock Agent received:"),
	)

	orchestrator.MaybeReply(context.Background(), 42, " hello ")

	select {
	case <-sender.called:
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for agent reply")
	}
	if sender.conversation != 42 {
		t.Fatalf("got conversation %d want 42", sender.conversation)
	}
	if sender.agentUserID != 10 {
		t.Fatalf("got agent user %d want 10", sender.agentUserID)
	}
	if sender.content != "Mock Agent received: hello" {
		t.Fatalf("got content %q", sender.content)
	}
}
