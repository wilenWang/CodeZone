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
	createErr  error
	listErr    error
	listUserID int64
}

func (r *recordingRepo) Create(ctx context.Context, input SendInput, contentPlain string) (Message, error) {
	if r.createErr != nil {
		return Message{}, r.createErr
	}
	return Message{ID: 1}, nil
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
