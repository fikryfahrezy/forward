package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/fikryfahrezy/forward/blog-api/internal/post"
)

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (post.PostWithAuthor, error) {
	p, err := s.repo.FindByIDWithAuthor(ctx, id)
	if err != nil {
		return post.PostWithAuthor{}, err
	}
	if p == (post.PostWithAuthor{}) {
		return post.PostWithAuthor{}, post.ErrPostNotFound
	}
	return p, nil
}
