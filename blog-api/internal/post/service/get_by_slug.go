package service

import (
	"context"

	"github.com/fikryfahrezy/forward/blog-api/internal/post"
)

func (s *Service) GetBySlug(ctx context.Context, slug string) (post.PostItem, error) {
	p, err := s.repo.FindBySlugWithAuthor(ctx, slug)
	if err != nil {
		return post.PostItem{}, err
	}
	if p == (post.PostWithAuthor{}) {
		return post.PostItem{}, post.ErrPostNotFound
	}
	return p.ToPostItem(), nil
}
