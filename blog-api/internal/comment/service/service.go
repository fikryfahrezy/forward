package service

import (
	"github.com/fikryfahrezy/forward/blog-api/internal/comment/repository"
)

type Service struct {
	repo *repository.Repository
}

func New(repo *repository.Repository) *Service {
	return &Service{repo: repo}
}
