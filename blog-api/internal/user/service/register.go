package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/fikryfahrezy/forward/blog-api/internal/user"
)

func (s *Service) Register(ctx context.Context, req user.RegisterRequest) (user.AuthResponse, error) {
	if err := req.Validate(); err != nil {
		return user.AuthResponse{}, err
	}

	existingUser, err := s.repo.FindByEmailOrUsername(ctx, req.Email, req.Username)
	if err != nil {
		return user.AuthResponse{}, err
	}
	if existingUser != (user.User{}) {
		if existingUser.Email == req.Email {
			return user.AuthResponse{}, user.ErrEmailExists
		}
		return user.AuthResponse{}, user.ErrUsernameExists
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return user.AuthResponse{}, err
	}

	now := time.Now()
	u := user.User{
		ID:        uuid.Must(uuid.NewV7()),
		Username:  req.Username,
		Email:     req.Email,
		Password:  string(hashedPassword),
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := s.repo.Create(ctx, u); err != nil {
		return user.AuthResponse{}, err
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
