package messages

import (
	"context"
	"errors"
	"testing"
)

func TestPlainTextFromMarkdown(t *testing.T) {
	input := "Hello **Bob**\n\n```go\nfmt.Println(\"x\")\n```"
	got := PlainTextFromMarkdown(input)
	want := "Hello Bob fmt.Println(\"x\")"
	if got != want {
		t.Fatalf("got %q want %q", got, want)
	}
}

type recordingRepo struct {
	createErr   error
	listErr     error
	listUserID  int64
	afterCreate func()
}

func (r *recordingRepo) Create(ctx context.Context, input SendInput, contentPlain string) (Message, error) {
	if r.createErr != nil {
		return Message{}, r.createErr
	}
	if r.afterCreate != nil {
		r.afterCreate()
	}
	return Message{
		ID:              1,
		ConversationID:  input.ConversationID,
		SenderID:        input.SenderID,
		ContentMarkdown: input.ContentMarkdown,
		ContentPlain:    contentPlain,
	}, nil
}

func (r *recordingRepo) ListBefore(ctx context.Context, conversationID int64, userID int64, beforeID int64, limit int) ([]Message, error) {
	r.listUserID = userID
	if r.listErr != nil {
		return nil, r.listErr
	}
	return []Message{}, nil
}

func (r *recordingRepo) MarkRead(ctx context.Context, conversationID int64, userID int64) error {
	return nil
}

func TestListBeforePassesUserID(t *testing.T) {
	repo := &recordingRepo{}
	service := NewService(repo)

	_, err := service.ListBefore(context.Background(), 42, 7, 0, 50)
	if err != nil {
		t.Fatalf("ListBefore returned error: %v", err)
	}
	if repo.listUserID != 7 {
		t.Fatalf("got userID %d want 7", repo.listUserID)
	}
}

func TestSendPropagatesNotFound(t *testing.T) {
	service := NewService(&recordingRepo{createErr: ErrNotFound})

	_, err := service.Send(context.Background(), SendInput{
		ConversationID:  42,
		SenderID:        7,
		ContentMarkdown: "hi",
	})

	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("got %v want ErrNotFound", err)
	}
}

type recordingNotifier struct {
	err                     error
	messageCreatedCalled    bool
	conversationUpdatedID   int64
	messageCreatedMessageID int64
	messageCreatedCtxErr    error
}

func (n *recordingNotifier) MessageCreated(ctx context.Context, message Message) error {
	n.messageCreatedCalled = true
	n.messageCreatedMessageID = message.ID
	n.messageCreatedCtxErr = ctx.Err()
	return n.err
}

func (n *recordingNotifier) ConversationUpdated(ctx context.Context, conversationID int64) error {
	n.conversationUpdatedID = conversationID
	return n.err
}

func TestSendNotifiesAfterSuccessfulCreateAndIgnoresNotifierErrors(t *testing.T) {
	notifier := &recordingNotifier{err: errors.New("websocket unavailable")}
	service := NewServiceWithNotifier(&recordingRepo{}, notifier)

	message, err := service.Send(context.Background(), SendInput{
		ConversationID:  42,
		SenderID:        7,
		ContentMarkdown: "hi",
	})
	if err != nil {
		t.Fatalf("Send returned error: %v", err)
	}
	if message.ID != 1 {
		t.Fatalf("got message ID %d want 1", message.ID)
	}
	if !notifier.messageCreatedCalled || notifier.messageCreatedMessageID != 1 {
		t.Fatalf("message notification not recorded correctly: %#v", notifier)
	}
	if notifier.conversationUpdatedID != 42 {
		t.Fatalf("got conversation update %d want 42", notifier.conversationUpdatedID)
	}
}

func TestSendUsesIndependentContextForPostCommitNotifications(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	notifier := &recordingNotifier{}
	service := NewServiceWithNotifier(&recordingRepo{afterCreate: cancel}, notifier)

	_, err := service.Send(ctx, SendInput{
		ConversationID:  42,
		SenderID:        7,
		ContentMarkdown: "hi",
	})
	if err != nil {
		t.Fatalf("Send returned error: %v", err)
	}
	if !notifier.messageCreatedCalled {
		t.Fatal("message notification was not recorded")
	}
	if notifier.messageCreatedCtxErr != nil {
		t.Fatalf("notification context was canceled: %v", notifier.messageCreatedCtxErr)
	}
}

type fakeResult struct {
	rowsAffected int64
}

func (r fakeResult) LastInsertId() (int64, error) {
	return 0, nil
}

func (r fakeResult) RowsAffected() (int64, error) {
	return r.rowsAffected, nil
}

func TestErrIfNoRowsAffected(t *testing.T) {
	if err := errIfNoRowsAffected(fakeResult{rowsAffected: 1}); err != nil {
		t.Fatalf("got %v want nil", err)
	}
	if err := errIfNoRowsAffected(fakeResult{rowsAffected: 0}); !errors.Is(err, ErrNotFound) {
		t.Fatalf("got %v want ErrNotFound", err)
	}
}
