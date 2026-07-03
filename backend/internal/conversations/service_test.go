package conversations

import (
	"context"
	"testing"
)

type fakeRepo struct {
	createdType string
	memberIDs   []int64
}

func (f *fakeRepo) Create(ctx context.Context, input CreateConversationInput) (Conversation, error) {
	f.createdType = input.Type
	f.memberIDs = input.MemberIDs
	return Conversation{ID: 42, Type: input.Type, Title: input.Title}, nil
}

func TestCreateGroupRequiresAtLeastThreeMembers(t *testing.T) {
	service := NewService(&fakeRepo{})
	_, err := service.Create(context.Background(), CreateConversationInput{
		WorkspaceID: 1,
		CreatedBy:   1,
		Type:        "group",
		Title:       ptr("Room"),
		MemberIDs:   []int64{1, 2},
	})
	if err != ErrGroupNeedsThreeMembers {
		t.Fatalf("got %v want %v", err, ErrGroupNeedsThreeMembers)
	}
}

func ptr(value string) *string {
	return &value
}
