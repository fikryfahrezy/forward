package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/fikryfahrezy/forward/blog-api/internal/user"
)

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (user.UserProfile, error) {
	u, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return user.UserProfile{}, err
	}
	if u == (user.User{}) {
		return user.UserProfile{}, user.ErrUserNotFound
	}
	return u.ToUserProfile(), nil
}
