package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/fikryfahrezy/forward/blog-api/internal/user"
)

func (s *Service) Register(ctx context.Context, req user.RegisterRequest) (user.AuthResponse, error) {
	// Validate input
	if req.Email == "" || req.Username == "" || req.Password == "" {
		return user.AuthResponse{}, user.ErrInvalidInput
	}

	// Check if email exists
	emailExists, err := s.repo.ExistsByEmail(ctx, req.Email)
	if err != nil {
		return user.AuthResponse{}, err
	}
	if emailExists {
		return user.AuthResponse{}, user.ErrEmailExists
	}

	// Check if username exists
	usernameExists, err := s.repo.ExistsByUsername(ctx, req.Username)
	if err != nil {
		return user.AuthResponse{}, err
	}
	if usernameExists {
		return user.AuthResponse{}, user.ErrUsernameExists
	}

	// Hash password
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
