package conversations

import (
	"context"
	"errors"
	"testing"
)

type fakeRepo struct {
	createdType        string
	memberIDs          []int64
	workspaceUserCount int
}

func (f *fakeRepo) Create(ctx context.Context, input CreateConversationInput) (Conversation, error) {
	f.createdType = input.Type
	f.memberIDs = input.MemberIDs
	return Conversation{ID: 42, Type: input.Type, Title: input.Title}, nil
}

func (f *fakeRepo) CountUsersInWorkspace(ctx context.Context, workspaceID int64, memberIDs []int64) (int, error) {
	if f.workspaceUserCount != 0 {
		return f.workspaceUserCount, nil
	}
	return len(memberIDs), nil
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

func TestCreateRejectsInvalidType(t *testing.T) {
	service := NewService(&fakeRepo{})
	_, err := service.Create(context.Background(), CreateConversationInput{
		WorkspaceID: 1,
		CreatedBy:   1,
		Type:        "channel",
		MemberIDs:   []int64{1, 2},
	})
	if !errors.Is(err, ErrInvalidConversationType) {
		t.Fatalf("got %v want ErrInvalidConversationType", err)
	}
}

func TestCreateRejectsNonPositiveMemberID(t *testing.T) {
	service := NewService(&fakeRepo{})
	_, err := service.Create(context.Background(), CreateConversationInput{
		WorkspaceID: 1,
		CreatedBy:   1,
		Type:        "direct",
		MemberIDs:   []int64{1, 0},
	})
	if !errors.Is(err, ErrInvalidMemberID) {
		t.Fatalf("got %v want ErrInvalidMemberID", err)
	}
}

func TestCreateDirectRequiresExactlyTwoMembers(t *testing.T) {
	service := NewService(&fakeRepo{})
	_, err := service.Create(context.Background(), CreateConversationInput{
		WorkspaceID: 1,
		CreatedBy:   1,
		Type:        "direct",
		MemberIDs:   []int64{1, 2, 3},
	})
	if !errors.Is(err, ErrDirectNeedsTwoMembers) {
		t.Fatalf("got %v want ErrDirectNeedsTwoMembers", err)
	}
}

func TestCreateRejectsMembersOutsideWorkspace(t *testing.T) {
	service := NewService(&fakeRepo{workspaceUserCount: 1})
	_, err := service.Create(context.Background(), CreateConversationInput{
		WorkspaceID: 1,
		CreatedBy:   1,
		Type:        "direct",
		MemberIDs:   []int64{1, 2},
	})
	if !errors.Is(err, ErrMembersOutsideWorkspace) {
		t.Fatalf("got %v want ErrMembersOutsideWorkspace", err)
	}
}

func TestCreateDeduplicatesMembersAndIncludesCreatorOnce(t *testing.T) {
	repo := &fakeRepo{}
	service := NewService(repo)

	_, err := service.Create(context.Background(), CreateConversationInput{
		WorkspaceID: 1,
		CreatedBy:   1,
		Type:        "direct",
		MemberIDs:   []int64{2, 1, 2},
	})
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}
	want := []int64{1, 2}
	if len(repo.memberIDs) != len(want) {
		t.Fatalf("got members %v want %v", repo.memberIDs, want)
	}
	for i := range want {
		if repo.memberIDs[i] != want[i] {
			t.Fatalf("got members %v want %v", repo.memberIDs, want)
		}
	}
}

func ptr(value string) *string {
	return &value
}
