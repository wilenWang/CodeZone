package conversations

import (
	"context"
	"errors"
)

var ErrGroupNeedsThreeMembers = errors.New("group conversations need at least three members")

type Conversation struct {
	ID            int64   `json:"id"`
	WorkspaceID   int64   `json:"workspaceId"`
	Type          string  `json:"type"`
	Title         *string `json:"title"`
	LastMessageID *int64  `json:"lastMessageId"`
	LastMessageAt *string `json:"lastMessageAt"`
	UnreadCount   int     `json:"unreadCount"`
}

type CreateConversationInput struct {
	WorkspaceID int64
	CreatedBy   int64
	Type        string
	Title       *string
	MemberIDs   []int64
}

type Repository interface {
	Create(ctx context.Context, input CreateConversationInput) (Conversation, error)
}

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, input CreateConversationInput) (Conversation, error) {
	input.MemberIDs = uniqueMemberIDs(input.CreatedBy, input.MemberIDs)
	if input.Type == "group" && len(input.MemberIDs) < 3 {
		return Conversation{}, ErrGroupNeedsThreeMembers
	}
	return s.repo.Create(ctx, input)
}

func uniqueMemberIDs(createdBy int64, memberIDs []int64) []int64 {
	seen := make(map[int64]bool, len(memberIDs)+1)
	out := make([]int64, 0, len(memberIDs)+1)
	for _, memberID := range append([]int64{createdBy}, memberIDs...) {
		if memberID == 0 || seen[memberID] {
			continue
		}
		seen[memberID] = true
		out = append(out, memberID)
	}
	return out
}
