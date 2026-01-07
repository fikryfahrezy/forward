package service

import (
	"context"

	"golang.org/x/crypto/bcrypt"

	"github.com/fikryfahrezy/forward/blog-api/internal/user"
)

func (s *Service) Login(ctx context.Context, req user.LoginRequest) (user.User, error) {
	if err := req.Validate(); err != nil {
		return user.User{}, err
	}

	u, err := s.repo.FindByEmail(ctx, req.Email)
	if err != nil {
		return user.User{}, err
	}
	if u == (user.User{}) {
		return user.User{}, user.ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(req.Password)); err != nil {
		return user.User{}, user.ErrInvalidCredentials
	}

	return u, nil
}
