package service

import (
	"github.com/fikryfahrezy/forward/blog-api/internal/server"
	"github.com/fikryfahrezy/forward/blog-api/internal/user/repository"
)

type Service struct {
	jwtGenerator *server.JWTGenerator
	repo         *repository.Repository
}

func New(
	jwtGenerator *server.JWTGenerator,
	repo *repository.Repository,
) *Service {
	return &Service{
		jwtGenerator: jwtGenerator,
		repo:         repo,
	}
}
