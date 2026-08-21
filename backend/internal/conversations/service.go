package conversations

import (
	"context"
	"errors"
	"fmt"
)

var ErrGroupNeedsThreeMembers = errors.New("group conversations need at least three members")
var ErrInvalidConversationType = errors.New("conversation type must be direct or group")
var ErrInvalidMemberID = errors.New("member ids must be positive")
var ErrDirectNeedsTwoMembers = errors.New("direct conversations need exactly two members")
var ErrMembersOutsideWorkspace = errors.New("all members must belong to the workspace")
var ErrDirectSelf = errors.New("direct conversation requires a different user")

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
	ListForUser(ctx context.Context, workspaceID int64, userID int64) ([]Conversation, error)
	CountUsersInWorkspace(ctx context.Context, workspaceID int64, memberIDs []int64) (int, error)
	EnsureDirect(ctx context.Context, workspaceID int64, createdBy int64, memberIDs []int64, pairKey string) (Conversation, error)
}

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, input CreateConversationInput) (Conversation, error) {
	memberIDs, err := uniqueMemberIDs(input.CreatedBy, input.MemberIDs)
	if err != nil {
		return Conversation{}, err
	}
	input.MemberIDs = memberIDs
	switch input.Type {
	case "direct":
		if len(input.MemberIDs) != 2 {
			return Conversation{}, ErrDirectNeedsTwoMembers
		}
	case "group":
		if len(input.MemberIDs) < 3 {
			return Conversation{}, ErrGroupNeedsThreeMembers
		}
	default:
		return Conversation{}, ErrInvalidConversationType
	}

	count, err := s.repo.CountUsersInWorkspace(ctx, input.WorkspaceID, input.MemberIDs)
	if err != nil {
		return Conversation{}, err
	}
	if count != len(input.MemberIDs) {
		return Conversation{}, ErrMembersOutsideWorkspace
	}
	return s.repo.Create(ctx, input)
}

func (s *Service) ListForUser(ctx context.Context, workspaceID int64, userID int64) ([]Conversation, error) {
	return s.repo.ListForUser(ctx, workspaceID, userID)
}

func (s *Service) EnsureDirect(ctx context.Context, workspaceID int64, userID int64, targetUserID int64) (Conversation, error) {
	if userID <= 0 || targetUserID <= 0 {
		return Conversation{}, ErrInvalidMemberID
	}
	if userID == targetUserID {
		return Conversation{}, ErrDirectSelf
	}
	createdBy := userID
	memberIDs := []int64{userID, targetUserID}
	count, err := s.repo.CountUsersInWorkspace(ctx, workspaceID, memberIDs)
	if err != nil {
		return Conversation{}, err
	}
	if count != len(memberIDs) {
		return Conversation{}, ErrMembersOutsideWorkspace
	}
	if userID > targetUserID {
		userID, targetUserID = targetUserID, userID
	}
	pairKey := fmt.Sprintf("%d:%d", userID, targetUserID)
	return s.repo.EnsureDirect(ctx, workspaceID, createdBy, memberIDs, pairKey)
}

func uniqueMemberIDs(createdBy int64, memberIDs []int64) ([]int64, error) {
	if createdBy <= 0 {
		return nil, ErrInvalidMemberID
	}
	seen := make(map[int64]bool, len(memberIDs)+1)
	out := make([]int64, 0, len(memberIDs)+1)
	candidates := append([]int64{createdBy}, memberIDs...)
	for _, memberID := range candidates {
		if memberID <= 0 {
			return nil, ErrInvalidMemberID
		}
		if seen[memberID] {
			continue
		}
		seen[memberID] = true
		out = append(out, memberID)
	}
	return out, nil
}

func isValidationError(err error) bool {
	return errors.Is(err, ErrInvalidConversationType) ||
		errors.Is(err, ErrInvalidMemberID) ||
		errors.Is(err, ErrDirectNeedsTwoMembers) ||
		errors.Is(err, ErrGroupNeedsThreeMembers) ||
		errors.Is(err, ErrMembersOutsideWorkspace) ||
		errors.Is(err, ErrDirectSelf)
}
