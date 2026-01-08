package service

import (
	"context"

	"golang.org/x/crypto/bcrypt"

	"github.com/fikryfahrezy/forward/blog-api/internal/user"
)

func (s *Service) Login(ctx context.Context, req user.LoginRequest) (user.AuthResponse, error) {
	if err := req.Validate(); err != nil {
		return user.AuthResponse{}, err
	}

	u, err := s.repo.FindByEmail(ctx, req.Email)
	if err != nil {
		return user.AuthResponse{}, err
	}
	if u == (user.User{}) {
		return user.AuthResponse{}, user.ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(req.Password)); err != nil {
		return user.AuthResponse{}, user.ErrInvalidCredentials
	}

	token, err := s.jwtGenerator.GenerateToken(u.ID.String(), u.Username, u.Email)
	if err != nil {
		return user.AuthResponse{}, err
	}

	return user.AuthResponse{
		Token: token,
		User:  u,
	}, nil
}
