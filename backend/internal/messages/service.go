package messages

import (
	"context"
	"errors"
	"strings"
)

var ErrEmptyMessage = errors.New("message content is empty")
var ErrNotFound = errors.New("conversation not found")

type Message struct {
	ID              int64  `json:"id"`
	ConversationID  int64  `json:"conversationId"`
	SenderID        int64  `json:"senderId"`
	ContentMarkdown string `json:"contentMarkdown"`
	ContentPlain    string `json:"contentPlain"`
	CreatedAt       string `json:"createdAt"`
}

type SendInput struct {
	ConversationID  int64
	SenderID        int64
	ContentMarkdown string
}

type Repository interface {
	Create(ctx context.Context, input SendInput, contentPlain string) (Message, error)
	ListBefore(ctx context.Context, conversationID int64, userID int64, beforeID int64, limit int) ([]Message, error)
	MarkRead(ctx context.Context, conversationID int64, userID int64) error
}

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Send(ctx context.Context, input SendInput) (Message, error) {
	if strings.TrimSpace(input.ContentMarkdown) == "" {
		return Message{}, ErrEmptyMessage
	}
	return s.repo.Create(ctx, input, PlainTextFromMarkdown(input.ContentMarkdown))
}

func (s *Service) ListBefore(ctx context.Context, conversationID int64, userID int64, beforeID int64, limit int) ([]Message, error) {
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	return s.repo.ListBefore(ctx, conversationID, userID, beforeID, limit)
}

func (s *Service) MarkRead(ctx context.Context, conversationID int64, userID int64) error {
	return s.repo.MarkRead(ctx, conversationID, userID)
}
