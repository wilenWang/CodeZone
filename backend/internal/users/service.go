package users

import "context"

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) List(ctx context.Context, workspaceID int64) ([]User, error) {
	return s.repo.List(ctx, workspaceID)
}

func (s *Service) FindByID(ctx context.Context, id int64) (User, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *Service) UpdateProfile(ctx context.Context, userID int64, displayName string, avatarURL *string) (User, error) {
	return s.repo.UpdateProfile(ctx, userID, displayName, avatarURL)
}
